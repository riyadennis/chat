FROM golang:alpine
RUN mkdir chat
ADD . /chat/
WORKDIR /chat
EXPOSE 8080
RUN go build -o chat .
CMD ["./chat"]