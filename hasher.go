package zzcache

import "hash/crc32"

// Hasher provides implementation of the hash function
type Hasher interface {
	Do([]byte) uint32
}

// CRC32 provides implementation of using golang
// crc32 algorithm
type CRC32 struct {
}

// Do returns CRC32 checksum
func (c *CRC32) Do(b []byte) uint32 {
	return crc32.ChecksumIEEE(b)
}
