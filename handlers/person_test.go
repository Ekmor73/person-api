package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "person-api/models"
)

// setupRouter создаёт тестовый роутер и базу данных
func setupRouter() (*gin.Engine, *gorm.DB) {
    gin.SetMode(gin.TestMode)
    // Используем SQLite в памяти для тестов
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.Person{})
    r := gin.Default()
    // Регистрируем маршруты
    r.POST("/people", CreatePerson(db))
    r.GET("/people", GetPeople(db))
    r.GET("/people/:id", GetPerson(db))
    r.PUT("/people/:id", UpdatePerson(db))
    r.DELETE("/people/:id", DeletePerson(db))
    return r, db
}

// TestCreatePerson тестирует создание человека
func TestCreatePerson(t *testing.T) {
    r, _ := setupRouter()
    payload := `{"name":"Дмитрий","surname":"Ушаков","patronymic":"Васильевич"}`
    req, _ := http.NewRequest("POST", "/people", bytes.NewBuffer([]byte(payload)))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    var person models.Person
    json.Unmarshal(w.Body.Bytes(), &person)
    assert.Equal(t, "Дмитрий", person.Name)
    assert.Equal(t, "Ушаков", person.Surname)
}

// TestGetPeople тестирует получение списка людей
func TestGetPeople(t *testing.T) {
    r, db := setupRouter()
    db.Create(&models.Person{Name: "Дмитрий", Surname: "Ушаков"})
    req, _ := http.NewRequest("GET", "/people?limit=10&skip=0", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    var people []models.Person
    json.Unmarshal(w.Body.Bytes(), &people)
    assert.Len(t, people, 1)
    assert.Equal(t, "Дмитрий", people[0].Name)
}

// TestGetPerson тестирует получение человека по ID
func TestGetPerson(t *testing.T) {
    r, db := setupRouter()
    db.Create(&models.Person{ID: 1, Name: "Дмитрий", Surname: "Ушаков"})
    req, _ := http.NewRequest("GET", "/people/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    var person models.Person
    json.Unmarshal(w.Body.Bytes(), &person)
    assert.Equal(t, "Дмитрий", person.Name)
}

// TestUpdatePerson тестирует обновление человека
func TestUpdatePerson(t *testing.T) {
    r, db := setupRouter()
    db.Create(&models.Person{ID: 1, Name: "Дмитрий", Surname: "Ушаков"})
    payload := `{"name":"Иван"}`
    req, _ := http.NewRequest("PUT", "/people/1", bytes.NewBuffer([]byte(payload)))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    var person models.Person
    json.Unmarshal(w.Body.Bytes(), &person)
    assert.Equal(t, "Иван", person.Name)
}

// TestDeletePerson тестирует удаление человека
func TestDeletePerson(t *testing.T) {
    r, db := setupRouter()
    db.Create(&models.Person{ID: 1, Name: "Дмитрий", Surname: "Ушаков"})
    req, _ := http.NewRequest("DELETE", "/people/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
    var response map[string]string
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "Удалён", response["message"])
}
