// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"runtime"
	"sync"
)

var (
	helperMutex sync.RWMutex
	helpers     map[uintptr]struct{}
)

// Helper passed 1 marks the caller as a helper function and will skip it when
// capturing the callsite location.
func Helper(skip int) {
	helper(skip + 1)
}

func helper(skip int) {
	callers := [1]uintptr{}
	if runtime.Callers(skip+2, callers[:]) == 0 {
		panic("failed to get caller information")
	}
	pc := callers[0]
	helperMutex.RLock()
	if _, ok := helpers[pc]; ok {
		helperMutex.RUnlock()
		return
	}
	helperMutex.RUnlock()
	helperMutex.Lock()
	defer helperMutex.Unlock()
	if helpers == nil {
		helpers = make(map[uintptr]struct{})
	}
	helpers[pc] = struct{}{}
}

// caller behaves like runtime.Caller but skips functions marked by helper.
func caller(skip int) (uintptr, string, int, bool) {
	pc := [8]uintptr{}
	n := runtime.Callers(skip+2, pc[:])
	if n == 0 {
		return 0, "", 0, false
	}
	helperMutex.RLock()
	pcs := pc[:]
	for i := n - 1; i >= 0; i-- {
		if _, ok := helpers[pc[i]]; ok {
			pcs = pc[i:]
			break
		}
	}
	helperMutex.RUnlock()
	frames := runtime.CallersFrames(pcs)
	if frame, ok := frames.Next(); ok {
		return frame.PC, frame.File, frame.Line, true
	}
	return 0, "", 0, false
}
