package structures

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	SenderID int
	Content  string
}

type Node struct {
	ID      int
	Total   int
	Channel chan Message
}

func (n *Node) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	message := <-n.Channel

	fmt.Printf("Node %d: %s\n", n.ID, message.Content)

	if n.ID < n.Total {
		n.Channel <- Message{SenderID: n.ID, Content: fmt.Sprintf("Hello from Node %d", n.ID)}
		time.Sleep(time.Millisecond * 100)
	} else {
		n.Channel <- Message{SenderID: n.ID, Content: fmt.Sprintf("Goodbye from Node %d", n.ID)}
	}
}