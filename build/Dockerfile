
FROM golang:1.24-alpine
EXPOSE 80 4000
WORKDIR /app
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
COPY . ./


RUN go mod download
# ENV GO111MODULE=off


RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o hello-app ./cmd/main.go

CMD [ "/go/bin/dlv", "--listen=:4000", "--headless=true", "--log-output=stack", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/app/hello-app" ]

