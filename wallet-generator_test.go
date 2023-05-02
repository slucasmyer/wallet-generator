package main

import (
	"fmt"
	"strings"
	"testing"
)


func TestGenerateMnemonic(t *testing.T) {

	userEntropy := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"
	mnemonic, err := generateMnemonic(24, userEntropy)
	if err != nil {
			t.Fatalf("Error generating mnemonic: %v", err)
	}

	words := strings.Split(mnemonic, " ")


	if len(words) != 24 {
			t.Errorf("Mnemonic should have 24 words, got %d", len(words))
	}
}

func TestMnemonicToSeed(t *testing.T) {

type test struct {
	name string
	passphrase string
	mnemonic string
	expected string
}

tests := []test{
	{
		"test1",
		"",
		"accuse differ elegant crop bonus civil broccoli mad become blouse galaxy accident leader simple crack decide airport road address deputy there candy knife fantasy",
		"2cf49b56ab67b80bdf996dc5dc53c869a94387efb188a2724db9ecc3fe720221910199e0d3a88751bd214ff19443e4389de876de01dd97c36a01d85e139f5421",
	},
	{
		"test2",
		"",
		"army hedgehog winner buddy base hover gauge pond bundle doll process cruise adjust avoid domain light away dismiss panic cattle put feel glory perfect",
		"e06f9990dcbf527fe30ba4ad95d6ccdf8c330545a5c1ec1c22afbc9c0173c8a",
	},
	{
		"test3",
		"satoshi",
		"biology wheat faculty blush uncle adult badge chimney bracket group debate crouch tonight ketchup agent jacket believe arctic leaf alert name humor cable tree",
		"7b5e5c04223a5e0fcc59887eedfbf3899a3364c3ead1db56fd491ccc1b0d8f3f435f4643fa78e1c4cae26a020a72458cac1c572c9677570078128faf28eb2c3d",
	},
	{
		"test4",
		"satoshi",
		"bread toward glass crop omit sheriff arctic surge black identify hurt apology scan since explain clay cabbage salon three angry junk shock blame slender",
		"49d5a47c707cbfb96198592cab721667f10732e106da08b4e22700ed28670ba2eb19dc7227fd63797bd14ce85a761cdfffc4bfa296f3085282e487d40db062b5",
	},
}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			seed := mnemonicToSeed(test.mnemonic, test.passphrase)
			if fmt.Sprintf("%x", seed) != test.expected {
					t.Errorf("mnemonicToSeed produced an incorrect seed. Expected %s, got %x", test.expected, seed)
			}
			
		})
	}
}
