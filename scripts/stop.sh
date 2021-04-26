#!/bin/bash

set -e

ps -ef | grep "obs" | grep -v grep | awk '{print $2}' | xargs kill -2
#ps -aux | grep "obs" | grep -v grep | awk '{print $2}' | xargs kill -2
#pgrep obs | xargs kill -2
