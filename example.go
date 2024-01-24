package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

var db *gorm.DB
var err error

type Country struct {
	Country_Id  string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Persons struct {
	gorm.Model
	Name    string    `json:"name"`
	Gender  string    `json:"gender"`
	Age     int       `json:"age"`
	Country []Country `json:"country"`
}

func main() {
	router := gin.Default()
	router.GET("/:name", getPerson)
	router.Run("localhost:8080")

}

func getPerson(c *gin.Context) {
	var person Persons
	name := c.Param("name")
	urlGender := fmt.Sprintf("https://api.genderize.io/?name=%s", string(name))
	urlAge := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	urlNation := fmt.Sprintf("https://api.nationalize.io/?name=%s", string(name))
	responeGender, err := http.Get(urlGender)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer responeGender.Body.Close()

	err = json.NewDecoder(responeGender.Body).Decode(&person)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	responseAge, err := http.Get(urlAge)
	if err != nil {
		panic(err)
	}
	defer responseAge.Body.Close()
	err = json.NewDecoder(responseAge.Body).Decode(&person)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	responseNation, err := http.Get(urlNation)
	if err != nil {
		fmt.Println("Error making JSON:", err)
		return
	}
	defer responseNation.Body.Close()

	err = json.NewDecoder(responseNation.Body).Decode(&person)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	c.IndentedJSON(http.StatusOK, person)

	db, err = gorm.Open("postgres", "user=postgres dbname=junior_example sslmode=disable password=123")
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	db.AutoMigrate(&person)
}

func findMaxProbabilityCountry(countries []Country) (string, float64) {
	if len(countries) == 0 {
		return "", 0.0
	}

	maxCountryID := countries[0].Country_Id
	maxProbability := countries[0].Probability

	for _, country := range countries {
		if country.Probability > maxProbability {
			maxProbability = country.Probability
			maxCountryID = country.Country_Id
		}
	}

	return maxCountryID, maxProbability
}
