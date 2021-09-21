# Capuchin

Distributed CSV Query Engine

Status: mostly nonsense

## Example 

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
