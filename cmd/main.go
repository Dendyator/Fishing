package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dendyator/Fishing/internal/config"
	"github.com/Dendyator/Fishing/internal/database"
	"github.com/Dendyator/Fishing/internal/handlers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Формирование DSN для миграций
	dsn := cfg.GetDatabaseDSN()

	// Выполнение миграций
	migrationURL := "file://migrations"
	m, err := migrate.New(migrationURL, dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации миграций: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}
	log.Println("Миграции успешно применены")

	// Инициализация базы данных
	db, err := database.NewDB(dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация базовых данных из config.yaml
	if err := db.InitializeData(cfg.FishingMethods, cfg.FishSpecies, cfg.Baits); err != nil {
		log.Fatalf("Ошибка инициализации данных: %v", err)
	}

	// Инициализация маршрутизатора
	router := mux.NewRouter()

	// Подключение обработчиков API
	handlers.SetupRoutes(router, db)

	// Обслуживание статических файлов
	staticDir := http.Dir("./static")
	fileServer := http.FileServer(staticDir)
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	// Запуск сервера
	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Запуск сервера на порту %s", cfg.Port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Сервер не смог запуститься: %v", err)
	}
}
