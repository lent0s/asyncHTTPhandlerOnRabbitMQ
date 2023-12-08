# asyncHTTPhandlerOnRabbitMQ
(Asynchronous processing engine HTTP requests on RabbitMQ) <a href="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/README.md">
                                                                   <img width="50" alt="switch to Russian" src="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/ru.jpg?raw=true">
                                                                 </a>

The project was released as a demonstration of the interaction between microservices and the RabbitMQ message broker.  
  
![Schema](https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/scheme.jpg?raw=true)  
Schematic diagram.

---

## Setup and launch
Setting up and running a project is done in several ways:
1. When compiling:
   1. [Api](#Api)
   2. [Config File](#Config-File)
2. [Executable File](#Executable-File)

### Api
To launch the client, just add to import
 ```Go
import (
    async "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/cmd"
)
 ```
call a function and pass configuration data to it
 ```Go
    async.RunService(LogMaxFileSizeKB uint16, LogPathFolder, LogFileName, RConnect, ServerHost, ServerPort string)
 ```
then send requests with tasks to ``http://ServerHost:ServerPort/cmd`` and receive responses.

### Config File
Fill [config file](#cоnfig-file) and start the project:
```
    go run async.go
```
or with import
 ```Go
import (
    async "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/cmd"
)
 ```
call
 ```Go
    async.Local()
 ```
then send requests with tasks to ``http://ServerHost:ServerPort/cmd`` and receive responses.

---

### Executable File
Compile the project with
```
    go build async.go
```
Start ```async``` with keys (or ```-help``` for reference):
- lpf - path to the directory where log files are stored
- lfn - log file name (without extension)
- lms - approximate log file size threshold in KB, after which a new file will be created
- rmq - address and authentication data to access to RabbitMQ
- rto - RabbitMQ server timeout
- sip - IP-address receiving requests
- spn - port receiving requests

Example: ```local -spn 9090``` - to run an application listening to requests for ```http://127.0.0.1:9090``` (ip-address default, port by key)

---

### cоnfig file
Config File ```./config/config.prop``` contains key:value pairs with data to launch the application. If there is an “unclear” situation with the configuration data, simply delete this file and, the next time you start the application, it will self-restore with the default values, after which you can change them correctly and restart the application.  
Default values:  
```text
##  logger
logPathFolder:              ./logs
logFileName:                log
logMaxFileSizeKB:           1024

##  rabbitMQ
rConnect:                   amqp://guest:guest@localhost:5672
rTimeout:                    5

##  server
serverHost:                 127.0.0.1
serverPort:                 9000
```
