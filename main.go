package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/magiconair/properties"
)

type Asset struct {
	Name         string  `json:"name"`
	Value        float32 `json:"value"`
	CoverageRisk float32 `json:"coverageRisk"`
	AssetClass   string  `json:"assetClass"`
}
type AssetsRequest struct {
	Luck   float32 `json:"luck"`
	Assets []Asset `json:"assets"`
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
}

type Scenario struct {
	ScenarioName string  `json:"scenario_name"`
	Cost         float32 `json:"cost"`
	AssetName    string  `json:"gap_name"`
	GapCost      float32 `json:"gap_cost"`
	AssetClass   string  `json:"assetClass"`
}

type GapsAndScenarios struct {
	Scenarios []Scenario `json:"scenarios"`
}

type AllDropDowns struct {
	Auto Auto `json:"auto"`
	Home Home `json:"home"`
}

var ddData *AllDropDowns
var gapsAndScenarios *GapsAndScenarios

const apiURL = "https://api.openai.com/v1/completions"

func main() {
	// init from a file
	p := properties.MustLoadFile("config.properties", properties.UTF8)

	router := gin.Default()
	router.GET("/alldropdowns", getDropdownData)
	router.GET("/score", getScore)
	router.GET("/add/:assetclass", getAssets)
	router.POST("/calculateRisk", getCalculateRisk)
	// router.GET("/albums/:id", getAlbumByID)

	if p.MustGetString("host") == "local" {
		router.Run("localhost:8080")
	} else {
		router.Run(":8080")
	}
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

	luck := assetsReq.Luck
	assets := assetsReq.Assets

	if len(assets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No assets provided"})
		return
	}

	// luck (single int), assets [{value: xxx, coverageRisk:1-3}]
	// Output -> coverageRisk, 3 is greatest risk
	// Scenario
	// Total payment per scenarios
	fmt.Println("Luck! " + fmt.Sprint(luck))
	maxGaps := 3
	// Can't easily use case statement w/Numbers
	if luck > 75 {
		maxGaps = 1
	} else if luck > 25 {
		maxGaps = 2
	} else {
		maxGaps = 3
	}

	fmt.Println(maxGaps)
	fileContent, err := os.Open("jsons/scenarios_and_gaps.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer fileContent.Close()

	byteResult, _ := io.ReadAll(fileContent)

	json.Unmarshal(byteResult, &gapsAndScenarios)
	sum := 0
	var vehicleArray []Scenario
	for _, scenario := range gapsAndScenarios.Scenarios {
		sum += int(scenario.GapCost)
		if scenario.AssetClass == "vehicle" {
			vehicleArray = append(vehicleArray, scenario)
		}
	}
	fmt.Println(sum)
	fmt.Println(vehicleArray)

	c.IndentedJSON(http.StatusOK, gapsAndScenarios)
}

// OPEN AI
func openAi() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := resty.New()

	response, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(`{
            "model": "text-davinci-003",
            "prompt": "Write a Go function to add two numbers",
            "max_tokens": 100
        }`).
		Post("https://api.openai.com/v1/completions")

	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// Print the response
	fmt.Println(response)
}
