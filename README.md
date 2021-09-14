## tesla_exporter
The tesla_exporter is work in China.
The tesla access_token is cached in local file. The tesla account email and password are self-hosting.

## build
```
go build -ldflags "-X 'main.email={YOUR_EMAIL}' -X 'main.password={YOUR_PWD}'" cmd/tesla_exporter.go
```
## start
```
# you can also start with supervisor
./tesla_exporter
```

## metrics
```
curl -v "http://127.1:9610/metrics"
```

## grafana
When the server is running, you can add exporter job into prometheus and `tesla.json` which is in grafana folder into your grafana web dashboard. Then you will see the images.
![tesla](/grafana/tesla.png)
