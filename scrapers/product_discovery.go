package scrapers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"nicheanal.com/config"

	"nicheanal.com/dal"
)

var _ = log.Print
var _ = bytes.HasPrefix
var _ = http.Get
var _ = ioutil.WriteFile

// ViralDiscoverydata for response type of product discovery
type ViralDiscoverydata struct {
	Data              []*dal.ProductDiscovery
	ContinuationToken *string
}

// PdPayload for requesting product discovery info
type PdPayload struct {
	AverageReviewCount   *dal.AppMinMax `json:"averageReviewCount"`
	AverageRevenue       *dal.AppMinMax `json:"averageRevenue"`
	AverageSales         *dal.AppMinMax `json:"averageSales"`
	ReviewRating         *dal.AppMinMax `json:"reviewRating"`
	SearchVolumeEstimate *dal.AppMinMax `json:"searchVolumeEstimate"`
	ContinuationToken    *string        `json:"continuationToken"`
	RootCategories       []string       `json:"rootCategories"`
	MarketPlace          string         `json:"marketplace"`
	Email                string         `json:"email"`
	ObjectID             string         `json:"objectID"`
	SalesPattern         string         `json:"salesPattern"`
	ShipSize             string         `json:"shipSize"`
}

// ScrapeProductDiscovery scrapes product discovery data with preset
func ScrapeProductDiscovery(pst *dal.AppPreset, delay []int) ([]ViralDiscoverydata, error) {
	collection := []ViralDiscoverydata{}
	for {
		pd := &PdPayload{
			AverageReviewCount:   pst.AverageReviewCount,
			AverageRevenue:       pst.AverageRevenue,
			AverageSales:         pst.AverageSales,
			ReviewRating:         pst.ReviewRating,
			SearchVolumeEstimate: pst.SearchVolumeEstimate,
			RootCategories:       pst.RootCategories,
			ContinuationToken:    pst.ContinuationToken,
			MarketPlace:          "US",
			Email:                config.Cfg.ViralLaunchEmail,
			ObjectID:             config.Cfg.ViralLaunchID,
			ShipSize:             "",
			SalesPattern:         "",
		}
		fmt.Println(pst.Preset, *pst.AverageRevenue.Lower, "-- starting")

		payload, _ := json.Marshal(pd)
		url := "https://viral-launch.com/sellers/assets/php/keyword-filter.php"
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Print("Error scraping product discovery: ", err)
			return collection, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		current := ViralDiscoverydata{}
		err = json.Unmarshal(body, &current)
		if err != nil {
			return collection, err
		}

		for i := range current.Data {
			current.Data[i].PresetRevenue = int64(*pst.AverageRevenue.Lower)
			current.Data[i].Preset = pst.Preset
		}

		collection = append(collection, current)
		pst.ContinuationToken = current.ContinuationToken
		if current.ContinuationToken == nil {
			break
		}
	}

	return collection, nil
}
