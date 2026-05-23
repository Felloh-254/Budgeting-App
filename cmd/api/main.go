package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Felloh-254/Budgeting-App/cmd/api/common"
	"github.com/Felloh-254/Budgeting-App/cmd/api/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type App struct {
	server  *echo.Echo
	handler *handlers.Handler
}

func main() {
	e := echo.New()

	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		e.Logger.Fatal("Error loading .env file", "error", err)
	}

	// Load the database configuration
	dbConfig, err := common.ConnectToDatabase()
	if err != nil {
		e.Logger.Fatal("Failed to connect to database", "error", err)
		return
	}
	defer dbConfig.DB.Close()

	h := &handlers.Handler{DB: dbConfig}

	app := &App{
		server:  e,
		handler: h,
	}

	fmt.Println(app)

	// Middleware

	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":" + os.Getenv("APP_PORT")); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
