#!/bin/bash

set -e

ps -ef | grep "obs" | grep -v grep | awk '{print $2}' | xargs kill -2
