package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptData struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        string    `json:"total"`
}

type Receipt struct {
	ID   string      `json:"id"`
	Data ReceiptData `json:"data"`
}

var receipts []Receipt

func processReceipts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var receiptData ReceiptData
	err := json.NewDecoder(r.Body).Decode(&receiptData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	receiptID := strconv.Itoa(rand.Intn(10000000))
	receipt := Receipt{
		ID:   receiptID,
		Data: receiptData,
	}
	receipts = append(receipts, receipt)

	response := map[string]string{"id": receiptID}
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(receiptData ReceiptData) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	retailerName := strings.ReplaceAll(receiptData.Retailer, "&", "")
	points += len(strings.ReplaceAll(retailerName, " ", ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, err := strconv.ParseFloat(receiptData.Total, 64)
	if err == nil && totalFloat == float64(int(totalFloat)) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	totalFloat *= 100 // Convert to cents for accurate comparison
	if math.Mod(totalFloat, 25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += len(receiptData.Items) / 2 * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range receiptData.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(priceFloat * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", receiptData.PurchaseDate+" "+receiptData.PurchaseTime)
	if err == nil && purchaseDateTime.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00 pm and before 4:00 pm.
	purchaseTime, err := time.Parse("15:04", receiptData.PurchaseTime)
	if err == nil && purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}


func getReceiptPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	receiptID := params["id"]

	var points int
	for _, receipt := range receipts {
		if receipt.ID == receiptID {
			points = calculatePoints(receipt.Data)
			break
		}
	}

	response := map[string]int{"points": points}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/receipts/process", processReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getReceiptPoints).Methods("GET")

	fmt.Println("Starting the server at port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
