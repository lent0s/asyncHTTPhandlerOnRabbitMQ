# asyncHTTPhandlerOnRabbitMQ
(Asynchronous processing engine HTTP requests on RabbitMQ) <a href="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/README_EN.md">
                                                                   <img width="50" alt="switch to English" src="https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/en.jpg?raw=true">
                                                                 </a>
  
Проект релизован в качестве демонстрации взаимодействия работы микросервисов и брокера сообщений RabbitMQ  
  
![Схема](https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/blob/main/doc/scheme.jpg?raw=true)
Принципиальная схема

---

## Подготовка к работе
Для корректной работы необходимо установить и запустить сервер RabbitMQ: <a href="https://www.rabbitmq.com/download.html">
<img width="200" alt="switch to English" src="https://www.rabbitmq.com/img/logo-rabbitmq.svg">
</a>  

---

## Настройка и запуск
Настройка и запуск проекта выполняются несколькими способами:
1. [При компиляции:](#При-компиляции:)
   1. [Api](#Api)
   2. [Файл конфигурации](#Файл-конфигурации)
2. [Исполняемый файл](#Исполняемый-файл)

---

## При компиляции:
- Добавьте в проект репозиторий командой  
```
go get github.com/lent0s/asyncHTTPhandlerOnRabbitMQ
```  
- Добавьте в импорт:

```Go
import (
    async "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/cmd"
)
```

### Api
- Вызовите функцию и передайте в неё конфигурационные данные

 ```Go
    async.RunService(LogMaxFileSizeKB uint16, RTimeout int, LogPathFolder, LogFileName, RConnect, ServerHost, ServerPort string)
 ```

где

```
- LogPathFolder - <b>существующий</b> путь к директории хранения лог-файлов
- LogFileName - имя лог-файла (без расширения)
- LogMaxFileSizeKB - примерное пороговое значение размера лог-файла в КБ, после которого будет создан новый файл
- RConnect - адрес и аутентификационные данные для доступа к RabbitMQ
- RTimeout - таймаут сервера RabbitMQ
- ServerHost - IP-адрес принимающий запросы
- ServerPort - порт принимающий запросы
```

- после чего скомпилируйте своё приложение и отправляйте POST запросы с заданиями на ``http://ServerHost:ServerPort/cmd`` и принимайте ответы

### Файл конфигурации
- Создайте директорию `config` в своём проекте и вызовите функцию `async.Local()`. При первой инициализации приложения в этой директории будет создан файл конфигурации
- Заполните [файл конфигурации](#файл-кoнфигурации)
- Запустите приложение, после чего можете отправлять запросы с заданиями на ``http://ServerHost:ServerPort/cmd`` и принимать ответы

---

## Исполняемый файл
- Скопируйте проект с помощью
```
git clone https://github.com/lent0s/asyncHTTPhandlerOnRabbitMQ
```  
- Заполните [файл конфигурации](#файл-кoнфигурации)
- Выполните сборку приложения `go build async.go`

Запустите ```async``` с ключами (или ```-help``` для справки):
- lpf - путь к директории хранения лог-файлов
- lfn - имя лог-файла (без расширения)
- lms - примерное пороговое значение размера лог-файла в КБ, после которого будет создан новый файл
- rmq - адрес и аутентификационные данные для доступа к RabbitMQ
- rto - таймаут сервера RabbitMQ
- sip - IP-адрес принимающий запросы
- spn - порт принимающий запросы

Пример: ```local -spn 9090``` - для запуска приложения с прослушиванием запросов на ```http://127.0.0.1:9090``` (ip-адрес по-умолчанию, порт по ключу)

---

### файл кoнфигурации
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
