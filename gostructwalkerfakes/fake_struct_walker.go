// Code generated by counterfeiter. DO NOT EDIT.
package gostructwalkerfakes

import (
	"sync"

	"github.com/DanLavine/gostructwalker"
)

type FakeStructWalker struct {
	WalkStub        func(interface{}) error
	walkMutex       sync.RWMutex
	walkArgsForCall []struct {
		arg1 interface{}
	}
	walkReturns struct {
		result1 error
	}
	walkReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeStructWalker) Walk(arg1 interface{}) error {
	fake.walkMutex.Lock()
	ret, specificReturn := fake.walkReturnsOnCall[len(fake.walkArgsForCall)]
	fake.walkArgsForCall = append(fake.walkArgsForCall, struct {
		arg1 interface{}
	}{arg1})
	stub := fake.WalkStub
	fakeReturns := fake.walkReturns
	fake.recordInvocation("Walk", []interface{}{arg1})
	fake.walkMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeStructWalker) WalkCallCount() int {
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	return len(fake.walkArgsForCall)
}

func (fake *FakeStructWalker) WalkCalls(stub func(interface{}) error) {
	fake.walkMutex.Lock()
	defer fake.walkMutex.Unlock()
	fake.WalkStub = stub
}

func (fake *FakeStructWalker) WalkArgsForCall(i int) interface{} {
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	argsForCall := fake.walkArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeStructWalker) WalkReturns(result1 error) {
	fake.walkMutex.Lock()
	defer fake.walkMutex.Unlock()
	fake.WalkStub = nil
	fake.walkReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStructWalker) WalkReturnsOnCall(i int, result1 error) {
	fake.walkMutex.Lock()
	defer fake.walkMutex.Unlock()
	fake.WalkStub = nil
	if fake.walkReturnsOnCall == nil {
		fake.walkReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.walkReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeStructWalker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeStructWalker) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ gostructwalker.StructWalker = new(FakeStructWalker)