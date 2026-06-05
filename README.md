# XLI Bot 🤖

Telegram-бот на Go с AI-агентом.

## Компиляция через GitHub Actions

1. Запушь этот репо на GitHub
2. Перейди в **Actions** → **Build XLI Bot**
3. Нажми **Run workflow**
4. Выбери OS и архитектуру
5. Скачай готовый бинарник из артефактов

## Локальная сборка

```bash
cp .env.example .env
# Отредактируй .env
go mod tidy
go build -o xli-bot ./cmd/bot
./xli-bot
```

## Для твоего сервера (Termux, ARM64)

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o xli-bot ./cmd/bot
```
