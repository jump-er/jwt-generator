## Introduction

### Сборка

Для Linux:
```
GOOS=linux GOARCH=amd64 go build
```

Для Mac:
```
GOOS=darwin GOARCH=arm64 go build
```

Готовый бинарь можно взять тут
https://nexus.infra.bingo.zone/#browse/browse:raw:jwt-generator

### Запуск

Для запуска генератора требуется экспортированная переменная `MAIN`:
```
export MAIN=/path/to/main
```

Запуск осуществляется из корня репозитория.

Даём права на исполнение.

Запускаем:
```
./jwt-generator
```
