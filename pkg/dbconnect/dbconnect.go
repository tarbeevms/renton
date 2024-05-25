package dbconnect

import (
	"database/sql"
	"fmt"
	"log"
	"myapp/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const maxRetries = 10
const retryDelay = 5 * time.Second

func ConnectToMySQL() (*sql.DB, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database)

	var db *sql.DB
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to connect to MySQL (attempt %d/%d): %v, %s", i+1, maxRetries, err, dsn)
			time.Sleep(retryDelay)
			continue
		}

		db.SetMaxOpenConns(10)

		if err := db.Ping(); err != nil {
			log.Printf("Failed to ping MySQL DB (attempt %d/%d): %v, %s", i+1, maxRetries, err, dsn)
			db.Close()
			time.Sleep(retryDelay)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("failed to connect to MySQL after %d attempts", maxRetries)
}
