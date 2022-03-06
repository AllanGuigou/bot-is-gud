# go-is-gud

## Build

```bash
go build -o bin/bot-is-gud
```

```bash
docker build -t go-is-gud .
```

## Usage

```bash
./bin/bot-is-gud -t DISCORD_TOKEN
```

`-t DISCORD_TOKEN` can be excluded if DISCORD_TOKEN is an environment variable

```bash
docker run -it --rm -e DISCORD_TOKEN=$DISCORD_TOKEN --name go-is-gud go-is-gud
```
