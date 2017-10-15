package main

import (
	"github.com/gin-gonic/gin"
	"github.com/aklyukin/ssh-manager-backend/structures"
	"github.com/aklyukin/ssh-manager-backend/database"
)

// ApiMiddleware will add the db connection to the context
func DbMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", database.DB)
		c.Next()
	}
}

func ManagerMiddleware(mchan chan structures.MMessage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("managerChannel", mchan)
		c.Next()
	}
}

func buildRoutes( mchan chan structures.MMessage) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(DbMiddleware())
	r.Use(ManagerMiddleware(mchan))

	api := r.Group("api")
	{
		api.POST("/servers", PostServer)
		api.GET("/servers", GetServers)
		api.GET("/servers/:id", GetServer)
		api.PUT("/servers/:id", UpdateServer)
		api.DELETE("/servers/:id", DeleteServer)
		api.POST("/users", PostUser)
		api.GET("/users", GetUsers)
		api.DELETE("/users/:id", DeleteUser)
	}

	return r
}