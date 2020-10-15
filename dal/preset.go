package dal

import (
	"github.com/lib/pq"
)

// Preset is the result type of the productDiscovery service
type Preset struct {
	// Name of the preset
	Preset string
	// Used to set the next revenue range for the preset
	RevenueRangeMove     int64
	AveragePrice         string // minmax type
	AverageRevenue       string // minmax type
	AverageReviewCount   string // minmax type
	AverageSales         string // minmax type
	AmazonFulfillment    string // minmax type
	AverageNetProfit     string // minmax type
	FbaFulfillment       string // minmax type
	FbmFulfillment       string // minmax type
	PriceChange          string // minmax type
	ProfitMargin         string // minmax type
	ReviewChange         string // minmax type
	ReviewRate           string // minmax type
	ReviewRating         string // minmax type
	SearchVolumeEstimate string // minmax type
	SalesPattern         string
	RootCategories       pq.StringArray `gorm:"type:text[]"`
	SalesChange          string         // minmax type
	SalesToReviews       string         // minmax type
	SalesYearOver        string         // minmax type
	ShipSize             string
	BestSalesPeriod      string
	Marketplace          string
	Email                string
	ObjectID             string
	ContinuationToken    string
}

// LoadData set all data of presetdata db model
func (pdb *Preset) LoadData(p AppPreset) {
	pdb.Preset = p.Preset
	pdb.RevenueRangeMove = IntPointerToInt(p.RevenueRangeMove)
	pdb.AveragePrice = ToString(p.AveragePrice)
	pdb.AverageRevenue = ToString(p.AverageRevenue)
	pdb.AverageReviewCount = ToString(p.AverageReviewCount)
	pdb.AverageSales = ToString(p.AverageSales)
	pdb.AmazonFulfillment = ToString(p.AmazonFulfillment)
	pdb.AverageNetProfit = ToString(p.AverageNetProfit)
	pdb.FbaFulfillment = ToString(p.FbaFulfillment)
	pdb.FbmFulfillment = ToString(p.FbmFulfillment)
	pdb.PriceChange = ToString(p.PriceChange)
	pdb.ProfitMargin = ToString(p.ProfitMargin)
	pdb.ReviewChange = ToString(p.ReviewChange)
	pdb.ReviewRate = ToString(p.ReviewRate)
	pdb.ReviewRating = ToString(p.ReviewRating)
	pdb.SearchVolumeEstimate = ToString(p.SearchVolumeEstimate)
	pdb.SalesChange = ToString(p.SalesChange)
	pdb.SalesToReviews = ToString(p.SalesToReviews)
	pdb.SalesYearOver = ToString(p.SalesYearOver)
	pdb.SalesPattern = StringPointerToString(p.SalesPattern)
	pdb.RootCategories = p.RootCategories
	pdb.ShipSize = StringPointerToString(p.ShipSize)
	pdb.BestSalesPeriod = StringPointerToString(p.BestSalesPeriod)
	pdb.Marketplace = p.Marketplace
	pdb.Email = p.Email
	pdb.ObjectID = p.ObjectID
	pdb.ContinuationToken = StringPointerToString(p.ContinuationToken)
	return
}

// ExportData converts pdb to viralPresetmedia type
func (pdb *Preset) ExportData() *AppPreset {
	preset := &AppPreset{}
	preset.Preset = pdb.Preset
	preset.RevenueRangeMove = &pdb.RevenueRangeMove
	preset.AveragePrice = ToStructMinMax(pdb.AveragePrice)
	preset.AverageRevenue = ToStructMinMax(pdb.AverageRevenue)
	preset.AverageReviewCount = ToStructMinMax(pdb.AverageReviewCount)
	preset.AverageSales = ToStructMinMax(pdb.AverageSales)
	preset.AmazonFulfillment = ToStructMinMax(pdb.AmazonFulfillment)
	preset.AverageNetProfit = ToStructMinMax(pdb.AverageNetProfit)
	preset.FbaFulfillment = ToStructMinMax(pdb.FbaFulfillment)
	preset.FbmFulfillment = ToStructMinMax(pdb.FbmFulfillment)
	preset.PriceChange = ToStructMinMax(pdb.PriceChange)
	preset.ProfitMargin = ToStructMinMax(pdb.ProfitMargin)
	preset.ReviewChange = ToStructMinMax(pdb.ReviewChange)
	preset.ReviewRate = ToStructMinMax(pdb.ReviewRate)
	preset.ReviewRating = ToStructMinMax(pdb.ReviewRating)
	preset.SearchVolumeEstimate = ToStructMinMax(pdb.SearchVolumeEstimate)
	preset.SalesChange = ToStructMinMax(pdb.SalesChange)
	preset.SalesToReviews = ToStructMinMax(pdb.SalesToReviews)
	preset.SalesYearOver = ToStructMinMax(pdb.SalesYearOver)
	preset.SalesPattern = &pdb.SalesPattern
	preset.RootCategories = pdb.RootCategories
	preset.ShipSize = &pdb.ShipSize
	preset.BestSalesPeriod = &pdb.BestSalesPeriod
	preset.Marketplace = pdb.Marketplace
	preset.Email = pdb.Email
	preset.ObjectID = pdb.ObjectID
	preset.ContinuationToken = &pdb.ContinuationToken

	return preset
}

// CreatePreset save new preset data to db
func CreatePreset(p AppPreset) error {
	pdb := &Preset{}
	pdb.LoadData(p)
	err := db.Where("preset=?", p.Preset).Delete(&Preset{}).Error
	if err != nil {
		return err
	}
	return db.Create(&pdb).Error
}

// GetPreset get one Preset with preset string
func GetPreset(p string) (*AppPreset, error) {
	pdb := Preset{}
	err := db.Where("preset=?", p).First(&pdb).Error
	if err != nil {
		return nil, err
	}
	return pdb.ExportData(), nil
}

// ListPreset returns all Preset
func ListPreset() ([]*AppPreset, error) {
	presets := []*AppPreset{}
	pdbs := []Preset{}
	err := db.Find(&pdbs).Error
	if err != nil {
		return presets, err
	}
	for _, v := range pdbs {
		presets = append(presets, v.ExportData())
	}
	return presets, nil
}

// RemovePreset remove one preset data
func RemovePreset(p string) error {
	return db.Where("preset=?", p).Delete(&Preset{}).Error
}

// RemoveAllPreset removes all preset data from db
func RemoveAllPreset() error {
	return db.Delete(&Preset{}).Error
}
