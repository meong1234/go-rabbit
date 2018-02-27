# go-rabbit
this rabbitMQ sample on golang, used on my presentation

### presentation
slide : https://docs.google.com/presentation/d/11vPV8Elw4S4aXZXdHoEPIXkOJAFvqsw-LO9qoUpV-cw/edit
youtube: https://www.youtube.com/watch?v=l_SgYpDu3JQ&t=2257s

### Dependencies

 - RabbitMq Client: [streadway amqp](http://github.com/streadway/amqp/)
 - Logger: [Logrus](github.com/sirupsen/logrus)
 - CLI: [Urfave Cli](github.com/urfave/cli)

### Run rabbitMQ using docker
 - `docker-compose up -d`

### Run the app using command line
 - make build

 - Run as publisher `./out/go-rabbit publish` in another terminal

 - Run as publisher `./out/go-rabbit subscribe` in another terminal
