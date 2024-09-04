package pricecontroller

import (
	"encoding/json"
	pricemodels "evaluation/my-go-project/models"
	priceservice "evaluation/my-go-project/service"
	"io"
	"net/http"
)

// Get handles the GET request for retrieving price data
func Get(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	if pair == "" {
		http.Error(w, "pair parameter is required", http.StatusBadRequest)
		return
	}
	duration := r.URL.Query().Get("duration")
	if duration == "" {
		http.Error(w, "duration parameter is required", http.StatusBadRequest)
		return
	}
	data, err := priceservice.GetPriceData(pair, duration)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"pair": pair,
		"data": data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

// Insert handles the POST request for inserting price data
func Insert(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var requestBody pricemodels.BitkubTradingViewRequestBody
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	priceData := priceservice.GetBitkubTradingView(pricemodels.BitkubTradingViewRequestBody{
		Symbol:     requestBody.Symbol,
		Resolution: requestBody.Resolution,
		FromTime:   requestBody.FromTime,
	})
	if err := priceservice.InsertPriceData(priceData); err != nil {
		http.Error(w, "Failed to insert data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data processed successfully"))
}
