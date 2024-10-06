
# Netro

**netro** is a versatile command-line tool for networking, diagnostics, and troubleshooting. It provides utilities for DNS lookups, network interface management, system diagnostics, and more. Built with Go and powered by the Cobra CLI library, Netro is designed to be lightweight and efficient, making it a go-to tool for developers, system administrators, and network engineers.

## Features

- **DNS Lookups**: Perform DNS queries (A, AAAA, CNAME, MX, NS, TXT records).
- **Network Interface Information**: Retrieve network interfaces and IP addresses (similar to `ifconfig`).
- **Network Statistics**: View active connections and network diagnostics (similar to `netstat`).
- **Extensible**: Modular and easily extendable with new commands.
- **Cross-Platform**: Works on Linux, macOS, and Windows.
- **Versioning**: Displays the version and build information dynamically.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
  - [Global Options](#global-options)
  - [Commands](#commands)
    - [curl](#curl)
    - [dig](#dig)
    - [ifconfig](#ifconfig)
    - [nc](#nc)
    - [netstat](#netstat)
    - [version](#version)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Build from Source

To build Netro from source, make sure you have Go installed (Go 1.18+ is recommended).

1. Clone the repository:

   ```
   git clone https://github.com/containeers/netro.git
   ```

2. Navigate to the project directory:

   ```
   cd netro
   ```

3. Build the Netro binary:

   ```
   go build -o netro
   ```

4. Move the binary to a location in your \`PATH\`, for example:

   ```
   sudo mv netro /usr/local/bin/
   ```

Now, you can run `netro` from anywhere in your terminal.

### Install via Go

Alternatively, you can install directly using `go install`:

```
go install github.com/containeers/netro@latest
```

## Usage

After installation, you can use `netro` with a variety of commands and options.

### Global Options

| Option         | Description                                |
|----------------|--------------------------------------------|
| `--help`       | Show help for any command                  |
| `--version`    | Show the version of the Netro CLI          |
| `-t, --toggle` | Enable or disable specific features        |

### Commands

#### `curl`

Perform HTTP requests, similar to curl, with proxy, method, and header support.

**Usage**:

```
netro curl [url] [flags]
```

**Examples**:

- Perform a simple GET request:

  ```  
  netro curl http://example.com
  ```

- Perform a POST request with data:

  ```
  netro curl http://example.com -X POST -d '{"name": "Netro"}' -H "Content-Type: application/json"
  ```

- Use a proxy for the request:

  ```
  netro curl http://example.com -x http://proxy.example.com:8080
  ```

#### `dig`

Perform DNS lookups for domain names.

**Usage**:

```
netro dig [domain] [flags]
```

**Examples**:

- Perform a DNS lookup for example.com:

  ```
  netro dig example.com
  ```

- Show only CNAME and IPs for a domain:

  ```
  netro dig example.com -s
  ``

#### `ifconfig`

Display network interface information (IP addresses, MAC addresses, MTU).

**Usage**:

```
netro ifconfig [interface]
```

**Examples**:

- Show all network interfaces:

  ```
  netro ifconfig
  ```

- Show details for a specific network interface:

  ```
  netro ifconfig eth0
  ```

#### `nc`

Netcat-like functionality for TCP and UDP connections, with listening mode, proxies, and timeouts.

**Usage**:

```
netro nc [host] [port] [flags]
```

**Examples**:

- Open a TCP connection to a remote server:

  ```
  netro nc example.com 80 -p tcp
  ```

- Start a TCP server and listen on port 8080:

  ```
  netro nc -l 8080 -p tcp
  ```

- Open a UDP connection:

  ```
  netro nc example.com 53 -p udp
  ```

- Open a TCP connection using a proxy:

  ```
  netro nc example.com 80 -p tcp -x http://proxy.example.com:8080
  ```

#### `netstat`

Display active network connections and socket statistics (TCP, UDP, UNIX).

**Usage**:

```bash
netro netstat
```

**Examples**:

- Show all active network connections:

  ```
  netro netstat
  ```

#### `version`

Display the current version and build information for Netro.

**Usage**:

```
netro version
```

**Example**:

```
Netro version: v1.0.0 (built on Oct 9 2024)
```

## Contributing

We welcome contributions! If you want to contribute to Netro, follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request, and ensure that your changes pass all tests.

Please make sure to follow the standard Go style and write tests for any new features.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
