# chuchu

A simple Go app to poll statuses and timetables of trains operating in Finnish railroad system.

It uses public API's from https://www.digitraffic.fi/rautatieliikenne/ as source.

The app supports CLI version and a lightweight HTTP server.


## CLI

### List stations

```
go run *.go -stations
```

### List all trains for a station

```
go run *.go -station <stationShortCode>
```

### List all trains 

```
go run *.go -all
```

### Show a single train

```
go run *.go -train <trainNumber>
```

### Watch a single train as it travels

```
go run *.go -watch <trainNumber>
```



## HTTP Server

```
go run *.go -server
```