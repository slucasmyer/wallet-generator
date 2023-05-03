# Sommelier Wallet Generator

A command-line tool that generates a new Sommelier.finance wallet, including a BIP39 mnemonic phrase, private key, public key, and address. Written in Go and currently uses the Cosmos SDK and btcsuite libraries.

## Prerequisites

To build and run the wallet generator securely, you'll need:

1. 2 computers: one internet-connected, and one air-gapped, both running identical Ubunutu Distros - *preferably ubuntu > 18.04* - with [Go](https://golang.org/dl/) installed on the internet-connected machine.
2. A sercure and encrypted transfer device, such as a USB drive.

## Usage

Foloo the steps below to generate a new wallet and securely store mnemonic and private key:

1. Prepare the internet-connected environment:
   1. On a clean, trusted, and internet-connected computer, running the same OS as our air-gapped machine - *preferably ubuntu > 18.04* - Clone the repository onto the internet-connected computer:  
    `$ git clone git@github.com:slucasmyer/wallet-generator.git`
   3. Verify the integrity of the software being installed on the internet-connected computer (check digital signatures, commit hashes, etc).
2. Build the wallet generator:
   1. Download the necessary dependencies, configure environment variables, and build the wallet generator:  
    `$ make all`
   2. Verify the integrity of the software being installed on the internet-connected computer (check digital signatures, commit hashes, etc).
3. Prepare the transfer device:
   1. Connect a clean, trusted, and preferably encrypted USB drive (or other hard-drive) to the internet-connected computer.
   2. Ensure that the USB drive does not contain any malicious software or files that could compromise the security of the air-gapped machine.
4. Prepare the air-gapped environment:
   1. Ensure that the computer being used to generate the keys is not connected to any network and has no means of connecting to the internet. All external devices or connections should be removed, and the computer should be running a clean operating system - *preferably ubuntu > 18.04*.
5. Transfer binaries to air-gapped environment:
   1. Move the USB drive to the air-gapped computer, ensuring that it remains secure during transportation. Connect the USB drive to the air-gapped machine and copy the source-code and compiled binaries to the machine's local storage.
   2. Verify the integrity of the software being installed on the air-gapped computer (check digital signatures, commit hashes, etc).
6. Generate a new wallet:
   1. Run the wallet generator:  
    `$ ./build/wallet-generator`
   2. Follow the prompts to generate a new wallet
   3. The program will output both to the CLI and to 5 encrypted key shares, any 3 of which are sufficient to reconstruct the original. To verify this, use the accompanying reconstruct tool:  
    `$ ./build/reconstruct`  
    then compare the output to the original mnemonic phrase and private key.
   4. Securely store the consituent shares of mnemonic phrase and private key in geogrpahically distributed safe places, preferably across multiple modalities (e.g. multi-cloud). The mnemonic phrase is required to recover your wallet, and the private key is required to sign transactions.

Example output:
```json
{
  "mnemonic": "artefact palace disorder ...",
  "privateKey": "base64_encoded_private_key",
  "publicKey": "base64_encoded_public_key",
  "address": "somm1abcd1234..."
}
```
## Notes

- The generated mnemonic phrase is 24 words long. Keep it safe, as it is required to recover your wallet.
- The tool derives the wallet using the default HD path for the Cosmos SDK.
- The generated address has a `somm` prefix, which is specific to the Sommelier network.

## Future Improvements

- Add support for BIP44 paths
- Add support for cross-compilation
- Automate the process of distributing and storing key shares
- Use a more secure method of generating entropy (YubiKey, etc.)
- Use proprietary library to generate key shares

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements or bug fixes.

## License

This project is open-source software, licensed under the [MIT License](LICENSE).