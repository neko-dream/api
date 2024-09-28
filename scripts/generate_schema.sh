#!/bin/bash

dbml2sql ./db/schema.dbml --postgres > ./internal/infrastructure/db/schema.sql
sqlc generate
