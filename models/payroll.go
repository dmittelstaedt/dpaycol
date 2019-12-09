package models

import "time"

// Payroll holds payroll information
type Payroll struct {
	ID               int       `json:"id"`
	Abrechnungskreis int       `json:"abrechnungskreis"`
	Abrechnungsmonat string    `json:"abrechnungsmonat"`
	Lauftermin       string    `json:"lauftermin"`
	Untertermin      string    `json:"untertermin"`
	Jobname          string    `json:"jobname"`
	Jobkette         string    `json:"jobkette"`
	IsEnd            bool      `json:"isEnd"`
	Timestamp        time.Time `json:"timestamp"`
	ServerID         int       `json:"serverId"`
}
