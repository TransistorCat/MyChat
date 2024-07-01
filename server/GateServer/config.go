package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Email struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	} `json:"email"`
	Mysql struct {
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Passwd string `json:"passwd"`
		User   string `json:"user"`
		Schema string `json:"schema"`
	} `json:"mysql"`
	Redis struct {
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Passwd string `json:"passwd"`
	} `json:"redis"`
}

var (
	EmailUser   string
	EmailPass   string
	MysqlHost   string
	MysqlPort   int
	MysqlPasswd string
	MysqlUser   string
	MysqlSchema string
	RedisHost   string
	RedisPort   int
	RedisPasswd string
	CodePrefix  = "code_"
)

func LoadConfig(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	EmailUser = config.Email.User
	EmailPass = config.Email.Pass
	MysqlHost = config.Mysql.Host
	MysqlPort = config.Mysql.Port
	MysqlUser = config.Mysql.User
	MysqlPasswd = config.Mysql.Passwd
	MysqlSchema = config.Mysql.Schema
	RedisHost = config.Redis.Host
	RedisPort = config.Redis.Port
	RedisPasswd = config.Redis.Passwd
}

func init() {
	LoadConfig("config.json")
	fmt.Println("Configuration loaded successfully")
}
