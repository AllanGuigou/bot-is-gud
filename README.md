# bot-is-gud

## Build

```bash
./gradlew shadowJar
```

```bash
docker build -t bot-is-gud .
```

## Run

```bash
NICKNAME_USERS="$(cat users)" DISCORD_TOKEN={TOKEN} java -jar build/libs/*.jar 
```

```bash
docker run --env-file .env bot-is-gud
```
