package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiveness(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	res := httptest.NewRecorder()

	s.Router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	body := strings.Trim(res.Body.String(), "\n")
	assert.Equal(t, `{"message":"Application healthy!"}`, body)
}
