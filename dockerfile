FROM golang

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

CMD [ "go", "run", "cmd/main.go" ]




