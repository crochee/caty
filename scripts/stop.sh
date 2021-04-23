#!/bin/bash

set -e

process_id=$(ps -ef | grep "obs" | grep -v grep | awk '{print $2}')

if [ -z "${process_id}" ]; then
  echo "process id is null."
else
  echo pid is ${process_id}
  kill -2 ${process_id}
fi