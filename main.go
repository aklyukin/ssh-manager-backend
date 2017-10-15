package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/aklyukin/ssh-manager-backend/database"
	"github.com/aklyukin/ssh-manager-backend/manager"
	"github.com/aklyukin/ssh-manager-backend/structures"
	"time"
	"fmt"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006/01/02 - 15:04:05") + " " + string(bytes))
}


func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	var err error
	database.DB, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer database.DB.Close()

	// Migrate the schema
	database.DB.AutoMigrate(&structures.Servers{})
	database.DB.AutoMigrate(&structures.ServerUsers{})
	database.DB.AutoMigrate(&structures.SshKeys{})
	database.DB.AutoMigrate(&structures.Users{})

	log.Printf("Start ssh-manager-backend")

	mchan := make(chan structures.MMessage)
	r := buildRoutes(mchan)
	go manager.Runmanager(mchan)

	r.Run(":8080")
}