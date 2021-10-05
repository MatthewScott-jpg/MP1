package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	name    int
	message int
}

func sendMessage(name int, out chan int) {
	out <- name
}

func receiveMessage(name int, out chan int) {
	fmt.Println(name, " received from", <-out)
}

func main() {
	n1 := Node{0, 0}
	n2 := Node{1, 0}
	n3 := Node{2, 0}

	nodes := []Node{n1, n2, n3}

	out := make(chan int, 3)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			n := nodes[i]
			sendMessage(n.name, out)
			time.Sleep(2 * time.Second)
			receiveMessage(n.name, out)
		}(i)
	}
	wg.Wait()
	fmt.Println("All done")
}
