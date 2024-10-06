/*
Copyright © 2024 Sandarsh Devappa <sd@containeers.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ncCmd represents the nc (Netcat) command
var ncCmd = &cobra.Command{
	Use:   "nc [host] [port]",
	Short: "Netro's implementation of Netcat (nc) for TCP and UDP connections",
	Long: `Netro's Netcat (nc) command supports TCP and UDP connections for interacting 
with remote servers. It can also listen for incoming connections using the -l flag.`,
	Args: cobra.RangeArgs(1, 2), // Accept one or two arguments (host is optional in listen mode)
	Run: func(cmd *cobra.Command, args []string) {
		var host, port string

		// In listen mode, we only need the port; otherwise, both host and port
		if len(args) == 1 {
			port = args[0]
		} else {
			host = args[0]
			port = args[1]
		}

		// Fetch flags
		protocol, _ := cmd.Flags().GetString("protocol")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		proxy, _ := cmd.Flags().GetString("proxy")
		listen, _ := cmd.Flags().GetBool("listen")

		// Execute the appropriate logic (listen mode or normal mode)
		if listen {
			err := executeNCListen(port, protocol)
			if err != nil {
				fmt.Printf("Error executing nc listen: %v\n", err)
				os.Exit(1)
			}
		} else {
			err := executeNC(host, port, protocol, timeout, proxy)
			if err != nil {
				fmt.Printf("Error executing nc: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ncCmd)

	// Define flags for the nc command
	ncCmd.Flags().StringP("protocol", "p", "tcp", "Specify the protocol to use (tcp or udp)")
	ncCmd.Flags().DurationP("timeout", "t", 5*time.Second, "Set timeout duration for the connection")
	ncCmd.Flags().StringP("proxy", "x", "", "Specify a TCP proxy URL for TCP connections (e.g., http://proxy.example.com:8080)")
	ncCmd.Flags().BoolP("listen", "l", false, "Listen for incoming connections on the specified port")
}

// executeNC handles TCP or UDP connections based on the provided protocol
func executeNC(host, port, protocol string, timeout time.Duration, proxy string) error {
	address := net.JoinHostPort(host, port)

	if protocol == "tcp" {
		// Handle TCP connection
		if proxy != "" {
			// Use proxy for TCP connection
			return executeTCPProxy(address, timeout, proxy)
		}
		return executeTCP(address, timeout)
	} else if protocol == "udp" {
		// Handle UDP connection
		return executeUDP(address, timeout)
	} else {
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// executeNCListen handles listening for incoming connections on the specified port
func executeNCListen(port, protocol string) error {
	address := net.JoinHostPort("", port) // Listen on all available interfaces

	if protocol == "tcp" {
		// Start TCP listener
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("failed to start TCP listener: %v", err)
		}
		defer listener.Close()

		fmt.Printf("Listening on %s (TCP)\n", address)

		// Accept incoming connections
		for {
			conn, err := listener.Accept()
			if err != nil {
				return fmt.Errorf("failed to accept connection: %v", err)
			}
			go handleTCPConnection(conn)
		}
	} else if protocol == "udp" {
		// Start UDP listener
		conn, err := net.ListenPacket("udp", address)
		if err != nil {
			return fmt.Errorf("failed to start UDP listener: %v", err)
		}
		defer conn.Close()

		fmt.Printf("Listening on %s (UDP)\n", address)

		// Handle UDP communication
		handleUDPConnection(conn)
	} else {
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}

	return nil
}

// handleTCPConnection handles an incoming TCP connection
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

	// Copy data between the connection and stdout/stderr
	go io.Copy(conn, os.Stdin) // Send data from stdin to the connection
	io.Copy(os.Stdout, conn)   // Receive data from the connection and print it
}

// handleUDPConnection handles UDP communication
func handleUDPConnection(conn net.PacketConn) {
	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error reading from UDP connection: %v\n", err)
			return
		}

		fmt.Printf("Received %d bytes from %s: %s\n", n, addr, strings.TrimSpace(string(buf[:n])))

		// Send response back
		_, err = conn.WriteTo([]byte("Message received"), addr)
		if err != nil {
			fmt.Printf("Error sending response: %v\n", err)
			return
		}
	}
}

// executeTCP establishes a TCP connection to the specified address
func executeTCP(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return fmt.Errorf("failed to establish TCP connection: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s (TCP)\n", address)
	return nil
}

// executeTCPProxy establishes a TCP connection through a proxy to the specified address
func executeTCPProxy(address string, timeout time.Duration, proxyURL string) error {

	// Parse the proxy URL
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %v", err)
	}

	// Connect to the proxy
	conn, err := net.DialTimeout("tcp", proxy.Host, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to proxy: %v", err)
	}
	defer conn.Close()

	// Send the HTTP CONNECT request to the proxy
	connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", address, address)
	_, err = conn.Write([]byte(connectReq))
	if err != nil {
		return fmt.Errorf("failed to send CONNECT request: %v", err)
	}

	// Read the proxy's response
	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return fmt.Errorf("failed to read proxy response: %v", err)
	}
	defer resp.Body.Close()

	// Check if the proxy successfully established the connection
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("proxy connection failed: %s", resp.Status)
	}

	fmt.Printf("Connected to %s through HTTP proxy %s\n", address, proxyURL)

	// You can now send and receive data over `conn`
	// This is where you'd typically implement the netcat-like functionality for communication
	// For example, using `conn.Read` and `conn.Write` to interact with the remote server

	return nil
}

// executeUDP establishes a UDP connection to the specified address
func executeUDP(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("udp", address, timeout)
	if err != nil {
		return fmt.Errorf("failed to establish UDP connection: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s (UDP)\n", address)
	return nil
}
