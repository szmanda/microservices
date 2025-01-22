FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# RUN go build -o app.bin ./main.go

EXPOSE 80
CMD ["sleep", "infinity"]
# CMD ["go build -o build.bin ./main.go; /app/build.bin"]
# CMD ["go run ./service/main.go"]