package transfer

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net"
	"sync"
	"time"
)

//implements rate tcp limit
type Reader struct {
	r       net.Conn
	limiter *rate.Limiter
}

type Writer struct {
	w       net.Conn
	limiter *rate.Limiter
}

const burstLimit = 1000 * 1000 * 1000

func NewReader(r net.Conn) *Reader {
	return &Reader{
		r: r,
	}
}

func NewWriter(w net.Conn) *Writer {
	return &Writer{
		w: w,
	}
}
func (s *Reader) SetRateLimit(bytesPerSec float64) {
	s.limiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.limiter.AllowN(time.Now(), burstLimit) // spend initial burst
}
func (s *Writer) SetRateLimit(bytesPerSec float64) {
	s.limiter = rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	s.limiter.AllowN(time.Now(), burstLimit) // spend initial burst
}

// Optimized GC memory
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

//todo thread pool
func handle(targetPort string, clientCon net.Conn) {
	defer clientCon.Close()
	targetConn, terr := net.Dial("tcp", targetPort)
	if terr != nil {
		log.Println("dialErr ", terr)
		clientCon.Close()
		return
	}
	go handleServerConn(NewWriter(targetConn), NewReader(clientCon))
	for {
		//not sure data length
		var buf = pool.Get().([]byte)
		defer pool.Put(buf)
		CreaDer := NewReader(clientCon)
		CreaDer.SetRateLimit(float64(500))
		num, readErr := CreaDer.r.Read(buf)
		if readErr != nil {
			log.Print("readErr ", readErr, CreaDer.r.RemoteAddr())
			CreaDer.r.Close()
			targetConn.Close()
			return
		}
		TarDer := NewWriter(targetConn)
		TarDer.SetRateLimit(float64(500))
		w, writeErr := TarDer.w.Write(buf[:num])
		if writeErr != nil {
			log.Print("writeErr ", writeErr, w)
			clientCon.Close()
			TarDer.w.Close()
			return
		}
		buf = buf[0:]
	}
}
func handleServerConn(targetConn *Writer, clientCon *Reader) {
	defer targetConn.w.Close()
	for {
		var buf = pool2.Get().([]byte)
		defer pool2.Put(buf)
		num, readErr := targetConn.w.Read(buf)
		targetConn.SetRateLimit(float64(500))
		if readErr != nil {
			log.Print("readErr ", readErr, targetConn.w.RemoteAddr())
			targetConn.w.Close()
			clientCon.r.Close()
			return
		}

		w, writeErr := clientCon.r.Write(buf[:num])
		clientCon.SetRateLimit(float64(500))
		if writeErr != nil {
			log.Print("writeErr ", writeErr, w)
			clientCon.r.Close()
			targetConn.w.Close()
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
