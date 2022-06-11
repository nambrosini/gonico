package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testPerson = Person{
    ID:        "1",
    FirstName: "Pinco",
    LastName:  "Pallino",
    Age:       25,
}

// HELPERS
func SetUpRouter() *gin.Engine {
    r := gin.Default()
    AssignRoutes(r)
    CreatePerson(r, testPerson)
    return r
}

func CreatePerson(r *gin.Engine, person Person) int {
    jsonValue, _ := json.Marshal(person)
	req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(jsonValue))
    w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
    return w.Code
}

// TESTS
// GET /
func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"Welcome to gonico!"}`
	
    r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

// GET /person/:id
func TestGetPersonHandler(t *testing.T) {
    r := SetUpRouter()

    code := CreatePerson(r, testPerson)
    assert.Equal(t, http.StatusCreated, code)

    req, _ := http.NewRequest("GET", "/person/" + testPerson.ID, nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    var person Person
    json.Unmarshal(w.Body.Bytes(), &person)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.NotEmpty(t, person)
}

func TestGetPersonNotFound(t *testing.T) {
    r := SetUpRouter()
    req, _ := http.NewRequest("GET", "/person/100", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}

// GET /people
func TestGetPeopleHandler(t *testing.T) {
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/people", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var people []Person
	json.Unmarshal(w.Body.Bytes(), &people)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, people)
}

// POST /person
func TestNewPersonHandler(t *testing.T) {
	r := SetUpRouter()

    person := testPerson
    person.ID = ""
	code := CreatePerson(r, person)

	assert.Equal(t, http.StatusCreated, code)
}

func TestNewPersonHandlerBadRequest(t *testing.T) {
    r := SetUpRouter()

    jsonValue, _ := json.Marshal("")
	req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(jsonValue))
    w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// PUT /person/:id
func TestUpdatePersonHandler(t *testing.T) {
	r := SetUpRouter()

    CreatePerson(r, testPerson)

	person := Person {
		ID:        testPerson.ID,
		FirstName: "John",
		LastName:  "Doe",
		Age:       23,
	}

	jsonValue, _ := json.Marshal(person)
	reqFound, _ := http.NewRequest("PUT", "/person/" + testPerson.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("PUT", "/person/12", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdatePersonBadRequest(t *testing.T) {
    r := SetUpRouter()

    jsonValue, _ := json.Marshal("")
	req, _ := http.NewRequest("PUT", "/person/1", bytes.NewBuffer(jsonValue))
    w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// DELETE /person/:id
func TestDeletePerson(t *testing.T) {
    r := SetUpRouter()

    req, _ := http.NewRequest("DELETE", "/person/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    var person Person
	json.Unmarshal(w.Body.Bytes(), &person)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, person)
}

func TestDeletePersonNotFound(t *testing.T) {
    r := SetUpRouter()

    req, _ := http.NewRequest("DELETE", "/person/10", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}