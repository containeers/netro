/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// Define variables to hold the real implementations of the net package functions
var getInterfaces = net.Interfaces
var getInterfaceByName = net.InterfaceByName

// ifconfigCmd represents the ifconfig command
var ifconfigCmd = &cobra.Command{
	Use:   "ifconfig [interface name]",
	Short: "Displays network interface information",
	Long:  `Displays network interface details. You can provide an interface name to show details of that specific interface, or leave it empty to show details for all interfaces.`,
	Args:  cobra.MaximumNArgs(1), // Allows 0 or 1 argument
	Run: func(cmd *cobra.Command, args []string) {
		// If an interface name is provided, filter by that name
		if len(args) == 1 {
			interfaceName := args[0]
			showInterfaceDetails(interfaceName)
		} else {
			// Otherwise, show details for all interfaces
			showAllInterfacesDetails()
		}
	},
}

func init() {
	rootCmd.AddCommand(ifconfigCmd)
}

// Function to show details of a specific interface
func showInterfaceDetails(interfaceName string) error {
	// Get the network interface by name
	iface, err := getInterfaceByName(interfaceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching interface %s: %v\n", interfaceName, err)
		return err
	}

	// Display interface information
	printInterfaceDetails(iface)
	return nil
}

// Function to show details of all interfaces
func showAllInterfacesDetails() {
	// Get a list of all network interfaces on the system
	interfaces, err := getInterfaces()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching interfaces: %v\n", err)
		os.Exit(1)
	}

	// Check if there are any interfaces
	if len(interfaces) == 0 {
		fmt.Println("No network interfaces found.")
		return
	}

	// Loop through each interface and display its information
	for _, iface := range interfaces {
		printInterfaceDetails(&iface)
	}
}

// Function to print the details of a given interface
func printInterfaceDetails(iface *net.Interface) {
	// Interface Name
	fmt.Printf("Interface: %s\n", iface.Name)

	// MAC Address (HardwareAddr)
	if len(iface.HardwareAddr) > 0 {
		fmt.Printf("  MAC Address: %s\n", iface.HardwareAddr)
	} else {
		fmt.Println("  MAC Address: N/A")
	}

	// MTU (Maximum Transmission Unit)
	fmt.Printf("  MTU: %d\n", iface.MTU)

	// Flags (Up, Loopback, etc.)
	fmt.Printf("  Flags: %s\n", iface.Flags)

	// Get and display IP addresses assigned to the interface
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error fetching addresses for interface %s: %v\n", iface.Name, err)
		return
	}

	if len(addrs) > 0 {
		fmt.Println("  IP Addresses:")
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok {
				// Print the IP address
				fmt.Printf("    - IP Address: %s\n", ipNet.IP.String())

				// Print the Netmask
				fmt.Printf("      Netmask: %s\n", net.IP(ipNet.Mask).String())
			} else {
				// If it's not an IPNet (rare case), print the address as it is
				fmt.Printf("    - %s\n", addr.String())
			}
		}
	} else {
		fmt.Println("  IP Addresses: None")
	}

	fmt.Println() // Add extra line for better readability
}
