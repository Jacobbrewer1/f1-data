# Data API

## Description

The Data API is a RESTful API that provides access to the F1 data stored in the database. The API allows users to query
the data using various filters and retrieve the results in JSON format; you are able to check out the documentation for
the API [here](https://github.com/Jacobbrewer1/f1-data-docs).

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

## Running the app

You can check out all the available commands by running the following command:

```shell
go build -o data && ./data --help
```

For each command you can check out each flag by running the following command:

```shell
./data <command> --help
```
