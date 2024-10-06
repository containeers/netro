/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// digCmd represents the dig command
var digCmd = &cobra.Command{
	Use:   "dig [domain]",
	Short: "Performs DNS lookups for the specified domain",
	Long: `Netro's dig command performs DNS lookups for the specified domain, 
similar to the 'dig' command in Unix. It supports querying for A, AAAA, MX, CNAME records, and prints the output in YAML format.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		simpleMode, _ := cmd.Flags().GetBool("s")
		queryDNS(domain, simpleMode)
	},
}

// Define the flag for simple mode
func init() {
	rootCmd.AddCommand(digCmd)
	digCmd.Flags().BoolP("s", "s", false, "Show only CNAME and A/AAAA IPs if available")
}

// DNSResults is a struct to hold all DNS query results in a structured format
type DNSResults struct {
	Domain string     `yaml:"domain"`
	A      []string   `yaml:"A,omitempty"`
	AAAA   []string   `yaml:"AAAA,omitempty"`
	CNAME  []string   `yaml:"CNAME,omitempty"` // Now supports multiple CNAMEs in the chain
	MX     []MXRecord `yaml:"MX,omitempty"`
	NS     []string   `yaml:"NS,omitempty"`
	TXT    []string   `yaml:"TXT,omitempty"`
}

type MXRecord struct {
	Host     string `yaml:"host"`
	Priority uint16 `yaml:"priority"`
}

// queryDNS performs DNS lookups and prints results in YAML, optionally with -s flag to show only CNAME and IPs
func queryDNS(domain string, simpleMode bool) {
	results := DNSResults{
		Domain: domain,
	}

	// A Record Lookup (NAME HERE <EMAIL ADDRESS>IPv4)
	aRecords, err := net.LookupIP(domain)
	if err == nil {
		for _, ip := range aRecords {
			if ip.To4() != nil {
				results.A = append(results.A, ip.String())
			}
		}
	}

	// AAAA Record Lookup (IPv6)
	for _, ip := range aRecords {
		if ip.To16() != nil && ip.To4() == nil {
			results.AAAA = append(results.AAAA, ip.String())
		}
	}

	// CNAME Lookup with chaining
	cnameChain := resolveCNAMEChain(domain)
	if len(cnameChain) > 0 {
		results.CNAME = cnameChain
	}

	// MX Record Lookup
	mxRecords, err := net.LookupMX(domain)
	if err == nil && !simpleMode { // Show MX records only in full mode
		for _, mx := range mxRecords {
			results.MX = append(results.MX, MXRecord{Host: mx.Host, Priority: mx.Pref})
		}
	}

	// NS Record Lookup (Name Servers)
	nsRecords, err := net.LookupNS(domain)
	if err == nil && !simpleMode { // Show NS records only in full mode
		for _, ns := range nsRecords {
			results.NS = append(results.NS, ns.Host)
		}
	}

	// TXT Record Lookup
	txtRecords, err := net.LookupTXT(domain)
	if err == nil && !simpleMode { // Show TXT records only in full mode
		results.TXT = append(results.TXT, txtRecords...)
	}

	// Handle printing results
	if simpleMode {
		// Only show CNAME and A/AAAA records in YAML
		printSimpleResults(results)
	} else {
		// Print all results in YAML format
		yamlOutput, err := yaml.Marshal(&results)
		if err != nil {
			fmt.Printf("Error marshaling to YAML: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(yamlOutput))
	}
}

// resolveCNAMEChain resolves a chain of CNAMEs starting from the initial domain
func resolveCNAMEChain(domain string) []string {
	var cnameChain []string

	for {
		cname, err := net.LookupCNAME(domain)
		if err != nil {
			break
		}

		// If the CNAME is the same as the domain, we've reached the final point
		if cname == domain {
			break
		}

		// Add the CNAME to the chain
		cnameChain = append(cnameChain, cname)

		// Continue resolving CNAME with the new domain name (next hop)
		domain = cname
	}

	return cnameChain
}

// printSimpleResults prints only CNAME and A/AAAA records in YAML format
func printSimpleResults(results DNSResults) {
	simpleResults := DNSResults{
		Domain: results.Domain,
		CNAME:  results.CNAME,
		A:      results.A,
		AAAA:   results.AAAA,
	}

	// Convert the simple results to YAML and print
	yamlOutput, err := yaml.Marshal(&simpleResults)
	if err != nil {
		fmt.Printf("Error marshaling to YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(yamlOutput))
}
