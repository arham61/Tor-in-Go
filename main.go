package main

import (
	"fmt"
	"math/rand"
	"sync"
)


type Message struct {
	SenderID int
	Content  string
	Path     []int 
}

type Node struct {
	ID        int
	Total     int
	Channel   chan Message
	WaitGroup *sync.WaitGroup
	Terminate chan struct{} 
	Processed bool      
}

func main() {

	numNodes := 5

	messageChannel := make(chan Message, numNodes) 

	var wg sync.WaitGroup
	terminateCh := make(chan struct{}) 
	for i := 1; i <= numNodes; i++ {
		wg.Add(1)
		node := Node{
			ID:        i,
			Total:     numNodes,
			Channel:   messageChannel,
			WaitGroup: &wg,
			Terminate: terminateCh,
		}
		go node.start()
	}


	startingNode := rand.Intn(numNodes) + 1
	fmt.Printf("Starting Node: %d\n", startingNode)

	initialMessage := Message{
		SenderID: startingNode,
		Content:  fmt.Sprintf("Message from Node %d", startingNode),
		Path:     []int{startingNode},
	}


	messageChannel <- initialMessage


	close(terminateCh)

	wg.Wait()

	close(messageChannel)
}

func (n *Node) start() {
	defer n.WaitGroup.Done()

	for {
		select {
		case message, ok := <-n.Channel:
			if !ok {
				return
			}

			fmt.Printf("Node %d: %s\n", n.ID, message.Content)

			if n.ID != message.SenderID && !n.Processed {
				forwardMessage(n, message)
			}

			n.Processed = true

		case <-n.Terminate:
			return
		}
	}
}

func forwardMessage(currentNode *Node, message Message) {
	nextNodeID := rand.Intn(currentNode.Total) + 1

	message.Path = append(message.Path, currentNode.ID)
	message.SenderID = currentNode.ID

	select {
	case currentNode.Channel <- message:
		fmt.Printf("Node %d forwarded the message to Node %d\n", currentNode.ID, nextNodeID)
	default:
	
	}
}