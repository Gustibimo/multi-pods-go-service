package routes

import (
	"bom-import-xls/internal/app/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", handlers.HelloHandler)
	e.POST("/upload", handlers.UploadHandler)
}
