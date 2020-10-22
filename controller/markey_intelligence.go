package controller

import (
	"log"

	"nicheanal.com/dal"
	"nicheanal.com/scrapers"
)

// MIScrape scraps market intelligence page and store to db
func MIScrape(phrase string, logger *log.Logger) error {
	res, err := scrapers.MarketIntelligence(phrase, logger)
	if err != nil {
		return err
	}
	err = dal.MarketIntelSave(res, phrase)
	return err
}
