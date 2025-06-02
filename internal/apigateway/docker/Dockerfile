FROM golang:1.23

WORKDIR /apigateway

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /apigateway.bin ./cmd/apigateway/apigateway.go

CMD ["/apigateway.bin"]