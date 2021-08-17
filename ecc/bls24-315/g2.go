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

package bls24315

import (
	"math"
	"math/big"
	"runtime"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	"github.com/consensys/gnark-crypto/ecc/bls24-315/internal/fptower"
	"github.com/consensys/gnark-crypto/internal/parallel"
)

// G2Affine point in affine coordinates
type G2Affine struct {
	X, Y fptower.E4
}

// G2Jac is a point with fptower.E4 coordinates
type G2Jac struct {
	X, Y, Z fptower.E4
}

//  g2JacExtended parameterized jacobian coordinates (x=X/ZZ, y=Y/ZZZ, ZZ**3=ZZZ**2)
type g2JacExtended struct {
	X, Y, ZZ, ZZZ fptower.E4
}

// g2Proj point in projective coordinates
type g2Proj struct {
	x, y, z fptower.E4
}

// -------------------------------------------------------------------------------------------------
// Affine

// Set sets p to the provided point
func (p *G2Affine) Set(a *G2Affine) *G2Affine {
	p.X, p.Y = a.X, a.Y
	return p
}

// ScalarMultiplication computes and returns p = a*s
func (p *G2Affine) ScalarMultiplication(a *G2Affine, s *big.Int) *G2Affine {
	var _p G2Jac
	_p.FromAffine(a)
	_p.mulGLV(&_p, s)
	p.FromJacobian(&_p)
	return p
}

// Add adds two point in affine coordinates.
// This should rarely be used as it is very inneficient compared to Jacobian
// TODO implement affine addition formula
func (p *G2Affine) Add(a, b *G2Affine) *G2Affine {
	var p1, p2 G2Jac
	p1.FromAffine(a)
	p2.FromAffine(b)
	p1.AddAssign(&p2)
	p.FromJacobian(&p1)
	return p
}

// Sub subs two point in affine coordinates.
// This should rarely be used as it is very inneficient compared to Jacobian
// TODO implement affine addition formula
func (p *G2Affine) Sub(a, b *G2Affine) *G2Affine {
	var p1, p2 G2Jac
	p1.FromAffine(a)
	p2.FromAffine(b)
	p1.SubAssign(&p2)
	p.FromJacobian(&p1)
	return p
}

// Equal tests if two points (in Affine coordinates) are equal
func (p *G2Affine) Equal(a *G2Affine) bool {
	return p.X.Equal(&a.X) && p.Y.Equal(&a.Y)
}

// Neg computes -G
func (p *G2Affine) Neg(a *G2Affine) *G2Affine {
	p.X = a.X
	p.Y.Neg(&a.Y)
	return p
}

// FromJacobian rescale a point in Jacobian coord in z=1 plane
func (p *G2Affine) FromJacobian(p1 *G2Jac) *G2Affine {

	var a, b fptower.E4

	if p1.Z.IsZero() {
		p.X.SetZero()
		p.Y.SetZero()
		return p
	}

	a.Inverse(&p1.Z)
	b.Square(&a)
	p.X.Mul(&p1.X, &b)
	p.Y.Mul(&p1.Y, &b).Mul(&p.Y, &a)

	return p
}

func (p *G2Affine) String() string {
	var x, y fptower.E4
	x.Set(&p.X)
	y.Set(&p.Y)
	return "E([" + x.String() + "," + y.String() + "]),"
}

// IsInfinity checks if the point is infinity (in affine, it's encoded as (0,0))
func (p *G2Affine) IsInfinity() bool {
	return p.X.IsZero() && p.Y.IsZero()
}

// IsOnCurve returns true if p in on the curve
func (p *G2Affine) IsOnCurve() bool {
	var point G2Jac
	point.FromAffine(p)
	return point.IsOnCurve() // call this function to handle infinity point
}

// IsInSubGroup returns true if p is in the correct subgroup, false otherwise
func (p *G2Affine) IsInSubGroup() bool {
	var _p G2Jac
	_p.FromAffine(p)
	return _p.IsOnCurve() && _p.IsInSubGroup()
}

// -------------------------------------------------------------------------------------------------
// Jacobian

// Set sets p to the provided point
func (p *G2Jac) Set(a *G2Jac) *G2Jac {
	p.X, p.Y, p.Z = a.X, a.Y, a.Z
	return p
}

// Equal tests if two points (in Jacobian coordinates) are equal
func (p *G2Jac) Equal(a *G2Jac) bool {

	if p.Z.IsZero() && a.Z.IsZero() {
		return true
	}
	_p := G2Affine{}
	_p.FromJacobian(p)

	_a := G2Affine{}
	_a.FromJacobian(a)

	return _p.X.Equal(&_a.X) && _p.Y.Equal(&_a.Y)
}

// Neg computes -G
func (p *G2Jac) Neg(a *G2Jac) *G2Jac {
	*p = *a
	p.Y.Neg(&a.Y)
	return p
}

// SubAssign substracts two points on the curve
func (p *G2Jac) SubAssign(a *G2Jac) *G2Jac {
	var tmp G2Jac
	tmp.Set(a)
	tmp.Y.Neg(&tmp.Y)
	p.AddAssign(&tmp)
	return p
}

// AddAssign point addition in montgomery form
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
func (p *G2Jac) AddAssign(a *G2Jac) *G2Jac {

	// p is infinity, return a
	if p.Z.IsZero() {
		p.Set(a)
		return p
	}

	// a is infinity, return p
	if a.Z.IsZero() {
		return p
	}

	var Z1Z1, Z2Z2, U1, U2, S1, S2, H, I, J, r, V fptower.E4
	Z1Z1.Square(&a.Z)
	Z2Z2.Square(&p.Z)
	U1.Mul(&a.X, &Z2Z2)
	U2.Mul(&p.X, &Z1Z1)
	S1.Mul(&a.Y, &p.Z).
		Mul(&S1, &Z2Z2)
	S2.Mul(&p.Y, &a.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U1.Equal(&U2) && S1.Equal(&S2) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &U1)
	I.Double(&H).
		Square(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &S1).Double(&r)
	V.Mul(&U1, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	S1.Mul(&S1, &J).Double(&S1)
	p.Y.Sub(&p.Y, &S1)
	p.Z.Add(&p.Z, &a.Z)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &Z2Z2).
		Mul(&p.Z, &H)

	return p
}

// AddMixed point addition
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-madd-2007-bl
func (p *G2Jac) AddMixed(a *G2Affine) *G2Jac {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.Z.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.Z.SetOne()
		return p
	}

	var Z1Z1, U2, S2, H, HH, I, J, r, V fptower.E4
	Z1Z1.Square(&p.Z)
	U2.Mul(&a.X, &Z1Z1)
	S2.Mul(&a.Y, &p.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U2.Equal(&p.X) && S2.Equal(&p.Y) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &p.X)
	HH.Square(&H)
	I.Double(&HH).Double(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &p.Y).Double(&r)
	V.Mul(&p.X, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	J.Mul(&J, &p.Y).Double(&J)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	p.Y.Sub(&p.Y, &J)
	p.Z.Add(&p.Z, &H)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &HH)

	return p
}

// Double doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G2Jac) Double(q *G2Jac) *G2Jac {
	p.Set(q)
	p.DoubleAssign()
	return p
}

// DoubleAssign doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G2Jac) DoubleAssign() *G2Jac {

	var XX, YY, YYYY, ZZ, S, M, T fptower.E4

	XX.Square(&p.X)
	YY.Square(&p.Y)
	YYYY.Square(&YY)
	ZZ.Square(&p.Z)
	S.Add(&p.X, &YY)
	S.Square(&S).
		Sub(&S, &XX).
		Sub(&S, &YYYY).
		Double(&S)
	M.Double(&XX).Add(&M, &XX)
	p.Z.Add(&p.Z, &p.Y).
		Square(&p.Z).
		Sub(&p.Z, &YY).
		Sub(&p.Z, &ZZ)
	T.Square(&M)
	p.X = T
	T.Double(&S)
	p.X.Sub(&p.X, &T)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M)
	YYYY.Double(&YYYY).Double(&YYYY).Double(&YYYY)
	p.Y.Sub(&p.Y, &YYYY)

	return p
}

// ScalarMultiplication computes and returns p = a*s
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G2Jac) ScalarMultiplication(a *G2Jac, s *big.Int) *G2Jac {
	return p.mulGLV(a, s)
}

func (p *G2Jac) String() string {
	if p.Z.IsZero() {
		return "O"
	}
	_p := G2Affine{}
	_p.FromJacobian(p)
	return "E([" + _p.X.String() + "," + _p.Y.String() + "]),"
}

// FromAffine sets p = Q, p in Jacboian, Q in affine
func (p *G2Jac) FromAffine(Q *G2Affine) *G2Jac {
	if Q.X.IsZero() && Q.Y.IsZero() {
		p.Z.SetZero()
		p.X.SetOne()
		p.Y.SetOne()
		return p
	}
	p.Z.SetOne()
	p.X.Set(&Q.X)
	p.Y.Set(&Q.Y)
	return p
}

// IsOnCurve returns true if p in on the curve
func (p *G2Jac) IsOnCurve() bool {
	var left, right, tmp fptower.E4
	left.Square(&p.Y)
	right.Square(&p.X).Mul(&right, &p.X)
	tmp.Square(&p.Z).
		Square(&tmp).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &bTwistCurveCoeff)
	right.Add(&right, &tmp)
	return left.Equal(&right)
}

// IsInSubGroup returns true if p is on the r-torsion, false otherwise.
// Z[r,0]+Z[-lambdaG2Affine, 1] is the kernel
// of (u,v)->u+lambdaG2Affinev mod r. Expressing r, lambdaG2Affine as
// polynomials in x, a short vector of this Zmodule is
// 1, x**4. So we check that p+x**4*phi(p)
// is the infinity.
func (p *G2Jac) IsInSubGroup() bool {

	var res G2Jac
	res.phi(p).
		ScalarMultiplication(&res, &xGen).
		ScalarMultiplication(&res, &xGen).
		ScalarMultiplication(&res, &xGen).
		ScalarMultiplication(&res, &xGen).
		AddAssign(p)

	return res.IsOnCurve() && res.Z.IsZero()

}

// mulWindowed 2-bits windowed exponentiation
func (p *G2Jac) mulWindowed(a *G2Jac, s *big.Int) *G2Jac {

	var res G2Jac
	var ops [3]G2Jac

	res.Set(&g2Infinity)
	ops[0].Set(a)
	ops[1].Double(&ops[0])
	ops[2].Set(&ops[0]).AddAssign(&ops[1])

	b := s.Bytes()
	for i := range b {
		w := b[i]
		mask := byte(0xc0)
		for j := 0; j < 4; j++ {
			res.DoubleAssign().DoubleAssign()
			c := (w & mask) >> (6 - 2*j)
			if c != 0 {
				res.AddAssign(&ops[c-1])
			}
			mask = mask >> 2
		}
	}
	p.Set(&res)

	return p

}

// psi(p) = u o frob o u**-1 where u:E'->E iso from the twist to E
func (p *G2Jac) psi(a *G2Jac) *G2Jac {
	p.Set(a)
	p.X.Frobenius(&p.X).Mul(&p.X, &endo.u)
	p.Y.Frobenius(&p.Y).Mul(&p.Y, &endo.v)
	p.Z.Frobenius(&p.Z)
	return p
}

// phi assigns p to phi(a) where phi: (x,y)->(ux,y), and returns p
func (p *G2Jac) phi(a *G2Jac) *G2Jac {
	p.Set(a)
	p.X.MulByElement(&p.X, &thirdRootOneG2)
	return p
}

// mulGLV performs scalar multiplication using GLV
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G2Jac) mulGLV(a *G2Jac, s *big.Int) *G2Jac {

	var table [15]G2Jac
	var zero big.Int
	var res G2Jac
	var k1, k2 fr.Element

	res.Set(&g2Infinity)

	// table[b3b2b1b0-1] = b3b2*phi(a) + b1b0*a
	table[0].Set(a)
	table[3].phi(a)

	// split the scalar, modifies +-a, phi(a) accordingly
	k := ecc.SplitScalar(s, &glvBasis)

	if k[0].Cmp(&zero) == -1 {
		k[0].Neg(&k[0])
		table[0].Neg(&table[0])
	}
	if k[1].Cmp(&zero) == -1 {
		k[1].Neg(&k[1])
		table[3].Neg(&table[3])
	}

	// precompute table (2 bits sliding window)
	// table[b3b2b1b0-1] = b3b2*phi(a) + b1b0*a if b3b2b1b0 != 0
	table[1].Double(&table[0])
	table[2].Set(&table[1]).AddAssign(&table[0])
	table[4].Set(&table[3]).AddAssign(&table[0])
	table[5].Set(&table[3]).AddAssign(&table[1])
	table[6].Set(&table[3]).AddAssign(&table[2])
	table[7].Double(&table[3])
	table[8].Set(&table[7]).AddAssign(&table[0])
	table[9].Set(&table[7]).AddAssign(&table[1])
	table[10].Set(&table[7]).AddAssign(&table[2])
	table[11].Set(&table[7]).AddAssign(&table[3])
	table[12].Set(&table[11]).AddAssign(&table[0])
	table[13].Set(&table[11]).AddAssign(&table[1])
	table[14].Set(&table[11]).AddAssign(&table[2])

	// bounds on the lattice base vectors guarantee that k1, k2 are len(r)/2 bits long max
	k1.SetBigInt(&k[0]).FromMont()
	k2.SetBigInt(&k[1]).FromMont()

	// loop starts from len(k1)/2 due to the bounds
	for i := int(math.Ceil(fr.Limbs/2. - 1)); i >= 0; i-- {
		mask := uint64(3) << 62
		for j := 0; j < 32; j++ {
			res.Double(&res).Double(&res)
			b1 := (k1[i] & mask) >> (62 - 2*j)
			b2 := (k2[i] & mask) >> (62 - 2*j)
			if b1|b2 != 0 {
				s := (b2<<2 | b1)
				res.AddAssign(&table[s-1])
			}
			mask = mask >> 2
		}
	}

	p.Set(&res)
	return p
}

// ClearCofactor maps a point in curve to r-torsion
func (p *G2Affine) ClearCofactor(a *G2Affine) *G2Affine {
	var _p G2Jac
	_p.FromAffine(a)
	_p.ClearCofactor(&_p)
	p.FromJacobian(&_p)
	return p
}

// ClearCofactor maps a point in curve to r-torsion
func (p *G2Jac) ClearCofactor(a *G2Jac) *G2Jac {
	// https://eprint.iacr.org/2017/419.pdf, section 4.2
	// multiply by (3x^4-3)*cofacor
	var xg, xxg, xxxg, xxxxg, res, t G2Jac
	xg.ScalarMultiplication(a, &xGen).Neg(&xg).SubAssign(a)
	xxg.ScalarMultiplication(&xg, &xGen).Neg(&xxg)
	xxxg.ScalarMultiplication(&xxg, &xGen).Neg(&xxxg)
	xxxxg.ScalarMultiplication(&xxxg, &xGen).Neg(&xxxxg)

	res.Set(&xxxxg).
		SubAssign(a)

	t.Set(&xxxg).
		psi(&t).
		AddAssign(&res)

	res.Set(&xxg).
		psi(&res).
		psi(&res).
		AddAssign(&t)

	t.Set(&xg).
		psi(&t).
		psi(&t).
		psi(&t).
		AddAssign(&res)

	res.Double(a).
		psi(&res).
		psi(&res).
		psi(&res).
		psi(&res).
		AddAssign(&t)

	p.Set(&res)

	return p

}

// -------------------------------------------------------------------------------------------------
// Jacobian extended

// Set sets p to the provided point
func (p *g2JacExtended) Set(a *g2JacExtended) *g2JacExtended {
	p.X, p.Y, p.ZZ, p.ZZZ = a.X, a.Y, a.ZZ, a.ZZZ
	return p
}

// setInfinity sets p to O
func (p *g2JacExtended) setInfinity() *g2JacExtended {
	p.X.SetOne()
	p.Y.SetOne()
	p.ZZ = fptower.E4{}
	p.ZZZ = fptower.E4{}
	return p
}

// fromJacExtended sets Q in affine coords
func (p *G2Affine) fromJacExtended(Q *g2JacExtended) *G2Affine {
	if Q.ZZ.IsZero() {
		p.X = fptower.E4{}
		p.Y = fptower.E4{}
		return p
	}
	p.X.Inverse(&Q.ZZ).Mul(&p.X, &Q.X)
	p.Y.Inverse(&Q.ZZZ).Mul(&p.Y, &Q.Y)
	return p
}

// fromJacExtended sets Q in Jacobian coords
func (p *G2Jac) fromJacExtended(Q *g2JacExtended) *G2Jac {
	if Q.ZZ.IsZero() {
		p.Set(&g2Infinity)
		return p
	}
	p.X.Mul(&Q.ZZ, &Q.X).Mul(&p.X, &Q.ZZ)
	p.Y.Mul(&Q.ZZZ, &Q.Y).Mul(&p.Y, &Q.ZZZ)
	p.Z.Set(&Q.ZZZ)
	return p
}

// unsafeFromJacExtended sets p in jacobian coords, but don't check for infinity
func (p *G2Jac) unsafeFromJacExtended(Q *g2JacExtended) *G2Jac {
	p.X.Square(&Q.ZZ).Mul(&p.X, &Q.X)
	p.Y.Square(&Q.ZZZ).Mul(&p.Y, &Q.Y)
	p.Z = Q.ZZZ
	return p
}

// add point in ZZ coords
// https://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#addition-add-2008-s
func (p *g2JacExtended) add(q *g2JacExtended) *g2JacExtended {
	//if q is infinity return p
	if q.ZZ.IsZero() {
		return p
	}
	// p is infinity, return q
	if p.ZZ.IsZero() {
		p.Set(q)
		return p
	}

	var A, B, X1ZZ2, X2ZZ1, Y1ZZZ2, Y2ZZZ1 fptower.E4

	// p2: q, p1: p
	X2ZZ1.Mul(&q.X, &p.ZZ)
	X1ZZ2.Mul(&p.X, &q.ZZ)
	A.Sub(&X2ZZ1, &X1ZZ2)
	Y2ZZZ1.Mul(&q.Y, &p.ZZZ)
	Y1ZZZ2.Mul(&p.Y, &q.ZZZ)
	B.Sub(&Y2ZZZ1, &Y1ZZZ2)

	if A.IsZero() {
		if B.IsZero() {
			return p.double(q)

		}
		p.ZZ = fptower.E4{}
		p.ZZZ = fptower.E4{}
		return p
	}

	var U1, U2, S1, S2, P, R, PP, PPP, Q, V fptower.E4
	U1.Mul(&p.X, &q.ZZ)
	U2.Mul(&q.X, &p.ZZ)
	S1.Mul(&p.Y, &q.ZZZ)
	S2.Mul(&q.Y, &p.ZZZ)
	P.Sub(&U2, &U1)
	R.Sub(&S2, &S1)
	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&U1, &PP)
	V.Mul(&S1, &PPP)

	p.X.Square(&R).
		Sub(&p.X, &PPP).
		Sub(&p.X, &Q).
		Sub(&p.X, &Q)
	p.Y.Sub(&Q, &p.X).
		Mul(&p.Y, &R).
		Sub(&p.Y, &V)
	p.ZZ.Mul(&p.ZZ, &q.ZZ).
		Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &q.ZZZ).
		Mul(&p.ZZZ, &PPP)

	return p
}

// double point in ZZ coords
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#doubling-dbl-2008-s-1
func (p *g2JacExtended) double(q *g2JacExtended) *g2JacExtended {
	var U, V, W, S, XX, M fptower.E4

	U.Double(&q.Y)
	V.Square(&U)
	W.Mul(&U, &V)
	S.Mul(&q.X, &V)
	XX.Square(&q.X)
	M.Double(&XX).
		Add(&M, &XX) // -> + a, but a=0 here
	U.Mul(&W, &q.Y)

	p.X.Square(&M).
		Sub(&p.X, &S).
		Sub(&p.X, &S)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M).
		Sub(&p.Y, &U)
	p.ZZ.Mul(&V, &q.ZZ)
	p.ZZZ.Mul(&W, &q.ZZZ)

	return p
}

// subMixed same as addMixed, but will negate a.Y
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#addition-madd-2008-s
func (p *g2JacExtended) subMixed(a *G2Affine) *g2JacExtended {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.ZZ.IsZero() {
		p.X = a.X
		p.Y.Neg(&a.Y)
		p.ZZ.SetOne()
		p.ZZZ.SetOne()
		return p
	}

	var P, R fptower.E4

	// p2: a, p1: p
	P.Mul(&a.X, &p.ZZ)
	P.Sub(&P, &p.X)

	R.Mul(&a.Y, &p.ZZZ)
	R.Neg(&R)
	R.Sub(&R, &p.Y)

	if P.IsZero() {
		if R.IsZero() {
			return p.doubleNegMixed(a)

		}
		p.ZZ = fptower.E4{}
		p.ZZZ = fptower.E4{}
		return p
	}

	var PP, PPP, Q, Q2, RR, X3, Y3 fptower.E4

	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&p.X, &PP)
	RR.Square(&R)
	X3.Sub(&RR, &PPP)
	Q2.Double(&Q)
	p.X.Sub(&X3, &Q2)
	Y3.Sub(&Q, &p.X).Mul(&Y3, &R)
	R.Mul(&p.Y, &PPP)
	p.Y.Sub(&Y3, &R)
	p.ZZ.Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &PPP)

	return p

}

// addMixed
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#addition-madd-2008-s
func (p *g2JacExtended) addMixed(a *G2Affine) *g2JacExtended {

	//if a is infinity return p
	if a.X.IsZero() && a.Y.IsZero() {
		return p
	}
	// p is infinity, return a
	if p.ZZ.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.ZZ.SetOne()
		p.ZZZ.SetOne()
		return p
	}

	var P, R fptower.E4

	// p2: a, p1: p
	P.Mul(&a.X, &p.ZZ)
	P.Sub(&P, &p.X)

	R.Mul(&a.Y, &p.ZZZ)
	R.Sub(&R, &p.Y)

	if P.IsZero() {
		if R.IsZero() {
			return p.doubleMixed(a)

		}
		p.ZZ = fptower.E4{}
		p.ZZZ = fptower.E4{}
		return p
	}

	var PP, PPP, Q, Q2, RR, X3, Y3 fptower.E4

	PP.Square(&P)
	PPP.Mul(&P, &PP)
	Q.Mul(&p.X, &PP)
	RR.Square(&R)
	X3.Sub(&RR, &PPP)
	Q2.Double(&Q)
	p.X.Sub(&X3, &Q2)
	Y3.Sub(&Q, &p.X).Mul(&Y3, &R)
	R.Mul(&p.Y, &PPP)
	p.Y.Sub(&Y3, &R)
	p.ZZ.Mul(&p.ZZ, &PP)
	p.ZZZ.Mul(&p.ZZZ, &PPP)

	return p

}

// doubleNegMixed same as double, but will negate q.Y
func (p *g2JacExtended) doubleNegMixed(q *G2Affine) *g2JacExtended {

	var U, V, W, S, XX, M, S2, L fptower.E4

	U.Double(&q.Y)
	U.Neg(&U)
	V.Square(&U)
	W.Mul(&U, &V)
	S.Mul(&q.X, &V)
	XX.Square(&q.X)
	M.Double(&XX).
		Add(&M, &XX) // -> + a, but a=0 here
	S2.Double(&S)
	L.Mul(&W, &q.Y)

	p.X.Square(&M).
		Sub(&p.X, &S2)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M).
		Add(&p.Y, &L)
	p.ZZ.Set(&V)
	p.ZZZ.Set(&W)

	return p
}

// doubleMixed point in ZZ coords
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-xyzz.html#doubling-dbl-2008-s-1
func (p *g2JacExtended) doubleMixed(q *G2Affine) *g2JacExtended {

	var U, V, W, S, XX, M, S2, L fptower.E4

	U.Double(&q.Y)
	V.Square(&U)
	W.Mul(&U, &V)
	S.Mul(&q.X, &V)
	XX.Square(&q.X)
	M.Double(&XX).
		Add(&M, &XX) // -> + a, but a=0 here
	S2.Double(&S)
	L.Mul(&W, &q.Y)

	p.X.Square(&M).
		Sub(&p.X, &S2)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M).
		Sub(&p.Y, &L)
	p.ZZ.Set(&V)
	p.ZZZ.Set(&W)

	return p
}

// -------------------------------------------------------------------------------------------------
// Homogenous projective

// Set sets p to the provided point
func (p *g2Proj) Set(a *g2Proj) *g2Proj {
	p.x, p.y, p.z = a.x, a.y, a.z
	return p
}

// FromJacobian converts a point from Jacobian to projective coordinates
func (p *g2Proj) FromJacobian(Q *G2Jac) *g2Proj {
	var buf fptower.E4
	buf.Square(&Q.Z)

	p.x.Mul(&Q.X, &Q.Z)
	p.y.Set(&Q.Y)
	p.z.Mul(&Q.Z, &buf)

	return p
}

// FromAffine sets p = Q, p in homogenous projective, Q in affine
func (p *g2Proj) FromAffine(Q *G2Affine) *g2Proj {
	if Q.X.IsZero() && Q.Y.IsZero() {
		p.z.SetZero()
		p.x.SetOne()
		p.y.SetOne()
		return p
	}
	p.z.SetOne()
	p.x.Set(&Q.X)
	p.y.Set(&Q.Y)
	return p
}

// BatchScalarMultiplicationG2 multiplies the same base (generator) by all scalars
// and return resulting points in affine coordinates
// uses a simple windowed-NAF like exponentiation algorithm
func BatchScalarMultiplicationG2(base *G2Affine, scalars []fr.Element) []G2Affine {

	// approximate cost in group ops is
	// cost = 2^{c-1} + n(scalar.nbBits+nbChunks)

	nbPoints := uint64(len(scalars))
	min := ^uint64(0)
	bestC := 0
	for c := 2; c < 18; c++ {
		cost := uint64(1 << (c - 1))
		nbChunks := uint64(fr.Limbs * 64 / c)
		if (fr.Limbs*64)%c != 0 {
			nbChunks++
		}
		cost += nbPoints * ((fr.Limbs * 64) + nbChunks)
		if cost < min {
			min = cost
			bestC = c
		}
	}
	c := uint64(bestC) // window size
	nbChunks := int(fr.Limbs * 64 / c)
	if (fr.Limbs*64)%c != 0 {
		nbChunks++
	}
	mask := uint64((1 << c) - 1) // low c bits are 1
	msbWindow := uint64(1 << (c - 1))

	// precompute all powers of base for our window
	// note here that if performance is critical, we can implement as in the msmX methods
	// this allocation to be on the stack
	baseTable := make([]G2Jac, (1 << (c - 1)))
	baseTable[0].Set(&g2Infinity)
	baseTable[0].AddMixed(base)
	for i := 1; i < len(baseTable); i++ {
		baseTable[i] = baseTable[i-1]
		baseTable[i].AddMixed(base)
	}

	pScalars := partitionScalars(scalars, c, false, runtime.NumCPU())

	// compute offset and word selector / shift to select the right bits of our windows
	selectors := make([]selector, nbChunks)
	for chunk := 0; chunk < nbChunks; chunk++ {
		jc := uint64(uint64(chunk) * c)
		d := selector{}
		d.index = jc / 64
		d.shift = jc - (d.index * 64)
		d.mask = mask << d.shift
		d.multiWordSelect = (64%c) != 0 && d.shift > (64-c) && d.index < (fr.Limbs-1)
		if d.multiWordSelect {
			nbBitsHigh := d.shift - uint64(64-c)
			d.maskHigh = (1 << nbBitsHigh) - 1
			d.shiftHigh = (c - nbBitsHigh)
		}
		selectors[chunk] = d
	}
	toReturn := make([]G2Affine, len(scalars))

	// for each digit, take value in the base table, double it c time, voila.
	parallel.Execute(len(pScalars), func(start, end int) {
		var p G2Jac
		for i := start; i < end; i++ {
			p.Set(&g2Infinity)
			for chunk := nbChunks - 1; chunk >= 0; chunk-- {
				s := selectors[chunk]
				if chunk != nbChunks-1 {
					for j := uint64(0); j < c; j++ {
						p.DoubleAssign()
					}
				}

				bits := (pScalars[i][s.index] & s.mask) >> s.shift
				if s.multiWordSelect {
					bits += (pScalars[i][s.index+1] & s.maskHigh) << s.shiftHigh
				}

				if bits == 0 {
					continue
				}

				// if msbWindow bit is set, we need to substract
				if bits&msbWindow == 0 {
					// add
					p.AddAssign(&baseTable[bits-1])
				} else {
					// sub
					t := baseTable[bits & ^msbWindow]
					t.Neg(&t)
					p.AddAssign(&t)
				}
			}

			// set our result point
			toReturn[i].FromJacobian(&p)

		}
	})
	return toReturn
}
