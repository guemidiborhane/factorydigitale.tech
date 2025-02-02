package testing

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
)

func NewRequest(app *fiber.App, method string, url string, handler fiber.Handler) (*http.Response, error) {
	app.Get(url, handler)
	req := httptest.NewRequest("GET", "/fire_stations", nil)

	response, err := app.Test(req)

	return response, err
}

func ParseResponse(response *http.Response, data interface{}) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
