# scripts

A Go mono-repo for small CLI tools and scripts.

## Scripts

| Script                | Description                                                                 |
| --------------------- | --------------------------------------------------------------------------- |
| `npm-min-release-age` | Filters `npm outdated` by release age, suggesting the latest "safe" version |
| `script-template`     | Placeholder / example for adding new scripts                                |

## Requirements

- [Go](https://go.dev/dl/) 1.26+

## Project Structure

```
.
├── cmd/                      # One directory per CLI binary
│   ├── npm-min-release-age/
│   │   └── main.go
│   └── script-template/
│       └── main.go
├── internal/                 # Private packages (only for this repo)
│   └── npm/
│       └── npm.go
├── pkg/                      # Public packages (reusable by external projects)
├── scripts/                  # Build / helper scripts
│   └── build.sh
├── go.mod                    # Single root module
└── README.md
```

## Build

### Build a single script

```bash
go build -o bin/npm-min-release-age ./cmd/npm-min-release-age
```

### Build all scripts

```bash
./scripts/build.sh        # outputs to ./bin/
./scripts/build.sh ~/bin  # outputs to ~/bin/
```

## Install

```bash
# Build everything into ~/bin
./scripts/build.sh ~/bin
export PATH="$HOME/bin:$PATH"
```

## Usage

### npm-min-release-age

```bash
# Default threshold: 7 days
npm-min-release-age

# Custom threshold: 14 days
npm-min-release-age 14

# Example output
✅ express: current=4.18.0 latest=5.2.1 (271 days old) SAFE
⏳ lodash: current=4.17.20 latest=4.18.1 (150 days old) NOT SAFE | latest safe: 4.17.21 (2016 days old)
```

- `✅` means the latest version is old enough to update.
- `⏳` means the latest version is too new. The "latest safe" version is the most recent one that meets your age threshold.

## Adding a New Script

1. Create a new directory under `cmd/`:

   ```bash
   mkdir cmd/my-new-script
   ```

2. Add a `main.go`:

   ```go
   package main
   import "fmt"
   func main() { fmt.Println("Hello!") }
   ```

3. Build it:
   ```bash
   go build -o bin/my-new-script ./cmd/my-new-script
   # or
   ./scripts/build.sh
   ```

## Shared Code

- Put **repo-private** packages in `internal/` (e.g., `internal/npm`).
- Put **reusable** packages in `pkg/` if you want them importable by other projects.

## Cross-Compile

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/npm-min-release-age-darwin ./cmd/npm-min-release-age

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/npm-min-release-age.exe ./cmd/npm-min-release-age
```
