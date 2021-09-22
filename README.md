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

Test dataset example:

```csv
order_id,user_id,date
abc123,abc123,2021-09-01
def456,def123,2021-09-02
abc123,abc123,2021-09-03
```

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

## Data Management 

1. Point your Capuchin cluster at your datalake, defining which column in your data is your date key. 
2. The Capuchin nodes will create a shard, loading batches of the data into memory using a date range. For example node 0 will load 20190101 to 20190801 into memory using a sharding algorithm to automatically divide the data.
3. 