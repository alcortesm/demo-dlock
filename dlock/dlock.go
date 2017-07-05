package dlock

// DLock defines the operations on a distributed lock.
//
// This is incomplete, it is missing leases and/or TTLs.
type DLock interface {
	TryLock() (bool, error)
	Unlock() error
}
