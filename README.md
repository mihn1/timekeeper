# TimeKeeper

- A free and open source time tracking application in the cooking
- At the moment TimeKeeper is an application usage tracker that monitors and categorizes your time spent in different applications and browser tabs.

## Features
### Implemented for now
- Tracks application switching and usage time
- Special browser tab tracking for web browsers (Chrome, Safari, and Brave)
- Categorizes usage into customizable categories
- Stores data in SQLite database or in-memory
### To be implemented
- Nice and simple GUI
- Sessions
- Goals
- Windows support
- Account then multiple devices support with data syncing

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/mihn1/timekeeper.git
cd timekeeper

# Build the project
go build -o timekeeper ./cmd/

# Run with SQLite storage
./timekeeper --db sqlite --dbpath ./db/timekeeper.db

# Use in-memory storage
./timekeeper --db inmem