FROM golang:latest as build

WORKDIR /build

COPY . .
RUN go build 

FROM build as release

RUN useradd -m app
WORKDIR /home/app
USER app

EXPOSE 6563
ENV METRICS_SERVER=localhost:6563

ENV DATABASE_URL="postgres://CHANGE_ME"
ENV BOT_TOKEN="my discord bot token"
ENV ENABLE_POLLING=true

EXPOSE 8080
COPY --from=build /build/aoc-dcbot .
CMD ["./aoc-dcbot"]
