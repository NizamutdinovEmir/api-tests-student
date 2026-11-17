# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите [Git](https://git-scm.com/install/) и [Go 1.25.1 или новее](https://go.dev/doc/install).
2. Склонируйте fork через HTTPS и в корне проекта скачайте зависимости.

```bash
go mod download
```

3. Подготовьте два разных аккаунта на [TestBoard](https://testboard.avito.com/). Продавец создаёт объявления, покупатель добавляет их в избранное. Создайте в корне проекта локальный файл `.env.override` с данными этих аккаунтов.

```
TEST_LOGIN=seller@example.com
TEST_PASSWORD=seller_password
TEST_BUYER_LOGIN=buyer@example.com
TEST_BUYER_PASSWORD=buyer_password
```

4. Для полного прогона создайте у продавца объявление «Пять ноль» с описанием и фотографией. Его используют сценарии первой практики. Вторая работа создаёт собственные тестовые объявления и удаляет их после проверок.

```bash
go test -v -p 1 -count=1 ./...
```

Чтобы запустить только вторую работу, выполните из корня проекта:

```bash
go test -v -p 1 -count=1 ./tests/scenarios/advertisementCreation ./tests/scenarios/favorites
```

Флаг `-p 1` запускает пакеты последовательно, чтобы не перегружать шлюз TestBoard.

## Логирование

Флаг `-v` выводит шаги тестов и результаты проверок. Для локальной диагностики HTTP-запросов можно установить
`DEBUG=true` в `.env.override`. Такой вывод содержит данные авторизации.

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
golangci-lint run
```
