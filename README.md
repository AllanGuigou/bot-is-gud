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
DISCORD_TOKEN={TOKEN} java -jar build/libs/*.jar 
```

```bash
docker run --env DISCORD_TOKEN={TOKEN} bot-is-gud
```
