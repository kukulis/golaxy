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

### Frontend Entity Classes Architecture

The JavaScript entity classes in `assets/js/entities/` combine both data storage and presentation logic:

- **Battle** - Manages battle state, shots, and fleets; includes ship lookup methods
- **Fleet** - Contains ships array and owner; maintains ship lookup map for efficient access
- **Ship** - Ship data with tech specs; includes SVG rendering (`creteShipSvg`) and click handling
- **Shot** - Shot data (source, destination, result); includes SVG building (`buildSvg`) for visualization

**Key Architecture Decision**: These are NOT pure DTOs. Each entity class serves as a view-model that combines:
1. **Backend data** - Properties received from API (id, name, tech, destroyed, etc.)
2. **Rendering properties** - Display-specific fields (battleX, battleY, svgElement)
3. **Presentation methods** - SVG generation, event handling, and drawing logic

**Why This Works**: Since these classes are specifically designed for game rendering, combining data with SVG generation logic keeps related functionality together. This is appropriate for a game UI where entities are responsible for both their state and visual representation.

**Data Update Pattern**: Each entity class includes an `updateFromDTO(data)` method that:
- Updates properties from backend data
- Returns `this` for method chaining
- Creates nested entity instances (e.g., Fleet creates Ships, Battle creates Fleets and Shots)
- Maintains relationships (e.g., Battle.fixShotsReferences links shots to ship instances)

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

### Frontend Visual Tests

- `assets/test_ship_designs.html` - Visual test for ship rendering. This is for manual visual verification only, not automated testing.

### Test-Driven Development (TDD)

**Project Decision**: Follow TDD principles where possible - write tests first before implementation.

**Workflow**:
1. Write the test first (which will fail)
2. Implement the minimum code to make the test pass
3. Refactor if needed while keeping tests green

**Benefits**:
- Tests drive the API design
- Ensures all code has test coverage
- Catches bugs early in development
- Makes refactoring safer

**Note**: While TDD is preferred, it may not be practical for all scenarios (e.g., exploratory prototyping, UI work). Use judgment, but default to test-first when feasible.

### Test Readability

- **Use "nice" numbers**: Choose test values that produce clean, easy-to-verify results. For example, use mass=64 (sqrt=8) instead of mass=60 (sqrt≈7.745) to get integer results.
- **Avoid obvious comments**: Do not comment self-explanatory code like "// When: Calculating ship tech" before a `CalculateShipTech()` call. Let the code speak for itself.
- **Keep tests concise**: Prefer minimal, focused test cases over verbose ones with excessive setup documentation.
