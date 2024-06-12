package main

import (
	"fmt"
	"sync"
)

// Unbuffered Channels -> There is no reserved space to allocate data, which means that a ready consumer is needed.
// In other words, the channel is always full and the consumer must be ready before the data is pushed into the channel.

// Wrong Aproach Example:
// The message will never be consumed.
// As the program gets stuck on c <- 10 statement the consumer never gets ready to consume it.
func unbuffered() {
	c := make(chan int)

	c <- 10 // blocking here

	msg := <-c
	fmt.Println("Message from channel: ", msg)
}

// Correct Aproach Example:
// Now the consumer will be ready before the channel receives a message, scheduling a go routine.
// After the message be consumed the program will not be blocked anymore.
func unbufferedCorrect() {
	c := make(chan int)

	// Preparing go routine that will schedule a function responsible for consuming the message.
	go func() {
		msg := <-c
		fmt.Println("Message from channel: ", msg)
	}()

	c <- 10 // blocking here
	// It still blocks here until someone reads the message but it will no longer cause a deadlock error
}

// Buffered Channels -> It is possible to define the length of the channel so it can receive more than one message.
// The channel does not need a ready consumer before the channel gets populated, to consume the message.
// In other words buffered channels works like a queue of messages with a defined length.
// If the buffered channel gets full it is going to cause a deadlock error.

// Wrong Aproach Example:
// In this case the message will be printed out, but it still causes a deadlock error because all goroutines are asleep after printing the messages.
// It happens because the program ends but the chan is still open.
func buffered() {
	c := make(chan int, 10)
	c <- 10
	c <- 20
	for msg := range c {
		fmt.Println("Message from channel: ", msg)
	}

}

// Correct Aproach Example:

// Adding wait groups so the result can be printed out correctly.
// In this case wg is needed because the code is only sending two messages to the channel and it`s size is 10.
// So it will not block until the code sends 10th message to it.
// Also, as it is not blocking, is not possible to guarantee that the consumer will read the message before the program is over, since the consumer func is scheduled.

func bufferedCorrect() {
	c := make(chan int, 10)

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		for msg := range c {
			fmt.Println("Message from channel: ", msg)
			wg.Done()
		}
	}(&wg)

	c <- 10
	c <- 20
	wg.Wait()
	fmt.Println("Exiting Program: Work is done!")
}

// Another way to do it without go routines is by simply closing the channel before ranging over it.
func bufferedClosingChannel() {
	c := make(chan int, 10)
	c <- 10
	c <- 20
	for msg := range c {
		fmt.Println("Message from channel: ", msg)
	}

	close(c)

}

func main() {
	// unbuffered()
	// unbufferedCorrect()
	buffered()
	// bufferedCorrect()
	// bufferedClosingChannel()
}
