## tesla_exporter
The tesla_exporter is work in China.
The tesla access_token is cached in local file. The tesla account email and password are self-hosting.

## build
```
go build -ldflags "-X 'main.email={YOUR_EMAIL}' -X 'main.password={YOUR_PWD}'" cmd/tesla_exporter.go
```


## metrics
```
curl -v "http://127.1:9610/metrics"
```
