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

package fp

// expBySqrtExp is equivalent to z.Exp(x, c19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52)
//
// uses github.com/mmcloughlin/addchain v0.4.0 to generate a shorter addition chain
func (z *Element) expBySqrtExp(x Element) *Element {
	// addition chain:
	//
	//	_10      = 2*1
	//	_11      = 1 + _10
	//	_101     = _10 + _11
	//	_110     = 1 + _101
	//	_111     = 1 + _110
	//	_1011    = _101 + _110
	//	_1100    = 1 + _1011
	//	_1101    = 1 + _1100
	//	_1111    = _10 + _1101
	//	_10001   = _10 + _1111
	//	_10011   = _10 + _10001
	//	_10111   = _110 + _10001
	//	_11001   = _10 + _10111
	//	_11011   = _10 + _11001
	//	_11111   = _110 + _11001
	//	_100011  = _1100 + _10111
	//	_100111  = _1100 + _11011
	//	_101001  = _10 + _100111
	//	_101011  = _10 + _101001
	//	_101101  = _10 + _101011
	//	_111001  = _1100 + _101101
	//	_1100000 = _100111 + _111001
	//	i46      = ((_1100000 << 5 + _11001) << 9 + _100111) << 8
	//	i62      = ((_111001 + i46) << 4 + _111) << 9 + _10011
	//	i89      = ((i62 << 7 + _1101) << 13 + _101001) << 5
	//	i109     = ((_10111 + i89) << 7 + _101) << 10 + _10001
	//	i130     = ((i109 << 6 + _11011) << 5 + _1101) << 8
	//	i154     = ((_11 + i130) << 12 + _101011) << 9 + _10111
	//	i179     = ((i154 << 6 + _11001) << 5 + _1111) << 12
	//	i198     = ((_101101 + i179) << 7 + _101001) << 9 + _101101
	//	i220     = ((i198 << 7 + _111) << 9 + _111001) << 4
	//	i236     = ((_101 + i220) << 7 + _1101) << 6 + _1111
	//	i265     = ((i236 << 5 + 1) << 11 + _100011) << 11
	//	i281     = ((_101101 + i265) << 4 + _1011) << 9 + _11111
	//	i299     = (i281 << 8 + _110 + _111001) << 7 + _101001
	//	return     2*i299
	//
	// Operations: 246 squares 54 multiplies

	// Allocate Temporaries.
	var (
		t0  = new(Element)
		t1  = new(Element)
		t2  = new(Element)
		t3  = new(Element)
		t4  = new(Element)
		t5  = new(Element)
		t6  = new(Element)
		t7  = new(Element)
		t8  = new(Element)
		t9  = new(Element)
		t10 = new(Element)
		t11 = new(Element)
		t12 = new(Element)
		t13 = new(Element)
		t14 = new(Element)
		t15 = new(Element)
		t16 = new(Element)
		t17 = new(Element)
		t18 = new(Element)
	)

	// var t0,t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11,t12,t13,t14,t15,t16,t17,t18 Element
	// Step 1: t4 = x⁰x2
	t4.Square(&x)

	// Step 2: t13 = x⁰x3
	t13.Mul(&x, t4)

	// Step 3: t8 = x⁰x5
	t8.Mul(t4, t13)

	// Step 4: t1 = x⁰x6
	t1.Mul(&x, t8)

	// Step 5: t9 = x⁰x7
	t9.Mul(&x, t1)

	// Step 6: t3 = x⁰xb
	t3.Mul(t8, t1)

	// Step 7: t0 = x⁰xc
	t0.Mul(&x, t3)

	// Step 8: t7 = x⁰xd
	t7.Mul(&x, t0)

	// Step 9: t6 = x⁰xf
	t6.Mul(t4, t7)

	// Step 10: t15 = x⁰x11
	t15.Mul(t4, t6)

	// Step 11: t16 = x⁰x13
	t16.Mul(t4, t15)

	// Step 12: t11 = x⁰x17
	t11.Mul(t1, t15)

	// Step 13: t10 = x⁰x19
	t10.Mul(t4, t11)

	// Step 14: t14 = x⁰x1b
	t14.Mul(t4, t10)

	// Step 15: t2 = x⁰x1f
	t2.Mul(t1, t10)

	// Step 16: t5 = x⁰x23
	t5.Mul(t0, t11)

	// Step 17: t17 = x⁰x27
	t17.Mul(t0, t14)

	// Step 18: z = x⁰x29
	z.Mul(t4, t17)

	// Step 19: t12 = x⁰x2b
	t12.Mul(t4, z)

	// Step 20: t4 = x⁰x2d
	t4.Mul(t4, t12)

	// Step 21: t0 = x⁰x39
	t0.Mul(t0, t4)

	// Step 22: t18 = x⁰x60
	t18.Mul(t17, t0)

	// Step 27: t18 = x⁰xc00
	for s := 0; s < 5; s++ {
		t18.Square(t18)
	}

	// Step 28: t18 = x⁰xc19
	t18.Mul(t10, t18)

	// Step 37: t18 = x⁰x183200
	for s := 0; s < 9; s++ {
		t18.Square(t18)
	}

	// Step 38: t17 = x⁰x183227
	t17.Mul(t17, t18)

	// Step 46: t17 = x⁰x18322700
	for s := 0; s < 8; s++ {
		t17.Square(t17)
	}

	// Step 47: t17 = x⁰x18322739
	t17.Mul(t0, t17)

	// Step 51: t17 = x⁰x183227390
	for s := 0; s < 4; s++ {
		t17.Square(t17)
	}

	// Step 52: t17 = x⁰x183227397
	t17.Mul(t9, t17)

	// Step 61: t17 = x⁰x30644e72e00
	for s := 0; s < 9; s++ {
		t17.Square(t17)
	}

	// Step 62: t16 = x⁰x30644e72e13
	t16.Mul(t16, t17)

	// Step 69: t16 = x⁰x1832273970980
	for s := 0; s < 7; s++ {
		t16.Square(t16)
	}

	// Step 70: t16 = x⁰x183227397098d
	t16.Mul(t7, t16)

	// Step 83: t16 = x⁰x30644e72e131a000
	for s := 0; s < 13; s++ {
		t16.Square(t16)
	}

	// Step 84: t16 = x⁰x30644e72e131a029
	t16.Mul(z, t16)

	// Step 89: t16 = x⁰x60c89ce5c26340520
	for s := 0; s < 5; s++ {
		t16.Square(t16)
	}

	// Step 90: t16 = x⁰x60c89ce5c26340537
	t16.Mul(t11, t16)

	// Step 97: t16 = x⁰x30644e72e131a029b80
	for s := 0; s < 7; s++ {
		t16.Square(t16)
	}

	// Step 98: t16 = x⁰x30644e72e131a029b85
	t16.Mul(t8, t16)

	// Step 108: t16 = x⁰xc19139cb84c680a6e1400
	for s := 0; s < 10; s++ {
		t16.Square(t16)
	}

	// Step 109: t15 = x⁰xc19139cb84c680a6e1411
	t15.Mul(t15, t16)

	// Step 115: t15 = x⁰x30644e72e131a029b850440
	for s := 0; s < 6; s++ {
		t15.Square(t15)
	}

	// Step 116: t14 = x⁰x30644e72e131a029b85045b
	t14.Mul(t14, t15)

	// Step 121: t14 = x⁰x60c89ce5c263405370a08b60
	for s := 0; s < 5; s++ {
		t14.Square(t14)
	}

	// Step 122: t14 = x⁰x60c89ce5c263405370a08b6d
	t14.Mul(t7, t14)

	// Step 130: t14 = x⁰x60c89ce5c263405370a08b6d00
	for s := 0; s < 8; s++ {
		t14.Square(t14)
	}

	// Step 131: t13 = x⁰x60c89ce5c263405370a08b6d03
	t13.Mul(t13, t14)

	// Step 143: t13 = x⁰x60c89ce5c263405370a08b6d03000
	for s := 0; s < 12; s++ {
		t13.Square(t13)
	}

	// Step 144: t12 = x⁰x60c89ce5c263405370a08b6d0302b
	t12.Mul(t12, t13)

	// Step 153: t12 = x⁰xc19139cb84c680a6e14116da0605600
	for s := 0; s < 9; s++ {
		t12.Square(t12)
	}

	// Step 154: t11 = x⁰xc19139cb84c680a6e14116da0605617
	t11.Mul(t11, t12)

	// Step 160: t11 = x⁰x30644e72e131a029b85045b68181585c0
	for s := 0; s < 6; s++ {
		t11.Square(t11)
	}

	// Step 161: t10 = x⁰x30644e72e131a029b85045b68181585d9
	t10.Mul(t10, t11)

	// Step 166: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb20
	for s := 0; s < 5; s++ {
		t10.Square(t10)
	}

	// Step 167: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f
	t10.Mul(t6, t10)

	// Step 179: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f000
	for s := 0; s < 12; s++ {
		t10.Square(t10)
	}

	// Step 180: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d
	t10.Mul(t4, t10)

	// Step 187: t10 = x⁰x30644e72e131a029b85045b68181585d9781680
	for s := 0; s < 7; s++ {
		t10.Square(t10)
	}

	// Step 188: t10 = x⁰x30644e72e131a029b85045b68181585d97816a9
	t10.Mul(z, t10)

	// Step 197: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d5200
	for s := 0; s < 9; s++ {
		t10.Square(t10)
	}

	// Step 198: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d
	t10.Mul(t4, t10)

	// Step 205: t10 = x⁰x30644e72e131a029b85045b68181585d97816a91680
	for s := 0; s < 7; s++ {
		t10.Square(t10)
	}

	// Step 206: t9 = x⁰x30644e72e131a029b85045b68181585d97816a91687
	t9.Mul(t9, t10)

	// Step 215: t9 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e00
	for s := 0; s < 9; s++ {
		t9.Square(t9)
	}

	// Step 216: t9 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e39
	t9.Mul(t0, t9)

	// Step 220: t9 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e390
	for s := 0; s < 4; s++ {
		t9.Square(t9)
	}

	// Step 221: t8 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e395
	t8.Mul(t8, t9)

	// Step 228: t8 = x⁰x30644e72e131a029b85045b68181585d97816a916871ca80
	for s := 0; s < 7; s++ {
		t8.Square(t8)
	}

	// Step 229: t7 = x⁰x30644e72e131a029b85045b68181585d97816a916871ca8d
	t7.Mul(t7, t8)

	// Step 235: t7 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a340
	for s := 0; s < 6; s++ {
		t7.Square(t7)
	}

	// Step 236: t6 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f
	t6.Mul(t6, t7)

	// Step 241: t6 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e0
	for s := 0; s < 5; s++ {
		t6.Square(t6)
	}

	// Step 242: t6 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e1
	t6.Mul(&x, t6)

	// Step 253: t6 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f0800
	for s := 0; s < 11; s++ {
		t6.Square(t6)
	}

	// Step 254: t5 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f0823
	t5.Mul(t5, t6)

	// Step 265: t5 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3951a78411800
	for s := 0; s < 11; s++ {
		t5.Square(t5)
	}

	// Step 266: t4 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3951a7841182d
	t4.Mul(t4, t5)

	// Step 270: t4 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3951a7841182d0
	for s := 0; s < 4; s++ {
		t4.Square(t4)
	}

	// Step 271: t3 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3951a7841182db
	t3.Mul(t3, t4)

	// Step 280: t3 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b600
	for s := 0; s < 9; s++ {
		t3.Square(t3)
	}

	// Step 281: t2 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f
	t2.Mul(t2, t3)

	// Step 289: t2 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f00
	for s := 0; s < 8; s++ {
		t2.Square(t2)
	}

	// Step 290: t1 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f06
	t1.Mul(t1, t2)

	// Step 291: t0 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f
	t0.Mul(t0, t1)

	// Step 298: t0 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3951a7841182db0f9f80
	for s := 0; s < 7; s++ {
		t0.Square(t0)
	}

	// Step 299: z = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3951a7841182db0f9fa9
	z.Mul(z, t0)

	// Step 300: z = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f52
	z.Square(z)

	return z
}

// expByLegendreExp is equivalent to z.Exp(x, 183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c3e7ea3)
//
// uses github.com/mmcloughlin/addchain v0.4.0 to generate a shorter addition chain
func (z *Element) expByLegendreExp(x Element) *Element {
	// addition chain:
	//
	//	_10       = 2*1
	//	_11       = 1 + _10
	//	_101      = _10 + _11
	//	_110      = 1 + _101
	//	_1000     = _10 + _110
	//	_1101     = _101 + _1000
	//	_10010    = _101 + _1101
	//	_10011    = 1 + _10010
	//	_10100    = 1 + _10011
	//	_10111    = _11 + _10100
	//	_11100    = _101 + _10111
	//	_100000   = _1101 + _10011
	//	_100011   = _11 + _100000
	//	_101011   = _1000 + _100011
	//	_101111   = _10011 + _11100
	//	_1000001  = _10010 + _101111
	//	_1010011  = _10010 + _1000001
	//	_1011011  = _1000 + _1010011
	//	_1100001  = _110 + _1011011
	//	_1110101  = _10100 + _1100001
	//	_10010001 = _11100 + _1110101
	//	_10010101 = _100000 + _1110101
	//	_10110101 = _100000 + _10010101
	//	_10111011 = _110 + _10110101
	//	_11000001 = _110 + _10111011
	//	_11000011 = _10 + _11000001
	//	_11010011 = _10010 + _11000001
	//	_11100001 = _100000 + _11000001
	//	_11100011 = _10 + _11100001
	//	_11100111 = _110 + _11100001
	//	i57       = ((_11000001 << 8 + _10010001) << 10 + _11100111) << 7
	//	i76       = ((_10111 + i57) << 9 + _10011) << 7 + _1101
	//	i109      = ((i76 << 14 + _1010011) << 9 + _11100001) << 8
	//	i127      = ((_1000001 + i109) << 10 + _1011011) << 5 + _1101
	//	i161      = ((i127 << 8 + _11) << 12 + _101011) << 12
	//	i186      = ((_10111011 + i161) << 8 + _101111) << 14 + _10110101
	//	i214      = ((i186 << 9 + _10010001) << 5 + _1101) << 12
	//	i236      = ((_11100011 + i214) << 8 + _10010101) << 11 + _11010011
	//	i268      = ((i236 << 7 + _1100001) << 11 + _100011) << 12
	//	i288      = ((_1011011 + i268) << 9 + _11000011) << 8 + _11100111
	//	return      (i288 << 7 + _1110101) << 5 + _11
	//
	// Operations: 246 squares 56 multiplies

	// Allocate Temporaries.
	var (
		t0  = new(Element)
		t1  = new(Element)
		t2  = new(Element)
		t3  = new(Element)
		t4  = new(Element)
		t5  = new(Element)
		t6  = new(Element)
		t7  = new(Element)
		t8  = new(Element)
		t9  = new(Element)
		t10 = new(Element)
		t11 = new(Element)
		t12 = new(Element)
		t13 = new(Element)
		t14 = new(Element)
		t15 = new(Element)
		t16 = new(Element)
		t17 = new(Element)
		t18 = new(Element)
		t19 = new(Element)
		t20 = new(Element)
	)

	// var t0,t1,t2,t3,t4,t5,t6,t7,t8,t9,t10,t11,t12,t13,t14,t15,t16,t17,t18,t19,t20 Element
	// Step 1: t8 = x⁰x2
	t8.Square(&x)

	// Step 2: z = x⁰x3
	z.Mul(&x, t8)

	// Step 3: t2 = x⁰x5
	t2.Mul(t8, z)

	// Step 4: t1 = x⁰x6
	t1.Mul(&x, t2)

	// Step 5: t3 = x⁰x8
	t3.Mul(t8, t1)

	// Step 6: t9 = x⁰xd
	t9.Mul(t2, t3)

	// Step 7: t6 = x⁰x12
	t6.Mul(t2, t9)

	// Step 8: t18 = x⁰x13
	t18.Mul(&x, t6)

	// Step 9: t0 = x⁰x14
	t0.Mul(&x, t18)

	// Step 10: t19 = x⁰x17
	t19.Mul(z, t0)

	// Step 11: t2 = x⁰x1c
	t2.Mul(t2, t19)

	// Step 12: t16 = x⁰x20
	t16.Mul(t9, t18)

	// Step 13: t4 = x⁰x23
	t4.Mul(z, t16)

	// Step 14: t14 = x⁰x2b
	t14.Mul(t3, t4)

	// Step 15: t12 = x⁰x2f
	t12.Mul(t18, t2)

	// Step 16: t15 = x⁰x41
	t15.Mul(t6, t12)

	// Step 17: t17 = x⁰x53
	t17.Mul(t6, t15)

	// Step 18: t3 = x⁰x5b
	t3.Mul(t3, t17)

	// Step 19: t5 = x⁰x61
	t5.Mul(t1, t3)

	// Step 20: t0 = x⁰x75
	t0.Mul(t0, t5)

	// Step 21: t10 = x⁰x91
	t10.Mul(t2, t0)

	// Step 22: t7 = x⁰x95
	t7.Mul(t16, t0)

	// Step 23: t11 = x⁰xb5
	t11.Mul(t16, t7)

	// Step 24: t13 = x⁰xbb
	t13.Mul(t1, t11)

	// Step 25: t20 = x⁰xc1
	t20.Mul(t1, t13)

	// Step 26: t2 = x⁰xc3
	t2.Mul(t8, t20)

	// Step 27: t6 = x⁰xd3
	t6.Mul(t6, t20)

	// Step 28: t16 = x⁰xe1
	t16.Mul(t16, t20)

	// Step 29: t8 = x⁰xe3
	t8.Mul(t8, t16)

	// Step 30: t1 = x⁰xe7
	t1.Mul(t1, t16)

	// Step 38: t20 = x⁰xc100
	for s := 0; s < 8; s++ {
		t20.Square(t20)
	}

	// Step 39: t20 = x⁰xc191
	t20.Mul(t10, t20)

	// Step 49: t20 = x⁰x3064400
	for s := 0; s < 10; s++ {
		t20.Square(t20)
	}

	// Step 50: t20 = x⁰x30644e7
	t20.Mul(t1, t20)

	// Step 57: t20 = x⁰x183227380
	for s := 0; s < 7; s++ {
		t20.Square(t20)
	}

	// Step 58: t19 = x⁰x183227397
	t19.Mul(t19, t20)

	// Step 67: t19 = x⁰x30644e72e00
	for s := 0; s < 9; s++ {
		t19.Square(t19)
	}

	// Step 68: t18 = x⁰x30644e72e13
	t18.Mul(t18, t19)

	// Step 75: t18 = x⁰x1832273970980
	for s := 0; s < 7; s++ {
		t18.Square(t18)
	}

	// Step 76: t18 = x⁰x183227397098d
	t18.Mul(t9, t18)

	// Step 90: t18 = x⁰x60c89ce5c2634000
	for s := 0; s < 14; s++ {
		t18.Square(t18)
	}

	// Step 91: t17 = x⁰x60c89ce5c2634053
	t17.Mul(t17, t18)

	// Step 100: t17 = x⁰xc19139cb84c680a600
	for s := 0; s < 9; s++ {
		t17.Square(t17)
	}

	// Step 101: t16 = x⁰xc19139cb84c680a6e1
	t16.Mul(t16, t17)

	// Step 109: t16 = x⁰xc19139cb84c680a6e100
	for s := 0; s < 8; s++ {
		t16.Square(t16)
	}

	// Step 110: t15 = x⁰xc19139cb84c680a6e141
	t15.Mul(t15, t16)

	// Step 120: t15 = x⁰x30644e72e131a029b850400
	for s := 0; s < 10; s++ {
		t15.Square(t15)
	}

	// Step 121: t15 = x⁰x30644e72e131a029b85045b
	t15.Mul(t3, t15)

	// Step 126: t15 = x⁰x60c89ce5c263405370a08b60
	for s := 0; s < 5; s++ {
		t15.Square(t15)
	}

	// Step 127: t15 = x⁰x60c89ce5c263405370a08b6d
	t15.Mul(t9, t15)

	// Step 135: t15 = x⁰x60c89ce5c263405370a08b6d00
	for s := 0; s < 8; s++ {
		t15.Square(t15)
	}

	// Step 136: t15 = x⁰x60c89ce5c263405370a08b6d03
	t15.Mul(z, t15)

	// Step 148: t15 = x⁰x60c89ce5c263405370a08b6d03000
	for s := 0; s < 12; s++ {
		t15.Square(t15)
	}

	// Step 149: t14 = x⁰x60c89ce5c263405370a08b6d0302b
	t14.Mul(t14, t15)

	// Step 161: t14 = x⁰x60c89ce5c263405370a08b6d0302b000
	for s := 0; s < 12; s++ {
		t14.Square(t14)
	}

	// Step 162: t13 = x⁰x60c89ce5c263405370a08b6d0302b0bb
	t13.Mul(t13, t14)

	// Step 170: t13 = x⁰x60c89ce5c263405370a08b6d0302b0bb00
	for s := 0; s < 8; s++ {
		t13.Square(t13)
	}

	// Step 171: t12 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f
	t12.Mul(t12, t13)

	// Step 185: t12 = x⁰x183227397098d014dc2822db40c0ac2ecbc000
	for s := 0; s < 14; s++ {
		t12.Square(t12)
	}

	// Step 186: t11 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b5
	t11.Mul(t11, t12)

	// Step 195: t11 = x⁰x30644e72e131a029b85045b68181585d97816a00
	for s := 0; s < 9; s++ {
		t11.Square(t11)
	}

	// Step 196: t10 = x⁰x30644e72e131a029b85045b68181585d97816a91
	t10.Mul(t10, t11)

	// Step 201: t10 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d5220
	for s := 0; s < 5; s++ {
		t10.Square(t10)
	}

	// Step 202: t9 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d
	t9.Mul(t9, t10)

	// Step 214: t9 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d000
	for s := 0; s < 12; s++ {
		t9.Square(t9)
	}

	// Step 215: t8 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e3
	t8.Mul(t8, t9)

	// Step 223: t8 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e300
	for s := 0; s < 8; s++ {
		t8.Square(t8)
	}

	// Step 224: t7 = x⁰x60c89ce5c263405370a08b6d0302b0bb2f02d522d0e395
	t7.Mul(t7, t8)

	// Step 235: t7 = x⁰x30644e72e131a029b85045b68181585d97816a916871ca800
	for s := 0; s < 11; s++ {
		t7.Square(t7)
	}

	// Step 236: t6 = x⁰x30644e72e131a029b85045b68181585d97816a916871ca8d3
	t6.Mul(t6, t7)

	// Step 243: t6 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e546980
	for s := 0; s < 7; s++ {
		t6.Square(t6)
	}

	// Step 244: t5 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e1
	t5.Mul(t5, t6)

	// Step 255: t5 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f0800
	for s := 0; s < 11; s++ {
		t5.Square(t5)
	}

	// Step 256: t4 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f0823
	t4.Mul(t4, t5)

	// Step 268: t4 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f0823000
	for s := 0; s < 12; s++ {
		t4.Square(t4)
	}

	// Step 269: t3 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b
	t3.Mul(t3, t4)

	// Step 278: t3 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b600
	for s := 0; s < 9; s++ {
		t3.Square(t3)
	}

	// Step 279: t2 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c3
	t2.Mul(t2, t3)

	// Step 287: t2 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c300
	for s := 0; s < 8; s++ {
		t2.Square(t2)
	}

	// Step 288: t1 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c3e7
	t1.Mul(t1, t2)

	// Step 295: t1 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f380
	for s := 0; s < 7; s++ {
		t1.Square(t1)
	}

	// Step 296: t0 = x⁰xc19139cb84c680a6e14116da060561765e05aa45a1c72a34f082305b61f3f5
	t0.Mul(t0, t1)

	// Step 301: t0 = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c3e7ea0
	for s := 0; s < 5; s++ {
		t0.Square(t0)
	}

	// Step 302: z = x⁰x183227397098d014dc2822db40c0ac2ecbc0b548b438e5469e10460b6c3e7ea3
	z.Mul(z, t0)

	return z
}
