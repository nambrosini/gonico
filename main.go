package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

var people = []Person{
	{ID: "1", FirstName: "Nico", LastName: "Ambrosini", Age: 24},
	{ID: "2", FirstName: "Elia", LastName: "Ambrosini", Age: 22},
	{ID: "3", FirstName: "Ari", LastName: "Ambrosini", Age: 57},
}

func main() {
	r := gin.Default()

    AssignRoutes(r)

	r.Run()
}

func AssignRoutes(r *gin.Engine) {
	r.GET("/", HomepageHandler)
	r.GET("/people", GetAllPeopleHandler)
	r.GET("/person/:id", GetPersonHandler)

	r.POST("/person", NewPersonHandler)

	r.PUT("/person/:id", UpdatePersonHandler)

	r.DELETE("/person/:id", DeletePersonHandler)
}

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to gonico!",
	})
}

func GetAllPeopleHandler(c *gin.Context) {
	c.JSON(http.StatusOK, people)
}

func GetPersonHandler(c *gin.Context) {
	id := c.Param("id")

	for i := 0; i < len(people); i++ {
		if people[i].ID == id {
			c.JSON(http.StatusOK, people[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Person not found.",
	})
}

func NewPersonHandler(c *gin.Context) {
	var newPerson Person

	if err := c.ShouldBindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

    if len(newPerson.ID) == 0 {
        newPerson.ID = xid.New().String()
    }

	people = append(people, newPerson)
	c.JSON(http.StatusCreated, newPerson)
}

func UpdatePersonHandler(c *gin.Context) {
	id := c.Param("id")
	var person Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	for i := 0; i < len(people); i++ {
		if people[i].ID == id {
			people[i] = person
			c.JSON(http.StatusOK, person)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Person not found",
	})
}

func DeletePersonHandler(c *gin.Context) {
	id := c.Param("id")

	for i := 0; i < len(people); i++ {
		if people[i].ID == id {
			person := people[i]
			people = append(people[:i], people[i+1:]...)
			c.JSON(http.StatusOK, person)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Person not found",
	})
}
