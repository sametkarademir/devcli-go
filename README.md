# DevKit CLI

A versatile command-line toolkit for developers and system administrators, built with Go.

## Features

- **Single Binary**: Distributed as a single executable, works cross-platform
- **Modular Design**: Organized command groups for different use cases
- **Multiple Output Formats**: Support for plain, JSON, and table formats
- **Pipe Support**: Works seamlessly with Unix pipes
- **Tab Completion**: Auto-completion support for bash, zsh, fish, and PowerShell

## Installation

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd DevCli

# Build for current platform
make build

# Build for macOS (current architecture)
make build-macos

# Build for macOS (both amd64 and arm64)
make build-macos-all

# Build for all platforms
make build-all
```

**Build Output:**
- `bin/devkit` - Current platform binary
- `bin/devkit-macos` - macOS binary (current architecture)
- `bin/devkit-darwin-amd64` - macOS Intel (amd64)
- `bin/devkit-darwin-arm64` - macOS Apple Silicon (arm64)

### Install

```bash
make install
```

### Shell Completion (Tab Auto-completion)

DevKit supports tab completion for bash, zsh, fish, and PowerShell.

#### Zsh (macOS/Linux)

```bash
# Generate completion script
devcli completion zsh > ~/.devcli-completion.zsh

# Add to your ~/.zshrc
echo "source ~/.devcli-completion.zsh" >> ~/.zshrc

# Reload shell (yeni terminal açın veya)
exec zsh
# veya
source ~/.zshrc
```

**Test Completion:**
```bash
# Tab tuşuna basarak test edin:
devcli <TAB>              # completion, dev, file, net, help gösterir
devcli dev <TAB>          # base64, hash, ulid, uuid gösterir
devcli dev hash <TAB>     # md5, sha1, sha256, sha512 gösterir
devcli dev base64 <TAB>   # encode, decode gösterir
devcli dev uuid --<TAB>   # --version, --count, --output gösterir
```

#### Bash

```bash
# Generate completion script
devcli completion bash > ~/.devcli-completion.bash

# Add to your ~/.bashrc or ~/.bash_profile
echo "source ~/.devcli-completion.bash" >> ~/.bashrc

# Reload shell
source ~/.bashrc
```

#### Fish

```bash
# Generate completion script
devkit completion fish > ~/.config/fish/completions/devkit.fish
```

After setup, you can use Tab to auto-complete commands, subcommands, and flags!

## Usage

### Developer Tools (`dev`)

#### UUID Generation

Generate UUID (v4 or v7) values:

```bash
# Generate UUID v4 (default)
devcli dev uuid

# Generate UUID v7
devcli dev uuid --version 7

# Generate multiple UUIDs
devcli dev uuid --count 5

# JSON output
devcli dev uuid --version 7 --count 3 --output json
```

#### ULID Generation

Generate ULID (Universally Unique Lexicographically Sortable Identifier) values:

```bash
# Generate a single ULID
devcli dev ulid

# Generate multiple ULIDs
devcli dev ulid --count 5

# JSON output
devcli dev ulid --count 3 --output json
```

#### Base64 Encode/Decode

Encode or decode base64 strings:

```bash
# Encode string
devcli dev base64 encode "hello world"

# Decode string
devcli dev base64 decode "aGVsbG8gd29ybGQ="

# Encode file
devcli dev base64 encode --file ./image.png

# From stdin
echo "test" | devcli dev base64 encode --stdin
echo "dGVzdA==" | devcli dev base64 decode --stdin

# JSON output
devcli dev base64 encode "hello world" --output json
```

#### JWT Operations

Decode and verify JWT tokens:

```bash
# Decode JWT token (without verification)
devcli dev jwt decode "eyJhbGciOiJIUzI1NiIs..."

# Decode from file
devcli dev jwt decode --file token.txt

# Decode from stdin
echo "eyJ..." | devcli dev jwt decode --stdin

# Verify JWT token
devcli dev jwt verify "eyJ..." --secret "my-secret-key"

# JSON output
devcli dev jwt decode "eyJ..." --output json
```

#### Hash Calculation

Calculate cryptographic hashes of strings or files:

```bash
# String hash
devcli dev hash sha256 "hello world"

# File hash
devcli dev hash md5 --file /path/to/file

# From stdin
echo "hello" | devcli dev hash sha256 --stdin

# JSON output
devcli dev hash sha256 "hello world" --output json
```

Supported algorithms: `md5`, `sha1`, `sha256`, `sha512`

#### URL Operations

URL encode, decode, and parse:

```bash
# Encode URL
devcli dev url encode "hello world"

# Decode URL
devcli dev url decode "hello%20world"

# Parse URL
devcli dev url parse "https://example.com/path?key=value"
```

#### HTML Entity Operations

HTML entity encode and decode:

```bash
# Encode HTML entities
devcli dev html encode "hello <world> & more"

# Decode HTML entities
devcli dev html decode "hello &lt;world&gt; &amp; more"
```

#### JSON Operations

JSON processing operations:

```bash
# Prettify JSON
devcli dev json prettify '{"name":"John","age":30}'

# Minify JSON
devcli dev json minify '{"name": "John", "age": 30}'

# Validate JSON
devcli dev json validate '{"name":"John","age":30}'

# Query JSON path
echo '{"name":"John"}' | devcli dev json path '$.name' --stdin
devcli dev json path '$.users[0].name' --file data.json
```

#### Epoch/Unix Timestamp

Convert between Unix timestamps and dates:

```bash
# Convert timestamp to date
devcli dev epoch 1699876543

# Convert date to timestamp
devcli dev epoch --to-unix "2024-01-15 10:30:00"

# Get current timestamp
devcli dev epoch now
```

#### Random Data Generation

Generate random strings, numbers, and passwords:

```bash
# Random string
devcli dev random string --length 16

# Random number
devcli dev random number --min 1 --max 100

# Random password
devcli dev random password --length 16 --symbols
```

#### Lorem Ipsum Generator

Generate placeholder text:

```bash
# Generate words
devcli dev lorem word --count 5

# Generate sentences
devcli dev lorem sentence --count 2

# Generate paragraphs
devcli dev lorem paragraph --count 1
```

#### Cron Expression Parser

Parse and explain cron expressions:

```bash
# Explain cron expression
devcli dev cron explain "0 9 * * 1-5"

# Get next run times
devcli dev cron next "0 9 * * 1-5" --count 5
```

#### Semantic Versioning

Compare and bump semantic versions:

```bash
# Compare versions
devcli dev semver compare "1.2.3" "1.2.4"

# Bump version
devcli dev semver bump major "1.2.3"
devcli dev semver bump minor "1.2.3"
devcli dev semver bump patch "1.2.3"
```

#### Environment File Management

Manage `.env` files:

```bash
# Set environment variable
devcli dev env set TEST_KEY=test_value --file .env

# Get environment variable
devcli dev env get TEST_KEY --file .env

# List all variables
devcli dev env list --file .env

# Unset environment variable
devcli dev env unset TEST_KEY --file .env
```

### File Operations (`file`)

#### File Statistics

Display detailed file information:

```bash
# File info
devcli file stat README.md

# Directory info
devcli file stat .

# JSON output
devcli file stat README.md --output json
```

#### Directory Tree

Display directory structure as a tree:

```bash
# Current directory
devcli file tree .

# With depth limit
devcli file tree . --depth 2

# Include hidden files
devcli file tree . --all
```

#### File Search

Search for text in files with colored output:

```bash
# Basic search
devcli file search "TODO" .

# Recursive search
devcli file search "function" ./src --recursive

# Filter by extension
devcli file search "error" . --extensions "go,js"

# Ignore directories
devcli file search "error" . --ignore "node_modules,vendor"

# Regex search
devcli file search "func.*main" . --regex --extensions "go"
```

#### Find and Replace

Find and replace text in multiple files:

```bash
# Basic find-replace
devcli file find-replace "old" "new" .

# Recursive
devcli file find-replace "TODO" "DONE" ./src --recursive

# Filter by extension
devcli file find-replace "error" "err" . --extensions "go" --dry-run

# Regex find-replace
devcli file find-replace "\\d+" "NUMBER" . --regex --dry-run
```

#### Bulk Rename

Rename files using patterns:

```bash
# Add prefix
devcli file rename --pattern "*.txt" --prefix "backup_" --path ./docs

# Replace pattern
devcli file rename --pattern "IMG_*.jpg" --replace "IMG_" "photo_" --path ./images

# Case conversion
devcli file rename --pattern "*.txt" --case upper --path ./docs --dry-run

# Add suffix
devcli file rename --pattern "*.txt" --suffix "_backup" --path ./docs
```

#### Format Conversion

Convert between file formats (JSON, YAML, TOML):

```bash
# JSON to YAML
devcli file convert config.json --to yaml

# YAML to JSON
devcli file convert config.yaml --to json --output config.json

# TOML to YAML
devcli file convert data.toml --to yaml
```

#### File Diff

Compare two files:

```bash
# Compare files
devcli file diff file1.txt file2.txt

# Unified diff format
devcli file diff file1.txt file2.txt --unified
```

#### Duplicate File Detection

Find and remove duplicate files:

```bash
# Find duplicates by hash
devcli file dedupe ./downloads --by hash

# Find duplicates by name
devcli file dedupe ./photos --by name --action list

# Delete duplicates (dry-run)
devcli file dedupe ./downloads --by hash --action delete --dry-run

# Recursive search
devcli file dedupe . --recursive --by hash
```

#### File Watching

Watch files for changes:

```bash
# Watch directory
devcli file watch ./src

# Execute command on change
devcli file watch ./src --on-change "go build"

# Watch specific pattern
devcli file watch . --pattern "*.go" --on-change "go test ./..."
```

### Network & System Operations (`net`)

#### Port Operations

Port scanning and status checking:

```bash
# Check port status
devcli net port check 8080
devcli net port check 80 --host google.com

# Scan port range
devcli net port scan localhost --range 1-1000
devcli net port scan 127.0.0.1 --range 80-443 --timeout 2

# List listening ports
devcli net port list
```

#### DNS Lookup

DNS lookup operations:

```bash
# A record (default)
devcli net dns lookup google.com

# Specific record types
devcli net dns lookup google.com --type MX
devcli net dns lookup google.com --type TXT
devcli net dns lookup google.com --type NS

# Reverse DNS
devcli net dns reverse 8.8.8.8
```

#### IP Information

Get IP address information:

```bash
# Public IP
devcli net ip

# Local IP
devcli net ip --local

# IP information
devcli net ip info 8.8.8.8
devcli net ip info 2001:4860:4860::8888
```

#### HTTP Requests

Send HTTP requests:

```bash
# GET request
devcli net http get https://api.github.com

# GET with headers
devcli net http get https://api.github.com --header "Accept: application/json"

# POST request
devcli net http post https://httpbin.org/post --data '{"name":"test"}'

# PUT request
devcli net http put https://httpbin.org/put --data '{"id":1}'

# DELETE request
devcli net http delete https://httpbin.org/delete
```

#### Ping

Ping a host with statistics:

```bash
# Basic ping
devcli net ping google.com

# Custom count
devcli net ping 8.8.8.8 --count 10

# Custom timeout
devcli net ping github.com --timeout 5
```

#### SSL Certificate

Check SSL certificate information:

```bash
# Certificate info
devcli net ssl check google.com
devcli net ssl check github.com:443

# Certificate expiry
devcli net ssl expiry google.com
devcli net ssl expiry example.com
```

#### Whois Lookup

Domain whois lookup:

```bash
# Whois query
devcli net whois google.com
devcli net whois github.com
```

#### Internet Speed Test

Test internet connection speed:

```bash
# Speed test
devcli net speed

# JSON output
devcli net speed --output json
```

#### System Information

Display system information:

```bash
# Full system info
devcli net sysinfo

# CPU only
devcli net sysinfo --cpu

# Memory only
devcli net sysinfo --memory

# Disk only
devcli net sysinfo --disk

# JSON output
devcli net sysinfo --output json
```

#### Process Management

List and manage processes:

```bash
# List processes
devcli net ps

# Sort by CPU
devcli net ps --sort cpu

# Sort by memory
devcli net ps --sort mem

# Filter processes
devcli net ps --filter "go" --limit 10

# Table output
devcli net ps --output table

# JSON output
devcli net ps --output json
```

#### Disk Usage

Analyze disk usage:

```bash
# Disk usage
devcli net disk /

# Current directory
devcli net disk .

# Table output
devcli net disk /home --output table
```

#### Network Interfaces

List network interfaces:

```bash
# List interfaces
devcli net interfaces

# JSON output
devcli net interfaces --output json
```

#### Open Ports

Show open ports and applications:

```bash
# List open ports
devcli net open-ports

# Table output
devcli net open-ports --output table
```

## Project Structure

```
devkit/
├── main.go                 # Application entry point
├── cmd/                    # Cobra command definitions
│   ├── root.go            # Root command
│   ├── dev/               # Developer tools
│   │   ├── dev.go         # Dev command group
│   │   ├── uuid.go        # UUID generation
│   │   ├── ulid.go        # ULID generation
│   │   ├── base64.go      # Base64 encode/decode
│   │   ├── jwt.go         # JWT operations
│   │   ├── hash.go        # Hash calculation
│   │   ├── url.go         # URL operations
│   │   ├── html.go        # HTML entity operations
│   │   ├── json.go        # JSON operations
│   │   ├── epoch.go       # Epoch/timestamp conversion
│   │   ├── random.go      # Random data generation
│   │   ├── lorem.go       # Lorem ipsum generator
│   │   ├── cron.go        # Cron expression parser
│   │   ├── semver.go      # Semantic versioning
│   │   └── env.go         # Environment file management
│   ├── file/              # File operations
│   │   ├── file.go        # File command group
│   │   ├── stat.go        # File statistics
│   │   ├── tree.go        # Directory tree
│   │   ├── search.go      # File search
│   │   ├── find-replace.go # Find and replace
│   │   ├── rename.go      # Bulk rename
│   │   ├── convert.go     # Format conversion
│   │   ├── diff.go        # File diff
│   │   ├── dedupe.go      # Duplicate detection
│   │   └── watch.go       # File watching
│   └── net/               # Network & system operations
│       ├── net.go         # Net command group
│       ├── port.go        # Port operations
│       ├── dns.go         # DNS lookup
│       ├── ip.go          # IP information
│       ├── http.go        # HTTP requests
│       ├── ping.go        # Ping
│       ├── ssl.go         # SSL certificate
│       ├── whois.go       # Whois lookup
│       ├── speed.go       # Speed test
│       ├── sysinfo.go     # System information
│       ├── ps.go          # Process management
│       ├── disk.go        # Disk usage
│       ├── interfaces.go  # Network interfaces
│       └── open-ports.go  # Open ports
├── internal/              # Internal packages
│   ├── output/            # Output formatting
│   ├── config/            # Configuration management
│   ├── utils/             # Utility functions
│   └── errors/            # Error handling
└── pkg/                   # Public packages
    └── version/           # Version information
```

## Development

### Prerequisites

- Go 1.22 or higher
- Make (optional, for Makefile commands)

### Building

```bash
# Build binary
make build

# Build for all platforms
make build-all

# Run tests
make test

# Lint code
make lint
```

## License

[Add your license here]
