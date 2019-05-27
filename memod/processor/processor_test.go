package processor

import "testing"

func TestFindBucketForFeerate(t *testing.T) {

	var feerateBuckets = [40]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 18, 22, 27, 33, 41, 50, 62, 76, 93, 114, 140, 172, 212, 261, 321, 395, 486, 598, 736, 905, 1113, 1369, 1684, 2071, 2547, 3133, 3854, 3855}

	// test a feerate of 0
	t1 := findBucketForFeerate(0)
	if t1 != 0 {
		t.Errorf("findBucketForFeerate(0) = %d; want 0", t1)
	}

	// test a feerate of 1
	t2 := findBucketForFeerate(1)
	if t2 != 0 {
		t.Errorf("findBucketForFeerate(1) = %d; want 0", t1)
	}

	// test a feerate of 2
	t3 := findBucketForFeerate(2)
	if t3 != 1 {
		t.Errorf("findBucketForFeerate(2) = %d; want 1", t3)
	}

	// test a feerate of 18
	t4 := findBucketForFeerate(18)
	if t4 != 12 {
		t.Errorf("findBucketForFeerate(18) = %d; want 12", t4)
	}

	// test a feerate of 19
	t5 := findBucketForFeerate(19)
	if t5 != 13 {
		t.Errorf("findBucketForFeerate(19) = %d; want 12, bucket = %d", t5, feerateBuckets[t5])
	}

	// test a feerate of -1
	t6 := findBucketForFeerate(-1)
	if t6 != 0 {
		t.Errorf("findBucketForFeerate(-1) = %d; want 0, bucket = %d", t6, feerateBuckets[t6])
	}

	// test a feerate of 19
	t10 := findBucketForFeerate(4000)
	if t10 != 39 {
		t.Errorf("findBucketForFeerate(4000) = %d; want 39", t10)
	}
}
