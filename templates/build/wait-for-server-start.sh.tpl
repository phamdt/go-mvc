#!/bin/bash

echo "Waiting for servers to start..."
attempts=1
while true; do
  docker exec -i {{Name}} curl -f http://localhost:8080/health > /dev/null 2> /dev/null
  if [ $? = 0 ]; then
    echo "Service started"
    break
  fi
  ((attempts++))
  if [[ $attempts == 5 ]]; then
    echo "5 attempts to check health failed"
    break
  fi
  sleep 10
  echo $attempts
done