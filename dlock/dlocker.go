package dlock

// DLock defines the operations on a distributed lock.
//
// This is incomplete, it is missing leases and/or TTLs.
type DLocker interface {
	Lock() error
	Unlock() error
}
