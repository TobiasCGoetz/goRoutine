package main

import (
	"fmt"
	"sync"
)

type ConcurrentQueue[T comparable] struct {
	items []T
	lock  sync.Mutex
	cond  *sync.Cond
}

func NewQueue[T comparable]() *ConcurrentQueue[T] {
	q := &ConcurrentQueue[T]{}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *ConcurrentQueue[T]) isEmpty() bool {
	return len(q.items) < 1
}

func (q *ConcurrentQueue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *ConcurrentQueue[T]) Dequeue() T {
	q.lock.Lock()
	defer q.lock.Unlock()
	for q.isEmpty() {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func worker(values []int, q *ConcurrentQueue[int], wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for _, val := range values {
		sum += val
	}
	fmt.Println("Worker started with ", values, " returned ", sum)
	q.Enqueue(sum)
}

func main() {
	fmt.Println("This is a test project that adds up arrays of integers in a go coroutine")
	fmt.Println("Init queue...")
	cq := NewQueue[int]()
	fmt.Println("Init wait group...")
	wg := sync.WaitGroup{}
	fmt.Println("Dispatch workers...")
	wg.Add(2)
	go worker([]int{1, 2, 3, 4, 5}, cq, &wg)
	go worker([]int{4, 5, 6, 7, 8}, cq, &wg)
	wg.Wait()
	fmt.Println("Grabbing results...")
	for !cq.isEmpty() {
		fmt.Println(cq.Dequeue())
	}

}
