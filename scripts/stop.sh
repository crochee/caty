#!/bin/bash

set -e

ps -ef | grep "obs" | grep -v grep | awk '{print $2}' | xargs kill -2
#ps -aux | grep "caty" | grep -v grep | awk '{print $2}' | xargs kill -2
#pgrep caty | xargs kill -2
