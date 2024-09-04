package priceservice

import (
	"context"
	"encoding/json"
	pricemodels "evaluation/my-go-project/models"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/influxdata/influxdb-client-go/v2"
)

const (
	myToken  = "f3c230012ee93b1ee7f64f137a0d7e57233d20d477d3cf0b3ea1eafefe49a0e3"
	myOrg    = "test123"
	myBucket = "test123"
)

var client = influxdb2.NewClient("http://localhost:8086", myToken)
var queryAPI = client.QueryAPI(myOrg)
var writeAPI = client.WriteAPIBlocking(myOrg, myBucket)

// FetchPriceData retrieves price data for a given pair and duration
func GetPriceData(pair, duration string) ([]map[string]interface{}, error) {
	query := `from(bucket:"` + myBucket + `") |> range(start: -` + duration + `) |> filter(fn: (r) => r._measurement == "` + pair + `")`
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer result.Close()

	var data []map[string]interface{}
	for result.Next() {
		record := result.Record()
		data = append(data, map[string]interface{}{
			"time":  record.Time().Format(time.RFC3339),
			"field": record.Field(),
			"value": record.Value(),
		})
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("error during result iteration: %w", result.Err())
	}
	return data, nil
}

// InsertPriceData inserts price data into InfluxDB
func InsertPriceData(priceData pricemodels.PriceDataDB) error {
	for _, data := range priceData.Data {
		timestamp := time.Unix(data.Timestamp, 0).UTC()
		p := influxdb2.NewPointWithMeasurement(priceData.Pair).
			AddField("price", data.Price).
			SetTime(timestamp)

		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			return fmt.Errorf("failed to write point: %w", err)
		}
	}
	return nil
}


// GetBitkubTradingView fetches trading data from Bitkub API
func GetBitkubTradingView(requestBody pricemodels.BitkubTradingViewRequestBody) pricemodels.PriceDataDB {
	host := "https://api.bitkub.com"
	path := "/tradingview/history"
	url := fmt.Sprintf("%s%s", host, path)
	toTime := time.Now().Unix()

	params := fmt.Sprintf("symbol=%s&resolution=%s&from=%d&to=%d", requestBody.Symbol, requestBody.Resolution, requestBody.FromTime, toTime)

	resp, err := http.Get(fmt.Sprintf("%s?%s", url, params))
	if err != nil {
		return pricemodels.PriceDataDB{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return pricemodels.PriceDataDB{}
	}

	var apiResponse pricemodels.BitkubTradingViewResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return pricemodels.PriceDataDB{}
	}

	if apiResponse.S != "ok" {
		return pricemodels.PriceDataDB{
			Pair: requestBody.Symbol,
			Data: []pricemodels.DataDB{},
		}
	}

	data := make([]pricemodels.DataDB, len(apiResponse.T))
	for i := range apiResponse.T {
		data[i] = pricemodels.DataDB{
			Timestamp: apiResponse.T[i],
			Price:     apiResponse.C[i],
		}
	}

	result := pricemodels.PriceDataDB{
		Pair: requestBody.Symbol,
		Data: data,
	}

	return result
}
