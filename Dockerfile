FROM golang:latest

WORKDIR /distilleries_of_scotland
COPY . ./

# Build the application
RUN go build -o distilleries_of_scotland

EXPOSE 8080

CMD ["./distilleries_of_scotland", "run", "0.0.0.0:8000"]
