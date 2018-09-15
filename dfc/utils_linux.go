// Package dfc is a scalable object-storage based caching system with Amazon and Google Cloud backends.
/*
 * Copyright (c) 2018, NVIDIA CORPORATION. All rights reserved.
 *
 */
package dfc

import (
	"fmt"
	"syscall"

	"github.com/NVIDIA/dfcpub/common"
)

// Get specific attribute for specified fqn.
func Getxattr(fqn string, attrname string) ([]byte, string) {
	data := make([]byte, maxAttrSize)
	read, err := syscall.Getxattr(fqn, attrname, data)
	common.Assert(read < maxAttrSize)
	if err != nil && err != syscall.ENODATA {
		return nil, fmt.Sprintf("Failed to get xattr %s for %s, err: %v", attrname, fqn, err)
	}
	if read > 0 {
		return data[:read], ""
	}
	return nil, ""
}

// Set specific named attribute for specific fqn.
func Setxattr(fqn string, attrname string, data []byte) (errstr string) {
	common.Assert(len(data) < maxAttrSize)
	err := syscall.Setxattr(fqn, attrname, data, 0)
	if err != nil {
		errstr = fmt.Sprintf("Failed to set extended attr for fqn %s attr %s, err: %v",
			fqn, attrname, err)
	}
	return
}

// Delete specific named attribute for specific fqn.
func Deletexattr(fqn string, attrname string) (errstr string) {
	err := syscall.Removexattr(fqn, attrname)
	if err != nil {
		errstr = fmt.Sprintf("Failed to remove extended attr for fqn %s attr %s, err: %v",
			fqn, attrname, err)
	}
	return
}

func TotalMemory() (mb uint64, err error) {
	sysinfo := &syscall.Sysinfo_t{}
	if err = syscall.Sysinfo(sysinfo); err != nil {
		return
	}
	mb = sysinfo.Totalram * uint64(sysinfo.Unit) / common.MiB
	return
}
