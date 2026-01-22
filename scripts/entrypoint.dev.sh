#!/bin/bash
set -e

echo "#################### starting deamon"
CompileDaemon --build="go build -o bin/main main.go" --command=./bin/main
