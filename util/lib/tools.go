package lib

import (
	"github.com/inconshreveable/go-update"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const VERSION string = "1"
const MIDURL string = "http://111.231.82.173/"
const MIDFILE string = "http://111.231.82.173/file/"
const MIDAUTH string = "http://111.231.82.173:9000/auth"
const MIDETCD string = "111.231.82.173:2379"

func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}

func DoUpdate() error {
	for {
		resp, err := http.Get(MIDFILE + "version.txt")
		if err != nil {
			return err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		current_file := strings.Split(os.Args[0], "\\")
		frpe, err := http.Get(MIDFILE + current_file[len(current_file)-1])
		if strings.TrimSpace(string(body)) != VERSION {
			err = update.Apply(frpe.Body, update.Options{TargetPath: os.Args[0]})
			if err != nil {
				// error handling
			}
		} else {
			//fmt.Println(string(body))
			time.Sleep(time.Duration(RandInt64(300, 1000)))
		}
		return err
	}
}

//func SystemCheck(){
//	switch runtime.GOOS {
//	case "windows":
//		current_file := strings.Split(os.Args[0], "\\")
//		c := exec.Command("cmd", "/C", "taskkill", "/IM",current_file[len(current_file)-1])
//		if err := c.Run(); err != nil {
//			fmt.Println("Error: ", err)
//		}
//	case "linux":
//	case "darwin":
//
//	case "freebsd":
//
//	}
//}
