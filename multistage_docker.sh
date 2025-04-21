docker build -t tg-bot:multistage -f Dockerfile.multistage .
docker run --name tg-bot-multistage -d tg-bot:multistage