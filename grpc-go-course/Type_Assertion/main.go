// package main

// import (
// 	"errors"
// 	"fmt"
// 	"sync"
// )

// type stack struct {
// 	lock    sync.Mutex
// 	element []int
// }

// func (s *stack) push(i int) {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	s.element = append(s.element, i)
// }

// func (s *stack) Pop() (int, error) {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	l := len(s.element)

// 	if l == 0 {
// 		return 0, errors.New("Stack is empty")
// 	}

// 	res := s.element[0]
// 	s.element = s.element[1:l]
// 	return res, nil
// }

// func main() {

// 	var stackobj stack
// 	stackobj.push(2)
// 	stackobj.push(3)
// 	stackobj.push(4)
// 	fmt.Println(stackobj.element)
// 	stackobj.Pop()
// 	fmt.Println(stackobj.element)

//

package main

import (
	"fmt"
)

// type item struct {
// 	value    string
// 	priority int
// 	index    int
// }

// type priorityqueue []*item

// func (pq priorityqueue) Len() int {
// 	return len(pq)
// }

// func (pq priorityqueue) Less(i, j int) bool {
// 	return pq[i].priority > pq[j].priority
// }

// func (pq priorityqueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// 	pq[i].index = i
// 	pq[j].index = j
// }

//var ch = make(chan string) //No wait untill channel is fill

//package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

func main() {
	const cap = 5
	ch := make(chan string, cap)

	go func() {
		for p := range ch {
			fmt.Println("employee : received :", p)
		}
	}()

	const work = 20
	for w := 0; w < work; w++ {
		select {
		case ch <- "paper":
			fmt.Println("manager : send ack")
		default:
			fmt.Println("manager : drop")
		}
	}

	close(ch)
}
