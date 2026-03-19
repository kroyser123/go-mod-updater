go-mod-updater
CLI‑утилита для анализа Go‑модулей и проверки обновлений зависимостей.

Функционал
Клонирует Git‑репозиторий (публичный или приватный)

Находит все go.mod файлы (поддержка монорепозиториев)

Анализирует текущие версии зависимостей

Проверяет наличие доступных обновлений

Показывает, какие зависимости можно обновить

Определяет тип обновления: patch / minor / major

Установка
bash
git clone https://github.com/kroyser123/go-mod-updater.git
cd go-mod-updater
go build -o modupdater cmd/modupdater/main.go
Использование
Базовый запуск
bash
./modupdater -repo https://github.com/go-kit/kit.git
Только устаревшие зависимости
bash
./modupdater -repo https://github.com/go-kit/kit.git -outdated
С учётом indirect‑зависимостей
bash
./modupdater -repo https://github.com/go-kit/kit.git -outdated -all
JSON‑вывод
bash
./modupdater -repo https://github.com/go-kit/kit.git -outdated -json
Флаги
Флаг	Описание
-repo	URL Git‑репозитория
-token	Токен для приватных репозиториев
-outdated	Показывать только устаревшие зависимости
-all	Включать indirect‑зависимости
-json	Вывод в JSON
-debug	Отладочные логи
Пример вывода (человекочитаемый формат)
MODULE: github.com/go-kit/kit
FILE:   C:\Users\9E8E~1\AppData\Local\Temp\modupdater-747241387\go.mod
GO:     1.17

github.com/go-kit/log                 v0.2.0 → v0.2.1   (patch)
github.com/go-logfmt/logfmt           v0.5.1 → v0.6.0   (minor)
github.com/nats-io/nats-server/v2     v2.8.4 → v2.12.5  (minor)
github.com/prometheus/client_golang   v1.11.1 → v1.23.2 (minor)
github.com/openzipkin/zipkin-go       v0.2.5 → v0.4.3   (minor)

Пример JSON‑вывода
{
  "Module": "github.com/go-kit/kit",
  "GoVersion": "1.17",
  "Statuses": [
    {
      "Path": "github.com/go-kit/log",
      "Current": "v0.2.0",
      "Latest": "v0.2.1",
      "Indirect": false,
      "NeedUpdate": true,
      "UpdateType": "patch"
    },
    {
      "Path": "github.com/golang-jwt/jwt/v4",
      "Current": "v4.0.0",
      "Latest": "v4.5.2",
      "Indirect": false,
      "NeedUpdate": true,
      "UpdateType": "minor"
    },
    {
      "Path": "github.com/nats-io/nats-server/v2",
      "Current": "v2.8.4",
      "Latest": "v2.12.5",
      "Indirect": false,
      "NeedUpdate": true,
      "UpdateType": "minor"
    },
    {
      "Path": "github.com/rabbitmq/amqp091-go",
      "Current": "v1.2.0",
      "Latest": "v1.10.0",
      "Indirect": false,
      "NeedUpdate": true,
      "UpdateType": "minor"
    }
  ]
}
Особенности
Shallow clone (--depth 1) для скорости

Работает с репозиториями без тегов

Параллельная обработка зависимостей

Поддержка приватных репозиториев (токен)

Определение типа обновлений через semver
