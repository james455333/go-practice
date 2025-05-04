# Overview

practice project for golang

# 📁 Project Structure

```
.
├── README.md
├── go.mod
└── go_routine
    └── main.go
```

# Contents

- [go-routine-practice](#go-routine-practice)
- [go-kit-example](#go-kit-practice)


## go-routine-practice

### Reference

- [how-to-use-goroutines-the-right-way-basis-to-advance-explained-simply](https://medium.com/@adityasinghrathore360/how-to-use-goroutines-the-right-way-basis-to-advance-explained-simply-9962f7c1b9e8)

## go-kit-practice

### Reference
- [go-kit example](https://gokit.io/examples/)

### stringsvc3

structure
```text
+--------------------------------------------------------------------+
|                          Transport Layer                          |
+--------------------------------------------------------------------+
| - Receives HTTP requests (/uppercase, /count)                     |
| - Uses go-kit's httptransport.NewServer                           |
| - Parses JSON into struct, passes to the endpoint layer           |
+--------------------------------------------------------------------+
                                 |
                                 v
+--------------------------------------------------------------------+
|                           Endpoint Layer                          |
+--------------------------------------------------------------------+
| - Defines makeUppercaseEndpoint / makeCountEndpoint               |
| - Converts request struct into method args                        |
| - Calls StringService.UpperCase / Count                           |
| - Wraps output and error into response struct                     |
+--------------------------------------------------------------------+
                                 |
                                 v
+--------------------------------------------------------------------+
|                          Middleware Layer                         |
+--------------------------------------------------------------------+
| - loggingMiddleware:       logs method, input, output, error, time|
| - instrumentingMiddleware: exposes Prometheus metrics             |
| - proxyingMiddleware:      optionally forwards certain requests   |
+-------------------------------+------------------------------------+
                                |
             +------------------+------------------+
             |                                     |
             v                                     v
+-------------------------------------+   +-------------------------------------+
|        Proxying Path (Optional)     |   |      Local Service Path (Default)   |
+-------------------------------------+   +-------------------------------------+
| - Forwards selected methods (e.g.,  |   | - Handles Count and other fallback  |
|   UpperCase) to remote instances    |   |   requests                          |
| - Uses endpoint.NewClient to send   |   | - Calls stringService directly      |
|   request to endpoint               |   |   (UpperCase, Count implementation) |
+-------------------------------------+   +-------------------------------------+
```


test commend
```shell
$ cd ./go_kit_example/stringsvc 
$ go build ./
$ ./stringsvc -listen :8081 &
$ ./stringsvc -listen :8082 &
$ ./stringsvc -listen :8083 &
$ ./stringsvc -listen :8080 -proxy=localhost:8081,localhost:8082,localhost:8083 &
$ sh test.sh
$ kill %1 %2 %3 %4
```

expected output
```text
$ cd ./go_kit_example/stringsvc 
$ go build ./
$ ./stringsvc -listen :8081 &
listen=:8081 caller=proxying.go:23 proxy_to=none
listen=:8081 caller=stringsvc.go:71 msg=HTTP addr=:8081
$ ./stringsvc -listen :8082 &
listen=:8082 caller=proxying.go:23 proxy_to=none
listen=:8082 caller=stringsvc.go:71 msg=HTTP addr=:8082
$ ./stringsvc -listen :8083 &
listen=:8083 caller=proxying.go:23 proxy_to=none
listen=:8083 caller=stringsvc.go:71 msg=HTTP addr=:8083
$ ./stringsvc -listen :8080 -proxy=localhost:8081,localhost:8082,localhost:8083 &
listen=:8080 caller=proxying.go:42 proxy_to="[localhost:8081 localhost:8082 localhost:8083]"
listen=:8080 caller=stringsvc.go:71 msg=HTTP addr=:8080
$ sh test.sh
listen=:8081 caller=middleware.go:22 method=uppercase input=foo output=FOO err=null took=559ns
listen=:8080 caller=middleware.go:22 method=uppercase input=foo output=FOO err=null took=886.081µs
{"v":"FOO"}
listen=:8082 caller=middleware.go:22 method=uppercase input=bar output=BAR err=null took=375ns
listen=:8080 caller=middleware.go:22 method=uppercase input=bar output=BAR err=null took=440.918µs
{"v":"BAR"}
listen=:8083 caller=middleware.go:22 method=uppercase input=baz output=BAZ err=null took=486ns
listen=:8080 caller=middleware.go:22 method=uppercase input=baz output=BAZ err=null took=601.592µs
{"v":"BAZ"}
$ kill %1 %2 %3 %4
```
