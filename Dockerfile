# Get node and build FE assets
FROM node:20.11.0-alpine3.10 as build_fe

RUN mkdir -p /app
COPY . /app
WORKDIR /app

RUN npm install
RUN npm run build

# Get golang and build BE assets
FROM golang:1.21.6-alpine3.19 as build_be
WORKDIR /app

COPY --link . /app

RUN go mod download

RUN go build -o server .

EXPOSE 8080

# Run
CMD ["/app/server"]