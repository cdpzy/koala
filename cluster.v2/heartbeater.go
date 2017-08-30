package cluster

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	client "github.com/coreos/etcd/clientv3"
)

var heartbeaterCtl chan struct{}

func heartbeater() {
	heartbeaterCtl = make(chan struct{})
	go func() {
		defer func() {
			log.Infoln("Node heartbeater stoped.")
		}()

		log.Infoln("Node heartbeater started.")
		ticker := time.After(5 * time.Second)
		for {
			select {
			case <-ticker:
				Nodes.Iterator(func(k string, v *Node) bool {
					heartbeaterTime, _ := v.Params.Int64("Heartbeater")
					s := time.Now().Sub(time.Unix(heartbeaterTime, 0)).Seconds()
					if s > 30 {
						Nodes.UnRegister(v.Name)
					}

					return true
				})

				if Local != nil {
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					etcdClient.Do(ctx, client.OpPut(Nodes.Path(Local)+"/Heartbeater", fmt.Sprint(time.Now().Unix())))
					cancel()
				}

				ticker = time.After(5 * time.Second)

			case <-heartbeaterCtl:
				return
			}
		}
	}()
}

func stopHeartbeater() {
	if heartbeaterCtl != nil {
		close(heartbeaterCtl)
	}
}
