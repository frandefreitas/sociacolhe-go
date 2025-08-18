FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /bin/sociacolhe ./cmd/server

FROM alpine:3.20
ENV PORT=8080
WORKDIR /srv
COPY --from=base /bin/sociacolhe /usr/local/bin/sociacolhe
EXPOSE 8080
CMD ["sociacolhe"]
