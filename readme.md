
# log test

The intended purpose of this software is to test logging an exceptionally large amount of data against logrotate to verify the time-slice that may result in lost messages when using copytruncate.


## setting up

To start, make sure golang is installed.

This test is dependent on linux for logrotate.  Copy the sample `test-logrotate.conf` file to `/etc/logrotate.d/` and adjust the path to the log file.

Grab the golang dependencies with:

    go get

You can build the raw executable and run it:

    go build
    log-test -t 30 2> test.log

Or you can run the software directly:

    go run main.go -t 30 2> test.log

_If you run a line-count against all the log files and it is different from the count the software printed, then logrotate lost messages._  At least, in theory.

It probably would make sense to experiment with different cutoff sizes, the default size is very small (300k) and may take less time to copy, therefore reducing the time-slice between copying and emptying the log file.


## operation

This software spits out messages onto stderr.

It requires a flag with the number of seconds to run.

It keeps track of the number of messages and prints that to stdout at the end.

This can then be compared by running a simple `wc -l *.log`.

