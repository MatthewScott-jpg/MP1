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
type NodeMap struct {
	nodes map[int]Node
	mu    sync.Mutex
}

var nm NodeMap
var nodes map[int]Node
var numInfected int64
var totalRuns int64

func (nm *NodeMap) pull(i int) {
	randnode := rand.Intn(3)
	fmt.Println(nm.nodes[i].name, "Pulling from ", randnode)
	nm.nodes[i].contactChannels[i] <- nm.nodes[randnode].message
	select {
	case x, ok := <-nm.nodes[randnode].contactChannels[i]:
		if ok {
			if nm.nodes[randnode].message == 1 {
				val := nm.nodes[randnode]
				val.message = x
				nm.nodes[i] = val
				atomic.AddInt64(&numInfected, 1)
				fmt.Println(i, "received", x)
			}
		}
	default:
		break //handles case where channel is empty
	}
	nm.mu.Unlock()
}

func (nm *NodeMap) pushSendMessage(i int) {
	randnode := rand.Intn(3)
	fmt.Println(nm.nodes[i].name, "Push: Sending to ", randnode)
	nm.nodes[i].contactChannels[randnode] <- nm.nodes[i].message
}

func (nm *NodeMap) pushReceiveMessage(i int) {
	nm.mu.Lock()
	select {
	case x, ok := <-nm.nodes[i].contactChannels[i]:
		if ok {
			if nm.nodes[i].message == 0 {
				val := nm.nodes[i]
				val.message = x
				nm.nodes[i] = val
				atomic.AddInt64(&numInfected, 1)
				fmt.Println(i, "received", x)
			}
		}
	default:
		break //handles case where channel is empty
	}
	nm.mu.Unlock()
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
	nm = NodeMap{nodes: nodes}
	nmlength := len(nm.nodes)

	totalRuns = 0

	//Implements push-pull protocol
	/*if numInfected > int64(nmlength)/2 {

	//} else {

	}*/

	//Implements push protocol
	for numInfected < 3 {
		var wg sync.WaitGroup
		//nm.mu.Lock()
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				n := nm.nodes[i]
				if n.message == 1 {
					nm.pushSendMessage(n.name)
				}
				time.Sleep(1 * time.Second)
				nm.pushReceiveMessage(n.name)
			}(i)
		}
		wg.Wait()
		//nm.mu.Unlock()
		atomic.AddInt64(&totalRuns, 1)
		time.Sleep(1 * time.Second)
	}

	//Implements pull protocol
	for numInfected < 3 {
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				n := nm.nodes[i]
				if n.message == 0 {
					nm.pull(n.name)
				}
			}(i)
		}
		wg.Wait()
		atomic.AddInt64(&totalRuns, 1)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("All done")
	fmt.Println("TR:", totalRuns)
}
