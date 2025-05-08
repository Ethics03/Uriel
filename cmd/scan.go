/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Your personal API-Security Tester",
	Long: `Uriel is an API-Security Testing CLI-Tool that uses 
	Llama-3 with tailored checks for seamless API performance.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := "Scan this API endpoint for basic security issues: https://api.example.com"

		data := map[string]interface{}{
			"model":  "llama3.2:latest",
			"prompt": prompt,
			"stream": false,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to send request to Ollama: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}

		fmt.Println("Uriel (Llama 3.2) says:")
		fmt.Println(string(body))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
