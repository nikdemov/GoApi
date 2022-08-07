//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"time"
)

func contextDemo(name string, ctx context.Context) {
	for {
		deadline, ok := ctx.Deadline()
		if ok {
			fmt.Println(name, "will expire at:", deadline)
		} else {
			fmt.Println(name, "has no deadline")
		}
		time.Sleep(time.Second)
	}
}

func main() {
	timeout := 3 * time.Second
	deadline := time.Now().Add(4 * time.Hour)
	timeOutContext, _ := context.WithTimeout(
		context.Background(), timeout)
	cancelContext, cancelFunc := context.WithCancel(
		context.Background())
	deadlineContext, _ := context.WithDeadline(
		cancelContext, deadline)

	go contextDemo("[timeoutContext]", timeOutContext)
	go contextDemo("[cancelContext]", cancelContext)
	go contextDemo("[deadlineContext]", deadlineContext)

	// Wait for the timeout to expire
	<-timeOutContext.Done()

	// This will cancel the deadline context as well as its
	// child - the cancelContext

	go func() {
		time.Sleep(time.Second * 7)
		fmt.Println("Cancelling the cancel context...")
		cancelFunc()
	}()

	<-cancelContext.Done()
	fmt.Println("The cancel context has been cancelled...")

	// Wait for both contexts to be cancelled
	<-deadlineContext.Done()
	fmt.Println("The deadline context has been cancelled...")
}
