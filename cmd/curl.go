/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl [URL]",
	Short: "Perform HTTP requests like curl",
	Long: `Netro's curl command lets you perform HTTP requests similar to the original curl utility. 
It supports proxies (-x), payloads (-d), multiple headers (-H), HTTP methods (-X), verbose output (-v), TLS details for HTTPS requests, and the ability to skip TLS verification (-k).`,
	Args: cobra.MinimumNArgs(1), // At least one argument is required (the URL)
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// Fetch flags
		proxy, _ := cmd.Flags().GetString("proxy")
		data, _ := cmd.Flags().GetString("data")
		headers, _ := cmd.Flags().GetStringArray("header")
		method, _ := cmd.Flags().GetString("method")
		verbose, _ := cmd.Flags().GetBool("verbose")
		insecure, _ := cmd.Flags().GetBool("insecure")

		// Execute the curl logic
		err := executeCurl(url, proxy, data, headers, method, verbose, insecure)
		if err != nil {
			fmt.Printf("Error executing curl: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(curlCmd)

	// Define flags for the curl command
	curlCmd.Flags().StringP("proxy", "x", "", "Specify a proxy to use")
	curlCmd.Flags().StringP("data", "d", "", "HTTP POST data (triggers POST request or other methods with -X)")
	curlCmd.Flags().StringArrayP("header", "H", []string{}, "Specify multiple headers (can be used multiple times)")
	curlCmd.Flags().StringP("method", "X", "GET", "Specify the HTTP method to use (GET, POST, PUT, DELETE, etc.)")
	curlCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output to show request and response details, including TLS details")
	curlCmd.Flags().BoolP("insecure", "k", false, "Allow insecure server connections when using SSL (skip TLS certificate verification)")
}

// executeCurl performs the HTTP request based on the provided flags
func executeCurl(urlStr, proxy, data string, headers []string, method string, verbose, insecure bool) error {
	// Create HTTP transport
	transport := &http.Transport{
		// Set TLS client configuration
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure, // Skip certificate verification if insecure mode is enabled
		},
	}

	// If a proxy is specified, set the proxy
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return fmt.Errorf("invalid proxy URL: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// Create HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}

	// Default to GET method if no method is specified
	if method == "" {
		method = "GET"
	}

	// Create the request, using the specified method
	var req *http.Request
	var err error
	if data != "" {
		req, err = http.NewRequest(method, urlStr, bytes.NewBuffer([]byte(data)))
	} else {
		req, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers to the request
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %s", header)
		}
		req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	// If verbose is enabled, print the request details
	if verbose {
		fmt.Println("----- Request -----")
		fmt.Printf("Method: %s\n", req.Method)
		fmt.Printf("URL: %s\n", req.URL)
		fmt.Println("Headers:")
		for key, value := range req.Header {
			fmt.Printf("  %s: %s\n", key, strings.Join(value, ", "))
		}
		if data != "" {
			fmt.Printf("Body: %s\n", data)
		}
		fmt.Println("-------------------")
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response body using io.ReadAll (instead of ioutil.ReadAll)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// If verbose is enabled, print the response details
	if verbose {
		fmt.Println("----- Response -----")
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Println("Headers:")
		for key, value := range resp.Header {
			fmt.Printf("  %s: %s\n", key, strings.Join(value, ", "))
		}

		// Print TLS details if the request was over HTTPS
		if resp.TLS != nil {
			printTLSDetails(resp.TLS)
		}
		fmt.Println("--------------------")
	}

	// Print the response body
	fmt.Printf("\nResponse Body:\n%s\n", string(body))

	return nil
}

// printTLSDetails prints TLS details from the response
func printTLSDetails(tlsState *tls.ConnectionState) {
	fmt.Println("----- TLS Information -----")
	fmt.Printf("Version: %s\n", tlsVersionToString(tlsState.Version))
	fmt.Printf("Cipher Suite: %s\n", tls.CipherSuiteName(tlsState.CipherSuite))
	fmt.Println("Server Certificates:")
	for _, cert := range tlsState.PeerCertificates {
		fmt.Printf("  Subject: %s\n", cert.Subject)
		fmt.Printf("  Issuer: %s\n", cert.Issuer)
		fmt.Printf("  Valid From: %s\n", cert.NotBefore.Format(time.RFC3339))
		fmt.Printf("  Valid Until: %s\n", cert.NotAfter.Format(time.RFC3339))
	}
	fmt.Println("----------------------------")
}

// tlsVersionToString converts the TLS version to a human-readable string
func tlsVersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return "Unknown TLS version"
	}
}
