package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	gin.SetMode(gin.ReleaseMode)
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
