# GO SSE EVENT

Simple app to send string object from server to client using HTTP 2 Server Push

## Requirement

- Go 1.12
- Go mod enabled

## How to run

- clone this repository
- run with `make run`
- server will run at `localhost:8080`, you can publish any string value to endpoint `/event` with method `POST`