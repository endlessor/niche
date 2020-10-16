package controller

import (
	"log"

	"nicheanal.com/dal"
	"nicheanal.com/scrapers"
)

// PdScrape scrapes and stores data from Product Discovery
func PdScrape(delay []int, revMin float64, revMax float64, preset *dal.AppPreset) error {
	if preset.AverageRevenue == nil {
		preset.AverageRevenue = &dal.AppMinMax{}
	}
	preset.AverageRevenue.Lower = &revMin
	preset.AverageRevenue.Upper = &revMax
	num := 0
	pds, err := scrapers.ScrapeProductDiscovery(preset, delay)
	for _, v := range pds {
		for _, c := range v.Data {
			err = dal.SaveProductDiscovery(c)
			if err != nil {
				return err
			}
			num++
		}
	}
	log.Printf("Scraped %d items", num)
	return nil
}

// ListAllProductDiscovery get all product discovery data
func ListAllProductDiscovery() ([]dal.ProductDiscovery, error) {
	return dal.ListAllProductDiscovery()
}

// // PdShowPreset shows a single preset
// func PdShowPreset(preset string) (*productdiscoverysvc.ViralPresetmedia, error) {
// 	return dal.GetPreset(preset)
// }

// // PdCreatePreset creates a new preset
// func PdCreatePreset(p *productdiscoverysvc.ApplicationViralPresetPayload) error {
// 	return dal.CreatePreset(p)
// }

// // PdListPresets lists saved presets
// func PdListPresets() (productdiscoverysvc.ViralPresetmediaCollection, error) {
// 	return dal.ListPreset()
// }

// // PdRemovePreset deletes a preset
// func PdRemovePreset(preset string) error {
// 	return dal.RemovePreset(preset)
// }

// // PdUpdatePreset updates a preset
// func PdUpdatePreset(p *productdiscoverysvc.ApplicationViralPresetPayload) error {
// 	return dal.CreatePreset(p)
// }

// // PdChangePName changes preset name
// func PdChangePName(old, new string) error {
// 	return dal.ChangePresetName(old, new)
// }
