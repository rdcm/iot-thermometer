## About

A pet project with an example of using the IOT. In this project, a scheme has been constructed for transmitting measurements from temperature and humidity sensors using the MQTT protocol. There is also a consumer who reads the received measurements and saves them in the Clickhouse. All obtained measurements can be visualized on a dashboard in Grafana.

## Scheme

## Architecture

## Up & Running

* - device according to the proposed scheme must be assembled independently

`docker-compose build --no-cache consumer && docker-compose up -d`

## Web interfaces
`http://localhost:15672/` - rabbitmq admin pannel, admin/admin
`http://localhost:3000/`  - grafana web interface, admin/admin