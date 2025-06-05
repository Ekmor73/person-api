package main

import (
    "os"
    "person-api/database"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
)

// main запускает приложение
func main() {
    // Загружаем переменные окружения из .env
    if err := godotenv.Load(); err != nil {
        logrus.Fatal("Ошибка загрузки .env: ", err)
    }
    // Настраиваем логирование
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetOutput(os.Stdout)
    logrus.SetLevel(logrus.InfoLevel)
    // Выполняем миграции базы данных
    logrus.Info("Запуск миграций...")
    if err := database.RunMigrations(); err != nil {
        logrus.Fatal("Ошибка миграций: ", err)
    }
    logrus.Info("Миграции выполнены")
}
