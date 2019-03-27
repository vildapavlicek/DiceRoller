package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
)

// returns maximum number found in provided slice
func maxInSlice(numbers []int) (int, error) {
	if numbers == nil {
		return 0, fmt.Errorf("nil numbers")
	}

	var max int = 0
	fmt.Printf("ranging values: %v\n", numbers)
	sort.Ints(numbers)
	length := len(numbers) - 1
	max = numbers[length]
	return max, nil
}

// finds maximum number in numbers array by breaking work into N pieces
// (where N is provided as a command line argument) and processing the
// pieces in parallel goroutines
func main() {
	parallelism, _ := strconv.Atoi(os.Args[1])
	fmt.Printf("ok, i'll use %d goroutines\n", parallelism)

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	var allMax int = 0
	perGoRoutine := len(numbers) / parallelism
	fmt.Printf("perGoRoutine: %v\n", perGoRoutine)
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			startIndex := perGoRoutine * i
			endIndex := perGoRoutine*i + perGoRoutine
			if i+1 == parallelism {
				endIndex = len(numbers)
			}
			fmt.Printf("ranging values %v in thread-%v\n", numbers[startIndex:endIndex], i)
			sliceMax, err := maxInSlice(numbers[startIndex:endIndex])
			mutex.Lock()
			if err != nil {
				fmt.Printf("error finding max for slice %d to %d, skipping this slice: %s\n", err)
				return
			}
			fmt.Printf("goroutine %d (slice %d to %d) found max %d\n", i, startIndex, endIndex, sliceMax)
			if sliceMax > allMax {
				allMax = sliceMax
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("maximum: %d\n", allMax)
}
