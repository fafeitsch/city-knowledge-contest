#!/bin/bash

./node_modules/@mockoon/cli/bin/run start --data ./nominatim-mock.json > /dev/null&
MOCKOON_PID=$!
trap 'kill ${MOCKOON_PID}; exit' INT QUIT TERM EXIT HUP
#waiting for mockoon to startup
sleep 3
playwright test

