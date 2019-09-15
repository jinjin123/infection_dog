package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Notifyer interface {
	Callback(*Config)
}
type Config struct {
	Url        string
	rwLock     sync.RWMutex
	notifyList []Notifyer
}

type PublicIp struct {
	RemoteAddr string `json:"remote"`
	Msg        string `json:"msg"`
	Code       int    `json:"code"`
}

//first load put config into the memory
func NewConfig() (conf *Config, err error) {
	var public PublicIp
	conf = &Config{}
	resp, err := http.PostForm("http://111.231.82.173:9000/auth", url.Values{"name": {"789"}, "ext": {"789"}, "auth": {"789"}})
	if err != nil {
		fmt.Printf("请检查网络")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal([]byte(body), &public); err == nil {
		if public.Code == -1 {
			fmt.Println(`Error:`, public.Msg)
		} else {
			conf.rwLock.Lock()
			conf.Url = public.RemoteAddr
			conf.rwLock.Unlock()
		}
	}
	// check the new config
	go conf.reload()
	return
}

//add the watcher
func (c *Config) AddObserver(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}

func (c *Config) reload() {
	ticker := time.NewTicker(time.Second * time.Duration(RandInt64(6, 18)))
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"111.231.82.173:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}
	fmt.Println("connect succ")
	defer cli.Close()
	for _ = range ticker.C {
		func() {
			rch := cli.Watch(context.Background(), "/url/ip/")
			for wresp := range rch {
				for _, ev := range wresp.Events {
					//fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					c.rwLock.Lock()
					c.Url = string(ev.Kv.Value)
					c.rwLock.Unlock()
					//notify watcher update
					for _, n := range c.notifyList {
						n.Callback(c)
					}
				}
			}
		}()
	}
}
func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}
