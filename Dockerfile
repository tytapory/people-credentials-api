FROM golang:1.24
RUN mkdir /app
WORKDIR /app
COPY . /app

EXPOSE 8080