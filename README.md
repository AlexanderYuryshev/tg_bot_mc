# Telegram Bot на Go: мастер-класс

Этот репозиторий содержит код Telegram-бота, написанного на Go, и материалы для его развертывания.

## 🛠 Подготовка к запуску

### Требования
- Установленный Go (версия 1.20+)
- Telegram API токен (получить у @BotFather)
- Для Docker-деплоя: установленный Docker

### Установка
1. Клонируйте этот репозиторий:
```bash
git clone https://github.com/yourusername/telegram-bot-go.git
cd telegram-bot-go
```
2. Установите зависимости:

```bash
go mod download
```
## 🚀 Запуск бота
### Локальный запуск
1. Создайте файл .env в корне проекта и добавьте туда ваш `TELEGRAM_BOT_TOKEN`, полученный из @BotFather (образец в [.env.example](./.env.example))
2. Запустите бота:

```bash
go run tg_bot_mc.go
```
### Деплой на VPS (для Linux-серверов типа CentOS или Ubuntu)
0. Убедитесь, что на сервере установлен Go (`go version`). Если нет, установите.
1. Создайте файл для systemd сервиса, образец лежит в [tg_bot_mc.service](./tg_bot_mc.service).
2. Скопируйте файлы бота на сервер, например, с помощью утилиты `scp` (в Windows доступна из консоли cmd):
```bash
scp -r C:/path/to/code server-user@your-server-ip:/path/to/code
```
3. Скопируйте файл сервиса на сервер в папку `/etc/systemd/system/`
```bash
scp -r C:/path/to/.service server-user@your-server-ip:/etc/systemd/system/
```
4. Включите и запустите сервис:
```bash
sudo systemctl enable tg_bot_mc
sudo systemctl start tg_bot_mc
```
### Docker-деплой (можно совместить с VPS)
В репозитории представлены два варианты - обычная сборка и multistage, которая сильно экономит место на диске.
Об образах Docker и их специфике можно почитать [по ссылке](https://hub.docker.com/_/golang).
1. Сборка и запуск обычного и multistage образов прописана в соответствующих скриптах:
```bash
sh simple_docker.sh
```
```bash
sh multistage_docker.sh
```
## 🔧 Возможности улучшения
### Доработка бота
1. Усовершенствовать функционал – например, хранить ссылки для каждого пользователя отдельно и проставлять домен в качестве категории по умолчанию, если она не была введена. Или выкинуть все это и сделать совсем другой функционал. Не стесняйтесь, реализуйте все ваши фантазии.
2. При деплое на сервер или через докер вынести `.env` из образа и папки с файлами и инициализировать из специального хранилища при запуске сервиса
3. Обработка ошибок – расширить логирование (с помощью log, zap) и реализовать механизмы повтора неуспешных операций.
4. Масштабирование – если подозреваете, что вашим ботом будут пользоваться тысячи людей, или просто хотите попрактиковаться в работе с нагрузкой, используйте очереди (RabbitMQ, Kafka) для высоконагруженных сервисов и Kubernetes для масштабируемости.

### CI/CD через GitHub Actions
1. Добавьте workflow для автоматического деплоя при push в main. Пример конфигурации:

```yaml
name: Deploy Bot
on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23.5'
      
      - name: Build and deploy
        env:
          TELEGRAM_BOT_TOKEN: ${{ secrets.BOT_TOKEN }}
        run: |
          scp -r . ${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }}:/path/to/bot
          ssh ${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }} "systemctl restart telegram-bot"
```
### Интеграция внешних API
1. Пример добавления OpenAI API:

```go
func generateAIResponse(prompt string) (string, error) {
    client := openai.NewClient("your-api-key")
    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT3Dot5Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: prompt,
                },
            },
        },
    )
    // обработка ответа
}
```