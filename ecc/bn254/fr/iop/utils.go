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

package iop

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func printVector(v []fr.Element) {
	fmt.Printf("[")
	for i := 0; i < len(v); i++ {
		fmt.Printf("Fr(%s), ", v[i].String())
	}
	fmt.Printf("]\n")
}

// func printPolynomials(p []*Polynomial) {
// 	fmt.Printf("[\n")
// 	for i := 0; i < len(p); i++ {
// 		printVector(p[i].Coefficients)
// 		fmt.Printf(",\n")
// 	}
// 	fmt.Printf("]\n")
// }

func printLayout(f Form) {

	if f.Basis == Canonical {
		fmt.Printf("CANONICAL")
	} else if f.Basis == LagrangeCoset {
		fmt.Printf("LAGRANGE_COSET")
	} else {
		fmt.Printf("LAGRANGE")
	}
	fmt.Println("")

	if f.Layout == Regular {
		fmt.Printf("REGULAR")
	} else {
		fmt.Printf("BIT REVERSED")
	}
	fmt.Println("")

}

//----------------------------------------------------
// exp functions until 5

func exp0(x fr.Element) fr.Element {
	var res fr.Element
	res.SetOne()
	return res
}

func exp1(x fr.Element) fr.Element {
	return x
}

func exp2(x fr.Element) fr.Element {
	return *x.Square(&x)
}

func exp3(x fr.Element) fr.Element {
	var res fr.Element
	res.Square(&x).Mul(&res, &x)
	return res
}

func exp4(x fr.Element) fr.Element {
	x.Square(&x).Square(&x)
	return x
}

func exp5(x fr.Element) fr.Element {
	var res fr.Element
	res.Square(&x).Square(&res).Mul(&res, &x)
	return res
}

// doesn't return any errors, it is a private method, that
// is assumed to be called with correct arguments.
func smallExp(x fr.Element, n int) fr.Element {
	if n == 0 {
		return exp0(x)
	}
	if n == 1 {
		return exp1(x)
	}
	if n == 2 {
		return exp2(x)
	}
	if n == 3 {
		return exp3(x)
	}
	if n == 4 {
		return exp4(x)
	}
	if n == 5 {
		return exp5(x)
	}
	return fr.Element{}
}
