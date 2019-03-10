package zzcache

// Hasher provides implementation of the hash function
type Hasher interface {
	Do(string) uint64
}
