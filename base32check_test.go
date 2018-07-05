// Copyright (c) 2018 The IPFN Developers
// Copyright (c) 2013-2018 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base32check

import (
	"testing"
)

var checkEncodingStringTests = []struct {
	in  string
	out string
}{
	{"", "00"},
	{" ", "yqzs"},
	{"-", "9hu0"},
	{"0", "xzss"},
	{"1", "xxms"},
	{"-1", "95c65"},
	{"11", "xyclw"},
	{"abc", "v93x8ss"},
	{"1234598760", "xyeqxrp48yuqwr3sj0"},
	{"abcdefghijklmnopqrstuvwxyz", "v93xxeq9venks6t2rrkx6mndwpchyum5w4m8w7qed27s"},
	{"00000000000000000000000000000000000000000000000000000000000000", "x0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cq0vpsx0cdx"},
}

func TestBase32Check(t *testing.T) {
	for x, test := range checkEncodingStringTests {
		// test encoding
		encoded := CheckEncodeToString([]byte(test.in))
		if test.out != encoded {
			t.Errorf("CheckEncode test #%d failed: got %s, want: %s", x, encoded, test.out)
		}

		// test decoding
		res, err := CheckDecodeString(encoded)
		if err != nil {
			t.Errorf("CheckDecodeString test #%d failed with err: %v", x, err)
		} else if string(res) != test.in {
			t.Errorf("CheckDecodeString test #%d failed: got: %s want: %s", x, res, test.in)
		}
	}

	// test the two decoding failure cases
	// case 1: checksum error
	_, err := CheckDecodeString("yqzx")
	if err != ErrChecksum {
		t.Error("Checkdecode test failed, expected ErrChecksum")
	}
	// case 2: invalid formats (string lengths below 5 mean the version byte and/or the checksum
	// bytes are missing).
	testString := ""
	for len := 0; len < 4; len++ {
		// make a string of length `len`
		_, err = CheckDecodeString(testString)
		if err != ErrInvalidFormat {
			t.Error("Checkdecode test failed, expected ErrInvalidFormat")
		}
	}
}

func TestBase32CheckZeros(t *testing.T) {
	res := CheckEncodeToString([]byte{0, 0, 0, 0, 0, 0, 0, 0, 123})
	if res != "00000000000007u6" {
		t.Errorf("CheckEncodeZeros failed: got %s", res)
	}
}
