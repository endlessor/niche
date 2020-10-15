package dal

// AppPreset for application preset model
type AppPreset struct {
	// Name of the preset
	Preset string `form:"preset,omitempty" json:"preset,omitempty" storm:"id"`
	// Used to set the next revenue range for the preset
	RevenueRangeMove   *int64
	AveragePrice       *AppMinMax
	AverageRevenue     *AppMinMax
	AverageReviewCount *AppMinMax
	AverageSales       *AppMinMax
	AmazonFulfillment  *AppMinMax
	AverageNetProfit   *AppMinMax
	FbaFulfillment     *AppMinMax
	FbmFulfillment     *AppMinMax
	PriceChange        *AppMinMax
	ProfitMargin       *AppMinMax
	ReviewChange       *AppMinMax
	ReviewRate         *AppMinMax
	ReviewRating       *AppMinMax
	SearchVolumeBroad  *AppMinMax
	SearchVolumeExact  *AppMinMax
	SalesPattern       *string
	RootCategories     []string
	SalesChange        *AppMinMax
	SalesToReviews     *AppMinMax
	SalesYearOver      *AppMinMax
	ShipSize           *string
	BestSalesPeriod    *string
	Marketplace        string
	Email              string
	ObjectID           string
	ContinuationToken  *string
}

// AppMinMax for min-max type
type AppMinMax struct {
	Lower *float64
	Upper *float64
}
