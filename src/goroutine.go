package main

import (
	"time"
	"fmt"
)

func main()  {
	var countNum int64
	countChan := make(chan int64,1000000000)
	go func() {
		for{
			go func() {
				countChan <- 1
			}()
		}
	}()
	go func() {
		for i := range countChan{
			countNum+=i
		}
	}()
	time.Sleep(1*time.Second)
	fmt.Println(countNum)
}
