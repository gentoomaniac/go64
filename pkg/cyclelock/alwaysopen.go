package cyclelock

// AlwaysOpenLock is a CycleLock for testing that doesn't lock at all so that tests can run with full speed
type AlwaysOpenLock struct {
	cycleCount int
}

// EnterCycle simulates waiting for a cpu cycle and does nothing
func (a *AlwaysOpenLock) EnterCycle() {
	a.cycleCount++
}

// ExitCycle simulates finishing the cycle and does nothing
func (a *AlwaysOpenLock) ExitCycle() {}

// CycleCount returns the number of cycles since the last ResetCycleCount()
func (a *AlwaysOpenLock) CycleCount() int {
	return a.cycleCount
}

// ResetCycleCount resets the cycle counter to 0
func (a *AlwaysOpenLock) ResetCycleCount() {
	a.cycleCount = 0
}

// WaitForLock waits for the next ExitCycle() call
func (a *AlwaysOpenLock) WaitForLock() {}

// Unlock simply unlocks and is supposed to be called from the controlling thread
func (a *AlwaysOpenLock) Unlock() {}
