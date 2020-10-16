package dal

import (
	"github.com/lib/pq"
)

// KeywordData is database model for keyword research result
type KeywordData struct {
	AveragePrice       float64
	AverageReviewCount float64
	AverageSales       float64
	DominantCategories pq.StringArray `gorm:"type:text[]"`
	Marketplace        string
	OpportunityScore   float64
	Phrase             string
	VolumeEstimate     float64
	VolumeEstimatedAt  string
	Score              float64
	OriginPhrase       string
}

// KeywordResearchSave create new keywork reseach scrape result searched
func KeywordResearchSave(kds []KeywordData) error {
	for _, v := range kds {
		err := db.Where("phrase=?", v.Phrase).Delete(&KeywordData{}).Error
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

// KeywordResearchGet get all keyword data from phrase
func KeywordResearchGet(phrase string) ([]KeywordData, error) {
	kds := []KeywordData{}
	err := db.Where("phrase=?", phrase).Find(&kds).Error
	return kds, err
}
