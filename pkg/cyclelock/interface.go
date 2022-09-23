package cyclelock

// CycleLock is used to enter and exit Cycles to get the emulation cycle correct
type CycleLock interface {
	EnterCycle()
	ExitCycle()
	CycleCount() int
	ResetCycleCount()
	Unlock()
	WaitForLock()
}
