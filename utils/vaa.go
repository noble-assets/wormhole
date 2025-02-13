package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
)

func testSigner() (*ecdsa.PrivateKey, common.Address) {
	// generate private key
	privateKey, _ := ecdsa.GenerateKey(ethcrypto.S256(), rand.Reader)

	return privateKey, ethcrypto.PubkeyToAddress(privateKey.PublicKey)
}

func CreateVAA(t *testing.T, payload string, sequence uint64) *vaautils.VAA {
	g1Pk, g1Addr := testSigner()
	g2Pk, g2Addr := testSigner()
	g3Pk, g3Addr := testSigner()
	_, g4Addr := testSigner()

	guardianAddresses := []common.Address{g1Addr, g2Addr, g3Addr, g4Addr}

	vaa := vaautils.VAA{
		Payload:  []byte(payload),
		Sequence: sequence,
	}

	vaa.AddSignature(g1Pk, 0)
	vaa.AddSignature(g2Pk, 1)
	vaa.AddSignature(g3Pk, 2)

	// verify signatures
	err := vaa.Verify(guardianAddresses)
	if err != nil {
		t.Errorf("verify failed: %s", err)
	}
	return &vaa
}
