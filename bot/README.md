# bot-is-gud

## Setup

- JDK 1.8

## Build

```bash
./gradlew shadowJar
```

```bash
docker build -t bot-is-gud .
```

## Run

```bash
DISCORD_TOKEN={TOKEN} java -jar build/libs/*.jar
```

```bash
docker run --env-file .env bot-is-gud
```
