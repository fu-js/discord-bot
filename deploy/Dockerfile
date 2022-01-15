FROM golang:1.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN make build_api

EXPOSE $PORT

CMD [ "./api" ]