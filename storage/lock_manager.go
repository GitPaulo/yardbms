package storage

import "sync"

type LockManager struct {
	tableLocks map[string]*sync.RWMutex
	lock       sync.Mutex
}

func NewLockManager() *LockManager {
	return &LockManager{
		tableLocks: make(map[string]*sync.RWMutex),
	}
}

func (lm *LockManager) getLock(tableName string) *sync.RWMutex {
	lm.lock.Lock()
	defer lm.lock.Unlock()
	if _, exists := lm.tableLocks[tableName]; !exists {
		lm.tableLocks[tableName] = &sync.RWMutex{}
	}
	return lm.tableLocks[tableName]
}

func (lm *LockManager) LockTable(tableName string) {
	lm.getLock(tableName).Lock()
}

func (lm *LockManager) UnlockTable(tableName string) {
	lm.getLock(tableName).Unlock()
}

func (lm *LockManager) RLockTable(tableName string) {
	lm.getLock(tableName).RLock()
}

func (lm *LockManager) RUnlockTable(tableName string) {
	lm.getLock(tableName).RUnlock()
}
