// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	vaautils "github.com/wormhole-foundation/wormhole/sdk/vaa"
)

type Guardian struct {
	PrivKey *ecdsa.PrivateKey
	Address common.Address
}

func GuardianSigner() Guardian {
	// generate private key
	privateKey, _ := ecdsa.GenerateKey(ethcrypto.S256(), rand.Reader)
	addr := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	return Guardian{
		PrivKey: privateKey,
		Address: addr,
	}
}

func CreateVAA(t *testing.T, guardians []Guardian, payload string, sequence uint64) *vaautils.VAA {
	vaa := vaautils.VAA{
		Payload:  []byte(payload),
		Sequence: sequence,
		Version:  vaautils.SupportedVAAVersion,
	}

	var guardianAddresses []common.Address
	for idx, guardian := range guardians {
		vaa.AddSignature(guardian.PrivKey, uint8(idx))
		guardianAddresses = append(guardianAddresses, guardian.Address)
	}

	// verify signatures
	err := vaa.Verify(guardianAddresses)
	if err != nil {
		t.Errorf("verify failed: %s", err)
	}
	return &vaa
}
