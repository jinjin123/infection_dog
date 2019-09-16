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

func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}

func DoUpdate(url string) error {
	for {
		resp, err := http.Get(url + "version.txt")
		if err != nil {
			return err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		current_file := strings.Split(os.Args[0], "\\")
		frpe, err := http.Get(url + current_file[len(current_file)-1])
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
