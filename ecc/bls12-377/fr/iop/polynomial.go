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
	"errors"
	"math/big"
	"math/bits"

	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr/fft"
)

//-----------------------------------------------------
// univariate polynomials

// Enum to tell in which basis a polynomial is represented.
type Basis int64

const (
	Canonical Basis = iota
	Lagrange
	LagrangeCoset
)

// Enum to tell if a polynomial is in bit reverse form or
// in the regular form.
type Layout int64

const (
	Regular Layout = iota
	BitReverse
)

// Form describes the form of a polynomial.
type Form struct {
	Basis  Basis
	Layout Layout
}

// Polynomial represents a polynomial, the vector of coefficients
// along with the basis and the layout.
type Polynomial struct {
	Coefficients []fr.Element
	Form
}

// return a copy of p
// return a copy of p
func (p *Polynomial) Copy() *Polynomial {
	size := len(p.Coefficients)
	var r Polynomial
	r.Coefficients = make([]fr.Element, size)
	copy(r.Coefficients, p.Coefficients)
	r.Form = p.Form
	return &r
}

// return an ID corresponding to the polynomial extra data
func getShapeID(p Polynomial) int {
	return int(p.Basis)*2 + int(p.Layout)
}

// WrappedPolynomial wrapps a polynomial so that it is
// interpreted as P'(X)=P(\omega^{s}X)
type WrappedPolynomial struct {
	*Polynomial
	Shift int
}

//----------------------------------------------------
// ToRegular

func (p *Polynomial) ToRegular(q *Polynomial) *Polynomial {

	if p != q {
		*p = *q.Copy()
	}
	if p.Layout == Regular {
		return p
	}
	fft.BitReverse(p.Coefficients)
	p.Layout = Regular
	return p
}

//----------------------------------------------------
// ToBitreverse

func (p *Polynomial) ToBitreverse(q *Polynomial) *Polynomial {

	if p != q {
		*p = *q.Copy()
	}
	if p.Layout == BitReverse {
		return p
	}
	fft.BitReverse(p.Coefficients)
	p.Layout = BitReverse
	return p
}

//----------------------------------------------------
// toLagrange

// the numeration corresponds to the following formatting:
// num = int(p.Basis)*2 + int(p.Layout)

// CANONICAL REGULAR
func (p *Polynomial) toLagrange0(d *fft.Domain) *Polynomial {
	p.Basis = Lagrange
	p.Layout = BitReverse
	d.FFT(p.Coefficients, fft.DIF)
	return p
}

// CANONICAL BITREVERSE
func (p *Polynomial) toLagrange1(d *fft.Domain) *Polynomial {
	p.Basis = Lagrange
	p.Layout = Regular
	d.FFT(p.Coefficients, fft.DIT)
	return p
}

// LAGRANGE REGULAR
func (p *Polynomial) toLagrange2(d *fft.Domain) *Polynomial {
	return p
}

// LAGRANGE BITREVERSE
func (p *Polynomial) toLagrange3(d *fft.Domain) *Polynomial {
	return p
}

// LAGRANGE_COSET REGULAR
func (p *Polynomial) toLagrange4(d *fft.Domain) *Polynomial {
	p.Basis = Lagrange
	p.Layout = Regular
	d.FFTInverse(p.Coefficients, fft.DIF, true)
	d.FFT(p.Coefficients, fft.DIT)
	return p
}

// LAGRANGE_COSET BITREVERSE
func (p *Polynomial) toLagrange5(d *fft.Domain) *Polynomial {
	p.Basis = Lagrange
	p.Layout = BitReverse
	d.FFTInverse(p.Coefficients, fft.DIT, true)
	d.FFT(p.Coefficients, fft.DIF)
	return p
}

// Set p to q in Lagrange form and returns it.
func (p *Polynomial) ToLagrange(q *Polynomial, d *fft.Domain) *Polynomial {
	id := getShapeID(*q)
	if q != p {
		*p = *q.Copy()
	}
	resize(p, d.Cardinality)
	switch id {
	case 0:
		return p.toLagrange0(d)
	case 1:
		return p.toLagrange1(d)
	case 2:
		return p.toLagrange2(d)
	case 3:
		return p.toLagrange3(d)
	case 4:
		return p.toLagrange4(d)
	case 5:
		return p.toLagrange5(d)
	default:
		panic("unknown ID")
	}
}

//----------------------------------------------------
// toCanonical

// CANONICAL REGULAR
func (p *Polynomial) toCanonical0(d *fft.Domain) *Polynomial {
	return p
}

// CANONICAL BITREVERSE
func (p *Polynomial) toCanonical1(d *fft.Domain) *Polynomial {
	return p
}

// LAGRANGE REGULAR
func (p *Polynomial) toCanonical2(d *fft.Domain) *Polynomial {
	p.Basis = Canonical
	p.Layout = BitReverse
	d.FFTInverse(p.Coefficients, fft.DIF)
	return p
}

// LAGRANGE BITREVERSE
func (p *Polynomial) toCanonical3(d *fft.Domain) *Polynomial {
	p.Basis = Canonical
	p.Layout = Regular
	d.FFTInverse(p.Coefficients, fft.DIT)
	return p
}

// LAGRANGE_COSET REGULAR
func (p *Polynomial) toCanonical4(d *fft.Domain) *Polynomial {
	p.Basis = Canonical
	p.Layout = BitReverse
	d.FFTInverse(p.Coefficients, fft.DIF, true)
	return p
}

// LAGRANGE_COSET BITREVERSE
func (p *Polynomial) toCanonical5(d *fft.Domain) *Polynomial {
	p.Basis = Canonical
	p.Layout = Regular
	d.FFTInverse(p.Coefficients, fft.DIT, true)
	return p
}

// ToCanonical Sets p to q, in canonical form and returns it.
func (p *Polynomial) ToCanonical(q *Polynomial, d *fft.Domain) *Polynomial {
	id := getShapeID(*q)
	if q != p {
		*p = *q.Copy()
	}
	resize(p, d.Cardinality)
	switch id {
	case 0:
		return p.toCanonical0(d)
	case 1:
		return p.toCanonical1(d)
	case 2:
		return p.toCanonical2(d)
	case 3:
		return p.toCanonical3(d)
	case 4:
		return p.toCanonical4(d)
	case 5:
		return p.toCanonical5(d)
	default:
		panic("unknown ID")
	}
}

//-----------------------------------------------------
// ToLagrangeCoset

func resize(p *Polynomial, newSize uint64) {
	z := make([]fr.Element, int(newSize)-len(p.Coefficients))
	p.Coefficients = append(p.Coefficients, z...)
}

// CANONICAL REGULAR
func (p *Polynomial) toLagrangeCoset0(d *fft.Domain) *Polynomial {
	p.Basis = LagrangeCoset
	p.Layout = BitReverse
	d.FFT(p.Coefficients, fft.DIF, true)
	return p
}

// CANONICAL BITREVERSE
func (p *Polynomial) toLagrangeCoset1(d *fft.Domain) *Polynomial {
	p.Basis = LagrangeCoset
	p.Layout = Regular
	d.FFT(p.Coefficients, fft.DIT, true)
	return p
}

// LAGRANGE REGULAR
func (p *Polynomial) toLagrangeCoset2(d *fft.Domain) *Polynomial {
	p.Basis = LagrangeCoset
	p.Layout = Regular
	d.FFTInverse(p.Coefficients, fft.DIF)
	d.FFT(p.Coefficients, fft.DIT, true)
	return p
}

// LAGRANGE BITREVERSE
func (p *Polynomial) toLagrangeCoset3(d *fft.Domain) *Polynomial {
	p.Basis = LagrangeCoset
	p.Layout = BitReverse
	d.FFTInverse(p.Coefficients, fft.DIT)
	d.FFT(p.Coefficients, fft.DIF, true)
	return p
}

// LAGRANGE_COSET REGULAR
func (p *Polynomial) toLagrangeCoset4(d *fft.Domain) *Polynomial {
	return p
}

// LAGRANGE_COSET BITREVERSE
func (p *Polynomial) toLagrangeCoset5(d *fft.Domain) *Polynomial {
	return p
}

// ToLagrangeCoset Sets p to q, in LagrangeCoset form and returns it.
func (p *Polynomial) ToLagrangeCoset(q *Polynomial, d *fft.Domain) *Polynomial {
	id := getShapeID(*q)
	if q != p {
		*p = *q.Copy()
	}
	resize(p, d.Cardinality)
	switch id {
	case 0:
		return p.toLagrangeCoset0(d)
	case 1:
		return p.toLagrangeCoset1(d)
	case 2:
		return p.toLagrangeCoset2(d)
	case 3:
		return p.toLagrangeCoset3(d)
	case 4:
		return p.toLagrangeCoset4(d)
	case 5:
		return p.toLagrangeCoset5(d)
	default:
		panic("unknown ID")
	}
}

//-----------------------------------------------------
// multivariate polynomials

// errors related to the polynomials.
var ErrInconsistentNumberOfVariable = errors.New("the number of variables is not consistent")

// Monomial represents a Monomial encoded as
// coeff*X₁^{i₁}*..*X_n^{i_n} if exponents = [i₁,..iₙ]
type Monomial struct {
	coeff     fr.Element
	exponents []int
}

// it is supposed that the number of variables matches
func (m Monomial) evaluate(x []fr.Element) fr.Element {

	var res, tmp fr.Element

	nbVars := len(x)
	res.SetOne()
	for i := 0; i < nbVars; i++ {
		if m.exponents[i] <= 5 {
			tmp = smallExp(x[i], m.exponents[i])
			res.Mul(&res, &tmp)
			continue
		}
		bi := big.NewInt(int64(i))
		tmp.Exp(x[i], bi)
		res.Mul(&res, &tmp)
	}
	res.Mul(&res, &m.coeff)

	return res

}

// reprensents a multivariate polynomial as a list of Monomial,
// the multivariate polynomial being the sum of the Monomials.
type MultivariatePolynomial struct {
	M []Monomial
	C fr.Element
}

// degree returns the total degree
func (m *MultivariatePolynomial) Degree() uint64 {
	r := 0
	for i := 0; i < len(m.M); i++ {
		t := 0
		for j := 0; j < len(m.M[i].exponents); j++ {
			t += m.M[i].exponents[j]
		}
		if t > r {
			r = t
		}
	}
	return uint64(r)
}

// AddMonomial adds a Monomial to m. If m is empty, the Monomial is
// added no matter what. But if m is already populated, an error is
// returned if len(e)\neq size of the previous list of exponents. This
// ensure that the number of variables is given by the size of any of
// the slices of exponent in any Monomial.
func (m *MultivariatePolynomial) AddMonomial(c fr.Element, e []int) error {

	// if m is empty, we add the first Monomial.
	if len(m.M) == 0 {
		r := Monomial{c, e}
		m.M = append(m.M, r)
		return nil
	}

	// at this stage all of exponennt in m are supposed to be of
	// the same size.
	if len(m.M[0].exponents) != len(e) {
		return ErrInconsistentNumberOfVariable
	}
	r := Monomial{c, e}
	m.M = append(m.M, r)
	return nil

}

// EvaluateSinglePoint a multivariate polynomial in x
// /!\ It is assumed that the multivariate polynomial has been
// built correctly, that is the sizes of the slices in exponents
// are the same /!\
func (m *MultivariatePolynomial) EvaluateSinglePoint(x []fr.Element) fr.Element {

	var res fr.Element

	for i := 0; i < len(m.M); i++ {
		tmp := m.M[i].evaluate(x)
		res.Add(&res, &tmp)
	}
	res.Add(&res, &m.C)
	return res
}

// EvaluatePolynomials evaluate h on x, interpreted as vectors.
// No transformations are made on the polynomials, it is assumed
// that they are all in the same format (e.g. all in LagrangeCoset,
// BitReverse).
func (m *MultivariatePolynomial) EvaluatePolynomials(x []WrappedPolynomial) (Polynomial, error) {

	var res Polynomial

	// check that the sizes are consistent
	nbPolynomials := len(x)
	nbElmts := len(x[0].Coefficients)
	for i := 0; i < nbPolynomials; i++ {
		if len(x[i].Coefficients) != nbElmts {
			return res, ErrInconsistentSize
		}
	}

	// check that the format are consistent
	expectedForm := x[0].Form
	for i := 0; i < nbPolynomials; i++ {
		if x[i].Form != expectedForm {
			return res, ErrInconsistentFormat
		}
	}

	res.Coefficients = make([]fr.Element, nbElmts)
	res.Form = expectedForm

	// two separate loops so the if is not inside the for loop
	if expectedForm.Layout == Regular {

		v := make([]fr.Element, nbPolynomials)

		for i := 0; i < nbElmts; i++ {

			for j := 0; j < nbPolynomials; j++ {

				v[j].Set(&x[j].Coefficients[(i+x[j].Shift)%nbElmts])

			}

			res.Coefficients[i] = m.EvaluateSinglePoint(v)
			res.Coefficients[i].Add(&res.Coefficients[i], &m.C)
		}

	} else {

		v := make([]fr.Element, nbPolynomials)
		nn := uint64(64 - bits.TrailingZeros(uint(nbElmts)))

		for i := 0; i < nbElmts; i++ {

			for j := 0; j < nbPolynomials; j++ {

				// take in account the fact that the polynomial mght be shifted...
				iRev := bits.Reverse64(uint64((i+x[j].Shift))%uint64(nbElmts)) >> nn
				v[j].Set(&x[j].Coefficients[iRev])

			}

			// evaluate h on x
			iRev := bits.Reverse64(uint64(i)) >> nn
			res.Coefficients[iRev] = m.EvaluateSinglePoint(v)
			res.Coefficients[iRev].Add(&res.Coefficients[iRev], &m.C)
		}

	}

	return res, nil

}
