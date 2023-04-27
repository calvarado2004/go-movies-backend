FROM docker.io/golang:1.20 as builder

RUN mkdir /app

COPY ./cmd /app/cmd
COPY ./internal /app/internal


WORKDIR /app

RUN go mod init github.com/calvarado2004/go-movies-backend && go get github.com/go-chi/chi/v5 && go get github.com/go-chi/cors && go get github.com/graphql-go/graphql && go get github.com/jackc/pgx/v4 && go get github.com/jackc/pgconn && go get github.com/golang-jwt/jwt/v4 && go get golang.org/x/crypto && go get golang.org/x/text

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o moviesApp ./cmd/api

RUN chmod +x /app/moviesApp

FROM alpine:latest 

RUN mkdir /app

COPY --from=builder /app/moviesApp /app

EXPOSE 8080

CMD [ "/app/moviesApp"]
