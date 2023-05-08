package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRedis(t *testing.T) {
	sourceFileName := "test/sample1.json"
	jsonFile, err := os.Open(sourceFileName)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", sourceFileName, err)
	}
	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", sourceFileName, err)
	}

	router := gin.Default()
	router.POST("/:tenant", PutRedis)
	router.GET("/:tenant", GetRedis)

	// Store the data in Redis using a POST request
	body := bytes.NewBuffer(content)
	req, _ := http.NewRequest("POST", "/tenant1", body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Fetch the data from Redis using a GET request
	req, _ = http.NewRequest("GET", "/tenant1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Compare the input JSON content with the actual REST request output
	responseBytes := resp.Body.Bytes()

	// Marshal the Go value back into JSON and remove possible extra whitespace
	var jsonData map[string]interface{}
	err = json.Unmarshal(responseBytes, &jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	expectedResponse := make(map[string]interface{})
	err = json.Unmarshal(content, &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal input JSON: %v", err)
	}

	assert.Equal(t, expectedResponse, jsonData)
}

func TestPutRedis(t *testing.T) {
	router := gin.Default()
	router.POST("/:tenant", PutRedis)
	router.GET("/:tenant", GetRedis)

	sourceFileName := "test/sample1.json"
	jsonFile, err := os.Open(sourceFileName)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", sourceFileName, err)
	}
	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", sourceFileName, err)
	}

	body := bytes.NewBuffer(content)
	req, _ := http.NewRequest("POST", "/tenant1", body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Fetch the data from Redis using a GET request
	req, _ = http.NewRequest("GET", "/tenant1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Compare the input JSON content with the actual REST request output
	responseBytes := resp.Body.Bytes()

	// Marshal the Go value back into JSON and remove possible extra whitespace
	var jsonData map[string]interface{}
	err = json.Unmarshal(responseBytes, &jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	expectedResponse := make(map[string]interface{})
	err = json.Unmarshal(content, &expectedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal input JSON: %v", err)
	}

	assert.Equal(t, expectedResponse, jsonData)
}
