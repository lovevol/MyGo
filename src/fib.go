package main

import (
	"fmt"
	"time"
)

func main()  {

	var a, b float64
	a = 0
	b = 1
	start := time.Now()
	fmt.Println(start.Second())
	for i := 1; i <= 1000; i++{
		a, b = b, a+b
	}
	fmt.Println(b)
	fmt.Print(time.Now().Nanosecond())
}