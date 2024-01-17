#!/usr/bin/env bash
#
# Install dependencies

go mod tidy
brew install golangci-lint
brew upgrade golangci-lint
brew install pre-commit
pre-commit install
pre-commit install --hook-type commit-msg
