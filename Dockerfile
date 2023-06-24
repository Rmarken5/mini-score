FROM golang:1.20

WORKDIR var

COPY . .

EXPOSE 8080

RUN go build -o ./mini-score ./service/cmd/server/main.go

CMD ["./mini-score"]
