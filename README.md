# Project L0

_______________________________________________________________
### ``` -1- БД находится в контейнере. Подключение по localhost:5433 ```
_______________________________________________________________
### ```-2- Сервис запускается на localhost:8080  ```
<br>у сервиса 2 основные ручки:
<br>GET "/api/orders" : получение данных заказа по order_uid
<br>POST "/api/orders":    запись данных в БД.
<br>-- данные сохраняются как в БД, так и в inmemory кэш.
<br>-- при запуске сервиса данные подтягиваются из БД в кэш
<br>-- дополнительно: в случае недоступности БД при попытке записи,
<br>--данные сохраняются в дополнительное inMemory хранилище, и
переносятся при дальнейшем удачном подкючении к БД.```
___________________________________________________________________

### ```-3- NATS в контейнере. |Подключение по localhost:4222 ```

______________________________________________________________________
### ```-4- Реализованы также publisher и consumer для взаимодействия с nats ```
#### <br>-4.1- publisher умеет генерировать JSON заказа с рандомными значениями полей нужного формата. order_uid задается вручную, для удобства тестирования.сгенерированный JSON публикуется в Jetstream с темой "orders"
#### -4.2- consumer подписан на JetStream тему "order" как только он получает сообщение по теме это сообщение отправляется POST запросом на localhost:8080/api/orders


______
## Getting Started

Инструкция по запуску.<br>
Nats-streaming сервер должен быть запущен на localhost:4222

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

run consumer

```bash
make run-consumer
```

run publisher

```bash
make run-publisher
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

clean up binary from the last build

```bash
make clean
```