# gb3

gb3 is a KISS cloud native tribune stack powered by [Go](https://golang.org/)

## gb0

gb0 is a guiless bouchot.

### Configuration

gb0 use the following environment variables.

#### GB0_LISTEN

Change default http listen address and port.

#### GB0_STORE_DIRECTORY

Post will be stored to disk only if this variable is a valid path to a writable directory.

#### GB0_MAX_POSTS_IN_HISTORY

Maximum number of posts returned in TSV backend.

#### GB0_MAX_DAYS_IN_HISTORY

Maximum number of days in the past to consider for building TSV backend.

### Endpoints

gb0 exposes two endpoint to retrieve and post messages.

#### /gb0/tsv

```bash
curl 'http://localhost:16667/gb0/tsv'
```

#### /gb0/post

```bash
curl -X POST -A "dave" -d "message=plop" 'http://localhost:16667/gb0/post'
```

or

```bash
curl -X POST -d "info=dave" -d "message=plop" 'http://localhost:16667/gb0/post'
```

## gb2c

gb2c is a webservice providing API for coincoins.

### Configuration

gb2c use the following environment variables.

#### GB2C_LISTEN

Change default http listen address and port. 

#### GB2C_VERBOSE

Set to true to enable full logging.

#### GB2C_LINUXFR_CLIENT_ID and GB2C_LINUXFR_CLIENT_SECRET

Oauth credentials for linuxfr API.

#### GB2C_GB0_TRIBUNES

List of gb0 tribunes url. Default value is:

```bash
GB2C_GB0_TRIBUNES=devnewton:https://gb3.devnewton.fr,gabuzomeu:https://gb3.plop.cc
```

#### GB2C_MIAOLI_TRIBUNES

List of miaoli tribunes url. Default value is:

```bash
GB2C_MIAOLI_TRIBUNES=reveildutroll:https://miaoli.im/tribune/LeR%C3%A9veilDuTroll
```

#### GB2C_MAX_POSTS_IN_HISTORY

Maximum number of posts returned by poll requests.

### Endpoints

gb2c exposes exports for listing tribunes then polling and posting messages.

#### /gb2c/list

```bash
curl 'http://localhost:16666/gb2c/list'
```

#### /gb2c/poll

```javascript
var postSource = new EventSource("http://localhost:16666/gb2c/poll");
postSource.onmessage = function (event) {
  console.log(event.data);
};
```

#### /gb2c/post

```bash
curl -X POST -d "info=dave" -d "message=plop" -d "tribune=euromussels" 'http://localhost:16666/gb2c/post'
```

### gc2

A simple coincoin frontend for gb2c.
