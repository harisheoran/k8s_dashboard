FROM golang:1.22.5 as build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/web/

#FROM scratch

#COPY --from=build /app/main .
#COPY ./ui ./

EXPOSE 1313

CMD [ "./main" ]