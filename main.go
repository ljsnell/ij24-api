package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type Albums struct {
	Albums []Album `json:"albums"`
}

type Album struct {
	ID     string  `json:"ID"`
	Title  string  `json:"Title"`
	Artist string  `json:"Artist"`
	Price  float32 `json:"Price"`
}

var albums *Albums

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/score", getScore)
	// router.GET("/albums/:id", getAlbumByID)

	router.Run(":8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	fileContent, err := os.Open("albums.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("The File is opened successfully...")

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	json.Unmarshal(byteResult, &albums)
	fmt.Println(albums.Albums)
	c.IndentedJSON(http.StatusOK, albums.Albums)
}

func getScore(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"score": 84})
}

// // getAlbumByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop over the list of albums, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }
