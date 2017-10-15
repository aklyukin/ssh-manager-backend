package structures

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Servers struct {
	Id          uint `gorm:"primary_key"`
	Hostname    string
	Port		string
	Ip          string
	Ip_status	string
	ServerUsers string
}

// Таблица пользователей на серверах и добавленных
type ServerUsers struct {
	ServerId 	uint `gorm:"unique_index:idx_serverusers"`
    UserName 	string `gorm:"unique_index:idx_serverusers"`
    UserId 		uint `gorm:"unique_index:idx_serverusers"`
}

// Таблица id реальных пользователей и их ключей
type SshKeys struct {
	UserId 		uint `gorm:"unique_index:idx_sshkeys"`
	SshKey		string `gorm:"unique_index:idx_sshkeys"`
	Comment		string
}

// Таблица реальных пользователей
type Users struct {
	Id    		uint `gorm:"primary_key"`
	UserName  	string
}

// Not database stuctures
// Manager message
type MMessage struct {
	Type string // user or server
	Id	uint // id of user/server
}
