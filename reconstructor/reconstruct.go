package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/hashicorp/vault/shamir"
)


func main() {
	// Set the number of shares required to reconstruct the sensitive data
	threshold := 3
	shares := []int{1, 2, 3, 4, 5}
	shareBytes := make([][]byte, threshold)
	// loop through shares removing a random element until length of shares == threshold
	for len(shares) > threshold {
		i := rand.Intn(len(shares))
		shares = append(shares[:i], shares[i+1:]...)
	}

	fmt.Println("shares used to reconstruct:", shares)
	// Read shares from files
	for i, v := range shares {
		filename := fmt.Sprintf("share_%d", v)
		filePath := fmt.Sprintf("shares/%s", filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading share from file:", err)
			os.Exit(1)
		}
		shareBytes[i] = data
	}

	// Combine the shares to reconstruct the sensitive data
	sensitiveDataJSON, err := shamir.Combine(shareBytes)
	if err != nil {
		fmt.Println("Error combining shares:", err)
		os.Exit(1)
	}

	// Parse the JSON to extract the mnemonic and private key
	var sensitiveData map[string]string
	err = json.Unmarshal(sensitiveDataJSON, &sensitiveData)
	if err != nil {
		fmt.Println("Error unmarshaling sensitive data:", err)
		os.Exit(1)
	}

	// Print the reconstructed sensitive data
	result := map[string]string{}
	for key, value := range sensitiveData {
		result[key] = value
	}
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Reconstructed sensitive data:\n", string(resultJSON))
}
