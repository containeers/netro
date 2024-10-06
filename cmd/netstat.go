/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/net"
	"github.com/spf13/cobra"
)

// netstatCmd represents the netstat command
var netstatCmd = &cobra.Command{
	Use:   "netstat",
	Short: "Displays network connections, routing tables, interface statistics, and process details.",
	Long:  `Netro's netstat command shows a list of active TCP and UDP connections, along with the process details (PID and process name) associated with each connection.`,
	Run: func(cmd *cobra.Command, args []string) {
		showNetstatWithProcesses()
	},
}

func init() {
	rootCmd.AddCommand(netstatCmd)
}

// showNetstatWithProcesses retrieves and prints active network connections along with associated processes
func showNetstatWithProcesses() {
	fmt.Println("Active Internet connections (servers and established)")
	fmt.Printf("%-7s %-56s %-56s %-11s\n", "Proto", "Local Address", "Foreign Address", "State")

	connections, err := net.Connections("all")
	if err != nil {
		log.Fatalf("Error retrieving network connections: %v", err)
	}

	for _, conn := range connections {
		protocol := getProtocolType(conn.Type) // Convert conn.Type to a string
		localAddr := fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port)
		remoteAddr := fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port)
		state := conn.Status

		// Display the connection details along with the process name and PID
		fmt.Printf("%-7s %-56s %-56s %-11s\n", protocol, localAddr, remoteAddr, state)
	}
}

// getProtocolType converts the protocol type from uint32 to a human-readable string
func getProtocolType(protocol uint32) string {
	switch protocol {
	case 1:
		return "tcp"
	case 2:
		return "udp"
	case 5:
		return "unix"
	default:
		return "unknown"
	}
}
