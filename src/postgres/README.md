# Migrations with SQLAlchemy + Alembic

While terraform handles the allocation of a Cloud SQL instance, there comes the issue of creating or migrating tables inside the newly created PostgreSQL database.

For this purpose, SQLAlchemy is used as an ORM and Alembic for migrations.

Inspiration was taken from [this medium post](https://benchling.engineering/move-fast-and-migrate-things-how-we-automated-migrations-in-postgres-d60aba0fc3d4).

## Set environment variables

Before running alembic, several environment variables must be set to authenticate to the PostgreSQL server:
* `POSTGRES_USERNAME` Username of the postgres service account.  Value is stored as a GCP Secret named`varity-postgres-username`.
* `POSTGRES_PASSWORD` Password of the postgres service account.  Value is stored as a GCP Secret named`varity-postgres-password`.
* `POSTGRES_HOST` IP of the postgres server.  This should be an output in Terraform named `postgres_ip`.
* `POSTGRES_DB` The postgres database name.  Currently, this should be `finance`.

## Running a migration

1. Edit or create a SQLAlchemy model in the `models` directory.
2. Autogenerate an Alembic revision.
```bash
alembic revision --autogenerate -m "My revision message"
```
3. Run the migration.
```bash
alembic upgrade head 
```

## TODO

* Validate schema changes with the SQLAlchemy Inspection API
