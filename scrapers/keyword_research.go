package scrapers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"nicheanal.com/config"
	"nicheanal.com/dal"
)

var kwInqURL = "https://viral-launch.com/sellers/assets/php/keyword-research-inquiry.php?market=US"
var kwmsURL = "https://viral-launch.com/sellers/assets/php/keyword-research-measurement.php?market=US"

// var spKwURL = "https://viral-launch.com/sellers/assets/php/keyword-filter.php"

// PhraseList for phrase list type
type PhraseList struct {
	At                string `json:"at"`
	RefreshPriority   int    `json:"refreshPriority"`
	MainKeywordsCount int    `json:"mainKeywordsCount"`
	UniqueWordsCount  int    `json:"uniqueWordsCount"`
	MainKeywords      []struct {
		Phrase string  `json:"phrase"`
		Score  float64 `json:"score"`
	} `json:"mainKeywords"`
}

// PhrasePayload for phrase payload type
type PhrasePayload struct {
	Phrases []string `json:"phrases"`
}

// KeywordPayload for keyword research payload type
type KeywordPayload struct {
	Phrase         string        `json:"phrase"`
	Marketplace    string        `json:"marketplace"`
	Email          string        `json:"email"`
	ObjectID       string        `json:"objectID"`
	SalesToReviews dal.AppMinMax `json:"salesToReviews"`
}

// KeywordResearch scrapes the viral launch app for keyword research
func KeywordResearch(query string) ([]dal.KeywordData, error) {
	data := []dal.KeywordData{}
	URL := fmt.Sprintf("%s&phrase=%s&email=%s&id=%s",
		kwInqURL, url.QueryEscape(query), url.QueryEscape(config.Cfg.ViralLaunchEmail), config.Cfg.ViralLaunchID)
	resPhrases, err := http.Get(URL)
	if err != nil {
		return data, err
	}

	body, _ := ioutil.ReadAll(resPhrases.Body)
	defer resPhrases.Body.Close()

	phraseList := PhraseList{}
	err = json.Unmarshal(body, &phraseList)
	if err != nil {
		return data, err
	}

	phrasePayload := PhrasePayload{}
	if phraseList.RefreshPriority == 0 {
		phrasePayload.Phrases = append(phrasePayload.Phrases, query)
		for i, v := range phraseList.MainKeywords {
			if i == 50 {
				break
			}
			phrasePayload.Phrases = append(phrasePayload.Phrases, v.Phrase)
		}

		payload, _ := json.Marshal(phrasePayload)
		URL = fmt.Sprintf("%s&email=%s&id=%s",
			kwmsURL, url.QueryEscape(config.Cfg.ViralLaunchEmail), config.Cfg.ViralLaunchID)
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
		if err != nil {
			return data, err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return data, err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		if len(data) > 0 {
			// add relevancy score
			for i := range data {
				data[i].Score = phraseList.MainKeywords[i].Score
				data[i].OriginPhrase = query
			}
		}
	} else {
		time.Sleep(5 * time.Second)
		KeywordResearch(query)
	}
	return data, nil
}

// // SpKeywordResearch scraps the viral launch app for keyword data
// func SpKeywordResearch(query string) (spkeywordresearchsvc.ViralSpkeyworddataCollection, error) {
// 	data := RespKwData{}
// 	var i float64 = 1
// 	p := KeywordPayload{
// 		Phrase:      query,
// 		Marketplace: "US",
// 		Email:       os.Getenv("VIRAL_LAUNCH_EMAIL"),
// 		ObjectID:    os.Getenv("VIRAL_LAUNCH_ID"),
// 		SalesToReviews: productdiscoverysvc.ApplicationViralMinMax{
// 			Lower: &i,
// 		},
// 	}
// 	payload, _ := json.Marshal(p)

// 	req, err := http.NewRequest("POST", spKwURL, bytes.NewBuffer(payload))
// 	if err != nil {
// 		return data.Data, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return data.Data, err
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	err = json.Unmarshal(body, &data)
// 	return data.Data, err
// }

//	payload, err := json.Marshal(preset) email=subs%40denevehome.com&id=kN1G4gsZEh
//    url := "https://viral-launch.com/sellers/assets/php/keyword-filter.php"
//    var jsonStr = payload
//    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
//    req.Header.Set("Content-Type", "application/json")
//
//    client := &http.Client{}
//    resp, err := client.Do(req)
//    if err != nil {
//        panic(err)
//    }
//    defer resp.Body.Close()
//
//    body, _ := ioutil.ReadAll(resp.Body)
//    data = body
//
//
//    current := app.DiscoveryData{}
//    err = json.Unmarshal(data, &current)
//    if err != nil {
//    	fmt.Println(err)
//	}
