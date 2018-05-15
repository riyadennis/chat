package handlers

import (
	"testing"
	"net/http"
	"strings"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"fmt"
	"mime/multipart"
)

func TestUploadHandlerServeHTTPWitInvalidForm(t *testing.T) {
	body := strings.NewReader("InvalidRequest")
	request, err := http.NewRequest("POST", "/handler", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)
	uploadHandler := uploadHandler{}
	writer := httptest.NewRecorder()
	uploadHandler.ServeHTTP(writer, request)
	assert.Equal(t, writer.Code, http.StatusInternalServerError)
	fmt.Printf("%+v", writer.Body)
	assert.Equal(t, "User not found in the request\n", writer.Body.String())
}

func TestUploadHandlerServeHTTPWitInvalidFileInForm(t *testing.T) {
	body := strings.NewReader("user_id=12")
	request, err := http.NewRequest("POST", "/handler", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	formValue := make(map[string][]string)
	fileValue := make(map[string][]*multipart.FileHeader)
	request.MultipartForm = &multipart.Form{Value:formValue,File: fileValue}
	assert.NoError(t, err)
	uploadHandler := uploadHandler{}
	writer := httptest.NewRecorder()
	uploadHandler.ServeHTTP(writer, request)
	assert.Equal(t, writer.Code, http.StatusInternalServerError)
	assert.Equal(t, "http: no such file\n", writer.Body.String())
}



