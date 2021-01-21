
# Shorts - URL shortener

Basic URL shortener in Go.

## API

| METHOD | URL          | Description
|:------:|--------------|-------------
| GET    |`/`           | Web Interface (not implemented yet)
| GET    |`/<id>`       | Redirect to URL
| GET    |`/shorts`     | Retrieve all Shorts
| POST   |`/shorts`     | Create a short
| DELETE |`/shorts/<id>`| Remove a short

## Commands:

```bash
go mod init github.com/vaidd4/go-shorts
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .
```
