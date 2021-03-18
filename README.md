# gb3

gb3 is a KISS cloud native tribune stack powered by [Go](https://golang.org/)

## gb0

gb0 is a guiless bouchot.

### Endpoints

gb0 exposes two endpoint to retrieve and post messages.

#### /tsv

```bash
curl 'http://localhost:6667/tsv'
```

#### /post

```bash
curl -X POST -d "message=plop" 'http://localhost:6667/post'
```

## gb2c

gb2c is a webservice providing API for coincoins.

### Endpoints

gb2c exposes two endpoint for polling and posting tribune messages.

#### /api/poll

```javascript
var postSource = new EventSource("http://localhost:6666/api/poll");
postSource.onmessage = function (event) {
  console.log(event.data);
};
```

#### /api/post

```bash
curl -X POST -d "message=plop" -d "tribune=euromussels" 'http://localhost:6666/api/post'
```

### gc2

A simple coincoin frontend embedded in gb2c.
