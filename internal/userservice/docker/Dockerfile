FROM golang:1.23

WORKDIR /userservice

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
ENV GRPC_VERBOSITY=DEBUG
ENV GRPC_TRACE=all
RUN CGO_ENABLED=0 GOOS=linux go build -o /userservice.bin ./cmd/userservice/userservice.go

CMD ["/userservice.bin"]