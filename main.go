package main

import (
    "os"
    "person-api/database"
    "person-api/handlers"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
)

// main запускает REST API сервер
func main() {
    // Загружаем переменные окружения
    if err := godotenv.Load(); err != nil {
        logrus.Fatal("Ошибка загрузки .env: ", err)
    }
    // Настраиваем логирование
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetOutput(os.Stdout)
    logrus.SetLevel(logrus.InfoLevel)
    // Выполняем миграции
    logrus.Info("Запуск миграций...")
    if err := database.RunMigrations(); err != nil {
        logrus.Fatal("Ошибка миграций: ", err)
    }
    // Инициализируем базу данных
    db := database.InitDB()
    // Настраиваем маршруты API
    r := gin.Default()
    r.POST("/people", handlers.CreatePerson(db))
    r.GET("/people", handlers.GetPeople(db))
    r.GET("/people/:id", handlers.GetPerson(db))
    r.PUT("/people/:id", handlers.UpdatePerson(db))
    r.DELETE("/people/:id", handlers.DeletePerson(db))
    // Определяем порт сервера
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    logrus.Infof("Сервер запускается на порту %s", port)
    // Запускаем сервер
    r.Run(":" + port)
}
