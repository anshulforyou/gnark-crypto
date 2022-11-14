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

package bn254

const MAX_BATCH_SIZE = 600

type batchOp struct {
	bucketID, pointID uint32
}

func (o batchOp) isNeg() bool {
	return o.pointID&1 == 1
}

// processChunkG1BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG1BatchAffine[B ibG1Affine, BS bitSet](chunk uint64,
	chRes chan<- g1JacExtended,
	c uint64,
	points []G1Affine,
	digits []uint32) {

	// init the buckets
	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// setup for the batch affine;
	// we do that instead of a separate object to give enough hints to the compiler to..
	// keep things on the stack.
	batchSize := len(buckets) / 20
	if batchSize > MAX_BATCH_SIZE {
		batchSize = MAX_BATCH_SIZE
	}
	if batchSize <= 0 {
		batchSize = 1
	}
	var bucketIds BS // bitSet to signify presence of a bucket in current batch
	cptAdd := 0      // count the number of bucket + point added to current batch
	cptSub := 0      // count the number of bucket - point added to current batch

	var P [MAX_BATCH_SIZE]*G1Affine // points to be added to R (buckets)
	var R [MAX_BATCH_SIZE]*G1Affine // bucket references

	canAdd := func(bID uint32) bool {
		return !bucketIds[bID]
	}

	isFull := func() bool {
		return (cptAdd + cptSub) == batchSize
	}

	executeAndReset := func() {
		if (cptAdd + cptSub) == 0 {
			return
		}
		batchAddG1Affine(R[:batchSize], P[:batchSize], cptAdd, cptSub)

		var tmp BS
		bucketIds = tmp
		cptAdd = 0
		cptSub = 0
	}

	add := func(op batchOp) {
		// CanAdd must be called before --> ensures bucket is not "used" in current batch

		BK := &buckets[op.bucketID]
		PP := &points[op.pointID>>1]
		if PP.IsInfinity() {
			return
		}
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			if op.isNeg() {
				BK.Neg(PP)
			} else {
				BK.Set(PP)
			}
			return
		}
		if BK.X.Equal(&PP.X) {
			if BK.Y.Equal(&PP.Y) {
				if op.isNeg() {
					// P + -P
					BK.setInfinity()
					return
				}
				// P + P: doubling, which should be quite rare -- may want to put it back in the batch add?
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				BK.Add(BK, BK)
				return
			}
			// b.Y == -p.Y
			if op.isNeg() {
				// doubling .
				BK.Add(BK, BK)
				return
			}
			BK.setInfinity()
			return
		}

		bucketIds[op.bucketID] = true
		if op.isNeg() {
			cptSub++
			R[batchSize-cptSub] = BK
			P[batchSize-cptSub] = PP
		} else {
			R[cptAdd] = BK
			P[cptAdd] = PP
			cptAdd++
		}

	}

	var queue [MAX_BATCH_SIZE]batchOp
	qID := 0

	processQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if canAdd(queue[i].bucketID) {
				add(queue[i])
				if isFull() {
					executeAndReset()
				}
				queue[i] = queue[qID-1]
				qID--
			}
		}
	}

	processTopQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if !canAdd(queue[i].bucketID) {
				return
			}
			add(queue[i])
			if isFull() {
				executeAndReset()
			}
			qID--
		}
	}

	for i, digit := range digits {

		if digit == 0 {
			continue
		}

		op := batchOp{pointID: uint32(i) << 1}
		// if msbWindow bit is set, we need to substract
		if digit&1 == 0 {
			// add
			op.bucketID = uint32((digit >> 1) - 1)
		} else {
			// sub
			op.bucketID = (uint32((digit >> 1)))
			op.pointID += 1
		}
		if canAdd(op.bucketID) {
			add(op)
			if isFull() {
				executeAndReset()
				processTopQueue()
			}
		} else {
			// put it in queue.
			queue[qID] = op
			qID++
			if qID == MAX_BATCH_SIZE-1 {
				executeAndReset()
				processQueue()
			}
			// queue = append(queue, op)
		}
	}

	for qID != 0 {
		processQueue()
		executeAndReset() // execute batch even if not full.
	}

	// flush items in batch.
	executeAndReset()

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g1JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].IsInfinity() {
			runningSum.addMixed(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total

}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketG1AffineC4 [1 << (4 - 1)]G1Affine
type bucketG1AffineC5 [1 << (5 - 1)]G1Affine
type bucketG1AffineC6 [1 << (6 - 1)]G1Affine
type bucketG1AffineC7 [1 << (7 - 1)]G1Affine
type bucketG1AffineC8 [1 << (8 - 1)]G1Affine
type bucketG1AffineC9 [1 << (9 - 1)]G1Affine
type bucketG1AffineC10 [1 << (10 - 1)]G1Affine
type bucketG1AffineC11 [1 << (11 - 1)]G1Affine
type bucketG1AffineC12 [1 << (12 - 1)]G1Affine
type bucketG1AffineC13 [1 << (13 - 1)]G1Affine
type bucketG1AffineC14 [1 << (14 - 1)]G1Affine
type bucketG1AffineC15 [1 << (15 - 1)]G1Affine
type bucketG1AffineC16 [1 << (16 - 1)]G1Affine

type ibG1Affine interface {
	bucketG1AffineC4 |
		bucketG1AffineC5 |
		bucketG1AffineC6 |
		bucketG1AffineC7 |
		bucketG1AffineC8 |
		bucketG1AffineC9 |
		bucketG1AffineC10 |
		bucketG1AffineC11 |
		bucketG1AffineC12 |
		bucketG1AffineC13 |
		bucketG1AffineC14 |
		bucketG1AffineC15 |
		bucketG1AffineC16
}

// processChunkG2BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG2BatchAffine[B ibG2Affine, BS bitSet](chunk uint64,
	chRes chan<- g2JacExtended,
	c uint64,
	points []G2Affine,
	digits []uint32) {

	// init the buckets
	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// setup for the batch affine;
	// we do that instead of a separate object to give enough hints to the compiler to..
	// keep things on the stack.
	batchSize := len(buckets) / 20
	if batchSize > MAX_BATCH_SIZE {
		batchSize = MAX_BATCH_SIZE
	}
	if batchSize <= 0 {
		batchSize = 1
	}
	var bucketIds BS // bitSet to signify presence of a bucket in current batch
	cptAdd := 0      // count the number of bucket + point added to current batch
	cptSub := 0      // count the number of bucket - point added to current batch

	var P [MAX_BATCH_SIZE]*G2Affine // points to be added to R (buckets)
	var R [MAX_BATCH_SIZE]*G2Affine // bucket references

	canAdd := func(bID uint32) bool {
		return !bucketIds[bID]
	}

	isFull := func() bool {
		return (cptAdd + cptSub) == batchSize
	}

	executeAndReset := func() {
		if (cptAdd + cptSub) == 0 {
			return
		}
		batchAddG2Affine(R[:batchSize], P[:batchSize], cptAdd, cptSub)

		var tmp BS
		bucketIds = tmp
		cptAdd = 0
		cptSub = 0
	}

	add := func(op batchOp) {
		// CanAdd must be called before --> ensures bucket is not "used" in current batch

		BK := &buckets[op.bucketID]
		PP := &points[op.pointID>>1]
		if PP.IsInfinity() {
			return
		}
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			if op.isNeg() {
				BK.Neg(PP)
			} else {
				BK.Set(PP)
			}
			return
		}
		if BK.X.Equal(&PP.X) {
			if BK.Y.Equal(&PP.Y) {
				if op.isNeg() {
					// P + -P
					BK.setInfinity()
					return
				}
				// P + P: doubling, which should be quite rare -- may want to put it back in the batch add?
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				BK.Add(BK, BK)
				return
			}
			// b.Y == -p.Y
			if op.isNeg() {
				// doubling .
				BK.Add(BK, BK)
				return
			}
			BK.setInfinity()
			return
		}

		bucketIds[op.bucketID] = true
		if op.isNeg() {
			cptSub++
			R[batchSize-cptSub] = BK
			P[batchSize-cptSub] = PP
		} else {
			R[cptAdd] = BK
			P[cptAdd] = PP
			cptAdd++
		}

	}

	var queue [MAX_BATCH_SIZE]batchOp
	qID := 0

	processQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if canAdd(queue[i].bucketID) {
				add(queue[i])
				if isFull() {
					executeAndReset()
				}
				queue[i] = queue[qID-1]
				qID--
			}
		}
	}

	processTopQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if !canAdd(queue[i].bucketID) {
				return
			}
			add(queue[i])
			if isFull() {
				executeAndReset()
			}
			qID--
		}
	}

	for i, digit := range digits {

		if digit == 0 {
			continue
		}

		op := batchOp{pointID: uint32(i) << 1}
		// if msbWindow bit is set, we need to substract
		if digit&1 == 0 {
			// add
			op.bucketID = uint32((digit >> 1) - 1)
		} else {
			// sub
			op.bucketID = (uint32((digit >> 1)))
			op.pointID += 1
		}
		if canAdd(op.bucketID) {
			add(op)
			if isFull() {
				executeAndReset()
				processTopQueue()
			}
		} else {
			// put it in queue.
			queue[qID] = op
			qID++
			if qID == MAX_BATCH_SIZE-1 {
				executeAndReset()
				processQueue()
			}
			// queue = append(queue, op)
		}
	}

	for qID != 0 {
		processQueue()
		executeAndReset() // execute batch even if not full.
	}

	// flush items in batch.
	executeAndReset()

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g2JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].IsInfinity() {
			runningSum.addMixed(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total

}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketG2AffineC4 [1 << (4 - 1)]G2Affine
type bucketG2AffineC5 [1 << (5 - 1)]G2Affine
type bucketG2AffineC6 [1 << (6 - 1)]G2Affine
type bucketG2AffineC7 [1 << (7 - 1)]G2Affine
type bucketG2AffineC8 [1 << (8 - 1)]G2Affine
type bucketG2AffineC9 [1 << (9 - 1)]G2Affine
type bucketG2AffineC10 [1 << (10 - 1)]G2Affine
type bucketG2AffineC11 [1 << (11 - 1)]G2Affine
type bucketG2AffineC12 [1 << (12 - 1)]G2Affine
type bucketG2AffineC13 [1 << (13 - 1)]G2Affine
type bucketG2AffineC14 [1 << (14 - 1)]G2Affine
type bucketG2AffineC15 [1 << (15 - 1)]G2Affine
type bucketG2AffineC16 [1 << (16 - 1)]G2Affine

type ibG2Affine interface {
	bucketG2AffineC4 |
		bucketG2AffineC5 |
		bucketG2AffineC6 |
		bucketG2AffineC7 |
		bucketG2AffineC8 |
		bucketG2AffineC9 |
		bucketG2AffineC10 |
		bucketG2AffineC11 |
		bucketG2AffineC12 |
		bucketG2AffineC13 |
		bucketG2AffineC14 |
		bucketG2AffineC15 |
		bucketG2AffineC16
}

type bitSetC4 [1 << (4 - 1)]bool
type bitSetC5 [1 << (5 - 1)]bool
type bitSetC6 [1 << (6 - 1)]bool
type bitSetC7 [1 << (7 - 1)]bool
type bitSetC8 [1 << (8 - 1)]bool
type bitSetC9 [1 << (9 - 1)]bool
type bitSetC10 [1 << (10 - 1)]bool
type bitSetC11 [1 << (11 - 1)]bool
type bitSetC12 [1 << (12 - 1)]bool
type bitSetC13 [1 << (13 - 1)]bool
type bitSetC14 [1 << (14 - 1)]bool
type bitSetC15 [1 << (15 - 1)]bool
type bitSetC16 [1 << (16 - 1)]bool

type bitSet interface {
	bitSetC4 |
		bitSetC5 |
		bitSetC6 |
		bitSetC7 |
		bitSetC8 |
		bitSetC9 |
		bitSetC10 |
		bitSetC11 |
		bitSetC12 |
		bitSetC13 |
		bitSetC14 |
		bitSetC15 |
		bitSetC16
}
