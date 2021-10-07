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
	contactChannels map[int]chan int
}

var nodes map[int]Node
var numInfected int64

func (n *Node) sendMessage() {
	randnode := rand.Intn(3)
	fmt.Println("Sending to ", randnode)
	n.contactChannels[randnode] <- n.message
}

func (n *Node) receiveMessage() {
	select {
	case x, ok := <-n.contactChannels[n.name]:
		if ok {
			if n.message == 0 {
				n.message = x
				atomic.AddInt64(&numInfected, 1)
				fmt.Println(n.name, "received", x)
			}
		}
	default:
		break //handles case where channel is empty
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numInfected = 1
	contacts := make(map[int]chan int)
	for i := 0; i < 3; i++ {
		contacts[i] = make(chan int, 3)
	}
	n0 := Node{0, 1, contacts}
	n1 := Node{1, 0, contacts}
	n2 := Node{2, 0, contacts}

	nodes = make(map[int]Node)
	nodes[0] = n0
	nodes[1] = n1
	nodes[2] = n2

	totalRuns := 0

	var wg sync.WaitGroup
	for numInfected < 3 {
		totalRuns += 1
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				n := nodes[i]
				if n.message == 1 {
					n.sendMessage()
				}
				time.Sleep(1 * time.Second)
				n.receiveMessage()
			}(i)
		}
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	fmt.Println("All done")
	fmt.Println("TR:", totalRuns)
}
