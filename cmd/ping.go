/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping [host]",
	Short: "Ping a host to measure network latency",
	Long: `Ping sends ICMP echo requests to network hosts to determine 
their availability and measure the time it takes for packets to travel to the host and back (round-trip time).`,
	Args: cobra.ExactArgs(1), // One argument required, the host to ping
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]

		// Fetch flags
		count, _ := cmd.Flags().GetInt("count")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		interval, _ := cmd.Flags().GetDuration("interval")

		// Execute ping logic
		err := executePing(host, count, timeout, interval)
		if err != nil {
			fmt.Printf("Error executing ping: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Define flags for the ping command
	pingCmd.Flags().IntP("count", "c", 4, "Number of packets to send")
	pingCmd.Flags().DurationP("timeout", "t", 5*time.Second, "Timeout duration for each ping request")
	pingCmd.Flags().DurationP("interval", "i", 1*time.Second, "Interval between successive packets")
}

// executePing sends ICMP ping packets to the specified host
func executePing(host string, count int, timeout, interval time.Duration) error {
	// Create a new ping instance
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return fmt.Errorf("failed to create pinger: %v", err)
	}

	// Set ping configuration
	pinger.Count = count
	pinger.Timeout = timeout
	pinger.Interval = interval
	pinger.SetPrivileged(true) // Required to send ICMP packets

	// Print ping result
	fmt.Printf("PING %s (%s): %d data bytes\n", pinger.Addr(), pinger.IPAddr(), 64)

	// Start pinging
	err = pinger.Run()
	if err != nil {
		return fmt.Errorf("failed to ping host: %v", err)
	}

	// Get ping statistics
	stats := pinger.Statistics()
	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Printf("%d packets transmitted, %d packets received, %.1f%% packet loss\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	fmt.Printf("round-trip min/avg/max/stddev = %.3f/%.3f/%.3f/%.3f ms\n",
		stats.MinRtt.Seconds()*1000, stats.AvgRtt.Seconds()*1000, stats.MaxRtt.Seconds()*1000, stats.StdDevRtt.Seconds()*1000)

	return nil
}
