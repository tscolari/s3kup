// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/tscolari/s3kup/fetch"
	"github.com/tscolari/s3kup/s3"
)

type FakeS3Client struct {
	ListStub        func(path string) (versions s3.Versions, err error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		path string
	}
	listReturns struct {
		result1 s3.Versions
		result2 error
	}
	GetStub        func(path string) ([]byte, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		path string
	}
	getReturns struct {
		result1 []byte
		result2 error
	}
}

func (fake *FakeS3Client) List(path string) (versions s3.Versions, err error) {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		path string
	}{path})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(path)
	} else {
		return fake.listReturns.result1, fake.listReturns.result2
	}
}

func (fake *FakeS3Client) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeS3Client) ListArgsForCall(i int) string {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return fake.listArgsForCall[i].path
}

func (fake *FakeS3Client) ListReturns(result1 s3.Versions, result2 error) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 s3.Versions
		result2 error
	}{result1, result2}
}

func (fake *FakeS3Client) Get(path string) ([]byte, error) {
	fake.getMutex.Lock()
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		path string
	}{path})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(path)
	} else {
		return fake.getReturns.result1, fake.getReturns.result2
	}
}

func (fake *FakeS3Client) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeS3Client) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].path
}

func (fake *FakeS3Client) GetReturns(result1 []byte, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

var _ fetch.S3Client = new(FakeS3Client)