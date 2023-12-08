# asyncHTTPhandlerOnRabbitMQ
(Asynchronous processing engine HTTP requests on RabbitMQ) <a href="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/README.md">
                                                                   <img width="50" alt="switch to Russian" src="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/ru.jpg?raw=true">
                                                                 </a>

The project was released as a demonstration of the interaction between microservices and the RabbitMQ message broker.  
  
![Schema](https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/scheme.jpg?raw=true)  
Schematic diagram

---

## Preparing for work
For correct operation you need to install and run the server RabbitMQ: <a href="https://www.rabbitmq.com/download.html">
<img width="200" alt="switch to English" src="https://www.rabbitmq.com/img/logo-rabbitmq.svg">
</a>  

---

## Setup and launch
Setting up and running a project is done in several ways:
1. When compiling:
   1. [Api](#Api)
   2. [Config File](#Config-File)
2. [Executable File](#Executable-File)

---

## When compiling:
- Add a repository to the project with the command  
```
go get github.com/lent0s/asyncHTTPhandlerOnRabbitMQ
```  
- Add to import:

```Go
import (
    async "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/cmd"
)
```

### Api
- Call the function and pass configuration data to it

 ```Go
    async.RunService(LogMaxFileSizeKB uint16, RTimeout int, LogPathFolder, LogFileName, RConnect, ServerHost, ServerPort string)
 ```

where

```
- LogPathFolder - <b>existing</b> path to the directory where log files are stored
- LogFileName - log file name (without extension)
- LogMaxFileSizeKB - approximate log file size threshold in KB, after which a new file will be created
- RConnect - address and authentication data to access to RabbitMQ
- RTimeout - RabbitMQ server timeout
- ServerHost - IP-address receiving requests
- ServerPort - port receiving requests
```

- then compile your application and send POST requests with tasks to ``http://ServerHost:ServerPort/cmd`` and accept responses

### Config File
- Create a `config` directory in your project and call the function `async.Local()`. When the application is initialized for the first time, a configuration file will be created in this directory
- Fill in [cоnfig file](#cоnfig-file)
- Launch the application, after which you can send requests with tasks to ``http://ServerHost:ServerPort/cmd`` and accept responses

---

### Executable File
- Copy the project using
```
git clone https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ
```  
- Fill in [cоnfig file](#cоnfig-file)
- Build the application `go build async.go`

Start ```async``` with keys (or ```-help``` for reference):
```
- lpf - path to the directory where log files are stored
- lfn - log file name (without extension)
- lms - approximate log file size threshold in KB, after which a new file will be created
- rmq - address and authentication data to access to RabbitMQ
- rto - RabbitMQ server timeout
- sip - IP-address receiving requests
- spn - port receiving requests
```

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
