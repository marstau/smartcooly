package model

import (
	"log"
	"strings"
	"time"
	"fmt"
	"os"

	"github.com/hprose/hprose-golang/io"
	"github.com/jinzhu/gorm"
	// for db SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/marstau/smartcooly/config"
)

var (
	// DB Database
	DB     *gorm.DB
	dbType string
	dbURL  string
	DBTYPES = map[string]string {
		"postgres" : "Postgres",
	}
)

func init() {

	if os.Getenv("DATABASE_URL") != "" {
		url := os.Getenv("DATABASE_URL")

		step := 0
		slashStep := 0
		var strs [6]string
		for i := 0; i < len(url); i++ {
			c := url[i]

			if c == ':' || c == '@' {
				step++
				continue
			} else if c == '/' {
				slashStep++
				continue
			}

			
			if slashStep < 3 {
				strs[step] += string(c)
			} else {
				strs[step + 1] += string(c)
			}
		}
		dbType = DBTYPES[strs[0]]
		dbURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", strs[3], strs[4], strs[1], strs[5], strs[2])
		
		log.Printf("dbURL=%v", dbURL)
	} else {
		dbType = config.String("dbtype")
		dbURL  = config.String("dburl")
	}

	io.Register((*User)(nil), "User", "json")
	io.Register((*Exchange)(nil), "Exchange", "json")
	io.Register((*Algorithm)(nil), "Algorithm", "json")
	io.Register((*Trader)(nil), "Trader", "json")
	io.Register((*Log)(nil), "Log", "json")
	var err error
	DB, err = gorm.Open(strings.ToLower(dbType), dbURL)
	if err != nil {
		log.Printf("Connect to %v database error: %v\n", dbType, err)
		dbType = "sqlite3"
		dbURL = "custom/data.db"
		DB, err = gorm.Open(dbType, dbURL)
		if err != nil {
			log.Fatalln("Connect to database error:", err)
		}
	}
	DB.AutoMigrate(&User{}, &Exchange{}, &Algorithm{}, &TraderExchange{}, &Trader{}, &Log{})
	users := []User{}
	DB.Find(&users)
	if len(users) == 0 {
		admin := User{
			Username: "admin",
			Password: "admin",
			Level:    99,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalln("Create admin error:", err)
		}
	}
	DB.LogMode(false)
	go ping()
}

func ping() {
	for {
		if err := DB.Exec("SELECT 1").Error; err != nil {
			log.Println("Database ping error:", err)
			if DB, err = gorm.Open(strings.ToLower(dbType), dbURL); err != nil {
				log.Println("Retry connect to database error:", err)
			}
		}
		time.Sleep(time.Minute)
	}
}

// NewOrm ...
func NewOrm() (*gorm.DB, error) {
	return gorm.Open(strings.ToLower(dbType), dbURL)
}
