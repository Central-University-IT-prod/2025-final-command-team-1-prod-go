FROM golang:1.23.6-alpine AS build

WORKDIR /go/src/app
COPY ./go.mod go.mod
COPY ./go.sum go.sum
COPY ./firebase_credentials.json firebase_credentials.json

RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app cmd/main.go
EXPOSE 8010

FROM alpine:latest
COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/scripts /scripts/
COPY --from=build /go/src/app/firebase_credentials.json /firebase_credentials.json

CMD ["/app"]