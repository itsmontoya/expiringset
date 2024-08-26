package expiringset

import (
	"sync"
	"time"
)

func New(ttl time.Duration) *ExpiringSet {
	var e ExpiringSet
	e.ttl = ttl
	e.byKey = make(map[string]time.Time)
	e.expirationList = make([]entry, 0, 32)
	return &e
}

type ExpiringSet struct {
	mux sync.RWMutex

	ttl time.Duration

	// Lookup table by key with a value of a unix timestamp representing the last accessed time
	byKey map[string]time.Time
	// Expiration list used to see what keys are expired
	expirationList []entry

	cancel *time.Timer
}

func (e *ExpiringSet) Set(key string) (inserted bool) {
	e.mux.Lock()
	defer e.mux.Unlock()

	inserted = e.set(key)
	e.refreshTimer()
	return
}

func (e *ExpiringSet) unset(key string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	lastAccessed, ok := e.byKey[key]
	if !ok {
		return
	}

	now := time.Now()
	if now.Sub(lastAccessed) < e.ttl {
		return
	}

	delete(e.byKey, key)
	e.cancel = nil
	e.refreshTimer()
}

func (e *ExpiringSet) set(key string) (inserted bool) {
	now := time.Now()
	expiresAt := now.Add(e.ttl)
	_, exists := e.byKey[key]
	inserted = !exists
	e.byKey[key] = now
	e.expirationList = append(e.expirationList, entry{expiresAt: expiresAt, key: key})
	return
}

func (e *ExpiringSet) refreshTimer() {
	if e.cancel != nil {
		return
	}

	entry, ok := e.shift()
	if !ok {
		return
	}

	now := time.Now()
	delta := entry.expiresAt.Sub(now)
	e.cancel = time.AfterFunc(delta, e.remove(entry))
}

func (e *ExpiringSet) shift() (out entry, ok bool) {
	if len(e.expirationList) == 0 {
		return
	}

	out = e.expirationList[0]
	e.expirationList = e.expirationList[1:]
	ok = true
	return
}

func (e *ExpiringSet) remove(in entry) func() {
	return func() {
		e.unset(in.key)
	}
}

type entry struct {
	expiresAt time.Time
	key       string
}
