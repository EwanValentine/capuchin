# Capuchin

Distributed CSV Query Engine

Status: mostly nonsense


## Install 

```bash 
$ git clone git@github.com:EwanValentine/capuchin.git
$ cd capuchin && go install
$ capuchin 
```

## Commands 

```bash
$ capuchin start // starts both gRPC and HTTP proxy server
$ capuchin grpc // starts just the gRPC server
$ capuchin http // starts just the http server, pointless as it uses the gRPC server 
```

## Query API

Using `httpie` using the test data set.

```bash 
$ http post localhost:9999/v1/query \
  select:='["user_id"]' \
  where="user_id = abc123" \
  source="./query/test-data.csv"
```

## Example (Library)

Example using Capuchin as a library:

```golang 
s := source.NewFileSource()
fileSource, err := s.Load("./query/test-data.csv")
if err != nil {
  log.Panic(err)
}

query := &query.Query{
  Select: []string{"user_id", "date"},
  Where:  "user_id = abc123",
}
query.Source(fileSource)

results, err := query.Exec()
if err != nil {
  log.Panic(err)
}

log.Println(results)
```
