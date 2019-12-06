FROM golang:alpine
RUN mkdir chat
ADD . /chat/
WORKDIR /chat
RUN go build -o chat .
CMD ["./chat"]