package handlers

import (
	"bom-import-xls/kafkaclient"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HelloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func UploadHandler(c echo.Context) error {
	kafkaclient.Publish(c.Request().Context(), "bom.xlsx", "file-bom-import-parsing")
	return c.String(http.StatusOK, "UploadHandler")
}
