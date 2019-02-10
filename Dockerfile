FROM golang:1.11

RUN mkdir /app

ADD dgpx /app/dgpx
ADD templates /app/templates

WORKDIR /app

EXPOSE 8080

ENTRYPOINT [ "/app/dgpx" ]
