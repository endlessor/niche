package scrapers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"nicheanal.com/config"
	"nicheanal.com/dal"
)

// Data for byte record
var Data []byte
var searchAsinsURL = "https://viral-launch.com/sellers/assets/php/market-intelligence/search-asins.php"
var searchListingURL = "https://viral-launch.com/sellers/assets/php/market-intelligence/search-listing.php"

// AsinPayload for requesting asin info
type AsinPayload struct {
	MarketPlace string `json:"marketplace"`
	Phrase      string `json:"phrase"`
	By          string `json:"by"`
	Source      string `json:"source"`
	ObjectID    string `json:"objectId"`
}

// AsinResp for asin response type
type AsinResp struct {
	PrimaryAsins []string `json:"primaryAsins"`
}

// MarketIntelligence scrapes market intelligence info from viral-launch-market site
func MarketIntelligence(phrase string, logger *log.Logger) ([]dal.MarketArgument, error) {
	p := &AsinPayload{
		MarketPlace: "US",
		Phrase:      phrase,
		By:          config.Cfg.ViralLaunchEmail,
		Source:      "viral-launch.com",
		ObjectID:    config.Cfg.ViralLaunchID,
	}
	asins, err := p.getAsins()
	if err != nil {
		return nil, err
	}
	logger.Println("asin---", len(asins))
	if len(asins) == 0 {
		logger.Panic("Exceed request limit for get asin")
		return nil, errors.New("Exceed request limit for get asin")
	}
	return p.getMarketData(asins, logger)
}

func (p *AsinPayload) getAsins() ([]string, error) {
	payload, _ := json.Marshal(p)
	req, err := http.NewRequest("POST", searchAsinsURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	asd := &AsinResp{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &asd)
	if err != nil {
		return nil, err
	}
	return asd.PrimaryAsins, nil
}

func (p *AsinPayload) getMarketData(asins []string, logger *log.Logger) ([]dal.MarketArgument, error) {
	vmd := []dal.MarketArgument{}
	payload, _ := json.Marshal(p)
	client := &http.Client{}

	for i := 0; i < len(asins); i++ {
		arg := dal.MarketArgument{}
		req, err := http.NewRequest("POST", fmt.Sprintf("%s?asin=%s", searchListingURL, asins[i]), bytes.NewBuffer(payload))
		if err != nil {
			logger.Println(err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &arg)
		if err != nil {
			logger.Println(err)
			continue
		}
		vmd = append(vmd, arg)
	}
	// time.Sleep(time.Duration(10) * time.Second)
	return vmd, nil
}
