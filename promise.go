package main

import (
	"fmt"
	"time"
)

//Promise prototype
type promisePrototype interface {
	then(onFulfilled func(int32), onRejected func(int32)) promisePrototype
	catch(onRejected func(int32)) promisePrototype
	finally(onFinally func(int32)) promisePrototype
}

// Promise which can resolve and reject a int32 value
type promise struct {
	executable func(func(int32), func(int32))
	reject     func(int32)
	resolve    func(int32)
	onSuccess  chan int32
	onError    chan int32
}

//then
func (p promise) then(onFulfilled func(int32), onRejected func(int32)) promisePrototype {
	go p.executable(p.resolve, p.reject)
	select {
	case success := <-p.onSuccess:
		if onFulfilled != nil {
			onFulfilled(success)
		}
	case error := <-p.onError:
		if onRejected != nil {
			onRejected(error)
		}
	}

	return p
}

//catch
func (p promise) catch(onRejected func(int32)) promisePrototype {
	go p.executable(p.resolve, p.reject)
	select {
	case <-p.onSuccess:
	case error := <-p.onError:
		if onRejected != nil {
			onRejected(error)
		}

	}
	return p
}

//finally
func (p promise) finally(onFinally func(int32)) promisePrototype {
	go p.executable(p.resolve, p.reject)
	select {
	case success := <-p.onSuccess:
		if onFinally != nil {
			onFinally(success)
		}
	case error := <-p.onError:
		if onFinally != nil {
			onFinally(error)
		}
	}
	return p
}

//create a promise
func newPromise(executable func(func(int32), func(int32))) promisePrototype {
	var p promise
	p.onSuccess = make(chan int32)
	p.onError = make(chan int32)

	p.resolve = func(value int32) {
		p.onSuccess <- value
	}
	p.reject = func(value int32) {
		p.onError <- value
	}
	p.executable = executable
	return p
}

func main() {
	// Promise 1: sleep for 3 second then reject with value 32
	p1 := newPromise(func(resolve func(int32), reject func(int32)) {
		time.Sleep(time.Millisecond * 200)
		resolve(32)
	})

	// Promise 2: sleep for 2 second then resolve with value 31
	p2 := newPromise(func(resolve func(int32), reject func(int32)) {
		time.Sleep(time.Millisecond * 50)
		reject(31)
	})

	go p1.then(func(value int32) { fmt.Println("P1 Then Resolved: ", value) }, nil).catch(func(value int32) { fmt.Println("P1 Catch: ", value) }).finally(func(value int32) { fmt.Println("P1 Finally", value) })
	go p2.then(nil, func(value int32) { fmt.Println("P2 Then Rejected: ", value) }).catch(func(value int32) { fmt.Println("P2 Catch: ", value) }).finally(func(value int32) { fmt.Println("P2 Finally", value) })

	//A promise is a async operation so let's sleep for a few seconds
	time.Sleep(time.Second)
}
