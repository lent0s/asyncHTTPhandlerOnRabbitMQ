# asyncHTTPhandlerOnRabbitMQ
(Asynchronous processing engine HTTP requests on RabbitMQ) <a href="https://github.com/asyncHTTPhandlerOnRabbitMQ/doc/README_EN.md">
                                                                   <img width="50" alt="switch to English" src="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/en.jpg?raw=true">
                                                                 </a>

Проект релизован в качестве демонстрации взаимодействия работы микросервисов и брокера сообщений RabbitMQ.  
![Схема](https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/scheme.jpg?raw=true)
Принципиальная схема.

## Настройка и запуск
Настройка и запуск проекта выполняются несколькими способами:
1. При компиляции:
   1. [Api](#Api)
   2. [Файл конфигурации](#Файл)
2. [Исполняемый файл](#Исполняемый)

### Api
Для запуска клиента достаточно добавить в импорт
 ```Go
import (
    async "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/cmd"
)
 ```
вызвать функцию и передать в неё конфигурационные данные
 ```Go
    async.RunService(LogMaxFileSizeKB uint16, LogPathFolder, LogFileName, RConnect, ServerHost, ServerPort string)
 ```
после чего отправлять запросы с заданиями на ``http://ServerHost:ServerPort/cmd`` и принимать ответы.

### Файл конфигурации
Заполнить [файл конфигурации](#файл) и запустить проект:
```
    go run async.go
```
или с импортом
 ```Go
import (
    async "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/cmd"
)
 ```
вызвать
 ```Go
    async.Local()
 ```
после чего отправлять запросы с заданиями на ``http://ServerHost:ServerPort/cmd`` и принимать ответы.

### Исполняемый файл
Скомпилируйте проект с помощью
```
    go build async.go
```
Запустите ```async``` с ключами (или ```-help``` для справки):
- lpf - путь к директории хранения лог-файлов
- lfn - имя лог-файла (без расширения)
- lms - примерное пороговое значение размера лог-файла в КБ, после которого будет создан новый файл
- rmq - адрес и аутентификационные данные для доступа к RabbitMQ
- rto - таймаут сервера RabbitMQ
- sip - IP-адрес принимающий запросы
- spn - порт принимающий запросы

Пример: ```local -spn 9090``` - для запуска приложения с прослушиванием запросов на ```http://127.0.0.1:9090``` (ip-адрес по-умолчанию, порт по ключу)

### файл конфигурации
Файл конфигурации ```./config/config.prop``` содержит пары ключ:значение с данными для запуска приложения. При "непонятной" ситуации с данными конфигурации - просто удалите этот файл и, при следующем запуске приложения, он самовосстановится со значениями по-умолчанию, после чего можете корректно их поменять и перезапустить приложение.  
Значения по умолчанию:  
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
