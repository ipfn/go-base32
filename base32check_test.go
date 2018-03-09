// Copyright (c) 2018 The IPFN Developers
// Copyright (c) 2013-2018 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base32check

import (
	"testing"
)

var checkEncodingStringTests = []struct {
	version byte
	in      string
	out     string
}{
	{20, "", "znts8576"},
	{20, " ", "zssfh0qjhy"},
	{20, "-", "zskjvbs4yv"},
	{20, "0", "zscbdprtm0"},
	{20, "1", "zsc5pjbhty"},
	{20, "-1", "zsknzyyc2x70"},
	{20, "11", "zscnz9a2ub50"},
	{20, "abc", "z3skyce97vrm5"},
	{20, "1234598760", "zscnyve5x5unsrekxpesys8w"},
	{20, "abcdefghijklmnopqrstuvwxyz", "z3skycmyv4nxw6qfrf4kcmtwrac8zunnw36hvamcd9af9qmy85"},
	{20, "00000000000000000000000000000000000000000000000000000000000000", "zscq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cqplklr6p0"},
}

func TestBase32Check(t *testing.T) {
	for x, test := range checkEncodingStringTests {
		// test encoding
		if res := CheckEncodeToString([]byte(test.in), test.version); test.out != res {
			t.Errorf("CheckEncode test #%d failed: got %s, want: %s", x, res, test.out)
		}

		// test decoding
		res, version, err := CheckDecodeString(test.out)
		if err != nil {
			t.Errorf("CheckDecodeString test #%d failed with err: %v", x, err)
		} else if version != test.version {
			t.Errorf("CheckDecodeString test #%d failed: got version: %d want: %d", x, version, test.version)
		} else if string(res) != test.in {
			t.Errorf("CheckDecodeString test #%d failed: got: %s want: %s", x, res, test.in)
		}
	}

	// test the two decoding failure cases
	// case 1: checksum error
	_, _, err := CheckDecodeString("znts8575")
	if err != ErrChecksum {
		t.Error("Checkdecode test failed, expected ErrChecksum")
	}
	// case 2: invalid formats (string lengths below 5 mean the version byte and/or the checksum
	// bytes are missing).
	testString := ""
	for len := 0; len < 4; len++ {
		// make a string of length `len`
		_, _, err = CheckDecodeString(testString)
		if err != ErrInvalidFormat {
			t.Error("Checkdecode test failed, expected ErrInvalidFormat")
		}
	}
}

func TestBase32CheckZeros(t *testing.T) {
	res := CheckEncodeToString([]byte{0, 0, 0, 0, 0, 0, 0, 0, 123}, 0x0)
	if res != "00000000000000qmdn28lms" {
		t.Errorf("CheckEncodeZeros failed: got %s", res)
	}
}
