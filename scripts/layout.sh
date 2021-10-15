#!/bin/bash

set -ex

awk '$1 == "module" {print $2}' ./go.mod | xargs fieldalignment ./pkg/service/account/account.go