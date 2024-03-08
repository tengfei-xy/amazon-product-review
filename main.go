package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/tengfei-xy/go-log"
	"gopkg.in/yaml.v3"
)

// https://www.amazon.de/dp/B0C9ZV7BX6

// https://www.amazon.de/Homimaster-Gaming-Stuhl-verstellbare-Belastbarkeit/product-reviews/B0C9ZV7BX6/ref=cm_cr_dp_d_show_all_btm?ie=UTF8&reviewerType=all_reviews
// https://www.amazon.de/product-reviews/B0C9ZV7BX6/ref=cm_cr_dp_d_show_all_btm?ie=UTF8&reviewerType=all_reviews

func init_config(flag flagStruct) {
	log.Infof("读取配置文件:%s", flag.config_file)

	yamlFile, err := os.ReadFile(flag.config_file)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &app)
	if err != nil {
		panic(err)
	}

}

func init_flag() flagStruct {
	var f flagStruct
	flag.StringVar(&f.asin, "asin", "", "指定asin，格式B开头")
	flag.StringVar(&f.config_file, "c", "config.yaml", "打开配置文件")
	flag.Parse()
	if f.asin == "" {
		log.Error("请指定asin")
		os.Exit(1)
	}
	return f
}
func init_mysql() {
	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", app.Mysql.Username, app.Mysql.Password, app.Mysql.Ip, app.Mysql.Port, app.Mysql.Database))
	if err != nil {
		panic(err)
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	log.Info("数据库已连接")
	app.db = DB
}

var app appConfig

func main() {
	flag := init_flag()
	init_config(flag)
	init_mysql()

	reviews_main(flag.asin)
}
func (app *appConfig) get_cookie() (string, error) {
	var cookie string

	if err := app.db.QueryRow("select cookie from cookie ORDER BY RAND() limit 1").Scan(&cookie); err != nil {
		return "", err
	}
	cookie = strings.TrimSpace(cookie)
	if app.cookie != cookie {
		log.Infof("使用新cookie: %s", cookie)
	}
	app.cookie = cookie
	return app.cookie, nil
}
