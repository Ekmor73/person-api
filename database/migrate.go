package database

import (
    "fmt"
    "os"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations выполняет миграции базы данных
func RunMigrations() error {
    // Формируем DSN для подключения к PostgreSQL
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
    // Инициализируем миграции из папки migrations
    m, err := migrate.New("file://migrations", dsn)
    if err != nil {
        return err
    }
    defer m.Close()
    // Применяем миграции
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}
