#!/bin/bash

dbml2sql ./db/schema.dbml --postgresql > ./internal/infrastructure/db/schema.sql
