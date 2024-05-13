FROM golang:1.20-alpine AS builder
LABEL authors="apartapatia"

WORKDIR /build

ADD go.mod .

COPY . .

COPY configs/ configs/

RUN go build -o computer_club_assistant cmd/computer_club_assistant/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/computer_club_assistant /build/computer_club_assistant
COPY --from=builder /build/configs /build/configs

CMD ./computer_club_assistant $FILE_NAME