# smartcooly

[![Travis](https://img.shields.io/travis/miaolz123/smartcooly.svg)](https://travis-ci.org/miaolz123/smartcooly) [![Go Report Card](https://goreportcard.com/badge/github.com/marstau/smartcooly)](https://goreportcard.com/report/github.com/marstau/smartcooly) [![Github All Releases](https://img.shields.io/github/downloads/miaolz123/smartcooly/total.svg)](https://github.com/marstau/smartcooly/releases) [![Gitter](https://img.shields.io/gitter/room/miaolz123/smartcooly.svg)](https://gitter.im/miaolz123-smartcooly/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link) [![Docker Pulls](https://img.shields.io/docker/pulls/miaolz123/smartcooly.svg)](https://hub.docker.com/r/miaolz123/smartcooly/) [![license](https://img.shields.io/github/license/miaolz123/smartcooly.svg)](https://github.com/marstau/smartcooly/blob/master/LICENSE)

[中文文档](http://smartcooly.marstau.com/docs/1.0/api-zh-cn/)

## Installation

You can install smartcooly from **installation package** or **Docker**.

The default username and password are `admin`, please modify them immediately after login!

### From installation package

1. Download the smartcooly installation package on [this page](https://github.com/marstau/smartcooly/releases)
2. Unzip the smartcooly installation package
3. Enter the extracted smartcooly installation directory
4. Run `smartcooly`

Then, smartcooly is running at `http://localhost:9876`.

**Linux & Mac user quick start command**

```shell
wget https://github.com/marstau/smartcooly/releases/download/v{{VERSION}}/smartcooly_{{OS}}_{{ARCH}}.tar.gz
tar -xzvf smartcooly_{{OS}}_{{ARCH}}.tar.gz
cd smartcooly_{{OS}}_{{ARCH}}
./smartcooly
```

Please replace *{{VERSION}}*, *{{OS}}*, *{{ARCH}}* first.

### by Docker

```shell
docker run --name=smartcooly -p 19876:9876 marstau/smartcooly
```

Then, smartcooly is running at `http://localhost:19876`.

### by heroku

```
heroku addons:create heroku-postgresql:hobby-basic -a smartcooly
heroku buildpacks:add https://github.com/debitoor/ssh-private-key-buildpack.git -a smartcooly
heroku config:set SSH_KEY=$(cat ~/.ssh/id_rsa | base64)  -a smartcooly
heroku buildpacks:add heroku/go -a smartcooly
```
## Usage

### Add an Exchange

![](https://raw.githubusercontent.com/miaolz123/smartcooly/master/docs/_media/add-exchange.png)

### Add an Algorithm

![](https://raw.githubusercontent.com/miaolz123/smartcooly/master/docs/_media/add-algorithm.png)

![](https://raw.githubusercontent.com/miaolz123/smartcooly/master/docs/_media/edit-algorithm.png)

### Deploy an Algorithm

![](https://raw.githubusercontent.com/miaolz123/smartcooly/master/docs/_media/add-trader.png)

### Run a Trader

![](https://raw.githubusercontent.com/miaolz123/smartcooly/master/docs/_media/run-trader.png)

## Commands

```
go build smartcooly.go
npm run dist
```

## Algorithm Reference

[Read Documentation](http://smartcooly.stockdb.org/#/#algorithm-reference)

## Contributing

Contributions are not accepted in principle until the basic infrastructure is complete.

However, the [ISSUE](https://github.com/marstau/smartcooly/issues) is welcome.

## Denote

eth:0x6E6dDDE24C79e94633CACAa442FadDcD41Af31Bd

btc:1KaJo5bFTpFzJXSYDiFjYkHELvcQdy9NSn


## License

Copyright (c) 2016 [marstau](https://github.com/marstau) by MIT

# Reference

* <https://github.com/ccxt/ccxt>
* <https://github.com/miaolz123/samaritan>
