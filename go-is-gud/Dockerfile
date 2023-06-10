FROM golang:1.20 AS build

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownling them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN CGO_ENABLED=0 go build -v -o /bot-is-gud

FROM alpine:latest
WORKDIR /
COPY --from=build /bot-is-gud /bot-is-gud

CMD ["./bot-is-gud"]
