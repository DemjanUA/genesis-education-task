FROM golang:1.20

WORKDIR /usr/src/app

COPY . ./

RUN go get -u github.com/gorilla/mux
RUN go get github.com/sendgrid/sendgrid-go
RUN go mod download
RUN go mod verify

RUN go build -o genesis-education-task

# app uses port 80
EXPOSE 80

CMD ["/usr/src/app/genesis-education-task"]