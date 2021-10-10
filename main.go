package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//Node Each has a name, message, and a corresponding channel map and name map\\
//Node Each node has its own channel map to ensure that no message is lost during Go Routines
type Node struct {
	name, message   int
	contactChannels map[int]chan int //for sending messages
	nameChannels    map[int]chan int //for sending name as a request
}

//NodeMap Structure for a map of nodes
type NodeMap struct {
	nodes map[int]Node
	mu    sync.Mutex
}

var nm NodeMap
var nodes map[int]Node
var numInfected int64
var totalRuns int64

//pullSendRequest
func (nm *NodeMap) pullSendRequest(i int) {
	if nm.nodes[i].message == 0 {
		randNode := rand.Intn(3)
		fmt.Println(nm.nodes[i].name, "Pulling from ", randNode)
		nm.nodes[i].nameChannels[randNode] <- nm.nodes[i].name
	}
}

func (nm *NodeMap) pullSendMessage(i int) {
	for j := 0; j < len(nm.nodes); j++ {
		select {
		case contactor, ok := <-nm.nodes[i].nameChannels[i]:
			if ok {
				fmt.Println(i, "sending message to", contactor)
				nm.nodes[i].contactChannels[contactor] <- nm.nodes[i].message
			}
		default:
			break //handles case where channel is empty
		}
	}
}

func (nm *NodeMap) pullReceiveMessage(i int) {
	nm.mu.Lock()
	if nm.nodes[i].message == 0 {
		for j := 0; j < len(nm.nodes); j++ {
			select {
			case newMessage, ok := <-nm.nodes[i].contactChannels[i]:
				if ok {
					if newMessage == 1 {
						tmpNode := nm.nodes[i]
						tmpNode.message = newMessage
						nm.nodes[i] = tmpNode
						atomic.AddInt64(&numInfected, 1)
					}
				}
			default:
				break //handles case where channel is empty
			}
		}
	}
	nm.mu.Unlock()
}

//Implements pull protocol
func pullProtocol(wg *sync.WaitGroup) {
	for numInfected < 3 {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				n := nm.nodes[i]
				nm.pullSendRequest(n.name)
				time.Sleep(1 * time.Second)
				nm.pullSendMessage(n.name)
				time.Sleep(1 * time.Second)
				nm.pullReceiveMessage(n.name)
			}(i)
		}
		wg.Wait()
		atomic.AddInt64(&totalRuns, 1)
		time.Sleep(1 * time.Second)
	}
}

func (nm *NodeMap) pushSendMessage(i int) {
	randNode := rand.Intn(3)
	fmt.Println(nm.nodes[i].name, "Push: Sending to ", randNode)
	nm.nodes[i].contactChannels[randNode] <- nm.nodes[i].message
}

func (nm *NodeMap) pushReceiveMessage(i int) {
	nm.mu.Lock()
	select {
	case newMessage, ok := <-nm.nodes[i].contactChannels[i]:
		if ok {
			if nm.nodes[i].message == 0 {
				tmpNode := nm.nodes[i]
				tmpNode.message = newMessage
				nm.nodes[i] = tmpNode
				atomic.AddInt64(&numInfected, 1)
				fmt.Println(i, "received", newMessage)
			}
		}
	default:
		break //handles case where channel is empty
	}
	nm.mu.Unlock()
}

//Implements push protocol
func pushProtocol(wg *sync.WaitGroup) {
	for numInfected < 3 {
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
		atomic.AddInt64(&totalRuns, 1)
		time.Sleep(1 * time.Second)
	}
}

//Runs Push or Pull protocol based on the ratio of numInfected and total nodes.
func pushPullProtocol(wg *sync.WaitGroup) {
	for numInfected < 3 {
		print(len(nm.nodes)/2)
		if numInfected < int64(len(nm.nodes)/2) {
			pushProtocol(wg)
		}else {
			pullProtocol(wg)
		}
		atomic.AddInt64(&totalRuns, 1)
		time.Sleep(1 * time.Second)
	}
}

//resets totalRuns and numInfected to starting value,
func resetVariables(totalNodes int) {
	totalRuns = 0
	numInfected = 1
	for i := 0; i < len(nm.nodes); i++ { //since we will set all nodes' channels to 0's, only need to clear 0's channels
		select {
		case _, ok := <-nm.nodes[0].contactChannels[i]:
			if ok {
				fmt.Println("Clearing Channel")
			}
		default:
			break //handles case where channel is empty
		}
	}
	//Resets nodes to original values
	for i := 1; i < totalNodes; i++ { //Node 0 is always infected, so it does not need to be reset
		tmpNode := Node{i, 0, nm.nodes[0].contactChannels, nm.nodes[0].nameChannels}
		nm.nodes[i] = tmpNode
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) //Ensure we get a different seed for random each time
	var wg sync.WaitGroup

	contacts := make(map[int]chan int)
	names := make(map[int]chan int)
	for i := 0; i < 3; i++ {
		contacts[i] = make(chan int, 3)
		names[i] = make(chan int, 3)
	}
	n0 := Node{0, 1, contacts, names}

	nodes = make(map[int]Node)
	nodes[0] = n0
	nm = NodeMap{nodes: nodes}

	resetVariables(3)
	//Calls push protocol function
	pushProtocol(&wg)
	fmt.Println("Push TR:", totalRuns)

	//reset variables for pull run
	resetVariables(3)

	//Calls pull protocol function
	pullProtocol(&wg)
	fmt.Println("Pull TR:", totalRuns)

	//reset variables for push-pull-run
	resetVariables(3)

	//Calls push-pull protocol function
	pushPullProtocol(&wg)
	fmt.Println("Push-Pull TR:", totalRuns)
}
