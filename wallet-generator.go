package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha512"
	_ "embed"
	"encoding/binary"
  "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
  "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hashicorp/vault/shamir"
	"golang.org/x/crypto/pbkdf2"
)

//go:embed wordList.txt
var wordList string

// setPrefix sets the bech32 prefix for the given network
func setPrefix(prefixAccAddr, prefixAccPub string) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(prefixAccAddr, prefixAccPub)
	config.Seal()
}

// getUserEntropy reads user input and returns a byte slice of the same length as the generated entropy.
func getUserEntropy() (string, error) {
	var userEntropy string
	fmt.Print("Please enter at least 33 random characters: ")
	_, err := fmt.Scan(&userEntropy)
	if err != nil {
		return "", err
	}
	if len(userEntropy) < 33 {
		return "", fmt.Errorf("user input must be at least 33 characters long")
	}
	return userEntropy, nil
}

// getUserPassphrase reads user input and returns a string.
func getUserPassphrase() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your passphrase (optional): ")
	passphrase, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error capturing passphrase:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(passphrase)
}

// generateMnemonic generates a BIP39 mnemonic phrase using the given number of words.
func generateMnemonic(numWords int, userEntropy string) (string, error) {
	words := strings.Split(wordList, "\n")

	if numWords%3 != 0 {
		return "", fmt.Errorf("number of words must be a multiple of 3")
	}

	entropy := make([]byte, numWords*11/8)
	_, err := rand.Read(entropy)
	if err != nil {
		return "", err
	}

	// Hash the user-provided entropy using SHA-512
	hash := sha512.New()
	hash.Write([]byte(userEntropy))
	userEntropyHash := hash.Sum(nil)

	// Combine the generated entropy with the hashed user-provided entropy using XOR
	for i := 0; i < len(entropy); i++ {
		entropy[i] ^= userEntropyHash[i % len(userEntropyHash)]
	}

	mnemonic := make([]string, numWords)
	for i := 0; i < numWords; i++ {
		startByte := (i * 11) / 8
		endByte := ((i+1) * 11) / 8
		if endByte == len(entropy) {
			endByte--
		}
		bits := (uint(binary.BigEndian.Uint16(entropy[startByte:endByte+1])) >> (8 - ((i * 11) % 8))) & 0x7FF
		mnemonic[i] = words[bits]
	}

	return strings.Join(mnemonic, " "), nil
}

func generateMnemonicWithUserInput(numWords int) string {
	userEntropy, err := getUserEntropy()
	if err != nil {
		fmt.Println("Error capturing user entropy:", err)
		os.Exit(1)
	}
	mnemonic, err := generateMnemonic(numWords, userEntropy)
	if err != nil {
		fmt.Println("Error generating mnemonic:", err)
		os.Exit(1)
	}
	return mnemonic
}

// mnemonicToSeed converts a BIP39 mnemonic phrase to a seed
func mnemonicToSeed(mnemonic, passphrase string) []byte {
	return pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"+passphrase), 4096, 64, sha512.New)
}

// derivePriv derives the private key and other information from the seed
func derivePriv(seed []byte) []byte {
	hdPath := hd.NewFundraiserParams(0, sdk.CoinType, 0)
	masterKey, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterKey, ch, hdPath.String())
	if err != nil {
		fmt.Println("Error deriving private key:", err)
		os.Exit(1)
	}
	return derivedPriv
}

// convertToPub converts the private key to a public key
func convertToPub(derivedPriv []byte) []byte {
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), derivedPriv)
	pubKey := privKey.PubKey()
	return pubKey.SerializeCompressed()
}

// keyToAddress converts the public key to an address
func keyToAddress(pubKeyBytes []byte) (secp256k1.PubKey, sdk.AccAddress) {
	pubKeyObj, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		fmt.Println("Error creating public key object:", err)
		os.Exit(1)
	}
	cosmosPubKey := secp256k1.PubKey{Key: pubKeyObj.SerializeCompressed()}
	addr := sdk.AccAddress(cosmosPubKey.Address().Bytes())
	return cosmosPubKey, addr
}

// func payloadToStrings converts the private key, public key, and address to strings
func payloadToStrings(derivedPriv, derivedPub []byte, addr sdk.AccAddress) (string, string, string) {
	derivedPrivString := base64.StdEncoding.EncodeToString(derivedPriv)
	pubKeyString := base64.StdEncoding.EncodeToString(derivedPub)
	addrString := addr.String()
	return derivedPrivString, pubKeyString, addrString
}

// createShares uses Shamir's Secret Sharing to split the private key into 5 shares (threshold 3)
func createShares(mnemonic string, privKeyStr string) [][]byte {
	// Prepare the sensitive data
	sensitiveData := map[string]string{
		"mnemonic":   mnemonic,
		"privateKey": privKeyStr,
	}

	// Convert the sensitive data to JSON
	sensitiveDataJSON, _ := json.Marshal(sensitiveData)

	// Split the sensitive data using Shamir's Secret Sharing
	numShares := 5
	threshold := 3
	shares, err := shamir.Split(sensitiveDataJSON, numShares, threshold)
	if err != nil {
		fmt.Println("Error splitting sensitive data:", err)
		os.Exit(1)
	}
	return shares
}

// storeShares stores the secret shares in files
func storeShares(secretShares [][]byte) {
	for i, share := range secretShares {
		filename := fmt.Sprintf("share_%d", i+1)
		filePath := fmt.Sprintf("shares/%s", filename)
		err := os.WriteFile(filePath, share, 0644)
		if err != nil {
			fmt.Println("Error writing share to file:", err)
			os.Exit(1)
		}
	}
}

// display displays the mnemonic, private key, public key, and address
func display(mnemonic, privKeyStr, pubKeyStr, addrStr string) {
	result := map[string]string{
		"mnemonic":   mnemonic,
		"privateKey": privKeyStr,
		"publicKey":  pubKeyStr,
		"address":    addrStr,
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(resultJSON))
}

func main() {
	
	// Set the Bech32 prefix to "somm"
	setPrefix("somm", "sommpub")

	// Generate a new BIP39 mnemonic phrase
	mnemonic := generateMnemonicWithUserInput(24)

	// Get the user's passphrase
	passphrase := getUserPassphrase()

	// Convert the mnemonic to a seed
	seed := mnemonicToSeed(mnemonic, passphrase)

	// Derive the private key and other information from the seed
	derivedPriv := derivePriv(seed)
	
	// Convert the private key to a public key
	pubKeyBytes := convertToPub(derivedPriv)

	// Convert the public key to an address
	cosmosPubKey, addr := keyToAddress(pubKeyBytes)

	// Encode the private key, public key, and address as strings
	privKeyStr, pubKeyStr, addrStr := payloadToStrings(derivedPriv, cosmosPubKey.Key, addr)

	// Split the sensitive data using Shamir's Secret Sharing
	secretShares := createShares(mnemonic, privKeyStr)

	// Save the shares to separate files
	storeShares(secretShares)

	// Print the mnemonic, private key, public key, and address as JSON
	display(mnemonic, privKeyStr, pubKeyStr, addrStr)

}