# go-dummer-simple

Generates dummy log data - a port of sonots/dummer 's dummer_simple

## Installation

`go get github.com/fujiwara/go-dummer-simple`

## Usage

```
Usage of go-dummer-simple:
  -i string
        Input file (Output messages by reading lines of the file in rotation)
  -m string
        Output message (default "time:2013-11-20 23:39:42 +0900\tlevel:ERROR\tmethod:POST\turi:/api/v1/people\treqtime:3.1983877060667103")
  -o string
        Output file (default "dummy.log")
  -r float
        Number of generating messages per second
  -s int
        Duration of running in second (default 1)
```

## LICENSE

The MIT License (MIT)

Copyright (c) 2016- FUJIWARA Shunichiro
