FROM golang:1.16-alpine AS build

RUN apk add --no-cache git

WORKDIR /tmp/app
COPY go.mod .
#COPY go.sum .

#RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .
#RUN go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=build /tmp/app/conf.json /app/conf.json
COPY --from=build /tmp/app/server /app/server

CMD ["./server"]