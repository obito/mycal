FROM golang:latest

WORKDIR /opt/
COPY *.go /opt/
COPY go.mod /opt/
COPY go.sum /opt/
RUN echo > calendar.ics
RUN go build -o /usr/bin/mycal

CMD [ "/usr/bin/mycal" ]
