package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	name            int
	message         int
	contactChannels map[int]chan int
}

func sendMessage(name int, out map[int]chan int) {
	randnode := rand.Intn(3)
	out[randnode] <- name
}

func receiveMessage(name int, out map[int]chan int) {
	for i := 0; i < 3; i++ {
		select {
		case x, ok := <-out[name]:
			if ok {
				fmt.Println(name, " received from", x)
			} else {
				break
			}
		default: //handles blocking when channels are emptied
		}
	}
}

func main() {
	contacts := make(map[int]chan int)
	for i := 0; i < 3; i++ {
		contacts[i] = make(chan int, 3)
	}
	n1 := Node{0, 0, contacts}
	n2 := Node{1, 0, contacts}
	n3 := Node{2, 0, contacts}

	nodes := []Node{n1, n2, n3}

	//out := make(chan int, 3)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			n := nodes[i]
			sendMessage(n.name, n.contactChannels)
			time.Sleep(1 * time.Second)
			receiveMessage(n.name, n.contactChannels)
		}(i)
	}
	wg.Wait()
	fmt.Println("All done")
}
