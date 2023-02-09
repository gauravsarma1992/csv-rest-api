# CSV REST APIs
REST APIs for CSV files

## Installation
```bash
cd src
go run run/server.go
```

## Fetch matching elements
```bash
curl http://localhost:9090/csv -XPOST -d '{"files": ["file1"], "filters": {"statuscode": "200", "requestname": "write"}}'
```
