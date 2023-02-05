package main

import (
	"math"
	"sync"
	"testing"
)

type Shape interface {
	Area() float64
}

type Rectangle struct {
	W, H float64
}

func (r Rectangle) Area() float64 {
	return r.W * r.H
}

type Circle struct {
	R float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.R, 2)
}

// BenchmarkInterface-8   	742840539	         1.493 ns/op
func BenchmarkInterface(b *testing.B) {
	var shape Shape = Rectangle{4, 2}
	for i := 0; i < b.N; i++ {
		shape.Area()
	}
}

// BenchmarkCast-8   	1000000000	         0.2912 ns/op
func BenchmarkCast(b *testing.B) {
	var shape Shape = Rectangle{4, 2}
	for i := 0; i < b.N; i++ {
		if rec, ok := shape.(Rectangle); ok {
			rec.Area()
		}
	}
}

type mySafeCounter struct {
	C  int
	mu sync.Mutex
	ch chan int
}

func (c *mySafeCounter) safeIncMutex() func(int) {
	var mu sync.Mutex
	return func(n int) {
		mu.Lock()
		c.C += n
		mu.Unlock()
	}
}

func (c *mySafeCounter) safeIncChan() func(int) {
	ch := make(chan int)

	go func() {
		for {
			c.C += <-ch
		}
	}()

	return func(n int) { ch <- n }
}

func runConcurrently(concurrency int, fn func()) {
	wg := sync.WaitGroup{}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			fn()
			wg.Done()
		}()
	}
	wg.Wait()
}

// Benchmark_Mutex-8   	     396	   2956641 ns/op
func Benchmark_Mutex(b *testing.B) {
	var counter int
	var mu sync.Mutex
	for i := 0; i < b.N; i++ {
		runConcurrently(10_000, func() {
			mu.Lock()
			counter += 1
			mu.Unlock()
		})
	}
}

// Benchmark_chan-8   	     181	   6066474 ns/op
func Benchmark_chan(b *testing.B) {
	var counter int
	ch := make(chan int)
	go func() {
		for {
			counter += <-ch
		}
	}()

	for i := 0; i < b.N; i++ {
		runConcurrently(10_000, func() {
			ch <- 1
		})
	}
}

// BenchmarkChan-8   	     325	   3968805 ns/op
