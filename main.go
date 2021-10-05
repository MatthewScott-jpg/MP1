package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	name, message   int
	infectionState  bool
	contactChannels map[int]chan int
	randval         int64
}

func (n *Node) sendMessage() {
	randnode := rand.Intn(3)
	n.contactChannels[randnode] <- n.message
}

func (n *Node) receiveMessage() {
	atomic.StoreInt64(&n.randval, 100)
	//n.randval = 100
	for i := 0; i < 3; i++ {
		select {
		case x, ok := <-n.contactChannels[n.name]:
			if ok {
				if !n.infectionState {
					n.message = x
					n.infectionState = true
					fmt.Println(n.name, "received", x)
				}
			} else {
				break
			}
		default: //handles blocking when channels are emptied
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	contacts := make(map[int]chan int)
	for i := 0; i < 3; i++ {
		contacts[i] = make(chan int, 3)
	}
	n0 := Node{0, 1, true, contacts, 0}
	n1 := Node{1, 0, false, contacts, 0}
	n2 := Node{2, 0, false, contacts, 0}

	nodes := []Node{n0, n1, n2}

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			n := nodes[i]
			if n.infectionState {
				n.sendMessage()
			}
			time.Sleep(1 * time.Second)
			n.receiveMessage()
		}(i)
	}
	wg.Wait()
	fmt.Println("All done")
	fmt.Println(n0)
	fmt.Println(n1)
	fmt.Println(n2)
}
