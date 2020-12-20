// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/20

// +build darwin

package uid

func readPlatformMachineID() (string, error) {
	return syscall.Sysctl("kern.uuid")
}
