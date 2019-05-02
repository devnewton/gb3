# gb3

gb3 is a KISS tribune powered by [Go](https://golang.org/)

## Endpoints

gb3 exposes two endpoint for polling and posting tribune messages.

### /api/poll

```javascript
var postSource = new EventSource("http://localhost:6666/api/poll");
postSource.onmessage = function (event) {
  console.log(event.data);
};
```

### /api/post

```bash
curl -X POST -d "message=plop" -d "tribune=euromussels" 'http://localhost:6666/api/post'
```
