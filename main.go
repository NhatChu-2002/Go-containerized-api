package main

import (
	"net/http"

	// "database/sql"
	// "encoding/json"
	"fmt"
	// _ "github.com/go-sql-driver/mysql"
	"log"
	// "github.com/gin-gonic/gin"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func main() {

	// r := gin.Default()

	// r.GET("/api/hello", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Hello world!",
	// 	})
	// })
	// r.Run()

	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
