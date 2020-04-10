###################
##  build stage  ##
###################
FROM golang:1.13.0-alpine as builder
WORKDIR /api-golang-kubernetes
COPY . .
RUN go build -v -o ./api-golang-kubernetes

##################
##  exec stage  ##
##################
FROM alpine:3.10.2
WORKDIR /app
COPY ./configs/config.json.default ./configs/config.json
COPY --from=builder /api-golang-kubernetes/api-golang-kubernetes /app/
CMD ["./api-golang-kubernetes"]
