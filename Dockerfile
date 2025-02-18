# Dockerfile

# Build the Go backend
FROM golang:1.20 AS builder
WORKDIR /app
COPY ./backend .
RUN go build -o main .

# Build the React frontend
FROM node:14 AS frontend
WORKDIR /app
COPY ./frontend .
RUN npm install && npm run build

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=frontend /app/build ./frontend/build
EXPOSE 8080
CMD ["./main"]