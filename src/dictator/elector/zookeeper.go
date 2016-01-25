package elector

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

func(ze *Elector) ZKConnect() (zk.State, error) {
	if ze.ZKConnection != nil {
		state := ze.ZKConnection.State()
		switch state {
			case zk.StateUnknown,zk.StateConnectedReadOnly,zk.StateExpired,zk.StateAuthFailed,zk.StateConnecting: {
				//Disconnect, and let Reconnection happen
				log.Info("Zookeeper Connection is in BAD State [",state,"] Reconnect")
				ze.ZKConnection.Close()
			}
			case zk.StateConnected, zk.StateHasSession: {
				log.Debug("Zookeeper Connection connected(",state,"), nothing to do.")

				return state, nil
			}
			case zk.StateDisconnected: {
				log.Info("Reporter Connection is Disconnected -> Reconnection")
			}
		}
	}
	conn, _, err := zk.Connect(ze.ZKHosts, 10 * time.Second)
	if err != nil {
		ze.ZKConnection = nil
		log.Info("Unable to Connect to ZooKeeper (",err,")")
		return zk.StateDisconnected, err
	}
	ze.ZKConnection = conn
	state := ze.ZKConnection.State()

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
