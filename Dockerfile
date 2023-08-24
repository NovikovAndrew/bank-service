# Build stage
FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM golang:1.20
WORKDIR /app
COPY  --from=builder /app/main .
COPY app.env .

EXPOSE 9091
CMD [ "/app/main " ]