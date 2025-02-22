// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package mimc

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
)

func TestDecompose(t *testing.T) {

	// create 10 random digits in basis r
	nbDigits := 10
	a := make([]fr.Element, nbDigits)
	for i := 0; i < nbDigits; i++ {
		a[i].SetRandom()
	}

	// create a big int whose digits in basis r are a
	m := fr.Modulus()
	var b, tmp big.Int
	for i := nbDigits - 1; i >= 0; i-- {
		b.Mul(&b, m)
		a[i].ToBigIntRegular(&tmp)
		b.Add(&b, &tmp)
	}

	// query the decomposition and compare to a
	bb := b.Bytes()
	d := Decompose(bb)
	for i := 0; i < nbDigits; i++ {
		if !d[i].Equal(&a[i]) {
			t.Fatal("error decomposition")
		}
	}

}
