# SIPq

[![Build Status](https://travis-ci.org/henryscala/sipq.svg?branch=master)](https://travis-ci.org/henryscala/sipq)

SIPq is a test tool for SIP protocol. Its ambition is to be a next generation of [SIPp](http://sipp.sourceforge.net/). There is for sure a long way to go. The primary purpose of it for now is for study, on how to design a SIP stack and test tool. If you want to do something serious, go for [SIPp](http://sipp.sourceforge.net/). 

SIPq is designed in mind as a SIP test tool, but components of it are also able to serve as a SIP stack. It should be practical to make it into use to construct a SIP application, since the source codes are organized in [golang](https://golang.org/) packages. 

Welcome any contributions! 

# References

[SIPp](http://sipp.sourceforge.net/) 

[SIP](https://tools.ietf.org/html/rfc3261) 

[Internet Message Format](https://tools.ietf.org/html/rfc2822)

[ABNF](https://tools.ietf.org/html/rfc5234)

[HTTP/1.1](https://tools.ietf.org/html/rfc2616)

[URI Gneric Syntax](https://tools.ietf.org/html/rfc2396)

# Dependence

[A golang js implementation,otto,](https://github.com/robertkrimen/otto) is used to write scenario script. 

# Build 

- setup golang and set a GOPATH environment variable 
- run the following commands

```
go get -u github.com/robertkrimen/otto
go get -u github.com/henryscala/sipq
cd $GOPATH/src/github.com/henryscala/sipq

go test ./...
go build
```

# Contributors 

[Tao Keqin](https://github.com/taokeqin)


