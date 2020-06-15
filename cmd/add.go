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
var worker string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add payroll information to REST API",
	Long: `Add payroll information to REST API including Abrechnungskreis, 
Abrechnungsmonat, Jobname and ID of the Jobkette.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		payroll.Timestamp = time.Now()
		hostname, _ := os.Hostname()
		payroll.ServerID = utils.GetServerID("http://"+configuration.APIEndpoint+"/servers", hostname)
		if cmd.Flags().Changed("worker") {
			payroll.Worker.String = worker
			payroll.Worker.Valid = true
		}

		if !payroll.CheckMonth() {
			return models.ErrInvalidMonth
		}
		if !payroll.CheckEnd() {
			return models.ErrEndCondition
		}
		payroll.WriteJSON(configuration, execDir, statsFile)

		body, err := json.Marshal(payroll)
		if err != nil {
			log.Printf("Error marshalling payroll: %+v", err)
			return models.ErrSendToAPI
		}
		respPayroll := utils.SendRequest("POST", "http://"+configuration.APIEndpoint+"/payrolls", body)
		defer respPayroll.Body.Close()

		var payrollResp models.Payroll
		if err := json.NewDecoder(respPayroll.Body).Decode(&payrollResp); err != nil {
			log.Printf("Error decoding payroll: %+v", err)
			return models.ErrSendToAPI
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
	addCmd.Flags().StringVarP(&worker, "worker", "w", "", "Worker")
	addCmd.MarkFlagRequired("jobkette")
	addCmd.Flags().BoolVarP(&payroll.IsEnd, "end", "e", false, "End flag")
	addCmd.Flags().IntVarP(&payroll.ReturnCode, "returncode", "r", -1, "Return Code of the job (required if end)")
}
