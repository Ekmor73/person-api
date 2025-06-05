package database

import (
    "fmt"
    "os"
    "github.com/sirupsen/logrus"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// InitDB инициализирует подключение к базе данных
func InitDB() *gorm.DB {
    // Формируем DSN для подключения
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
    logrus.Info("Подключение к базе данных...")
    // Открываем соединение с базой
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        logrus.Fatal("Ошибка подключения к базе: ", err)
    }
    logrus.Info("База данных подключена")
    return db
}
