package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_Shorter(t *testing.T) {
	handlers := Handlers{}
	app := fiber.New()
	app.Post("/shorten", handlers.ShortenUrl)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).
		Encode(ShortUrlRequest{LongUrl: "https://dashboard.heroku.com/apps/e7ast1c-shrty"})
	if err != nil {
		t.Error(fmt.Sprintf("app.Test error: %s", err))
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", &buf)
	req.Header = map[string][]string{"Content-Type": {"application/json"}}
	resp, err := app.Test(req)
	if err != nil {
		t.Error(fmt.Sprintf("app.Test error: %s", err))
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(fmt.Sprintf("io.ReadAll error: %s", err))
	}
	log.Println(string(all))
}
