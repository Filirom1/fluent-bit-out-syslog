# Output Syslog Plugin for Fluent-Bit

Work In Progress: do not use yet

## Options

**Network:** If network is empty (default), it will connect to the local syslog server. Otherwise tcp and udp is supported.

**Address:** Remote syslog address. Format: IP:PORT or HOSTNAME:PORT. Only use if network is not empty

**Severity:** Default info. Syslog severity as defined by RFC 3164:
* emerg 
* alert 
* crit 
* err 
* warning
* notice
* info
* debug

**Facility:** Default local0. Syslog facility as defined by RFC 3164:
* kern
* user
* mail
* daemon
* auth
* syslog
* lpr
* news
* uucp
* cron
* authpriv
* ftp
* local0
* local1
* local2
* local3
* local4
* local5
* local6
* local7

**Tag:** By default use tag provided by FluentBit.

## Configuration example

```
[OUTPUT]
    Name       syslog
    Match      *
    Network    udp
    Address    localhost:514
    Severity   info
    Facility   ftp
    Tag        my-tag
```

## Build

```
$ make
go build -buildmode=c-shared -o out_syslog.so .
```

## Usage

```
$ td-agent-bit -v -e ./out_syslog.so -i cpu -o syslog
```

## Version

Tested with [fluent-bit-go v0.12](https://github.com/fluent/fluent-bit-go/tree/api-0.12)

