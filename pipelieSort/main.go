package main

import (
	"os"
	"MyGo/pipelieSort/pipeline"
	"bufio"
	"strconv"
	"fmt"
)

const (
	resultFileName    = "result"
	size              = 100000
)

func main() {
	demo(4)
	//print()
}

func demo(fileCount int) {
	var resultFileNames []string
	//创建n个包含原始数据的文件 并且分别排序后写入对应的结果文件
	for i := 0; i < fileCount; i ++ {
		//使用闭包 在for 中使用 defer
		func() {
			resultFile, err := os.Create(resultFileName + strconv.Itoa(i) + ".result")
			if err != nil {
				panic(err)
			}
			defer resultFile.Close()

			writer := bufio.NewWriter(resultFile)
			pipeline.WriteSink(writer, pipeline.MemorySort(pipeline.RandomSource(size)))
			writer.Flush()
			resultFileNames = append(resultFileNames, resultFile.Name())
		}()
	}

	var resultChans []<-chan int
	var files []*os.File
	//获取到n个结果文件的输出
	for _, name := range resultFileNames {
		//使用闭包 在for 中使用 defer
		func() {
			resultFile, err := os.Open(name)
			if err != nil {
				panic(err)
			}
			files = append(files, resultFile)
			resultChans = append(resultChans, pipeline.ReaderSource(bufio.NewReader(resultFile)))
		}()
	}

	//归并后写入最终结果文件
	resultFile, err := os.Create(resultFileName)
	if err != nil {
		panic(err)
	}
	defer resultFile.Close()
	writer := bufio.NewWriter(resultFile)
	pipeline.WriteSink(writer, pipeline.MergeManay(resultChans...))
	writer.Flush()

	for _,file := range files{
		func(){
			defer file.Close()
		}()
	}
	print()
}

//输出一部分数据查看
func print() {
	resultFile, err := os.Open(resultFileName)
	if err != nil {
		panic(err)
	}
	defer resultFile.Close()
	count := 0
	for v := range pipeline.ReaderSource(bufio.NewReader(resultFile)) {
		fmt.Println(v)
		count ++
		if count > 100 {
			break
		}
	}
}
