// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/20

package uid

import (
	"encoding/base32"
	"encoding/binary"
	"math"
	"sync"
	"time"
)

type UUID [22]byte

var (
	version   byte = 0x01 // set version to 1
	pid            = getPid()
	clockSeq       = randUint64()
	machineID      = readMachineID()
	mutex     sync.Mutex
)

func New() UUID {
	return fromTime(time.Now())
}

// 8 bytes of time (ns) + 1 bytes of version +
// first 3 bytes of md5(Machine hostname) + 2 byes of pid + 8 random bytes
func fromTime(aTime time.Time) UUID {
	var id UUID

	utcTime := aTime.In(time.UTC)
	// Timestamp ns, 8 bytes, big endian
	binary.BigEndian.PutUint64(id[:], uint64(utcTime.UnixNano()))
	// version, 1 bytes
	id[8] = version
	// Machine, first 3 bytes of md5(hostname)
	id[9] = machineID[0]
	id[10] = machineID[1]
	id[11] = machineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	id[12] = byte(pid >> 8)
	id[13] = byte(pid)
	// random, 8 bytes big endian
	mutex.Lock()
	clockSeq = (clockSeq + 1) & math.MaxUint64
	mutex.Unlock()
	id[14] = byte(clockSeq >> 56)
	id[15] = byte(clockSeq >> 48)
	id[16] = byte(clockSeq >> 40)
	id[17] = byte(clockSeq >> 32)
	id[18] = byte(clockSeq >> 24)
	id[19] = byte(clockSeq >> 16)
	id[20] = byte(clockSeq >> 8)
	id[21] = byte(clockSeq)
	return id
}

var encode = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv").WithPadding(base32.NoPadding)

func (u UUID) String() string {
	return encode.EncodeToString(u[:])
}
