FROM golang:1.18-bullseye as base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify
COPY . ./
RUN go test ./... 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM gcr.io/distroless/static-debian11
COPY --from=base /main .
ENTRYPOINT ["./main"]