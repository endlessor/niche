package dal

// ProductDiscovery for product discovery type
type ProductDiscovery struct {
	Phrase                           string `form:"phrase,omitempty" json:"phrase,omitempty" storm:"id"`
	Preset                           string
	PresetRevenue                    int64
	AverageNetProfit                 float64
	AveragePrice                     float64
	AverageProfitMargin              float64
	AverageRevenue                   float64
	AverageReviewCount               float64
	AverageReviewRate                float64
	AverageReviewRating              float64
	AverageSales                     float64
	AverageSalesToReviews            float64
	TotalReviewCount                 int64
	AveragePriceChangeLastNinetyDays float64
	AverageReviewCountChangeMonthly  float64
	AverageSalesChangeLastNinetyDays float64
	BestSalesPeriod                  string
	SalesPattern                     string
	SalesYearOverYear                float64
	StarRating                       float64
	VolumeEstimate                   float64
	VolumeEstimatedAt                string
	Category                         string
	DominantShipmentSizeTier         string
	Marketplace                      string
	IsPopular                        bool
}

// SaveProductDiscovery save product discovery data
func SaveProductDiscovery(p *ProductDiscovery) error {
	err := db.Where("phrase=? AND preset=?", p.Phrase, p.Preset).Delete(&ProductDiscovery{}).Error
	if err != nil {
		return err
	}
	return db.Create(&p).Error
}

// ListDiscoveryData returns all phraseData from product discovery with preset
func ListDiscoveryData(preset string) ([]ProductDiscovery, error) {
	pds := []ProductDiscovery{}
	err := db.Where("preset=?", preset).Find(&pds).Error
	return pds, err
}

// ListAllProductDiscovery returns all product discovery
func ListAllProductDiscovery() ([]ProductDiscovery, error) {
	pds := []ProductDiscovery{}
	err := db.Find(&pds).Error
	return pds, err
}
