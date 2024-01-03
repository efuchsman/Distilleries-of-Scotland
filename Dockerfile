FROM golang:latest

WORKDIR /distilleries_of_scotland
COPY . ./

RUN go build -o distilleries_of_scotland

EXPOSE 8080

CMD ["./distilleries_of_scotland"]
