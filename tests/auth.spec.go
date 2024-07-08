package tests

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"h-two/internal/dto"
	"h-two/internal/models"
	"h-two/internal/server"
	"h-two/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

type MockOrganizationRepository struct {
	mock.Mock
}

func (m *MockOrganizationRepository) CreateOrganization(org *models.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepository) GetOrganizationsByUser(userId string) ([]*models.Organization, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) GetOrganizationById(userId string, orgId string) (*models.Organization, error) {
	args := m.Called(userId, orgId)
	return args.Get(0).(*models.Organization), args.Error(1)
}

func (m *MockOrganizationRepository) AddUserToOrganization(orgId string, userId string) error {
	args := m.Called(orgId, userId)
	return args.Error(0)
}

func (m *MockOrganizationRepository) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)

}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) (*dto.UserResponse, error) {
	args := m.Called(user)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(userId string) (*models.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func (m *MockOrganizationRepository) IsUserInOrganization(userId string, orgId string) (bool, error) {
	args := m.Called(userId, orgId)
	return args.Bool(0), args.Error(1)
}

func setupServer() *server.Server {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)

	}
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	log.Println("PORT: ", port)
	db, _, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Create a gorm instance from the sqlmock instance
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//dbInstance := database.New()
	mockRepo := new(MockOrganizationRepository)
	mockRepo.On("Begin").Return(gdb)
	mockRepo.On("GetOrganizationsByUser", mock.AnythingOfType("string")).Return([]*models.Organization{}, nil)
	mockRepo.On("CreateOrganization", mock.AnythingOfType("*models.Organization")).Return(nil)
	organizationService := services.NewOrganizationService(mockRepo)

	userRepo := new(MockUserRepository) // Pass the dbInstance to the UserRepository
	// Create a UserResponse that contains the user details
	userResponse := &dto.UserResponse{
		UserId:    "some-user-id", // Replace with an actual user ID
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "1234567890",
	}

	// Set up the CreateUser method to return the UserResponse
	userRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(userResponse, nil)
	h, _ := services.HashPassword("password123")
	user := &models.User{
		UserId:    "some-user-id", // Replace with an actual user ID
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  h, // This should be the hashed password
		Phone:     "1234567890",
	}

	// Set up the GetUserByEmail method to return the User
	userRepo.On("GetUserByEmail", "john.doe@example.com").Return(user, nil)
	userRepo.On("Begin").Return(gdb)
	authService := services.NewAuthService(userRepo, organizationService) // Pass the UserRepository to the AuthService
	userService := services.NewUserService(userRepo)                      // Assuming you have a function to create a new AuthService
	return &server.Server{
		Port:                port,
		AuthService:         authService,
		UserService:         userService,
		OrganizationService: organizationService,
	}
}

//func TestHelloWorldHandler(t *testing.T) {
//	s := &server.Server{}
//	r := gin.New()
//	r.GET("/", s.HelloWorldHandler)
//	// Create a test HTTP request
//	req, err := http.NewRequest("GET", "/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	// Create a ResponseRecorder to record the response
//	rr := httptest.NewRecorder()
//	// Serve the HTTP request
//	r.ServeHTTP(rr, req)
//	// Check the status code
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
//	}
//	// Check the response body
//	expected := "{\"message\":\"Hello World\"}"
//	if rr.Body.String() != expected {
//		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
//	}
//}

func TestTokenGeneration(t *testing.T) {
	userID := "1"
	token, err := services.GenerateJWT(userID)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(services.SecretKey), nil
	})
	if err != nil {
		t.Fatalf("Error parsing JWT: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("Token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Claims are not of type jwt.MapClaims")
	}

	if claims["userId"] != userID {
		t.Fatalf("Expected userId to be %s, got %s", userID, claims["userId"])
	}

	expirationTime := claims["exp"].(float64)
	if time.Now().Add(services.TokenDuration).Sub(time.Unix(int64(expirationTime), 0)) > time.Minute {
		t.Fatal("Token expiration time is not within expected duration")
	}
}

func TestOrganizationAccessControl(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.GET("/", s.GetOrganizationsHandler)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code to be %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "{\"status\":\"success\",\"message\":\"Organizations retrieved successfully\",\"data\":{\"organisations\":[]}}"
	if rr.Body.String() != expected {
		t.Fatalf("Expected body to be %s, got %s", expected, rr.Body.String())
	}
}

func TestCreateOrganizationHandler(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.POST("/api/organization", s.CreateOrganizationHandler)

	// Define a CreateOrganizationRequest
	reqBody := &dto.CreateOrganizationRequest{
		Name:        "Test Organization",
		Description: "This is a test organization",
	}

	// Encode the request body into JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Error encoding request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	req, err := http.NewRequest("POST", "/api/organization", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status code to be %d, got %d", http.StatusCreated, rr.Code)
	}

	// Define the expected response
	expected := dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Organization created successfully",
		Data: &dto.GetOrganizationResponse{
			OrgId:       "",
			Name:        reqBody.Name,
			Description: reqBody.Description,
		},
	}

	// Convert the expected response to JSON
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Error encoding expected response to JSON: %v", err)
	}

	// Compare the actual response with the expected response
	if rr.Body.String() != string(expectedJSON) {
		t.Fatalf("Expected body to be %s, got %s", string(expectedJSON), rr.Body.String())
	}
}

func TestRegisterUserWithDefaultOrganization(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.POST("/api/register", s.RegisterHandler)

	h, _ := services.HashPassword("password123")
	// Define a CreateUserRequest
	reqBody := &dto.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  h,
		Phone:     "1234567890",
	}

	// Encode the request body into JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Error encoding request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status code to be %d, got %d", http.StatusCreated, rr.Code)
	}
	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Extract the accessToken from the map
	accessToken, ok := responseBody["data"].(map[string]interface{})["accessToken"].(string)
	if !ok || accessToken == "" {
		t.Fatal("Expected accessToken to be present and not empty")
	}

	// Define the expected response
	expected := dto.ApiSuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data: map[string]interface{}{
			"accessToken": accessToken, // You can't predict the actual access token value, so leave it as an empty string
			"user": &dto.UserResponse{
				UserId:    "some-user-id",
				FirstName: reqBody.FirstName,
				LastName:  reqBody.LastName,
				Email:     reqBody.Email,
				Phone:     reqBody.Phone,
			},
		},
	}

	// Convert the expected response to JSON
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Error encoding expected response to JSON: %v", err)
	}

	// Compare the actual response with the expected response
	if rr.Body.String() != string(expectedJSON) {
		t.Fatalf("Expected body to be %s, got %s", string(expectedJSON), rr.Body.String())
	}
}

func TestLoginUserSuccess(t *testing.T) {
	s := setupServer()
	r := gin.New()
	r.POST("/api/login", s.LoginHandler)

	// Define a LoginRequest
	reqBody := &dto.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	// Encode the request body into JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Error encoding request body: %v", err)
	}

	// Create a new HTTP request with the JSON body
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code to be %d, got %d", http.StatusOK, rr.Code)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	// Check if the accessToken field is present and not empty
	accessToken, ok := responseBody["data"].(map[string]interface{})["accessToken"].(string)
	if !ok || accessToken == "" {
		t.Fatal("Expected accessToken to be present and not empty")
	}
}
