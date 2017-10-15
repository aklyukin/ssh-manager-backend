package manager

import (
	"github.com/aklyukin/ssh-manager-backend/structures"
	"github.com/aklyukin/ssh-manager-backend/database"
	"net"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/aklyukin/ssh-manager-backend/sshcmd"
	"log"
	"strings"

)


func Runmanager(mchan chan structures.MMessage) {

	log.Printf("Start manager...")
	for {
		select {
		case message := <-mchan:
			if message.Type == "server" {
				log.Printf("Run refresh server for id: " + message.Type + "\n")
				//log.Printf(message.Id)
				log.Printf("\n")
				go refreshServer(message.Id)
			}
		}
	}
}

func refreshServer( id uint) {
	var server structures.Servers
	database.DB.First(&server, id)
	server.Ip = getIpForHost(server.Hostname)
	server.ServerUsers = strings.Join(sshcmd.GetUsers(server.Hostname),", ")
	database.DB.Save(&server)
	//sshcmd.RunCmd(server.Hostname)
}

func getIpForHost(hostname string) string {
	addr,err := net.LookupIP(hostname)
	if err != nil {
		return "Unknown host"
	} else {
		return net.IP.String(addr[0])
	}
	return net.IP.String(addr[0])
}