#!/usr/bin/env bash
# Runs every reference implementation against vector.json inside docker,
# so the snippets shown on the documentation site are known to work as pasted.
set -euo pipefail
cd "$(dirname "$0")"

run() {
  local name=$1 image=$2
  shift 2
  echo "== $name"
  docker run --rm -v "$PWD:/verify" -w "/verify/$name" "$image" "$@"
}

run go     golang:1.26-alpine            go run .
run node   node:22-alpine                node verify.js
run php    php:8.3-cli-alpine            php verify.php
run python python:3.12-alpine            python3 verify.py
run java   eclipse-temurin:21-jdk-alpine java Verify.java
echo "all implementations agree"
