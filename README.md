# TeamCity Golang API Wrapper

TeamCity API wrapper with some extra tools for TeamCity management.

## Testing

`cp .env-sample .env` to create a `.env` file, and configure to point to your UAT TeamCity instance.

`export $(<.env)` to export these variables in your shell.

Tests do not execute any modifying actions against the TeamCity server.

`go test ./..` to test all packages.
