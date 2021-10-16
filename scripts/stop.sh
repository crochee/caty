#!/bin/bash

set -e

ps -ef | grep "obs" | grep -v grep | awk '{print $2}' | xargs kill -2
#ps -aux | grep "cca" | grep -v grep | awk '{print $2}' | xargs kill -2
#pgrep cca | xargs kill -2
