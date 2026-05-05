// ABOUTME: Deterministic byte serialization helpers for hashing and
// ABOUTME: signing, producing canonical output across Go versions.
package quark

import (
	"encoding/binary"
	"hash"
)

func writeBytes(h hash.Hash, b []byte) {
	var l [4]byte
	binary.BigEndian.PutUint32(l[:], uint32(len(b)))
	h.Write(l[:])
	h.Write(b)
}

func writeString(h hash.Hash, s string) {
	writeBytes(h, []byte(s))
}

func writeInt64(h hash.Hash, v int64) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(v))
	h.Write(b[:])
}

func writeInt32(h hash.Hash, v int32) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(v))
	h.Write(b[:])
}
