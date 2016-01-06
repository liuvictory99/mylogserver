package main

import(
	"fmt"
	"os"
	"os/signal"
	"net"
	"sync"
	"log/logutil"
)
const (
  maxBufPoolSize = 100
)
var service = ":1213"
var bufs [][]byte 
var lock sync.Mutex

func main (){

    bufs = make([][]byte,0,maxBufPoolSize)

	udpAddr,err := net.ResolveUDPAddr("udp4",service)
	checkError(err)
	conn, err := net.ListenUDP("udp",udpAddr)
	checkError(err)

	for i:=0; i<100;i++ {
	   go func(){
			for {
				handleClient(conn)
			}	
		}()
	}
       
	var ch chan os.Signal= make(chan os.Signal)
	var chh os.Signal
	go func(){
	    signal.Notify(ch,os.Interrupt,os.Kill)
	}()
	
	select {
		case chh = <-ch:
		     switch chh {
				case os.Interrupt:
				     fmt.Println("Interrupt")
			    case os.Kill:
				     fmt.Println("kill")
				default:
				     fmt.Println("Unknow");
			}
	}
}


//从缓冲区抛出一块内存
func popBuf() []byte {
     lock.Lock()
     var buf []byte
     if len(bufs) < 2 {
        buf = make([]byte,4096)
	 	bufs = append(bufs, buf)
     } else {
        buf = bufs[len(bufs)-1]
	    bufs = bufs[0:len(bufs)-1]
     }
     lock.Unlock()
     return buf
}


//回收内存到缓冲区
func putBuf(buf []byte){
     lock.Lock()
     if len(bufs) < maxBufPoolSize {
        bufs = append(bufs, buf[0:])
     }
	 
     lock.Unlock()
}

//处理日志请求
func handleClient(conn *net.UDPConn){
	buf := popBuf()
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil || n == 0{
		return
	}

	logutil.Write(buf[0:n])
	putBuf(buf)
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr,"Fatal error ",err.Error())
		os.Exit(1)
	}
}