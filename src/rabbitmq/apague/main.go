package main

import (
	"fmt"
	"sync"
	"time"
)

func SomaNSemThread(n int) int {
	s := 0

	for i := 0; i < n; i++ {
		s++
	}
	return s
}

func SomaNComThread(n int, t int) int {
	wg := sync.WaitGroup{}
	s := make([]int, t)

	cSize := n / t

	for idx := 0; idx < t; idx++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			b := cSize * j
			e := b + cSize
			fmt.Println(b, e)
			for i := b; i < e; i++ {
				s[j] = s[j] + 1
			}
		}(idx)
	}
	wg.Wait()

	sTotal := 0
	for idx := 0; idx < t; idx++ {
		sTotal = sTotal + s[idx]
	}

	// adjust for some cases, e.g., n = 1000, t = 3
	sTotal = sTotal + n - cSize*t

	return sTotal
}

func main() {
	t1 := time.Now()
	s1 := SomaNSemThread(1000)
	t2 := time.Now()

	t3 := time.Now()
	s2 := SomaNComThread(1000, 2)
	t4 := time.Now()

	fmt.Println(s1, s2, t2.Sub(t1), t4.Sub(t3))
}
