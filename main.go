package main

import (
    "fmt"
    "time"

    "github.com/cdelorme/go-log"
)

func main() {

    // establish logger
    logger := log.Logger{Level:log.Debug}

    // acquire timeout
    timeout := 2

    // controls
    after := time.After(time.Duration(timeout) * time.Second)
    do := true

    // counter
    count := 0

    for do {
        select {
        case <-after:
                do = false
        default:
            logger.Debug("message")
            count = count + 1
        }
    }

    fmt.Printf("Total Messages: %d\n", count)
}
