## Description
This is a simple scanning application to detect sensitive keywords in public Git repositories.

## Packages
The server contains several Go packages:
* `api` - API handlers
* `app` - Application layer that contains business logic
* `docs` - Generated Swagger API documents
* `jobs` - Job server
* `migration` - Database migration files
* `model` - Definition and method for data model
* `store` - Storage layer
* `util` - Utility functions

## Installation

* Cloning the repository and build
```
git clone https://github.com/sonda2208/grc.git
cd grc
go build
```
* Place `config.json` file in the same directory of the executable file or set `CONFIG` environment variable with path to your configuration file.
* Sample `config.json` file content:
```json
{
  "SQLSetting": {
    "DataSource": "user=postgres password= database=gr host=localhost sslmode=disable",
    "MaxIdleConnections": 8,
    "MaxOpenConnections": 8
  }
}
```
* Then run
```
./grc
```