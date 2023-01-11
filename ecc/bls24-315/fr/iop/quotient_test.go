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

package iop

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	"github.com/consensys/gnark-crypto/ecc/bls24-315/fr/fft"
)

// computes x₃ in h(x₁,x₂,x₃) = x₁^{2}*x₂ + x₃ - x₁^{3}
// from x₁ and x₂.
func computex3(x []fr.Element) fr.Element {

	var a, b fr.Element
	a.Square(&x[0]).Mul(&a, &x[1])
	b.Square(&x[0]).Mul(&b, &x[0])
	a.Sub(&b, &a)
	return a

}

func buildPoly(size int, form Form) WrappedPolynomial {
	var f Polynomial
	f.Coefficients = make([]fr.Element, size)
	f.Basis = form.Basis
	f.Layout = form.Layout
	return WrappedPolynomial{Polynomial: &f, Shift: 0}
}

func evalCanonical(p Polynomial, x fr.Element) fr.Element {

	var res fr.Element
	for i := len(p.Coefficients) - 1; i >= 0; i-- {
		res.Mul(&res, &x)
		res.Add(&res, &p.Coefficients[i])
	}
	return res
}

func TestQuotient(t *testing.T) {

	// create the multivariate polynomial h
	// h(x₁,x₂,x₃) = x₁^{2}*x₂ + x₃ - x₁^{3}
	nbEntries := 3

	var h MultivariatePolynomial
	var one, minusOne fr.Element
	one.SetOne()
	minusOne.SetOne().Neg(&minusOne)
	h.AddMonomial(one, []int{2, 1, 0})
	h.AddMonomial(one, []int{0, 0, 1})
	h.AddMonomial(minusOne, []int{3, 0, 0})

	// create an instance (f_i) where h holds
	sizeSystem := 8

	form := Form{Basis: Lagrange, Layout: Regular}

	entries := make([]WrappedPolynomial, nbEntries)
	entries[0] = buildPoly(sizeSystem, form)
	entries[1] = buildPoly(sizeSystem, form)
	entries[2] = buildPoly(sizeSystem, form)

	for i := 0; i < sizeSystem; i++ {

		entries[0].Coefficients[i].SetRandom()
		entries[1].Coefficients[i].SetRandom()
		tmp := computex3(
			[]fr.Element{entries[0].Coefficients[i],
				entries[1].Coefficients[i]})
		entries[2].Coefficients[i].Set(&tmp)

		x := []fr.Element{
			entries[0].Coefficients[i],
			entries[1].Coefficients[i],
			entries[2].Coefficients[i],
		}
		tmp = h.Evaluate(x)
		if !tmp.IsZero() {
			t.Fatal("system does not vanish on x^n-1")
		}
	}

	printVector(entries[0].Coefficients)
	fmt.Println("--")

	// compute the quotient q
	expectedForm := Form{Basis: Canonical, Layout: Regular}
	domains := [2]*fft.Domain{nil, nil}
	quotient, err := ComputeQuotient(entries, h, expectedForm, domains)
	if err != nil {
		t.Fatal(err)
	}

	// checks that h(f_i) = (x^n-1)*q by evaluating the relation
	// at a random point
	var c fr.Element
	c.SetRandom()

	nbElmtsExtended := ecc.NextPowerOfTwo(h.Degree() * uint64(sizeSystem))
	domainExtended, _ := buildDomain(int(nbElmtsExtended), domains[1])

	entries[0].ToCanonical(entries[0].Polynomial, domainExtended)
	printVector(entries[0].Coefficients)

	entries[1].ToCanonical(entries[1].Polynomial, domainExtended)
	entries[1].ToRegular(entries[1].Polynomial)

	entries[2].ToCanonical(entries[2].Polynomial, domainExtended)
	entries[2].ToRegular(entries[2].Polynomial)

	x := []fr.Element{
		evalCanonical(*entries[0].Polynomial, c),
		evalCanonical(*entries[1].Polynomial, c),
		evalCanonical(*entries[2].Polynomial, c),
	}
	l := h.Evaluate(x)
	var xnminusone fr.Element
	xnminusone.Set(&c).
		Square(&xnminusone).
		Square(&xnminusone).
		Square(&xnminusone).
		Sub(&xnminusone, &one)
	r := evalCanonical(quotient, c)
	r.Mul(&r, &xnminusone)

	if !r.Equal(&l) {
		t.Fatal("error quotient")
	}

}
