FROM golang as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app cmd/main.go
#CMD [ "./bin/app" ]


FROM alpine

WORKDIR /app
COPY --from=builder /app/bin/app .

COPY /database/migrations /app/database/migrations

#RUN chmod 777 /app
EXPOSE 8081
CMD [ "./app" ]