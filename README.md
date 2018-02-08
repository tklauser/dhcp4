# dhcp4

[![CircleCI](https://circleci.com/gh/u-root/dhcp4.svg?style=svg)](https://circleci.com/gh/u-root/dhcp4) [![Go Report Card](https://goreportcard.com/badge/github.com/u-root/dhcp4)](https://goreportcard.com/report/github.com/u-root/dhcp4) [![GoDoc](https://godoc.org/github.com/u-root/dhcp4?status.svg)](https://godoc.org/github.com/u-root/dhcp4) [![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/dhcp4/blob/master/LICENSE)

Package `dhcp4` implements IPv4 DHCP packet and option encoding and decoding as
described in RFC 2131, 2132, and 3396.

Option parsing is in the `dhcp4opts` sub-package, and a simple client is
included in `dhcp4client`. Some day, there may be a server.
