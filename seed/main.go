package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"nicheanal.com/config"
	"nicheanal.com/controller"
	"nicheanal.com/dal"
)

func main() {
	config.LoadConfig("../config/config.json")
	dal.LoadDB()

	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "nicheanalysis", log.LstdFlags)

	// logger.Println("Starting db seeding...")
	// seedPreset(logger)
	// logger.Println("DB seeding completed")

	// logger.Println("Starting scraping product discovery...")
	// scrapProductDiscovery(logger)
	// logger.Println("Complete scraping product discovery")

	logger.Println("Starting scraping amazon keywords...")
	scrapKeywords(logger)
	logger.Println("Complete scraping amazon keywords")

	// if err != nil {
	// 	params := &logger.LogParams{}
	// 	params.Add("reason:", err)
	// 	params.Add("requestBody:", requestBody)
	// 	logger.ErrorP("unable to parse requestBody:", params)

	// 	return
	// }

}

func seedPreset(logger *log.Logger) {
	ps := []dal.AppPreset{}
	pj, err := os.Open("preset.json")
	if err != nil {
		logger.Fatal("Faile to load seed data", err)
	}
	defer pj.Close()

	b, _ := ioutil.ReadAll(pj)
	err = json.Unmarshal(b, &ps)
	if err != nil {
		logger.Fatal("Failed to parse preset data", err)
	}

	err = dal.RemoveAllPreset()
	if err != nil {
		logger.Fatal("Failed to remove all preset data", err)
	}
	for _, p := range ps {
		err = controller.CreatePreset(p)
		if err != nil {
			logger.Fatal("Failed to add preset data", err)
		}
	}
}

func scrapProductDiscovery(logger *log.Logger) {
	aps, err := controller.ListPresets()
	if err != nil {
		logger.Println(err)
		return
	}

	for k := 0; k < len(aps); k++ {
		err = controller.PdScrape([]int{}, float64(200000), float64(20000000), aps[k])
		if err != nil {
			logger.Println(err)
			return
		}

		for i := 1; i < 8; i++ {
			err = controller.PdScrape([]int{}, float64(200000-25000*i), float64(200000-25000*(i-1)), aps[k])
			if err != nil {
				logger.Println(err)
				continue
			}
		}
	}
}

func scrapKeywords(logger *log.Logger) {
	apds, err := controller.ListAllProductDiscovery()
	if err != nil {
		logger.Println(err)
		return
	}
	for _, pd := range apds {
		phrase := pd.Phrase
		err = controller.KeywordResearchSave(phrase)
		if err != nil {
			logger.Println(err)
			continue
		}
	}
	// for _, pm := range pl {
	// 	phl, err := dal.ListPhraseDataByPhrase(pm.Preset)
	// 	if err != nil {
	// 		logger.Println(err)
	// 		return
	// 	}
	// 	for i, phm := range phl {
	// 		phrase := phm.Phrase
	// 		if !dal.CheckPhraseExist(phrase) {
	// 			if phrase != "" {
	// 				continue
	// 			}
	// 			_, err = controller.AmazonScrape(phrase)
	// 			if err != nil {
	// 				logger.Println(err)
	// 				continue
	// 			}
	// 			logger.Println("searching items on", i)
	// 			_, err = controller.MIScrape(phrase, logger)
	// 			if err != nil {
	// 				logger.Println(err)
	// 				continue
	// 			}
	// 			_, err = controller.KeywordResearchSave(phrase)
	// 			if err != nil {
	// 				logger.Println(err)
	// 				continue
	// 			}
	// 			err = dal.SavePhraseModel(phrase)
	// 			if err != nil {
	// 				logger.Println(err)
	// 				continue
	// 			}
	// 		}
	// 		err = controller.SpPhraseSave(pm.Preset, phrase)
	// 		if err != nil {
	// 			logger.Println(err)
	// 			continue
	// 		}
	// 	}
	// }
}

// // CleanPhraseModel cleans phrase modal table
// func CleanPhraseModel(logger *log.Logger) {
// 	pds, err := dal.GetAllMiData()
// 	if err != nil {
// 		logger.Panic(err)
// 	}
// 	for _, v := range pds {
// 		phrase := v.Phrase
// 		if dal.CheckPhraseExist(phrase) {
// 			continue
// 		}
// 		fmt.Println("save--", phrase)
// 		err = dal.SavePhraseModel(phrase)
// 		if err != nil {
// 			logger.Println(err)
// 			continue
// 		}
// 	}
// }
