package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	generateFiles(200)
}

func calculate(target int, sum int, start int, result *[]int, delta *int) {
	if target == sum {
		*delta = *delta + 1
		log.Printf("%d @ %d = %v", target, *delta, *result)
	}

	for i := start; i < target; i++ {
		tmpSum := sum + i
		if tmpSum <= target {
			*result = append(*result, i)
			calculate(target, tmpSum, i+1, result, delta)
			*result = (*result)[:len(*result)-1]
		} else {
			break
		}
	}
}

func calculateToFile(file *os.File, target int, sum int, start int, result *[]int, delta *int) {
	if target == sum {
		*delta = *delta + 1
		log.Printf("%d @ %d = %v", target, *delta, *result)
		buf, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		buf = append(buf, '\n')
		if _, err := file.Write(buf); err != nil {
			log.Println(err)
		}
	}

	for i := start; i < target; i++ {
		tmpSum := sum + i
		if tmpSum <= target {
			*result = append(*result, i)
			calculateToFile(file, target, tmpSum, i+1, result, delta)
			*result = (*result)[:len(*result)-1]
		} else {
			break
		}
	}
}

func generateFiles(n int) {
	var wg sync.WaitGroup
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(target int) {
			buffer := make([]int, 0)
			totalCombinations := 0
			generateFile(target, &buffer, &totalCombinations)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func generateFile(target int, buffer *[]int, totalCombinations *int) {
	f, err := os.OpenFile(fmt.Sprintf("%d.txt", target), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	calculateToFile(f, target, 0, 1, buffer, totalCombinations)
}
