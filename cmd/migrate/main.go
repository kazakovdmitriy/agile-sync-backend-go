package main

import (
	"backend_go/internal/infrastructure/config"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Путь к миграциям относительно этого файла
const migrationDir = "./migrations"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("использование: go run cmd/migrate/main.go up|down|status|version")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml" // fallback для локальной разработки
	}

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Database.URL == "" {
		log.Fatal("ошибка: DATABASE_URL не задан в конфигурации")
	}

	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		log.Fatal("не удалось открыть подключение к БД:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("не удалось подключиться к БД:", err)
	}

	switch os.Args[1] {
	case "up":
		err = goose.Up(db, migrationDir)
		fmt.Println("миграции успешно применены")
	case "down":
		err = goose.Down(db, migrationDir)
		fmt.Println("последняя миграция откачена")
	case "status":
		err = goose.Status(db, migrationDir)
	case "version":
		v, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Текущая версия БД: %d\n", v)
		return
	default:
		log.Fatalf("неизвестная команда: %s", os.Args[1])
	}

	if err != nil {
		log.Fatal("ошибка миграции:", err)
	}
}
