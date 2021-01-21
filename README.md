
# Shorts - URL shortener

## API

- GET `/`: Web Interface (NotImplemented)
- GET/POST `/shorts`: retrieve/create Shorts
- DELETE `/shorts/<id>`: delete a Short
- GET `/<id>`: redirect to URL

## Commands:

```bash
go mod init github.com/vaidd4/go-shorts
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .
```
