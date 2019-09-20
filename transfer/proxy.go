package transfer

import (
	"fmt"
	"log"
	"net"
	"sync"
)

var pool = &sync.Pool{
	New: func() interface{} {
		log.Println("new 1")
		return make([]byte, 9192)
	},
}

var pool2 = &sync.Pool{
	New: func() interface{} {
		log.Println("new 22222222")
		return make([]byte, 9192)
	},
}

func handle(targetPort string, clientCon net.Conn) {
	defer clientCon.Close()
	targetConn, terr := net.Dial("tcp", targetPort)
	if terr != nil {
		log.Println("dialErr ", terr)
		clientCon.Close()
		return
	}
	go handleServerConn(targetConn, clientCon)
	for {
		var buf = pool.Get().([]byte)
		defer pool.Put(buf)
		num, readErr := clientCon.Read(buf)
		if readErr != nil {
			log.Print("readErr ", readErr, clientCon.RemoteAddr())
			clientCon.Close()
			targetConn.Close()
			return
		}
		w, writeErr := targetConn.Write(buf[:num])
		if writeErr != nil {
			log.Print("writeErr ", writeErr, w)
			clientCon.Close()
			targetConn.Close()
			return
		}
		buf = buf[0:]
	}
}
func handleServerConn(targetConn, clientCon net.Conn) {
	defer targetConn.Close()
	for {
		var buf = pool2.Get().([]byte)
		defer pool2.Put(buf)
		num, readErr := targetConn.Read(buf)
		if readErr != nil {
			log.Print("readErr ", readErr, targetConn.RemoteAddr())
			targetConn.Close()
			clientCon.Close()
			return
		}
		w, writeErr := clientCon.Write(buf[:num])
		if writeErr != nil {
			log.Print("writeErr ", writeErr, w)
			clientCon.Close()
			targetConn.Close()
			return
		}
		buf = buf[0:]
	}

}
func InitCfg(target string, localAddr string) {
	localConn, err := net.Listen("tcp", localAddr)
	defer localConn.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		c, err := localConn.Accept()
		if err != nil {
			log.Println("接受连接失败 ", err)
		} else {
			go handle(target, c)
		}
	}
}
