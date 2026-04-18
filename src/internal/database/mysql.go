package database

import (
	"database/sql"
	"fmt"
	"inventory-app/internal/config"
	"log"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLDB(conf config.MysqlDB) *sql.DB {
	user := url.QueryEscape(conf.User)
	pass := url.QueryEscape(conf.Pass)
	tz := url.QueryEscape(conf.Tz)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		user,
		pass,
		conf.Host,
		conf.Port,
		conf.Name,
		tz,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB Open err:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB Ping err:", err)
	}

	log.Println("✅ MySQL Connected")

	return db
}
