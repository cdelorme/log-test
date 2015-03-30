
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

This can then be compared by running a simple `wc -l *.log*`.


## use case

The current project for which this test-code was created has a situation in which we need logrotation to prevent filling up the disks

Our application uses basic redirection of stdout/stderr, and logrotate will break this stream unless it has `copytruncate`.

We are forwarding our logs to an external collection and indexing tool using rsyslog.

To read data from a file into rsyslog we are using the `imfile` module.


**We then face two major concerns:**

Unfortunately, the imfile module breaks when copytruncate'd logrotation runs, because it cannot track state.

There is also the fact that logrotate itself claims that a short "time slice" may occur between the copy and truncate operations resulting in a loss of messages.


**Resolutions:**

The options available to us include:

1. temporarily turning off rsyslog for _everything_ in a postrotate operation to restart rsyslog

2. use named pipes to redirect output to a temporary stream which we can control from a small script

3. change our applications to write to syslog directly, or to our external collector

4. use a separate client application to collect and forward logs


Why these all fail to meet our needs:

4. fails in the same way as imfile, logrotate breaks the state of any forwarding applications

3. we cannot control third-party dependencies which also rely on basic redirection for output

2. creates multiple new points of failure in logging for every application using it

1. probably not ideal, but _this is the solution supplied on support sites_


Right now we are going with option **2**, but I propose switching to option **1**.  This code is a test-case to validate the potential loss of messages from logrotate in a totally theoretical (and unlikely) scenario where the _only_ task of this application is to log messages as quickly as possible.

Our goal for our pipelines is distributed operations across cloud, handling upwards of 75,000 messages per second.  My goal is to see what the cut-off per-machine might be, and whether any loss of messages is realistic.


## results

    TODO

