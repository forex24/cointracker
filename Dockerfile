FROM golang:1.18-alpine as builder
RUN apk --no-cache --update add git gcc musl-dev linux-headers ca-certificates tzdata bash

RUN wget -q https://github.com/markbates/refresh/releases/download/v1.4.11/refresh_1.4.11_linux_amd64.tar.gz \
    && tar -xzf refresh_1.4.11_linux_amd64.tar.gz && mv refresh /usr/local/bin/refresh && chmod u+x /usr/local/bin/refresh

WORKDIR /opt/cointracker
COPY ./cointracker .

RUN go build -o ./bin/cointracker .


CMD ["refresh", "-c", "refresh.yml"]

FROM alpine:3.12 as dist
RUN apk --no-cache --update add ca-certificates tzdata bash

WORKDIR /opt/cointracker
COPY --from=builder /opt/cointracker/bin/cointracker ./bin/cointracker

CMD ["/opt/cointracker/bin/cointracker"]