package controller

import (
	"nicheanal.com/dal"
	"nicheanal.com/scrapers"
)

// KeywordResearchSave scrape keyword research page and store to db
func KeywordResearchSave(phrase string) error {
	res, err := scrapers.KeywordResearch(phrase)
	if err != nil {
		return err
	}
	return dal.KeywordResearchSave(res)
}
