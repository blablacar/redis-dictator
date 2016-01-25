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

func(re *Node) Initialize(Name string, Host string, Port int) (error) {
	re.Name = Name
	re.Host = Host
	re.Port = Port
	re.Role = "UNKNOWN"
	return nil
}

func(re *Node) SlaveOf(host string, port string) (error) {
	 client := redis.NewClient(&redis.Options{
        Addr:     re.Host + ":" + strconv.Itoa(re.Port),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    slaveOf := client.SlaveOf(host, port)
    if slaveOf.Val() != "OK"{
    	return slaveOf.Err()
    }
    return nil
}


func(re *Node) SetRole(role string, master *Node) (error) {
	switch role {
	case "MASTER":
		re.Role = "MASTER"
		err := re.SlaveOf("NO", "ONE")
		if err != nil {
			return err
		}
	case "SLAVE":
		if master != nil {
			err := re.SlaveOf(master.Host, strconv.Itoa(master.Port))
			if err != nil {
				return err
			}
			re.Role = "SLAVE"
		}else{
			return errors.New("Master is empty!")
		}
	default:
	}
    return nil
}
