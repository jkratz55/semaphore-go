# Semaphore GO

A simple semaphore implementation for Golang.

A semaphore is a variable or abstract data type used to control access to a common resource by multiple threads and avoid critical section problems in a concurrent system such as a multitasking operating system. Semaphores are a type of synchronization primitive.

## Usage

The Semaphore type has a very simple API consisting of three methods: Acquire, TryAcquire, and Release. Acquire and TryAcquire are used to acquire the semaphore. Acquire is blocking but accepts a context so the operation can be cancelled or timeout, and TryAcquire is non-blocking. In contrast Release is used to release the Semaphore after it's been acquired.

Calls to release should always be paired with successful acquisitions of the Semaphore to release the Semaphore. Calling Release when there wasn't a successful acquire can/will lead to a deadlock. 

### Example

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jkratz55/semaphore-go"
)

func main()  {
	sem := semaphore.New(5) // create a semaphore with max concurrency of 5

	for i:=0; i < 20; i++ {
		go func(id int) {
			_ = sem.Acquire(context.Background())
			defer sem.Release()

			fmt.Printf("ID %d\n", id)
			time.Sleep(time.Second * 5)
		}(i)
	}

	if sem.TryAcquire() {
		fmt.Printf("I tried and I made it!\n")
		sem.Release() // Release should only be called if acquiring the Semaphore succeeded
	} else {
		fmt.Printf("Ahhhh snap! I didn't make it\n")
	}

	time.Sleep(time.Second * 30) // don't do this at home, this is just lazy way to ensure the program finishes
}
```