package dal

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

// MarketArgument for database model
type MarketArgument struct {
	Phrase                        string
	Preset                        string
	Asin                          string
	At                            string
	Title                         string
	Description                   string
	Price                         float64
	ParentASIN                    string
	Brand                         string
	Bsr                           int64
	Category                      string
	CategoryID                    int64
	Categories                    pq.StringArray `gorm:"type:text[]"`
	Features                      pq.StringArray `gorm:"type:text[]"`
	ImageUrls                     pq.StringArray `gorm:"type:text[]"`
	NetProfit                     float64
	Revenue                       float64
	ReviewCount                   int64
	ReviewRate                    float64
	ReviewRating                  float64
	PriceAmazon                   float64
	PriceNew                      float64
	ProductGroup                  string
	ProfitMargin                  float64
	Sales                         int64
	SalesLastYear                 int64
	SalesToReviews                float64
	SellerCount                   int64
	UnitMargin                    float64
	StarRating                    float64
	BestSalesPeriod               string
	IsNameBrand                   bool
	PriceChangeLastNinetyDays     float64
	ReviewCountChangeMonthly      float64
	SalesChangeLastNinetyDays     float64
	SalesPattern                  string
	SalesYearOverYear             float64
	InitialCost                   float64
	InitialNetProfit              float64
	InitialOrganicSalesProjection int64
	InitialUnitsToOrder           int64
	OngoingOrganicSalesProjection int64
	OngoingUnitsToOrder           int64
	PromotionDuration             int64
	PromotionUnitsDaily           int64
	Fulfillment                   string
	PromotionUnitsTotal           int64
	IsVariationWithSharedBSR      bool
	OfferCountNew                 int64
	OfferCountUsed                int64
	PackageHeight                 float64
	PackageLength                 float64
	PackageQuantity               int64
	PackageWeight                 float64
	PackageWidth                  float64
}

// MarketIntelSave saves market intelligence scraping result
func MarketIntelSave(res []MarketArgument, phrase string) error {
	for _, v := range res {
		v.Phrase = phrase
		err := db.Where("asin=?", v.Asin).Delete(&MarketArgument{}).Error
		if err != nil {
			return err
		}
		err = db.Create(&v).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// GetMIDataByPhrase get all mi data with phrase
func GetMIDataByPhrase(phrase string) ([]MarketArgument, error) {
	mis := []MarketArgument{}
	err := db.Where("phrase=?", phrase).Find(&mis).Error
	return mis, err
}

// GetMIData get a market argument object from asin
func GetMIData(asin string) (*MarketArgument, error) {
	res := &MarketArgument{}
	err := db.Where("asin=?", asin).First(&res).Error
	return res, err
}

// GetAllMiData get all mi data
func GetAllMiData() ([]MarketArgument, error) {
	mis := []MarketArgument{}
	now := time.Now().Format("2006-01-02")
	now = now[0:7]
	fmt.Println("getmarket data---")
	err := db.Where("at LIKE ?", now+"%").Find(&mis).Error
	return mis, err
}
