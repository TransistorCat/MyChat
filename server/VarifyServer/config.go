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
		Host string `json:"host"`
		Port int    `json:"port"`
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
	RedisHost = config.Redis.Host
	RedisPort = config.Redis.Port
	RedisPasswd = config.Redis.Passwd
}

func init() {
	LoadConfig("config.json")
	fmt.Println("Configuration loaded successfully")
}
