package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	logFile, err := os.OpenFile("./im_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetPrefix("[im_server]")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	pid := os.Getpid()
	log.Printf("imserver start run in pid:%d,log setup ok", pid)

}
func main() {
	pid := os.Getpid()
	log.Printf("-----imserver start run in pid:%d-----", pid)
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
