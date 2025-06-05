package handlers

import (
	"encoding/json"
	"net/http"
	"person-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PersonCreate определяет структуру для создания человека
type PersonCreate struct {
	Name       string `json:"name" binding:"required"`    // Имя (обязательное)
	Surname    string `json:"surname" binding:"required"` // Фамилия (обязательная)
	Patronymic string `json:"patronymic"`                 // Отчество (опционально)
}

// PersonUpdate определяет структуру для обновления человека
type PersonUpdate struct {
	Name        *string `json:"name"`        // Имя (опционально)
	Surname     *string `json:"surname"`     // Фамилия (опционально)
	Patronymic  *string `json:"patronymic"`  // Отчество (опционально)
	Age         *int    `json:"age"`         // Возраст (опционально)
	Gender      *string `json:"gender"`      // Пол (опционально)
	Nationality *string `json:"nationality"` // Национальность (опционально)
}

// GenderizeResponse структура ответа от Genderize.io
type GenderizeResponse struct {
	Gender      string  `json:"gender"`      // Пол
	Probability float64 `json:"probability"` // Вероятность
}

// NationalizeResponse структура ответа от Nationalize.io
type NationalizeResponse struct {
	Country []struct {
		CountryID   string  `json:"country_id"`  // Код страны
		Probability float64 `json:"probability"` // Вероятность
	} `json:"country"`
}

// @Summary Создание нового человека
// @Description Принимает имя, фамилию и отчество. Определяет пол и национальность с помощью внешних API.
// @Tags people
// @Accept json
// @Produce json
// @Param person body PersonCreate true "Данные для создания"
// @Success 200 {object} models.Person
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /people [post]
func CreatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input PersonCreate
		// Валидируем входные данные
		if err := c.ShouldBindJSON(&input); err != nil {
			logrus.Errorf("Ошибка ввода: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Создаём модель человека
		person := models.Person{
			Name:       input.Name,
			Surname:    input.Surname,
			Patronymic: input.Patronymic,
		}
		// Запрашиваем пол через Genderize.io
		if resp, err := http.Get("https://api.genderize.io?name=" + person.Name); err == nil {
			defer resp.Body.Close()
			var genderResp GenderizeResponse
			if err := json.NewDecoder(resp.Body).Decode(&genderResp); err == nil && genderResp.Probability > 0.7 {
				person.Gender = genderResp.Gender
			}
		}
		// Запрашиваем национальность через Nationalize.io
		if resp, err := http.Get("https://api.nationalize.io?name=" + person.Name); err == nil {
			defer resp.Body.Close()
			var natResp NationalizeResponse
			if err := json.NewDecoder(resp.Body).Decode(&natResp); err == nil && len(natResp.Country) > 0 && natResp.Country[0].Probability > 0.3 {
				person.Nationality = natResp.Country[0].CountryID
			}
		}
		logrus.Infof("Создание: %s %s", person.Name, person.Surname)
		// Сохраняем в базе
		if err := db.Create(&person).Error; err != nil {
			logrus.Errorf("Ошибка создания: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать"})
			return
		}
		logrus.Infof("Создан ID: %d", person.ID)
		c.JSON(http.StatusOK, person)
	}
}

// @Summary Получить список людей
// @Description Возвращает список людей с фильтрацией и пагинацией
// @Tags people
// @Produce json
// @Param name query string false "Фильтр по имени"
// @Param surname query string false "Фильтр по фамилии"
// @Param age query int false "Фильтр по возрасту"
// @Param gender query string false "Фильтр по полу"
// @Param nationality query string false "Фильтр по национальности"
// @Param skip query int false "Смещение (пагинация)"
// @Param limit query int false "Ограничение (пагинация)"
// @Success 200 {array} models.Person
// @Failure 500 {object} models.ErrorResponse
// @Router /people [get]
func GetPeople(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var people []models.Person
		query := db
		// Фильтры по параметрам запроса
		if name := c.Query("name"); name != "" {
			query = query.Where("name ILIKE ?", "%"+name+"%")
		}
		if surname := c.Query("surname"); surname != "" {
			query = query.Where("surname ILIKE ?", "%"+surname+"%")
		}
		if age := c.Query("age"); age != "" {
			ageInt, _ := strconv.Atoi(age)
			query = query.Where("age = ?", ageInt)
		}
		if gender := c.Query("gender"); gender != "" {
			query = query.Where("gender = ?", gender)
		}
		if nationality := c.Query("nationality"); nationality != "" {
			query = query.Where("nationality = ?", nationality)
		}
		// Пагинация
		skip, _ := strconv.Atoi(c.DefaultQuery("skip", "0"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		query = query.Offset(skip).Limit(limit)
		// Выполняем запрос
		if err := query.Find(&people).Error; err != nil {
			logrus.Errorf("Ошибка получения: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить"})
			return
		}
		logrus.Infof("Получено: %d записей", len(people))
		c.JSON(http.StatusOK, people)
	}
}

// @Summary Получить человека по ID
// @Description Возвращает информацию о человеке по его ID
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} models.Person
// @Failure 404 {object} models.ErrorResponse
// @Router /people/{id} [get]
func GetPerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var person models.Person
		// Ищем запись
		if err := db.First(&person, id).Error; err != nil {
			logrus.Errorf("Не найден ID=%d: %v", id, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найден"})
			return
		}
		logrus.Infof("Получен ID: %d", id)
		c.JSON(http.StatusOK, person)
	}
}

// @Summary Обновить данные человека
// @Description Обновляет существующего человека по ID. Принимает только те поля, которые нужно изменить.
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body PersonUpdate true "Обновляемые данные"
// @Success 200 {object} models.Person
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /people/{id} [put]
func UpdatePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var person models.Person
		// Ищем запись
		if err := db.First(&person, id).Error; err != nil {
			logrus.Errorf("Не найден ID=%d: %v", id, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найден"})
			return
		}
		var input PersonUpdate
		// Валидируем входные данные
		if err := c.ShouldBindJSON(&input); err != nil {
			logrus.Errorf("Ошибка ввода: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Обновляем поля
		if input.Name != nil {
			person.Name = *input.Name
		}
		if input.Surname != nil {
			person.Surname = *input.Surname
		}
		if input.Patronymic != nil {
			person.Patronymic = *input.Patronymic
		}
		if input.Age != nil {
			person.Age = input.Age
		}
		if input.Gender != nil {
			person.Gender = *input.Gender
		}
		if input.Nationality != nil {
			person.Nationality = *input.Nationality
		}
		// Сохраняем изменения
		if err := db.Save(&person).Error; err != nil {
			logrus.Errorf("Ошибка обновления ID=%d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить"})
			return
		}
		logrus.Infof("Обновлён ID: %d", id)
		c.JSON(http.StatusOK, person)
	}
}

// @Summary Удалить человека
// @Description Удаляет человека по ID
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /people/{id} [delete]
func DeletePerson(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		// Удаляем запись
		if err := db.Delete(&models.Person{}, id).Error; err != nil {
			logrus.Errorf("Ошибка удаления ID=%d: %v", id, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Не найден"})
			return
		}
		logrus.Infof("Удалён ID: %d", id)
		c.JSON(http.StatusOK, gin.H{"message": "Удалён"})
	}
}
