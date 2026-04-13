package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRespondErrorEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/e", func(c *gin.Context) {
		RespondError(c, http.StatusBadRequest, "erro teste")
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/e", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if body["mensagem"] != "erro teste" {
		t.Fatalf("mensagem = %v", body["mensagem"])
	}
	if v, ok := body["data"]; !ok || v != nil {
		t.Fatalf("expected data=nil, got %v", v)
	}
}

func TestRespondOKEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/ok", func(c *gin.Context) {
		RespondOK(c, http.StatusOK, gin.H{"x": 1})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ok", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if body["mensagem"] != "ok" {
		t.Fatalf("mensagem = %v", body["mensagem"])
	}
}

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/health", HealthHandler)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if body["mensagem"] != "ok" {
		t.Fatalf("mensagem = %v", body["mensagem"])
	}
}
