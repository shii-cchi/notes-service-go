FROM golang:1.22.5-alpine
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ARG DB_URI
ENV DB_URI=$DB_URI
WORKDIR notes-service-go
COPY . .
RUN go mod download
RUN go build -o notes_server cmd/main.go
CMD ["sh", "-c", "goose -dir internal/database/migrations postgres ${DB_URI} up && ./notes_server"]