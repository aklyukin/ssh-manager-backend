package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/aklyukin/ssh-manager-backend/structures"
	"log"
)

type LoginJSON struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostServer(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	mChannel := c.MustGet("managerChannel").(chan structures.MMessage)

	var server structures.Servers
	c.Bind(&server)

	if server.Hostname != "" {
		if server.Port == "" {
			server.Port = "22"
		}
		dbConn.Create(&server)
		c.JSON(201, gin.H{"success": server})
		mChannel <- structures.MMessage{
			Type: "server",
			Id: server.Id }
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/users
}

func GetServers(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	var servers []structures.Servers
	dbConn.Find(&servers)

	c.JSON(200, servers)

	// curl -i http://localhost:8080/api/servers
}

func GetServer(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	id := c.Params.ByName("id")
	var server structures.Servers
	dbConn.First(&server, id)
	if server.Id != 0 {
		c.JSON(200, server)
	} else {
		c.JSON(404, gin.H{"error": "Server not found"})
	}

	// curl -i http://localhost:8080/api/servers/3
}

func UpdateServer(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	id := c.Params.ByName("id")
	var server structures.Servers
	dbConn.First(&server, id)

	if server.Hostname != "" {
		if server.Id != 0 {
			var newServer structures.Servers
			c.Bind(&newServer)
			result := structures.Servers{
				Id:          server.Id,
				Hostname:    newServer.Hostname,
				Ip:          newServer.Ip,
				ServerUsers: newServer.Ip,
			}
			dbConn.Save(&result)
			c.JSON(200, gin.H{"success": result})
		} else {
			c.JSON(404, gin.H{"error": "Server not found"})
		}
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
}

func DeleteServer(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	id := c.Params.ByName("id")
	var server structures.Servers
	dbConn.First(&server, id)

	if server.Id != 0 {
		dbConn.Delete(&server)
		c.JSON(200, gin.H{"success": "Server #" + id + " deleted"})
	} else {
		c.JSON(404, gin.H{"error": "Server not found"})
	}
}

func PostUser(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	var user structures.Users

	var sshkey string
	c.BindJSON(&sshkey)
	log.Printf(sshkey)

	c.Bind(&user)
	if user.UserName != "" {
		dbConn.Create(&user)
		c.JSON(201, gin.H{"success": user})
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

func GetUsers(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	var users []structures.Users
	dbConn.Find(&users)

	c.JSON(200, users)
}

func DeleteUser(c *gin.Context) {
	dbConn, ok := c.MustGet("databaseConn").(*gorm.DB)
	if !ok {
		c.JSON(422, gin.H{"error": "database error"})
	}

	id := c.Params.ByName("id")
	var user structures.Users
	dbConn.First(&user, id)

	if user.Id != 0 {
		dbConn.Delete(&user)
		c.JSON(200, gin.H{"success": "User #" + id + " deleted"})
	} else {
		c.JSON(404, gin.H{"error": "User not found"})
	}
}