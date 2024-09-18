package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AssetsRequest struct {
	Assets []string `json:"assets"`
}

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
	// scenarios
}

type Scenario struct {
	ScenarioName string  `json:"scenario_name"`
	Cost         float32 `json:"cost"`
}

type Gap struct {
	AssetName string  `json:"asset_name"`
	Quote     float32 `json:"quote"`
}

type GapsAndScenarios struct {
	Scenario []Scenario `json:"scenarios"`
	Gap      []Gap      `json:"gaps"`
}
type AllDropDowns struct {
	Auto Auto `json:"auto"`
	Home Home `json:"home"`
}

var ddData *AllDropDowns
var gapsAndScenarios *GapsAndScenarios

func main() {
	router := gin.Default()
	router.GET("/alldropdowns", getDropdownData)
	router.GET("/score", getScore)
	router.GET("/add/:assetclass", getAssets)
	router.POST("/calculateRisk", getCalculateRisk)
	// router.GET("/albums/:id", getAlbumByID)

	router.Run(":8080")
}

// getDropdownData responds with the list of all dropdown fields the mobile app will need as JSON.
func getDropdownData(c *gin.Context) {
	fileContent, err := os.Open("jsons/all-dropdowns-data.json")

	if err != nil {
		log.Fatal(err)
		return
	}

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

func getCalculateRisk(c *gin.Context) {
	// Gets a list of assets. Returns a Json of scenariosAndGaps
	var assetsReq AssetsRequest
	if err := c.ShouldBindJSON(&assetsReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	assets := assetsReq.Assets

	if len(assets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No assets provided"})
		return
	}

	fileContent, err := os.Open("jsons/scenarios_and_gaps.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer fileContent.Close()

	byteResult, _ := io.ReadAll(fileContent)

	json.Unmarshal(byteResult, &gapsAndScenarios)
	// fmt.Println(albums.Albums)
	c.IndentedJSON(http.StatusOK, gapsAndScenarios)

}
