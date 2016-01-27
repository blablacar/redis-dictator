package elector

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

type ZKDebugLogger struct {}

func(ZKDebugLogger) Printf(format string, a ...interface{}) {
	log.Debug(format, a)
}

func(ze *Elector) ZKConnect() (zk.State, error) {
	if ze.ZKConnection != nil {
		state := ze.ZKConnection.State()
		switch state {
			case zk.StateUnknown,zk.StateConnectedReadOnly,zk.StateExpired,zk.StateAuthFailed,zk.StateConnecting: {
				//Disconnect, and let Reconnection happen
				log.Warn("Zookeeper connection is in BAD State [",state,"] reconnect")
				ze.ZKConnection.Close()
			}
			case zk.StateConnected, zk.StateHasSession: {
				log.Debug("Zookeeper connection connected(",state,"), nothing to do.")

				return state, nil
			}
			case zk.StateDisconnected: {
				log.Warn("Reporter connection is Disconnected -> Reconnection")
			}
		}
	}
	conn, ev, err := zk.Connect(ze.ZKHosts, 10 * time.Second)
	if err != nil {
		ze.ZKConnection = nil
		log.Warn("Unable to connect to ZooKeeper (",err,")")
		return zk.StateDisconnected, err
	}

	ze.ZKConnection = conn

	var zkLogger ZKDebugLogger
	ze.ZKConnection.SetLogger(zkLogger)

	state := ze.ZKConnection.State()

	ze.ZKEvent = ev
	return state, nil
}

func(ze *Elector) ZKCreateNode(path string, data string, flag int32)(string, error){
	p, err := ze.ZKConnection.Create(path, []byte(data), flag, zk.WorldACL(zk.PermAll))
	return p, err
}

func(ze *Elector) ZKCreatePath(path string) error {
	paths := strings.Split(path, "/")
	full := ""
	for i, node := range paths {
		if i > 0 {
			full +=  "/"
		}
		full += node
		exists, _, err := ze.ZKConnection.Exists(full)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		_, err = ze.ZKConnection.Create(full, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	return nil
}
