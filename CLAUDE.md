# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Galaktika is a web-based galaxy game focused on building and competing fleets. The MVP allows players to build fleets using resources and compete against other fleets.

## Project Structure

```
galaktika/
├── cmd/
│   └── server/           # HTTP server entry point
│       └── main.go       # Gin-based web server
├── internal/             # Private application code
│   ├── api/              # API handlers (future)
│   └── game/             # Game logic (future)
├── pkg/                  # Public libraries
│   ├── ship/             # Ship and flight mechanics
│   └── util/             # Generic utility functions
├── assets/               # Static files served to browser
│   ├── js/
│   │   └── game.js       # SVG-based game rendering
│   ├── css/
│   │   └── style.css
│   └── index.html
├── go.mod
└── CLAUDE.md
```

## Architecture

The project follows standard Go web application conventions:

- **cmd/server** - Main application entry point using Gin framework
- **pkg/ship** - Ship and flight mechanics (glaktika.eu/galaktika/pkg/ship)
- **pkg/util** - Generic utility functions (glaktika.eu/galaktika/pkg/util)
- **assets/** - Frontend assets served statically via Gin

The frontend uses vanilla JavaScript with SVG for rendering. The game is turn-based, so rendering occurs on state changes rather than continuous animation.

### Ship Package Architecture

The ship package models space vessels and fleets:

- **Ship** - Individual vessel with position (X, Y), cargo (people, materials), and technology specs
- **ShipTech** - Technology specifications (Attack, Guns, Defense, Speed, CargoCapacity, Mass)
- **Flight** - Collection of ships that move together; flight speed is limited by the slowest ship

Key design pattern: Flight speed calculation demonstrates two approaches - imperative loop (Speed) and functional using util helpers (Speed2).

### Util Package

Provides generic functional programming utilities:
- **ArrayMap** - Transform elements with a function
- **ArrayFilter** - Select elements matching a predicate
- **ArrayReduce** - Reduce elements to a single value with accumulator

These are used throughout the codebase for cleaner functional-style operations.

## Development Commands

### Running the Server

```bash
# Run the web server (serves on http://localhost:8080)
go run cmd/server/main.go

# Build and run
go build -o galaktika cmd/server/main.go
./galaktika
```

### Running Tests

```bash
# Run all tests in the project
go test ./...

# Run tests in a specific package
go test ./pkg/ship
go test ./pkg/util

# Run a specific test
go test ./pkg/ship -run TestFlightSpeed0
```

### Building

```bash
# Build the server binary
go build -o galaktika cmd/server/main.go

# Build with module awareness
go build -mod=readonly -o galaktika cmd/server/main.go
```

### Module Management

```bash
# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Testing Conventions

- Test files follow Go naming: `*_test.go`
- Tests use standard `testing` package
- Test function naming: `Test<Feature><Variation>` (e.g., `TestFlightSpeed0`, `TestFlightSpeed1`)
