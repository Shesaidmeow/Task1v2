package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gopkg.in/reform.v1"
	mysqld "gopkg.in/reform.v1/dialects/mysql"
)

// Config структура для хранения конфигурации
type Config struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	HTTPPort   string
}

// News структура для хранения новости
//go:generate reform

// reform:News
type News struct {
	ID         int64   `reform:"Id,pk"`
	Title      string  `reform:"Title"`
	Content    string  `reform:"Content"`
	Categories []int64 `json:"Categories"`
}

var db *reform.DB

func main() {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Чтение конфигурации
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
	)
	//fmt.Println(dsn)

	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	
	pool.SetMaxIdleConns(3)
 	pool.SetMaxOpenConns(3)
	
	pingErr := pool.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// Создание таблиц
	_, err = pool.Exec(`
		CREATE TABLE IF NOT EXISTS News (
			Id bigint NOT NULL AUTO_INCREMENT,
			Title tinytext NOT NULL,
			Content longtext NOT NULL,
			PRIMARY KEY (Id)
		);
	`)

	if err != nil {
		log.Fatal("cannot create News table", err)
	}

	_, err = pool.Exec(`
		CREATE TABLE IF NOT EXISTS NewsCategories (
			NewsId bigint NOT NULL,
			CategoryId bigint NOT NULL,
			PRIMARY KEY (NewsId, CategoryId)
		);
	`)

	if err != nil {
		log.Fatal("cannot create NewsCategories table", err)
	}

	// Инициализация reform с использованием подключения к базе данных
	db = reform.NewDB(pool, mysqld.Dialect, nil)

	// Создание Fiber приложения
	app := fiber.New()

	// Регистрация обработчиков
	app.Post("/edit/:Id", editNews)
	app.Get("/list", listNews)

	// Запуск сервера
	port := viper.GetString("HTTP_PORT")
	log.Fatal(app.Listen(":" + port))
}
