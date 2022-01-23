FROM openjdk:8-jre-alpine
RUN apk update && apk add aspell && apk add aspell-en

ENV APPLICATION_USER ktor
RUN adduser -D -g '' $APPLICATION_USER

RUN mkdir /app
RUN chown -R $APPLICATION_USER /app

USER $APPLICATION_USER

COPY ./build/libs/bot-is-gud-1.0-SNAPSHOT.jar /app/bot-is-gud.jar
WORKDIR /app

# dump master list of words and filter out words that end with 's
RUN aspell -l en dump master | grep -v "'s$" > words

CMD ["java", "-server", "-XX:+UnlockExperimentalVMOptions", "-XX:+UseCGroupMemoryLimitForHeap", "-XX:InitialRAMFraction=2", "-XX:MinRAMFraction=2", "-XX:MaxRAMFraction=2", "-XX:+UseG1GC", "-XX:MaxGCPauseMillis=100", "-XX:+UseStringDeduplication", "-jar", "bot-is-gud.jar"]
