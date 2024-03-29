package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"url-shortener/config"
	"url-shortener/internal/repository"
	"url-shortener/internal/usecase"
)

func TestHandler_GetLinkHandler(t *testing.T) {
	tests := []struct {
		name                 string
		target               string
		expectedStatusCode   int
		expectedResponseHead string
	}{
		{
			name:                 "Ok",
			target:               "/zE",
			expectedStatusCode:   307,
			expectedResponseHead: "http://zrnzruvv7qfdy.ru/hlc65i",
		},
		{
			name:                 "Err",
			target:               "/IVI1",
			expectedStatusCode:   400,
			expectedResponseHead: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &repository.Config{
				DriverName:     "map",
				DataSourcePath: "test",
			}
			repo, err := repository.New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			logic := usecase.New(repo)

			repo.AddLink("http://zrnzruvv7qfdy.ru/hlc65i", "zE", "df")

			conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
			handler := Handler{conf: conf, logic: logic}

			req := httptest.NewRequest("GET", test.target,
				nil)
			w := httptest.NewRecorder()

			router := gin.Default()
			router.GET("/:id", handler.GetLinkHandler)

			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Header().Get("Location"), test.expectedResponseHead)
		})
	}
}

func TestHandler_CreateLinkHandler(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			inputBody:            "vk.com/gasayminajj",
			expectedStatusCode:   201,
			expectedResponseBody: `zE`,
		},
		{
			name:                 "server error",
			inputBody:            "q",
			expectedStatusCode:   500,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &repository.Config{
				DriverName:     "map",
				DataSourcePath: "test",
			}
			repo, err := repository.New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			logic := usecase.New(repo)

			conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
			handler := Handler{conf: conf, logic: logic}

			req := httptest.NewRequest("POST", "/",
				bytes.NewBufferString(test.inputBody))
			w := httptest.NewRecorder()
			// определяем хендлер
			router := gin.Default()
			router.Use(handler.CreateLinkHandler)

			router.ServeHTTP(w, req)
			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_APICreateLinkHandler(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Ok",
			inputBody:            `{"url":"vk.com/gasayminajj"}`,
			expectedStatusCode:   201,
			expectedResponseBody: `{"result":"zE"}`,
		},
		{
			name:                 "server error",
			inputBody:            "q",
			expectedStatusCode:   500,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &repository.Config{
				DriverName:     "map",
				DataSourcePath: "test",
			}
			repo, err := repository.New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			logic := usecase.New(repo)

			conf := &config.Config{Host: "127.0.0.1", DBConfig: cfg}
			handler := Handler{conf: conf, logic: logic}

			req := httptest.NewRequest("POST", "/api/shorten",
				bytes.NewBufferString(test.inputBody))
			w := httptest.NewRecorder()

			router := gin.Default()
			router.Use(handler.APICreateLinkHandler)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
