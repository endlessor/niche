package dal

import (
	"encoding/json"
	"math"
)

// BoolPointerToBool returns origin boolean type from reference
func BoolPointerToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// StringPointerToString returns origin string type from reference
func StringPointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// IntPointerToInt returns origin int64 type from reference
func IntPointerToInt(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// FloatPointerToFloat returns origin float64 type from reference
func FloatPointerToFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return math.Round(*f*100) / 100
}

// ToString returns string of marshaled object
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// ContainsStringArray check it straing array contains a specific string
func ContainsStringArray(sa []string, s string) bool {
	for _, ss := range sa {
		if ss == s {
			return true
		}
	}
	return false
}

// ToStructMinMax returns ViralMinMax struct from string
func ToStructMinMax(s string) *AppMinMax {
	m := &AppMinMax{}
	json.Unmarshal([]byte(s), &m)
	return m
}

// ToStructPhraseData returns ApplicationViralPhraseData struct from string
// func ToStructPhraseData(s string) *phrasesvc.ApplicationViralPhraseData {
// 	m := &phrasesvc.ApplicationViralPhraseData{}
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// // ToStructMarketData returns ViralMarketdata struct from string
// func ToStructMarketData(s string) *phrasesvc.ViralMarketdata {
// 	m := &phrasesvc.ViralMarketdata{}
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// // ToStructKRData returns ViralKeyworddata struct from string
// func ToStructKRData(s string) []*phrasesvc.ViralKeyworddata {
// 	var m []*phrasesvc.ViralKeyworddata
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// // ToBestSellerRankings returns array of ApplicationViralBSR from string
// func ToBestSellerRankings(s string) []*marketintelligencesvc.ApplicationViralBSR {
// 	var m []*marketintelligencesvc.ApplicationViralBSR
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// // ToFees returns Fees from string
// func ToFees(s string) *marketintelligencesvc.ApplicationViralFees {
// 	var m *marketintelligencesvc.ApplicationViralFees
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// // ToOffers returns offers from string
// func ToOffers(s string) []*marketintelligencesvc.ApplicationViralOffer {
// 	var m []*marketintelligencesvc.ApplicationViralOffer
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// // ToDaily returns daily from string
// func ToDaily(s string) *marketintelligencesvc.ApplicationViralDaily {
// 	var m *marketintelligencesvc.ApplicationViralDaily
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }
