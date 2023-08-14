package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"janjiss.com/rest/users"
)

// This is a mock implementation of an in-memory user service
type InMemoryUserService struct {
	storage []*users.User
}

func (us *InMemoryUserService) CreateUser(name, email string) (*users.User, error) {
	user := &users.User{Name: name, Email: email}
	us.storage = append(us.storage, user)
	return user, nil
}

func (us *InMemoryUserService) GetAllUsers(cursor string) ([]users.User, error) {
	var users []users.User

	for _, user := range us.storage {
		users = append(users, *user)
	}

	return users, nil
}

func (us *InMemoryUserService) Login(email string) (string, error) {
	return "JWT_TOKEN", nil
}

func (us *InMemoryUserService) FindOneByEmail(email string) (*users.User, error) {

	for _, user := range us.storage {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, nil
}

func TestBuildCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	userService := &InMemoryUserService{storage: make([]*users.User, 0)}

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

		expectedResponse := map[string]interface{}{
			"user": map[string]interface{}{
				"name":  "TestName",
				"email": "test@email.com",
				"id":    "00000000-0000-0000-0000-000000000000",
			},
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("Expected response and received response are not equal: \nGot:  %#v \nWant: %#v", expectedResponse, response)
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

		expectedResponse := map[string]interface{}{
			"error": "invalid character 'i' looking for beginning of value",
		}

		json.Unmarshal(w.Body.Bytes(), &response)

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("Expected response and received response are not equal: \nGot:  %#v \nWant: %#v", expectedResponse, response)
		}
	})
}

func TestBuildGetAllUsersHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// Setup
	userService := &InMemoryUserService{storage: make([]*users.User, 0)}

	router := gin.Default()
	router.GET("/users", BuildGetAllUsersHandler(userService))

	t.Run("success", func(t *testing.T) {
		userService.CreateUser("TestName", "email@user.com")

		body, _ := json.Marshal(&AllUsers{Cursor: "00000000-0000-0000-0000-000000000000"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", bytes.NewBuffer(body))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK but got %v", w.Code)
		}

		expectedResponse := map[string]interface{}{
			"cursor": "AAAAAAAAAAAAAAAAAAAAAA==",
			"users": []interface{}{
				map[string]interface{}{
					"name":  "TestName",
					"email": "email@user.com",
					"id":    "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("Expected response and received response are not equal: \nGot:  %#v \nWant: %#v", expectedResponse, response)
		}

	})
}

func TestBuildLoginHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// Setup
	userService := &InMemoryUserService{storage: make([]*users.User, 0)}

	router := gin.Default()
	router.POST("/login", BuildLoginHandler(userService))

	t.Run("success", func(t *testing.T) {
		userService.CreateUser("TestName", "email@user.com")

		body, _ := json.Marshal(&Login{Email: "email@user.com"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK but got %v", w.Code)
		}

		expectedResponse := map[string]interface{}{
			"token": "JWT_TOKEN",
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if !reflect.DeepEqual(expectedResponse, response) {
			t.Errorf("Expected response and received response are not equal: \nGot:  %#v \nWant: %#v", expectedResponse, response)
		}

	})
}
