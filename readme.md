
# log test

The intended purpose of this software is to test logging an exceptionally large amount of data against logrotate to verify the time-slice that may result in lost messages when using copytruncate.


## setting up

To start, make sure golang is installed.

This test is dependent on linux for logrotate.  Copy the sample `test-logrotate.conf` file to `/etc/logrotate.d/` and adjust the path to the log file.

To run the test you should temporarily add the `logrotate` script (ex. `/etc/cron.daily/logrotate`) to be called every minute via the system crontab.

Grab the golang dependencies with:

    go get

You can build the raw executable and run it:

    go build
    log-test -t 120 2> test.log

Or you can run the software directly:

    go run main.go -t 120 2> test.log

This will take up a sizable amount of space on disk, as it will write to the log file as quickly as possible.  **That is the softwares only purpose.**

_If you run a line-count against all the log files and it is different from the count the software printed, then logrotate lost messages._  At least, in theory.

The logrotate size configuration value is only considered when logrotate itself is executed via crontab.  Therefore the file may exceed the size limit by a significant amount, unless you force the logrorate operation to occur more regularly.  A more regular interval with smaller sizes may result in less data lost during the copytruncate "time slice", and is worth investigating.


## operation

This software spits out messages onto stderr.

It requires a flag with the number of seconds to run.

It keeps track of the number of messages and prints that to stdout at the end.

This can then be compared by running a simple `wc -l *.log*`.


## use case

A project I am working on has a situation in which we need logrotation to prevent filling up the disks.

Our application uses basic redirection of stdout/stderr, and logrotate will break this stream unless it has `copytruncate`.

We are forwarding our logs to an external collection and indexing tool using rsyslog.

To read data from our applications log file(s) into rsyslog we are using the `imfile` module.  _Some of the applications are third party tools._


**We encountered major conflicts:**

The `imfile` module breaks and ceases to forward when logrotate runs with `copytruncate`.  Presumably because it cannot track state anymore.

The logrotate documentation itself claims that a short "time slice" may occur between the copy and truncate operations which may result in data loss.


**Limited Options:**

We found four possible resolutions around the web for this:

1. temporarily turning off rsyslog for _everything_ in a postrotate operation to restart rsyslog

2. use named pipes to redirect output to a temporary stream which we can control from a small script

3. change our applications to write to syslog directly, or to our external collector

4. use a separate client application to collect and forward logs


None of these truly meet our needs, and here is why:

1. Stopping rsyslog means temporarily loosing any logs going through syslog itself.

2. named pipes add additional points of failure and complexity to what should otherwise be simple output handling.

3. We cannot change all the applications we depend on which also do not directly use syslog.

4. A third party client very likely suffers from the same conflicts as imfile does with logrotation, and we'd rather use the onboard tools (for not).


We are going with the first option for now, and this project is intended to either proove or disprove whether or not the time-slice loss is a real concern.

