/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/dmittelstaedt/dpaycol/models"
	"github.com/dmittelstaedt/dpaycol/utils"

	"github.com/spf13/cobra"
)

var payroll models.Payroll

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add payroll information to REST API",
	Long: `Add payroll information to REST API including Abrechnungskreis, 
Abrechnungsmonat, Jobname and ID of the Jobkette.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		payroll.Timestamp = time.Now()
		hostname, _ := os.Hostname()
		apiURL := "http://" + configuration.APIEndpoint + "/servers?" + hostname
		resp := utils.SendRequest("GET", apiURL)
		defer resp.Body.Close()

		var servers []models.Server
		if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
			log.Println(err)
		}

		payroll.ServerID = servers[0].ID
		if !payroll.CheckMonth() {
			return models.ErrInvalidMonth
		}
		if !payroll.CheckEnd() {
			return models.ErrEndCondition
		}
		payroll.WriteJSON(configuration, execDir, statsFile)
		apiURLPayrolls := "http://" + configuration.APIEndpoint + "/payrolls"
		respPayroll := utils.SendRequest("POST", apiURLPayrolls)
		defer respPayroll.Body.Close()

		var payrolls []models.Payroll
		if err := json.NewDecoder(resp.Body).Decode(&payrolls); err != nil {
			log.Println(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringVarP(&payroll.Kunde, "kunde", "c", "", "Kunde intern")
	addCmd.Flags().IntVarP(&payroll.Abrechnungskreis, "abrechnungskreis", "k", -1, "Abrechnungskreis (required)")
	addCmd.MarkFlagRequired("abrechnungskreis")
	addCmd.Flags().StringVarP(&payroll.Abrechnungsmonat, "abrechnungsmonat", "m", "", "Abrechnungsmonat (required)")
	addCmd.MarkFlagRequired("abrechnungsmonat")
	addCmd.Flags().StringVarP(&payroll.Lauftermin, "lauftermin", "l", "", "Lauftermin")
	addCmd.Flags().StringVarP(&payroll.Untertermin, "untertermin", "u", "", "Untertermin")
	addCmd.Flags().StringVarP(&payroll.Jobname, "jobname", "n", "", "Jobname (required)")
	addCmd.MarkFlagRequired("jobname")
	addCmd.Flags().StringVarP(&payroll.Jobkette, "jobkette", "i", "", "ID of the Jobkette (required)")
	addCmd.MarkFlagRequired("jobkette")
	addCmd.Flags().BoolVarP(&payroll.IsEnd, "end", "e", false, "End flag")
	addCmd.Flags().IntVarP(&payroll.ReturnCode, "returncode", "r", -1, "Return Code of the job (required if end)")
}
