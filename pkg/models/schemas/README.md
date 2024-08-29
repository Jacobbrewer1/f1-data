# Schemas

In here go all the table schema definitions for the database. This is a good place to define the tables and their
columns, as well as any constraints that should be enforced on the data.

# Generation

The schemas are generated using [GOSCHEMA](https://github.com/Jacobbrewer1/goschema). This is a tool that reads the
schema files and generates the SQL to create. You can use the following command from the root of the repo after adding
GOSCHEMA to your path:

```bash
goschema generate --templates=./pkg/models/templates/*tmpl --out=./pkg/models --sql=./pkg/models/schemas/*.sql
```
