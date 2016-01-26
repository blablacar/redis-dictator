package node

import (
	"gopkg.in/redis.v3"
	"strconv"
	"errors"
)

type Node struct {
	Name string
	Host string
	Port int
	Role string
}

func(rn *Node) Initialize(Name string, Host string, Port int) (error) {
	rn.Name = Name
	rn.Host = Host
	rn.Port = Port
	rn.Role = "UNKNOWN"
	return nil
}

func(rn *Node) SlaveOf(host string, port string) (error) {
	 client := redis.NewClient(&redis.Options{
        Addr:     rn.Host + ":" + strconv.Itoa(rn.Port),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    slaveOf := client.SlaveOf(host, port)
    if slaveOf.Val() != "OK"{
    	return slaveOf.Err()
    }
    return nil
}

func(rn *Node) SetRole(role string, master *Node) (error) {
	switch role {
	case "MASTER":
		rn.Role = "MASTER"
		err := rn.SlaveOf("NO", "ONE")
		if err != nil {
			return err
		}
	case "SLAVE":
		if rn.Is(master) {
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

func(rn *Node) Is(n *Node) (bool) {
	if rn.Host == n.Host && rn.Port == n.Port {
		return true
	}else{
		return false
	}
}