package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	targetURL       string
	authCheck       bool
	headers         []string
	requestBodyFile string
)

var requestBody string

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Your personal API-Security Tester",
	Long: `Uriel is an API-Security Testing CLI-Tool that uses 
	Llama-3 with tailored checks for seamless API performance.`,
	Run: func(cmd *cobra.Command, args []string) {

		if targetURL == "" {
			log.Fatalf("Please provide a URL using the --url flag.")
		}

		var req *http.Request
		var err error

		if requestBodyFile != "" {
			content, err := os.ReadFile(requestBodyFile)
			if err != nil {
				log.Fatalf("Failed to read body file: %v", err)
			}
			requestBody = string(content)
		}

		if requestBody != "" {
			req, err = http.NewRequest("POST", targetURL, bytes.NewBuffer([]byte(requestBody)))
		} else {
			req, err = http.NewRequest("GET", targetURL, nil)
		}
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		for _, header := range headers {
			parts := strings.SplitN(header, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(parts[0], parts[1])
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to perform HTTP request: %v", err)
		}
		defer resp.Body.Close()

		headersJSON, _ := json.MarshalIndent(resp.Header, "", "  ")

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)
		if len(bodyStr) > 300 {
			bodyStr = bodyStr[:300] + "...(truncated)"
		}
		prompt := fmt.Sprintf(`Scan this API endpoint for authentication or security issues.

					 URL: %s
					 Response status: %d
					 Headers: %s
					 Body (truncated): %s

					Please analyze this response and suggest what kind of authentication mechanism is expected, and any potential security flaws.
					make it really short and simple not more than 200 words
						`, targetURL, resp.StatusCode, string(headersJSON), bodyStr)

		if authCheck {
			prompt = "check for authentication-related flaws such as missing tokens, insecure login flows, or improper access control. Not more than 200 words"
		}

		data := map[string]interface{}{
			"model":  "llama3",
			"prompt": prompt,
			"stream": false,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Failed to marshal request: %v", err)
		}

		Ollamaresp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to send request to Ollama: %v", err)
		}
		defer Ollamaresp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(Ollamaresp.Body).Decode(&result); err != nil {
			log.Fatalf("Failed to decode JSON response: %v", err)
		}

		if responseText, ok := result["response"].(string); ok {
			fmt.Println("Uriel (Llama 3.2) Scan Result:")
			fmt.Println(responseText)
		} else {
			fmt.Println("Unexpected response format from Ollama:")
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&targetURL, "url", "u", "", "Target API URL to scan (required)")

	scanCmd.Flags().BoolVar(&authCheck, "auth", false, "Check for authentication flaws")

	scanCmd.Flags().StringArrayVar(&headers, "header", []string{}, "Custom headers in 'Key: Value' format")

	scanCmd.Flags().StringVar(&requestBody, "body", "", "Optional JSON body to include in the request")

	scanCmd.Flags().StringVar(&requestBodyFile, "body-file", "", "Path to JSON file to use as request body")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
