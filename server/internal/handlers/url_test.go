package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlers_Shorter(t *testing.T) {
	handlers := Handlers{}
	app := fiber.New()
	app.Get("/live", handlers.HealthCheck)

	req := httptest.NewRequest("GET", "/live", nil)

	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != fiber.StatusOK {
		t.Error(fmt.Sprintf("app.Test error: %s", err))
	}
	err = resp.Write(os.Stdout)
}
