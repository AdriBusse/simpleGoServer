FROM golang:1.20-alpine

WORKDIR /app

COPY . ./
RUN go mod download


EXPOSE 3333

CMD [ "go", "run", "main.go" ]