// Package cmn provides common low-level types and utilities for all aistore projects
/*
 * Copyright (c) 2018, NVIDIA CORPORATION. All rights reserved.
 */
package cmn

import (
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
)

const (
	BadCksumPrefix = "BAD CHECKSUM:"
	MLCG32         = 1103515245 // xxhash seed
)

func NewCRC32C() hash.Hash {
	return crc32.New(crc32.MakeTable(crc32.Castagnoli))
}

type (
	Cksummer interface {
		Get() (string, string)
		String() string
	}
	cksum struct {
		ty    string
		value string
	}
)

func NewCksum(ty, value string) Cksummer {
	if ty == "" || value == "" {
		return nil
	}
	if ty != ChecksumXXHash && ty != ChecksumMD5 && ty != ChecksumCRC32C {
		AssertMsg(false, fmt.Sprintf("invalid checksum type: %s (with value of: %s)", ty, value))
	}
	return cksum{ty, value}
}

func HashToStr(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

func EqCksum(a, b Cksummer) bool {
	if a == nil || b == nil {
		return false
	}
	t1, v1 := a.Get()
	t2, v2 := b.Get()
	return t1 == t2 && v1 == v2
}

func BadCksum(a, b Cksummer) string {
	if a != nil && b == nil {
		return fmt.Sprintf("%s (%s != %v)", BadCksumPrefix, a, b)
	} else if a == nil && b != nil {
		return fmt.Sprintf("%s (%v != %s)", BadCksumPrefix, a, b)
	} else if a == nil && b == nil {
		return fmt.Sprintf("%s (nil != nil)", BadCksumPrefix)
	}
	t1, v1 := a.Get()
	t2, v2 := b.Get()
	if t1 == t2 {
		return fmt.Sprintf("%s %s(%s != %s)", BadCksumPrefix, t1, v1, v2)
	}
	return fmt.Sprintf("%s %s != %s", BadCksumPrefix, a, b)
}

func (v cksum) Get() (string, string) { return v.ty, v.value }
func (v cksum) Type() string          { return v.ty }
func (v cksum) Value() string         { return v.value }
func (v cksum) String() string {
	return fmt.Sprintf("(%s,%s...)", v.ty, v.value[:Min(10, len(v.value))])
}
