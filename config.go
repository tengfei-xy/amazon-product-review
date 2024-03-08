package main

import "database/sql"

type appConfig struct {
	Mysql      `yaml:"mysql"`
	Basic      `yaml:"basic"`
	Proxy      `yaml:"proxy"`
	db         *sql.DB
	cookie     string
	primary_id int64
}

type Basic struct {
	Test   bool   `yaml:"test"`
	Domain string `yaml:"domain"`
}
type Proxy struct {
	Enable bool     `yaml:"enable"`
	Sockc5 []string `yaml:"socks5"`
}
type Mysql struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
type flagStruct struct {
	asin        string
	config_file string
}
