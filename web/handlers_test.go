package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"janjiss.com/rest/users"
)

// This is a mock implementation of an in-memory user service
type InMemoryUserService struct {
	storage map[string]*users.User
}

func (us *InMemoryUserService) CreateUser(name, email string) (*users.User, error) {
	user := &users.User{Name: name, Email: email}
	// For simplicity, use email as the key
	us.storage[email] = user
	return user, nil
}

func (us *InMemoryUserService) GetAllUsers(cursor string) ([]users.User, error) {
	return []users.User{}, nil
}

func (us *InMemoryUserService) Login(email string) (string, error) {
	return "", nil
}

func (us *InMemoryUserService) FindOneByEmail(email string) (*users.User, error) {
	return us.storage[email], nil
}

func TestBuildCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	userService := &InMemoryUserService{storage: make(map[string]*users.User)}

	router := gin.Default()
	router.POST("/users", BuildCreateUserHandler(userService))

	t.Run("success", func(t *testing.T) {
		body, _ := json.Marshal(&CreateUser{Name: "TestName", Email: "test@email.com"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK but got %v", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		userResponse, ok := response["user"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected user in response")
			return
		}
		if userResponse["name"] != "TestName" {
			t.Errorf("Expected TestName but got %v", userResponse["Name"])
		}
	})

	t.Run("bind error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte("invalid json")))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status BadRequest but got %v", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		if _, ok := response["error"]; !ok {
			t.Errorf("Expected error in response")
		}
	})
}
