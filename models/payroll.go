package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"time"
)

// Payroll holds payroll information
type Payroll struct {
	ID               int            `json:"id"`
	Kunde            string         `json:"kunde"`
	Abrechnungskreis int            `json:"abrechnungskreis"`
	Abrechnungsmonat string         `json:"abrechnungsmonat"`
	Lauftermin       string         `json:"lauftermin"`
	Untertermin      string         `json:"untertermin"`
	Jobname          string         `json:"jobname"`
	Jobkette         string         `json:"jobkette"`
	Worker           sql.NullString `json:"worker"`
	IsEnd            int            `json:"isEnd"`
	ReturnCode       int            `json:"returnCode"`
	Timestamp        time.Time      `json:"timestamp"`
	HeapXmx          int            `json:"heapXmx"`
	Threads          sql.NullString `json:"threads"`
	ServerID         int            `json:"serverId"`
}

// WriteJSON writes the run statistics to a file in JSON encoding.
func (payroll *Payroll) WriteJSON(c Configuration, execDir string, statsFile string) {
	bytes, err := json.Marshal(payroll)
	if err != nil {
		log.Fatal(err)
	}

	bytes = append(bytes, "\n"...)

	var statsLogPath string
	if c.StatsPath == "" {
		statsLogPath = execDir + "/" + statsFile
	} else {
		statsLogPath = c.StatsPath
	}

	file, err := os.OpenFile(statsLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

// CheckMonth checks whether a payroll struct has correct month or not
func (payroll *Payroll) CheckMonth() bool {
	r, err := regexp.Compile("^\\d{4}(0[1-9]|1[0-2])$")
	if err != nil {
		return false
	}
	return r.MatchString(payroll.Abrechnungsmonat)
}

// CheckEnd checks whether a payroll struct has correct end condition or not
func (payroll *Payroll) CheckEnd() bool {
	if payroll.IsEnd == 0 && payroll.ReturnCode < 0 {
		return true
	}
	if payroll.IsEnd == 1 && payroll.ReturnCode >= 0 {
		return true
	}
	return false
}
