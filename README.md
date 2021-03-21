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

### Endpoints

gb0 exposes two endpoint to retrieve and post messages.

#### /tsv

```bash
curl 'http://localhost:16667/tsv'
```

#### /gb0/post

```bash
curl -X POST -d "message=plop" 'http://localhost:16667/gb0/post'
```

## gb2c

gb2c is a webservice providing API for coincoins.

### Configuration

gb2c use the following environment variables.

#### GB2C_LISTEN

Change default http listen address and port. 

#### GB2C_VERBOSE

Set to true to enable full logging.

#### GB2C_FULLTEXT_SEARCH_INDEX_PATH

Fulltext search is only available if this variable is a valid path to a writable directory.

#### GB2C_LINUXFR_CLIENT_ID and GB2C_LINUXFR_CLIENT_SECRET

Oauth credentials for linuxfr API.

### Endpoints

gb2c exposes two endpoint for polling and posting tribune messages.

#### /gb2c/poll

```javascript
var postSource = new EventSource("http://localhost:16666/gb2c/poll");
postSource.onmessage = function (event) {
  console.log(event.data);
};
```

#### /gb2c/post

```bash
curl -X POST -d "message=plop" -d "tribune=euromussels" 'http://localhost:16666/gb2c/post'
```

### gc2

A simple coincoin frontend embedded in gb2c.
