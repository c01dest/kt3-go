package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{2, 4, 6, 8, 10}
	results := make(chan int, len(numbers))
	var wg sync.WaitGroup

	for _, n := range numbers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			results <- num * num
		}(n)
	}

	// Ждем окончания вычислений в горутинах и закрываем канал
	wg.Wait()
	close(results)

	var totalSum int
	for square := range results {
		totalSum += square
	}

	fmt.Printf("Сумма квадратов: %d\n", totalSum)
}