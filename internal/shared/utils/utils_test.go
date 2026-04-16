package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

	body := DecodeBodyMap(t, w.Body.Bytes())
	AssertMessageEquals(t, body, "erro teste")
	AssertDataNil(t, body)
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

	body := DecodeBodyMap(t, w.Body.Bytes())
	AssertMessageEquals(t, body, "ok")
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

	body := DecodeBodyMap(t, w.Body.Bytes())
	AssertMessageEquals(t, body, "ok")
}

func TestHttpResponseHelpers(t *testing.T) {
	t.Run("DecodeBodyMap success", func(t *testing.T) {
		body := DecodeBodyMap(t, []byte(`{"mensagem":"ok","data":null}`))
		if body["mensagem"] != "ok" {
			t.Fatalf("mensagem = %v", body["mensagem"])
		}
	})

	t.Run("DecodeBodyMap fatal path", func(t *testing.T) {
		oldFatalf := testFatalf
		defer func() { testFatalf = oldFatalf }()

		called := false
		msg := ""
		testFatalf = func(_ testing.TB, format string, args ...any) {
			called = true
			msg = format
		}

		body := DecodeBodyMap(t, []byte(`{invalid`))
		if !called {
			t.Fatal("expected fatal path to be called")
		}
		if !strings.Contains(msg, "unmarshal error") {
			t.Fatalf("unexpected fatal message format: %s", msg)
		}
		if body != nil {
			t.Fatalf("expected nil map on decode error, got %+v", body)
		}
	})

	t.Run("AssertMessagePresent success", func(t *testing.T) {
		AssertMessagePresent(t, map[string]any{"mensagem": "ok", "data": nil})
	})

	t.Run("AssertMessagePresent fatal path", func(t *testing.T) {
		oldFatalf := testFatalf
		defer func() { testFatalf = oldFatalf }()

		called := false
		testFatalf = func(_ testing.TB, _ string, _ ...any) { called = true }
		AssertMessagePresent(t, map[string]any{"data": nil})
		if !called {
			t.Fatal("expected fatal path to be called")
		}
	})

	t.Run("AssertMessageEquals success", func(t *testing.T) {
		AssertMessageEquals(t, map[string]any{"mensagem": "ok"}, "ok")
	})

	t.Run("AssertMessageEquals fatal path", func(t *testing.T) {
		oldFatalf := testFatalf
		defer func() { testFatalf = oldFatalf }()

		called := false
		testFatalf = func(_ testing.TB, _ string, _ ...any) { called = true }
		AssertMessageEquals(t, map[string]any{"mensagem": "x"}, "ok")
		if !called {
			t.Fatal("expected fatal path to be called")
		}
	})

	t.Run("AssertDataNil success", func(t *testing.T) {
		AssertDataNil(t, map[string]any{"data": nil})
	})

	t.Run("AssertDataNil fatal path", func(t *testing.T) {
		oldFatalf := testFatalf
		defer func() { testFatalf = oldFatalf }()

		called := false
		testFatalf = func(_ testing.TB, _ string, _ ...any) { called = true }
		AssertDataNil(t, map[string]any{"data": 1})
		if !called {
			t.Fatal("expected fatal path to be called")
		}
	})

	t.Run("AssertDataArray success", func(t *testing.T) {
		data := AssertDataArray(t, map[string]any{"data": []interface{}{1.0, "x"}})
		if len(data) != 2 {
			t.Fatalf("expected len(data)=2, got %d", len(data))
		}
	})

	t.Run("AssertDataArray fatal path", func(t *testing.T) {
		oldFatalf := testFatalf
		defer func() { testFatalf = oldFatalf }()

		called := false
		testFatalf = func(_ testing.TB, _ string, _ ...any) { called = true }
		data := AssertDataArray(t, map[string]any{"data": "not-array"})
		if !called {
			t.Fatal("expected fatal path to be called")
		}
		if data != nil {
			t.Fatalf("expected nil data on type mismatch, got %+v", data)
		}
	})
}
