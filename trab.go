package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func isPrime(p int) bool {
	if p%2 == 0 {
		return false
	}
	for i := 3; i*i <= p; i += 2 {
		if p%i == 0 {
			return false
		}
	}
	return true
}

func contaPrimosSeq(slice []int) int {
	count := 0
	for _, value := range slice {
		if isPrime(value) {
			count++
		}
	}
	return count
}

func contaPrimosConc(slice []int, numProcs int) int {
	count := 0
	var wg sync.WaitGroup
	ch := make(chan int, numProcs)

	for i := 0; i < numProcs; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localCount := 0
			for j := workerID; j < len(slice); j += numProcs {
				if isPrime(slice[j]) {
					localCount++
				}
			}
			ch <- localCount
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for localCount := range ch {
		count += localCount
	}

	return count
}

func main() {
	numProcs := 8
	runtime.GOMAXPROCS(numProcs)

	slices := [][]int{
		{101, 883, 359, 941, 983, 859, 523, 631, 181, 233},
		{547369, 669437, 683251, 610279, 851117, 655439, 937351, 419443, 128467, 316879},
		{550032733, 429415309, 109543211, 882936113, 546857209, 756170741, 699422809, 469062577, 117355333, 617320027},
		{7069402558433, 960246047869, 5738081989711, 5358141480883, 2569391599009, 4135462531597, 7807787948171, 130788041233,
			2708131414819, 1571981553097},
		{383376390724197361, 882611655919772761, 533290385325847007, 17969611178168479, 903013501582628521, 541906710014517121,
			281512690206248899, 403936627075987639, 775148726422474717, 942319117335957539},
	}

	for _, slice := range slices {
		firstNumber := slice[0]
		fmt.Println("------ Conta primos de tamanho", len(fmt.Sprint(firstNumber)), "-------")
		start := time.Now()
		p := contaPrimosSeq(slice)
		fmt.Printf("-> Sequencial ------ Segundos:%f\n", time.Since(start).Seconds())
		fmt.Println("------ N primos:", p)

		//for numWorkers := 2; numWorkers <= 16; numWorkers += 1 {
		numWorkers := numProcs
		start2 := time.Now()
		p = contaPrimosConc(slice, numWorkers)
		fmt.Printf("-> Concorrente (%d workers) ------ Segundos: %f\n", numWorkers, time.Since(start2).Seconds())
		fmt.Println("------ N primos:", p)
		//}
	}
}
