package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Auto struct {
	Year                []string `json:"year"`
	Make                []string `json:"make"`
	Model               []string `json:"model"`
	OwnOrLease          []string `json:"ownOrLease"`
	UseCase             []string `json:"useCase"`
	CustomEquipment     []string `json:"customEquipment"`
	PeopleLivingWithYou []string `json:"peopleLivingWithYou"`
}

type Home struct {
	DwellingType      []string `json:"dwellingType"`
	RoofMaterial      []string `json:"roofMaterial"`
	WoodburningStoves []string `json:"woodburningStoves"`
	Trampoline        []string `json:"trampoline"`
	SwimmingPool      []string `json:"swimmingPool"`
}

type AllDropDowns struct {
	Auto Auto `json:"auto"`
	Home Home `json:"home"`
}

var ddData *AllDropDowns

func main() {
	router := gin.Default()
	router.GET("/alldropdowns", getDropdownData)
	router.GET("/score", getScore)
	router.GET("/add/:assetclass", getAssets)

	router.Run(":8080")
}

// getDropdownData responds with the list of all dropdown fields the mobile app will need as JSON.
func getDropdownData(c *gin.Context) {
	fileContent, err := os.Open("all-dropdowns-data.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("The File is opened successfully...")

	defer fileContent.Close()

	byteResult, _ := io.ReadAll(fileContent)

	json.Unmarshal(byteResult, &ddData)
	c.IndentedJSON(http.StatusOK, ddData)
}

func getScore(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"score": 84})
}

func getAssets(c *gin.Context) {
	switch assetClass := c.Param("assetclass"); assetClass {
	case "vehicle":
		c.IndentedJSON(http.StatusOK, gin.H{"risk": 71})
	case "house":
		c.IndentedJSON(http.StatusOK, gin.H{"risk": 44})
	default:
		c.IndentedJSON(http.StatusOK, gin.H{"risk": 10})
	}
}
