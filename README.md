# Bart

[![Build Status](https://travis-ci.org/mp4096/bart.svg?branch=master)](https://travis-ci.org/mp4096/bart)
[![Go Report Card](https://goreportcard.com/badge/github.com/mp4096/bart)](https://goreportcard.com/report/github.com/mp4096/bart)


## Security

`bart` uses the `net/smtp` package,
[which uses TLS if possible](https://golang.org/pkg/net/smtp/#SendMail).
Still, just to be safe, I explicitly discourage using `bart` for mission-critical information.

## Installation

### From source on Linux

```
$ make install
```

### Binaries

You can get them from GitHub releases.
