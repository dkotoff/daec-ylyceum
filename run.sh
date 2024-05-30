#!/bin/sh

export TIME_DIVISION_MS=$(awk '/^TIME_DIVISION_MS/{print $3}' .config)
export TIME_MULTIPLICATION_M=$(awk '/^TIME_MULTIPLICATION_MS/{print $3}' .config)
export TIME_SUBTRACTION_MS=$(awk '/^TIME_SUBTRACTION_MS/{print $3}' .config)
export TIME_ADDITION_MS=$(awk '/^TIME_ADDITION_MS/{print $3}' .config)

go build -C ./server/cmd  
go build -C ./agent/cmd

./server/cmd/cmd &
./agent/cmd/cmd

wait