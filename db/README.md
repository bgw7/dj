# PostgreSQL
# Project for managing local postgres DB

1. Uses environment variables supported by psql & PostgreSQL
    - https://www.postgresql.org/docs/current/libpq-envars.html
2. `make sql-init` or `./init.sh` runs ./*.sql in schemas with `PGDATABASE` and `DB_SUPER_USER`