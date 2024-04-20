FROM golang:latest as builder
ENV LOCALE=nb_NO
#ENV LOCALE=en_US
WORKDIR /app
COPY go.* ./
COPY vendor/ ./vendor/
COPY ttf/ ./ttf/
COPY *.go ./
COPY cmd/server/ ./cmd/server/
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN cd cmd/server && go build -v -mod=vendor --tags $LOCALE -o ../../server
FROM debian:buster-slim
COPY --from=builder /app/server /server
COPY cmd/server/static/ /static/
ENV PORT 8080
CMD ["/server"]
