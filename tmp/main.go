package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m sync.Mutex
	var rw sync.RWMutex

	// Using sync.Mutex: One exclusive straw
	fmt.Println("Using sync.Mutex: One exclusive straw")
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		fmt.Println("Friend 1 wants to drink...")
		m.Lock()
		fmt.Println("Friend 1 is drinking")
		time.Sleep(2 * time.Second)
		fmt.Println("Friend 1 is done drinking")
		m.Unlock()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second) // Ensure Friend 1 starts first
		fmt.Println("Friend 2 wants to drink...")
		m.Lock()
		fmt.Println("Friend 2 is drinking")
		time.Sleep(2 * time.Second)
		fmt.Println("Friend 2 is done drinking")
		m.Unlock()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second) // Ensure Friend 2 starts second
		fmt.Println("Friend 3 wants to drink...")
		m.Lock()
		fmt.Println("Friend 3 is drinking")
		time.Sleep(2 * time.Second)
		fmt.Println("Friend 3 is done drinking")
		m.Unlock()
	}()

	wg.Wait()
	fmt.Println()

	// Using sync.RWMutex: Many reader straws, one exclusive writer straw
	fmt.Println("Using sync.RWMutex: Many reader straws, one exclusive writer straw")
	wg.Add(3)

	go func() {
		defer wg.Done()
		fmt.Println("Reader 1 wants to drink...")
		rw.RLock()
		fmt.Println("Reader 1 is drinking")
		time.Sleep(2 * time.Second)
		fmt.Println("Reader 1 is done drinking")
		rw.RUnlock()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second) // Ensure Reader 1 starts first
		fmt.Println("Writer wants to drink alone...")
		rw.Lock()
		fmt.Println("Writer is drinking alone")
		time.Sleep(2 * time.Second)
		fmt.Println("Writer is done drinking")
		rw.Unlock()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second) // Ensure Writer starts second
		fmt.Println("Reader 2 wants to drink...")
		rw.RLock()
		fmt.Println("Reader 2 is drinking")
		time.Sleep(2 * time.Second)
		fmt.Println("Reader 2 is done drinking")
		rw.RUnlock()
	}()

	wg.Wait()
	fmt.Println("All done!")
}
