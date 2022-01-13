package utils

import (
	"encoding/json"
	"io"
)

type Price struct {
	Kind  string
	Delta int
	Price float64
}

type Station struct {
	Name   string
	Prices map[string]Price
}

func NewStation(name string, prices map[string]Price) Station {
	return Station{name, prices}
}

func Parse(body io.Reader) ([]Station, error) {

	var jsonData map[string]interface{}

	err := json.NewDecoder(body).Decode(&jsonData)

	if err != nil {
		return nil, err
	}

	stations := jsonData["data"].(map[string]interface{})["stations"].([]interface{})

	result := make([]Station, 0, len(stations))

	for _, station := range stations {

		name := station.(map[string]interface{})["name"].(string)

		prices := station.(map[string]interface{})["prices"].(map[string]interface{})

		p := map[string]Price{}

		_, ok := prices["diesel"]
		if ok {
			price := Price{Kind: "diesel", Price: prices["diesel"].(float64), Delta: int(prices["diesel_ch"].(float64))}
			p["diesel"] = price
		}

		_, ok = prices["e10"]
		if ok {
			price := Price{Kind: "e10", Price: prices["e10"].(float64), Delta: int(prices["e10_ch"].(float64))}
			p["e10"] = price
		}

		_, ok = prices["e5"]
		if ok {
			price := Price{Kind: "e5", Price: prices["e5"].(float64), Delta: int(prices["e5_ch"].(float64))}
			p["e5"] = price
		}

		result = append(result, NewStation(name, p))
	}

	return result, nil
}
