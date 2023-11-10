package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Receipt struct to represent the JSON receipt
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var processedReceipts = make(map[string]int)

func main() {
	r := gin.Default()

	r.POST("/receipts/process", ProcessReceipt)
	r.GET("/receipts/:id/points", GetPoints)

	log.Fatal(r.Run(":8080"))
}

func ProcessReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	points := CalculatePoints(receipt)
	id := GenerateReceiptID()

	processedReceipts[id] = points

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func GetPoints(c *gin.Context) {
	id := c.Param("id")
	points := LookupPointsByID(id)
	c.JSON(http.StatusOK, gin.H{"points": points})
}

func CalculatePoints(receipt Receipt) int {
	// Implement the point calculation logic
	points := 0

	//add points for retailer alphanumeric characters
	regex, _ := regexp.Compile("[^a-zA-Z0-9]+")
	points = points + len(regex.ReplaceAllString(receipt.Retailer, ""))

	//add 50 points if total is round dollar with no cents
	totalCost, _ := strconv.ParseFloat(receipt.Total, 64)
	if math.Mod(totalCost, 1.00) == 0.0 {
		points = points + 50
	}
	//add 25 points if total is multiple of .25
	if math.Mod(totalCost, 0.25) == 0.0 {
		points = points + 25
	}

	//add 5 points for every two items on receipt
	points = points + 5*(len(receipt.Items)/2)

	//if trimmed(remove leading and trailing white space) length is a multiple of 3
	//then add the price multiplied by .2 and rounded up to the nearest integer
	for _, item := range receipt.Items {
		points = points + CalculateItemPoints(item)
	}

	//add 6 points if the purchase date is odd
	dateString := receipt.PurchaseDate[strings.LastIndex(receipt.PurchaseDate, "-")+1 : len(receipt.PurchaseDate)]
	date, _ := strconv.ParseInt(dateString, 10, 64)
	if date%2 == 1 {
		points = points + 6
	}

	//add 10 points if time of purchase is after 2:00pm and before 4:00pm
	hourString := receipt.PurchaseTime[0:strings.Index(receipt.PurchaseTime, ":")]
	hour, _ := strconv.ParseInt(hourString, 10, 64)
	if hour >= 14 && hour < 16 {
		points = points + 10
	}

	return points
}

func CalculateItemPoints(item Item) int {
	//if trimmed(remove leading and trailing white space) length is a multiple of 3
	//then add the price multiplied by .2 and rounded up to the nearest integer
	points := 0
	trimmedDescription := strings.TrimSpace(item.ShortDescription)
	if len(trimmedDescription)%3 == 0 {
		itemTotal, _ := strconv.ParseFloat(item.Price, 64)
		pointsAdded := itemTotal * 0.2
		points = points + int(pointsAdded)
		if math.Mod(pointsAdded, 1.0) != 0 {
			//round up if pointsAdded isn't a whole number
			points = points + 1
		}
	}

	return points
}

func GenerateReceiptID() string {
	// Generate a UUID for the receipt
	return uuid.New().String()
}

func LookupPointsByID(id string) int {
	// Look up the points for a receipt by ID
	return processedReceipts[id]
}
