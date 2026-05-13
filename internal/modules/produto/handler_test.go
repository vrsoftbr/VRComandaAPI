package produto

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"vrcomandaapi/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type serviceStub struct {
	listFn func(ctx context.Context, req ListProdutosRequest) (interface{}, error)
}

func (s serviceStub) List(ctx context.Context, req ListProdutosRequest) (interface{}, error) {
	return s.listFn(ctx, req)
}

func TestHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns 500 when service fails", func(t *testing.T) {
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListProdutosRequest) (interface{}, error) {
			if req.IDLoja != 5 || req.CodigoBarras != "123" || req.DescricaoCompleta != "REPOLHO" || req.DescricaoCupom != "KG" {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return nil, errors.New("boom")
		}})

		r := gin.New()
		r.GET("/produtos", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/produtos?idLoja=5&codigoBarras=123&descricaocompleta=REPOLHO&descricaocupom=KG", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", w.Code)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertDataNil(t, body)
	})

	t.Run("returns 200 with data", func(t *testing.T) {
		called := 0
		h := NewHandler(serviceStub{listFn: func(_ context.Context, req ListProdutosRequest) (interface{}, error) {
			called++
			if req.IDLoja != 7 || req.CodigoBarras != "789" || req.DescricaoCompleta != "ROXO" || req.DescricaoCupom != "KG" {
				t.Fatalf("unexpected request passed to service: %+v", req)
			}
			return ProdutosPaginatedResponse{
				Items: []ProdutoResponse{{ID: "1", IDProduto: 2, IDLoja: 7, DescricaoCompleta: "REPOLHO ROXO KG", DescricaoCupom: "REPOLHO ROXO KG"}},
				Page:  1,
				Limit: 20,
				Total: 1,
				Pages: 1,
			}, nil
		}})

		r := gin.New()
		r.GET("/produtos", h.List)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/produtos?idLoja=7&codigoBarras=789&descricaocompleta=ROXO&descricaocupom=KG", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		if called != 1 {
			t.Fatalf("service calls = %d", called)
		}

		body := utils.DecodeBodyMap(t, w.Body.Bytes())
		utils.AssertMessageEquals(t, body, "ok")
		data := body["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		if len(items) != 1 {
			t.Fatalf("expected one item in data, got %d", len(items))
		}
	})
}
