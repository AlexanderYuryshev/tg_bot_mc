docker build -t tg-bot:simple -f Dockerfile .
docker run --name tg-bot-simple -d tg-bot:simple