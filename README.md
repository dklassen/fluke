# [STATUS::Pre-Alpha] - Still messing around

## ClickHouse Schema Migrations

At some point, every data engineer thinks: “Why don’t we just write a SQL parser?” It’s a rite of passage—and a sign you haven’t quite hit Staff Engineer yet. So let’s grab the ClickHouse grammar files and make some bold, probably bad, but definitely fun decisions.

This project is one of those experiments: a tool that generates ClickHouse SQL migration files by comparing two schemas. No more hand-crafting scripts—just feed in two schema versions and let the tool detect changes and spit out the SQL.


### 💡 Goals and Requirements

The goal: Compare two ClickHouse SQL schemas (i.e. sets of CREATE statements) and generate migration files containing the SQL needed to move from one schema to another.

The tool covers creating, dropping, and altering tables/columns. It also deals with engine changes and structural clause updates using configurable strategies.

We’ve also got environment-specific quirks—like different engine setups in dev vs prod—so we want a reliable way to reflect and manage those variations in generated migrations.


### 📝 Supported SQL Parsing

✅ CREATE TABLE
 - Detects new tables and generates CREATE TABLE statements accordingly.

✅ ALTER TABLE
- Emits statements to:
- Add new columns
- Drop columns
- Modify existing columns

### 🚩 Structural Change Strategies

Certain structural changes can’t be applied with ALTER and require a full table rebuild. These include modifications to:
	•	ORDER BY
	•	PRIMARY KEY
	•	PARTITION BY

When these are detected, the tool follows one of two strategies:

1.	Nuclear — Just blow it up:
Emits DROP TABLE + CREATE TABLE.
	
2.	Sane — Be careful:
Logs a warning and suggests potential non-destructive migration paths.


## 🗺️ **Schema Parsing**

The tool accurately, or will when finished, parses and capture the following information from each table schema:

- **Table Name**
- **Column Names and Types**
- **Engine Type**
- **Structural Clauses:**
  - `ORDER BY`
  - `PRIMARY KEY`
  - `PARTITION BY`

This data is represented using Go structs, making it easy to compare schemas and generate migration statements.


## 💡 **Example Usage**

Given the following old schema:
```sql
CREATE TABLE users (
    id UInt32,
    name String
) ENGINE = MergeTree()
ORDER BY id;
```

And a new schema:
```sql
CREATE TABLE users (
    id UInt32,
    name String,
    age Int32
) ENGINE = MergeTree()
ORDER BY (id, name);
```

The migration file generated could be in destructive mode:
```sql
DROP TABLE IF EXISTS users;
CREATE TABLE users (id UInt32, name String, age Int32) ENGINE = MergeTree() ORDER BY (id, name);
```
