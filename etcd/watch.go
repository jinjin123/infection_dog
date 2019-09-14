package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"111.231.82.173:2379"}, //etcd集群三个实例的端口
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")

	defer cli.Close()

	for true {
		rch := cli.Watch(context.Background(), "/url/ip/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q:%q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				fmt.Printf("%q", ev.Kv.Value)
			}
		}
	}
}
