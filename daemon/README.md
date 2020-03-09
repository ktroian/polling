This is a daemon (micro)service, it's main goal is to poll the rest API and put its result into a database.
Currently, only postgres is supported.

Building:

``` go build main.go db.go poll.go company.go```

CLI arguments:
```
  -host string
    	hostname for the service API (default "localhost")
  -port int
    	a port of the service API (default 8080)
  -interval int
    	an interval to check the API (default 30)
  -log string
    	path to a file to .log file to store logs
      if none logs will be written to stdout
  -password string
    	a password for the database
      (Yes, this is not secure)
  -username string
    	username to log into the database
```

Known bugs:
  - Program is panicing in case of unsuccessful connection
  - Also, it may panic on invalid arguments (validation is not implemented)
  - Program will not delete data that it didn't put in database during this session
  - Daemonizing is not fare
  
Usage notes:
  - Program will create a new table with name 'companies'
    in case it is already exists, data will not be removed
  - Postgres Database is expected to run on default PostgreSQL port - 5432
    Unfortunately, config files are not implemented yet

Usage:
  Program can be used in 2 modes, as a separate process or as a deamon (to be fare, just detached process for now)
  
  As a process:
 
  ``` ./main [host|port|interval|log|username|password]```
  
  As a daemon:
  
  ``` ./main [host|port|interval|log|username|password] & disown```
  
  Alternatively, other external tools may be used to start a daemon:
  
  nohup
  
  ``` nohup ./main [args...]```
  
  daemonize
  
  ``` daemonize -p /var/run/main.pid -l /var/lock/subsys/main -u nobody ./main [args...]```
  
TODO:
  - Implement starting daemon from the process
  - Fixes of crushes
  - Test coverage
