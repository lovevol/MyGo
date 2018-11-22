package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
)

func Source(num ...int) <-chan int {
	numChan := make(chan int)
	go func() {
		for _, n := range num {
			numChan <- n
		}
		close(numChan)
	}()
	return numChan
}

func ReaderSource(reader io.Reader) <- chan int{
	outChan := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		for  {
			n,err := reader.Read(buffer)
			if n > 0 {
				num := int(binary.BigEndian.Uint64(buffer))
				outChan <- num
			}
			if err != nil{
				break
			}
		}
		close(outChan)
	}()
	return outChan
}

func RandomSource(size int) <- chan int {
	outChan := make(chan int)
	go func() {
		for i := 0; i < size; i ++ {
			outChan <- rand.Int()
		}
		close(outChan)
	}()
	return outChan
}

func MemorySort(numChan <-chan int) chan int {
	var nums []int
	var num int
	var ok bool
	for {
		num, ok = <-numChan
		if !ok {
			break
		}
		nums = append(nums, num)
	}
	sort.Ints(nums)
	outChan := make(chan int)
	go func() {
		for _, n := range nums {
			outChan <- n
		}
		close(outChan)
	}()
	return outChan
}

func Merge(c1, c2 <-chan int) <-chan int {
	resultChan := make(chan int)
	go func() {
		num1, ok1 := <-c1
		num2, ok2 := <-c2
		for {
			if ok1 && ok2 {
				if num1 < num2 {
					resultChan <- num1
					num1, ok1 = <-c1
				} else {
					resultChan <- num2
					num2, ok2 = <-c2
				}
			} else if ok1 {
				resultChan <- num1
				num1, ok1 = <-c1
			} else if ok2 {
				resultChan <- num2
				num2, ok2 = <-c2
			} else {
				break
			}
		}
		close(resultChan)
	}()
	return resultChan
}

func MergeManay(chs ...<-chan int) <-chan int {
	if len(chs) == 1 {
		return chs[0]
	}else if len(chs) == 2 {
		return Merge(chs[0], chs[1])
	}else if len(chs) > 2{
		middle := len(chs) / 2
		return Merge(MergeManay(chs[:middle]...), MergeManay(chs[middle:]...))
	}else{
		return make(chan int)
	}
}

func WriteSink(writer io.Writer, in <-chan int)  {
	for v := range in{
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}
