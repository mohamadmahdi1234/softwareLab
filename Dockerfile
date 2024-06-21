FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://repo.snapp.tech/repository/goproxy/,goproxy.cn,goproxy.io,direct
RUN go mod download && go mod verify

COPY . .
ENV DB_USERNAME=helloDocker
ENV DB_PASSWORD=secret

CMD ["go","run","main.go","serve"]