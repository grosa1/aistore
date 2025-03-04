// Package cos provides common low-level types and utilities for all aistore projects
/*
 * Copyright (c) 2018-2022, NVIDIA CORPORATION. All rights reserved.
 */
package cos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NVIDIA/aistore/cmn/debug"
)

func IsParseBool(s string) bool {
	yes, err := ParseBool(s)
	_ = err // error means false
	return yes
}

// ParseBool converts string to bool (case-insensitive):
//
//	y, yes, on -> true
//	n, no, off, <empty value> -> false
//
// strconv handles the following:
//
//	1, true, t -> true
//	0, false, f -> false
func ParseBool(s string) (value bool, err error) {
	if s == "" {
		return
	}
	s = strings.ToLower(s)
	switch s {
	case "y", "yes", "on":
		return true, nil
	case "n", "no", "off":
		return false, nil
	}
	return strconv.ParseBool(s)
}

func StringSliceToIntSlice(strs []string) ([]int64, error) {
	res := make([]int64, 0, len(strs))
	for _, s := range strs {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func StrToSentence(str string) string {
	if str == "" {
		return ""
	}
	capitalized := CapitalizeString(str)
	if !strings.HasSuffix(capitalized, ".") {
		capitalized += "."
	}
	return capitalized
}

func ConvertToString(value any) (valstr string, err error) {
	switch v := value.(type) {
	case string:
		valstr = v
	case bool, int, int32, int64, uint32, uint64, float32, float64:
		valstr = fmt.Sprintf("%v", v)
	default:
		debug.FailTypeCast(value)
		err = fmt.Errorf("failed to assert type: %v(%T)", value, value)
	}
	return
}
