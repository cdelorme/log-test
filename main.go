package main

import (
    "fmt"
    "time"
    "os"

    "github.com/cdelorme/go-log"
    "github.com/cdelorme/go-option"
    "github.com/cdelorme/go-maps"
)

func main() {

    // establish logger
    logger := log.Logger{Level:log.Debug}

    // register options
    appOptions := option.App{Description: "log-as-fast-as-you-can"}
    appOptions.Flag("timeout", "timeout (in seconds) to print log messages", "-t", "--time")
    appOptions.Example("log-test -t 5 2>logs/test.log")
    flags := appOptions.Parse()

    // acquire timeout
    timeout, _ := maps.Int(&flags, 0, "timeout")
    if (timeout <= 0) {
        logger.Error("please supply a timeout...")
        os.Exit(1)
    }

    // controls
    after := time.After(time.Duration(timeout) * time.Second)
    do := true

    // counter
    count := 0

    // loop until timeout
    for do {
        select {
        case <-after:
                do = false
        default:
            logger.Debug("a json blob: %+v", appOptions)
            count = count + 1
        }
    }

    // print total to stdout
    fmt.Printf("Total Messages: %d\n", count)
}
