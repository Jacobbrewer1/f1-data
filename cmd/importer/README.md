# Importer

## Description

The Importer app is a simple command line tool that scrapes the F1 archive and imports the data into a database.

## Setup

### Configuration

The app requires a configuration file to be supplied. The configuration file is a JSON file with the following
structure:

```json
{
  "vault": {
    "address": "https://example.com:8200",
    "auth": {
      "username": "example-username",
      "password": "example-password"
    },
    "database": {
      "path": "example/database/creds/example-role"
    }
  },
  "database": {
    "host": "localhost:5432",
    "schema": "f1data"
  },
  "importer": {
    "f1_base_url": "https://www.formula1.com",
    "from_year": 1950,
    "to_year": -1
  }
}
```

* The `vault` section contains the configuration for the Vault server. The `address` field is the address of the Vault
  server.
    * The `auth` section contains the credentials for the Vault server. The `username` and `password` fields are the
      credentials for the Vault server.
    * The `database` section contains the path to the database credentials in the Vault server. The `path` field is the
      path to the database credentials in the Vault server.
* The `database` section contains the configuration for the database.
    * The `host` field is the address of the database server. (address and port separated by a colon)
    * The `schema` field is the schema to use in the database.
* The `importer` section contains the configuration for the importer.
    * The `f1_base_url` field is the base URL for the F1 archive.
    * The `from_year` field is the year to start importing data from.
    * The `to_year` field is the year to stop importing data at. If the value is -1, the importer will import data
      up to the current year.

## Running the app

You can check out all the available commands by running the following command:

```shell
go build -o importer && ./importer --help
```

For each command you can check out each flag by running the following command:

```shell
./importer <command> --help
```
