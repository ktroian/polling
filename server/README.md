This is a simple (micro)service, it's main goal is to serve .csv data via simple REST API.

Building:

``` go build main.go writer.go api.go company.go```

CLI arguments:
```
  -port int
    	a port of the server API (default 8080)
  -file string
      path to the file.csv
```

Usage:
 
  ``` ./main [port|file.csv]```
  
  
TODO:
  - New router for better parallelism
  - Test coverage
  - Refactor
