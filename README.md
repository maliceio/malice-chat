# malice-api
Malice API

### Requirements
 - Docker
 - NodeJS/npm

### Install
```bash
$ git clone https://github.com/maliceio/malice-api
$ cd malice-api
$ docker-compose up -d db
$ docker-compose up api
$ cd static
$ npm install
$ webpack-dev-server --port 8888
```

Navigate to `$(docker-machine env malice):8888`
