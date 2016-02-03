package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v3"
	"strconv"
	"errors"
	"time"
)

type Redis struct {
	Name string
	Host string
	Port int
	Role string
	LoadingTimeout int
	Conn *redis.Client
}

func(rn *Redis) Initialize(Name string, Host string, Port int, LoadingTimeout int) (error) {
	rn.Name = Name
	rn.Host = Host
	rn.Port = Port
	rn.LoadingTimeout = LoadingTimeout
	rn.Role = "UNKNOWN"
	return nil
}

func (rn *Redis) Connect() (error){
	var err error
    for i := 0; i < rn.LoadingTimeout; i++ {
		rn.Conn = redis.NewClient(&redis.Options{
	        Addr:     rn.Host + ":" + strconv.Itoa(rn.Port),
	        Password: "", // no password set
	        DB:       0,  // use default DB
	    })

	    err = rn.Conn.Ping().Err()
	    if err != nil {
			log.WithError(err).Debug("Wait for Redis Server...")
			time.Sleep(time.Second)
		}else{
			err = nil
			break
		}
	}
	return err
}

func(rn *Redis) SlaveOf(host string, port string) (error) {
	err := rn.Connect()
	if err != nil {
    	return errors.New("Can't connect to Redis.")
    }

    slaveOf := rn.Conn.SlaveOf(host, port)
    if slaveOf.Val() != "OK"{
    	err = slaveOf.Err()
    }

    err = rn.Conn.Close()

    return err
}

func(rn *Redis) Is(n *Redis) (bool) {
	if rn.Host == n.Host && rn.Port == n.Port {
		return true
	}else{
		return false
	}
}

func(rn *Redis) SetRole(role string, master *Redis) (error) {
	switch role {
	case "MASTER":
		rn.Role = "MASTER"
		err := rn.SlaveOf("NO", "ONE")
		if err != nil {
			return err
		}
	case "SLAVE":
		if rn.Is(master) {
			rn.Role = "MASTER"
			return errors.New("I can't be slave of myself...")
		}
		if master != nil {
			err := rn.SlaveOf(master.Host, strconv.Itoa(master.Port))
			if err != nil {
				return err
			}
			rn.Role = "SLAVE"
		}else{
			return errors.New("Master is empty!")
		}
	default:
		return errors.New("Role Unknown")
	}
	return nil
}

