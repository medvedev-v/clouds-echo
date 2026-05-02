# Clouds Echo
Реализация https://www.cloudping.info/ на Go, плюс проверки для запрещенных и разрешенных ресурсов

## Запросы
GET http://localhost:8080/echo/clouds

GET http://localhost:8080/echo/forbidden

GET http://localhost:8080/echo/allowed

## Ответ
```json 
[
  {
    "url": "https://dynamodb.us-east-1.amazonaws.com/ping",
    "ping": 776,
    "responsecode": "200 OK"
  },
...
 ``` 
