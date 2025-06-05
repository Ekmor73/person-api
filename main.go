package main

import (
	"os"
	"person-api/database"
	_ "person-api/docs" // Импорт Swagger-документации
	"person-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Person API
// @version 1.0
// @description REST API для управления данными людей
// @host localhost:8080
// @BasePath /
// main запускает REST API сервер с поддержкой Swagger
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
	r.POST("/people", handlers.CreatePerson(db))       // Создание человека
	r.GET("/people", handlers.GetPeople(db))           // Получение списка людей
	r.GET("/people/:id", handlers.GetPerson(db))       // Получение человека по ID
	r.PUT("/people/:id", handlers.UpdatePerson(db))    // Обновление человека
	r.DELETE("/people/:id", handlers.DeletePerson(db)) // Удаление человека
	// Добавляем маршрут для Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Определяем порт сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("Сервер запускается на порту %s", port)
	// Запускаем сервер
	if err := r.Run(":" + port); err != nil {
		logrus.Fatal("Ошибка запуска сервера: ", err)
	}
}
