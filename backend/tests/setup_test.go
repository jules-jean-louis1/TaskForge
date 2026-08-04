package tests

import (
	"os"
	"time"

	"cmd/api/internal/db"
	"cmd/api/internal/middleware"
	"cmd/api/internal/models"
	"cmd/api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	_ = godotenv.Load("../docker/.env")
	_ = godotenv.Load(".env")

	if os.Getenv("JWT_SECRET") == "" {
		_ = os.Setenv("JWT_SECRET", "super-secret-key-for-test")
	}

	router := gin.Default()

	dsn := os.Getenv("DATABASE_TEST_URL")
	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5433/taskforge_test_db?sslmode=disable"
	}

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Échec de connexion à la BDD de test: " + err.Error())
	}
	db.SetDB(testDB)

	_ = testDB.AutoMigrate(&models.User{})

	apiGroup := router.Group("/api/v1")
	routes.UserRoutes(apiGroup)
	routes.TicketRoutes(apiGroup)

	return router
}

func GenerateTestToken(userIDStr string, role string) string {
	_ = godotenv.Load("../docker/.env")
	_ = godotenv.Load(".env")

	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		secretStr = "super-secret-key-for-test"
		_ = os.Setenv("JWT_SECRET", secretStr)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		userID = uuid.New()
	}

	var u models.User
	email := "test@example.com"
	if err := testDB.First(&u, "id = ?", userID).Error; err == nil {
		email = u.Email
	}

	claims := middleware.Claims{
		ID:        userID.String(),
		Firstname: "Test",
		Lastname:  "User",
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretStr))
	if err != nil {
		panic("Échec de signature du token de test: " + err.Error())
	}

	return tokenStr
}
