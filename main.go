package main

import (
	"sync"
	"go-tor/structures"
)


func main() {
	numNodes := 5

	messageChannel := make(chan structures.Message)

	var wg sync.WaitGroup
	for i := 2; i <= numNodes; i++ {
		wg.Add(1)
		node := structures.Node{
			ID:      i,
			Total:   numNodes,
			Channel: messageChannel,
		}
		go node.Start(&wg)
	}

	messageChannel <- structures.Message{SenderID: 1, Content: "Hello from Node 1"}

	go func() {
		wg.Wait()
		close(messageChannel) 
	}()

	// Block the main goroutine until all nodes are done
	for range messageChannel {
		// Keep the main goroutine alive until the messageChannel is closed
	}
}

