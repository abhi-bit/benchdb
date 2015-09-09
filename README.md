### benchdb

Store Golang project benchmarks in Couchbase and then tag git SHAs against perf dip/jump.

### Sample run:

```
➜  benchdb git:(master) ✗ ./benchdb -conn http://localhost:8091/default 
PASS
BenchmarkFib100          10000            152631 ns/op              32 B/op          1 allocs/op
ok      github.com/abhi-bit/fib       2.032s
```

It would store following blob in Couchbase:

```
Key: 1, value: {“id”:1,”batch_id”:”118d631c9b98ba33caa0c0aabd57ff7d”,”latest_sha”:”59dedc23e6”,”datetime”:”2015-09-09T13:44:22+05:30”,”name”:”Fib100”,”n”:10000,”ns_op”:152631,”allocated_bytes_op”:32,”allocs_op”:1}
```

