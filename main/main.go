package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

var ip_list []string

func WritLog(log_file string, ip string) {
	open, err := os.OpenFile(log_file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("创建或写入文件失败，报错为：", err)
		return
	}
	writer := bufio.NewWriter(open)

	writer.Write([]byte(ip))
	writer.WriteByte('\n')

	writer.Flush()

	defer open.Close()
}

func ReadVm(file string) {
	open, err := os.Open(file)
	if err != nil {
		fmt.Println("打开文件失败：", err)
		return
	}
	defer open.Close()

	reader := bufio.NewReader(open)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		ip_list = append(ip_list, string(line))
	}
}

func Network() {
	//从命令行获取检测的端口
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	wg.Add(len(ip_list))
	port := os.Args[1]

	//遍历ip
	for i := 0; i < len(ip_list); i++ {

		address := ip_list[i] + ":" + port

		//开启协程，测试端口是否通
		go func() {
			_, err := net.DialTimeout("tcp", address, 3*time.Second)

			if err != nil {
				fmt.Printf("%v端口不通,原因为：%v\n", address, err)
				WritLog("./端口未开放的ip.txt", address)
			} else {
				fmt.Printf("%v端口通\n", address)
				WritLog("./端口开放的ip.txt", address)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func main() {
	start := time.Now()

	ReadVm("./vm_list.txt")
	Network()

	fmt.Println("检测结束，结果请查看同级目录输出文件")

	fmt.Println("程序耗时：", time.Since(start))
}
