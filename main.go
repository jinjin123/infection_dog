package main

import (
	"infection/etcd"
	//"infection/machineinfo"
	//"strings"
	"sync/atomic"
	"systray"
	//"time"
	"infection/util/icon"
	"infection/util/lib"
)

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
	systray.SetIcon(icon.Data)
	systray.SetTitle("freedom")
	systray.SetTooltip("running...")
	mQuit := systray.AddMenuItem("Quit", "Quit freedom")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}
func onExit() {
	// clean up here
}
func main() {
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)
	//machineinfo.MachineSend(conf.Url)
	go lib.DoUpdate(conf.Url)
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
