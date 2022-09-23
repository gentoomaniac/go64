package cyclelock

// ChannelLock uses a sync.mutex to lock the ressource
type ChannelLock struct {
	cycleCount int
	lock       chan bool
}

// Init initialises the lock
func (l *ChannelLock) Init() {
	l.lock = make(chan bool, 1)
}

// EnterCycle simulates waiting for a cpu cycle and does nothing
func (l *ChannelLock) EnterCycle() {
	<-l.lock
	l.cycleCount++
}

// ExitCycle simulates finishing the cycle and does nothing
func (l *ChannelLock) ExitCycle() {
	l.lock <- true
}

// CycleCount returns the number of cycles since the last ResetCycleCount()
func (l *ChannelLock) CycleCount() int {
	return l.cycleCount
}

// ResetCycleCount resets the cycle counter to 0
func (l *ChannelLock) ResetCycleCount() {
	l.cycleCount = 0
}

// Unlock simply unlocks and is supposed to be called from the controlling thread
func (l *ChannelLock) Unlock() {
	l.lock <- true
}

// WaitForLock waits for the next ExitCycle() call
func (l *ChannelLock) WaitForLock() {
	<-l.lock
}
