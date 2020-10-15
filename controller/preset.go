package controller

import (
	"nicheanal.com/dal"
)

// CreatePreset save new preset data to db
func CreatePreset(p dal.AppPreset) error {
	return dal.CreatePreset(p)
}

// ListPresets lists saved presets
func ListPresets() ([]*dal.AppPreset, error) {
	return dal.ListPreset()
}
