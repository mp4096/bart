# Bart

[![Build Status](https://travis-ci.org/mp4096/bart.svg?branch=master)](https://travis-ci.org/mp4096/bart)
[![Go Report Card](https://goreportcard.com/badge/github.com/mp4096/bart)](https://goreportcard.com/report/github.com/mp4096/bart)


## Security

`bart` uses the `net/smtp` package,
[which uses TLS if possible](https://golang.org/pkg/net/smtp/#SendMail).
Still, just to be safe, I explicitly discourage using `bart` for mission-critical information.

## Installation

### From source

```
$ make install
```

### Binaries

You can get them from [GitHub releases](https://github.com/mp4096/bart/releases).


## Usage example

You want to send an HTML email stored as a [Mustache template](https://mustache.github.io/) in `template.mustache`.
Your configuration (see below) is defined in `config.bart.yml`.

To preview the rendered emails, call `bart` without the send `-s` flag:

```
$ bart -t template.mustache -c config.bart.yml
Hello, Jane Doe
Send flag not set: opening preview in "chromium-browser"
Send flag not set: opening preview in "chromium-browser"
```

To send the email, add the send flag `-s`:

```
$ bart -t template.mustache -c config.bart.yml -s
Hello, Jane Doe
Please enter your credentials for "smtpserver.xyz.com"
Login: janedoe
Password:
Will send to [john.doe@abc.com jane.doe@xyz.com]
Will send to [max.mustermann@def.de jane.doe@xyz.com]
```

Since this email will not appear in your provider's `Sent` folder,
`bart` will send you a BCC copy.

For help, call `bart -h`.

## How to configure

Here's an example config file:

```yaml
author:
  name: Jane Doe
  email: jane.doe@xyz.com
  browser: chromium-browser

email_server:
    hostname: smtpserver.xyz.com
    port: 123

global_context:
  subject: Global subject

recipients:
  john.doe@abc.com:
    salutation: Hi John
    subject: Local subject, overrides global one
  max.mustermann@def.de:
    salutation: Hello Max
```

Basically, `recipients` is a hashmap keyed by recipient email addresses;
the values are local contexts specific for each recipient.
Notice that local context overrides global one!
