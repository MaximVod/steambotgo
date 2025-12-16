# SteamBotGo

A Go application for interacting with Steam store APIs.

## Features
- Search for games on Steam
- Get pricing information for games across different regions
- Clean architecture with separation of concerns

## Architecture
This project follows a clean architecture pattern with:
- **cmd**: Entry point of the application
- **internal/adapters**: External service adapters (Steam API)
- **internal/entities**: Domain entities
- **internal/interfaces**: Port interfaces
- **internal/usecases**: Business logic

## Setup
```bash
go mod tidy
```

## License
MIT
