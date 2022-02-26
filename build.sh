#!/bin/bash

# compile graphql schema from files
cat internal/graph/schema/*.gql > internal/graph/schema.graphql

# build go code from graph
go generate ./internal/graph