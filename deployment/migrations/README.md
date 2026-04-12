# Database Migrations

SQL migrations for the portfolio database schema.

## Naming Convention

```
{number}_{description}.up.sql    - Apply migration
{number}_{description}.down.sql  - Rollback migration
```

Example:
- `000001_initial_schema.up.sql`
- `000001_initial_schema.down.sql`

## Running Migrations

### Option 1: Using golang-migrate CLI

Install:
```bash
# Mac
brew install golang-migrate

# Windows (with scoop)
scoop install migrate
```

Run against local database:
```bash
migrate -path deployment/migrations -database "mysql://root:devpassword@tcp(localhost:3307)/portfolio" up
```

Run against Cloud SQL:
```bash
migrate -path deployment/migrations -database "mysql://root:PASSWORD@tcp(CLOUD_SQL_IP)/portfolio" up
```

Rollback:
```bash
migrate -path deployment/migrations -database "mysql://..." down 1
```

### Option 2: Manual

Connect to MySQL and run the SQL files directly:
```bash
mysql -h localhost -P 3307 -u root -p portfolio < deployment/migrations/000001_initial_schema.up.sql
```

## Current Migrations

| Number | Description | Status |
|--------|-------------|--------|
| 000001 | Initial schema (art_tiles, categories, art_categories) | Ready |

## Adding New Migrations

1. Create two files with the next number:
   - `000002_your_description.up.sql`
   - `000002_your_description.down.sql`

2. The `.up.sql` applies changes
3. The `.down.sql` reverses them

4. Test locally before running in production!
