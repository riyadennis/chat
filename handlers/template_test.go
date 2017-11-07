package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTemplateHandler(t *testing.T) {
	templateHandler := NewTemplateHandler("testfile")
	assert.IsType(t, &TemplateHandler{}, templateHandler)
}
func TestTemplateHandlerServeHTTP(t *testing.T) {
	templateHandler := NewTemplateHandler("testfile")
	req, err := http.NewRequest("GET", "/room", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	templateHandler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}
