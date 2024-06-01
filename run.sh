#!/bin/sh

export TIME_DIVISION_MS=$(awk '/^TIME_DIVISION_MS/{print $3}' .config)
export TIME_MULTIPLICATION_M=$(awk '/^TIME_MULTIPLICATION_MS/{print $3}' .config)
export TIME_SUBTRACTION_MS=$(awk '/^TIME_SUBTRACTION_MS/{print $3}' .config)
export TIME_ADDITION_MS=$(awk '/^TIME_ADDITION_MS/{print $3}' .config)
export COMPUTING_POWER=$(awk '/^COMPUTING_POWER/{print $3}' .config)
export SERVER_PORT=$(awk '/^SERVER_PORT/{print $3}' .config)

go build -C ./server/cmd  
go build -C ./agent/cmd

cd ./server
go build ./cmd  
./cmd/cmd &
pid1=$!

cd ../agent
go build ./cmd
./cmd/cmd &
pid2=$!

cd ..

cleanup() {
    echo "Cleaning up"
    kill $pid1
    kill $pid2
    rm ./server/cmd/cmd
    rm ./agent/cmd/cmd
    wait
    exit 0
}

trap cleanup INT

wait