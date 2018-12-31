// Copyright © 2018 The IPFN Developers. All Rights Reserved.
// Copyright © 2013-2014 The btcsuite developers. All Rights Reserved.
//
// Use of this source code is governed by an ISC license.
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package base32i

import (
	"errors"
	"hash/crc32"
)

// ErrChecksum indicates that the checksum of a check-encoded string does not verify against
// the checksum.
var ErrChecksum = errors.New("checksum error")

// ErrInvalidFormat indicates that the check-encoded string has an invalid format.
var ErrInvalidFormat = errors.New("invalid format: checksum bytes missing")

const cSize = 1

func checksum(input []byte) (cksum byte) {
	vint := encodeVarint(uint64(crc32.ChecksumIEEE(input)))
	return vint[len(vint)-1]
}

// source: github.com/gogo/protobuf/proto
func encodeVarint(x uint64) []byte {
	var buf [10]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

// CheckEncode prepends and appends a four byte checksum.
func CheckEncode(input []byte) []byte {
	return Encode(checkBuffer(input))
}

// CheckEncodeToString is CheckEncode to string.
func CheckEncodeToString(input []byte) string {
	return EncodeToString(checkBuffer(input))
}

// CheckEncodePrefixed is CheckEncode to string with multibase prefix 'i'.
func CheckEncodePrefixed(input []byte) string {
	input = checkBuffer(input)
	buf := make([]byte, Encoding.EncodedLen(len(input))+1)
	buf[0] = 'i'
	Encoding.Encode(buf[1:], input)
	return string(buf)
}

// CheckDecodeString decodes a string that was encoded with CheckEncode and verifies the checksum.
func CheckDecodeString(input string) (result []byte, err error) {
	decoded, err := DecodeString(input)
	if err != nil {
		return
	}
	if len(decoded) < 1 {
		err = ErrInvalidFormat
		return
	}
	cksum := decoded[len(decoded)-1]
	result = decoded[:len(decoded)-1]
	if checksum(result) != cksum {
		err = ErrChecksum
		return
	}
	return
}

func checkBuffer(input []byte) []byte {
	return append(input, checksum(input))
}
