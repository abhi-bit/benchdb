### benchdb

Store Golang project benchmarks in Couchbase and then tag git SHAs against perf dip/jump.

### Sample run:

```
➜  fib git:(master) ../benchdb/benchdb -conn http://localhost:8091/default 
PASS
BenchmarkFib1   1000000000               2.14 ns/op            0 B/op          0 allocs/op
BenchmarkFib2   200000000                6.86 ns/op            0 B/op          0 allocs/op
BenchmarkFib3   100000000               11.7 ns/op             0 B/op          0 allocs/op
BenchmarkFib10   3000000               398 ns/op               0 B/op          0 allocs/op
BenchmarkFib20     30000             49429 ns/op               0 B/op          0 allocs/op
BenchmarkFib40         2         750234949 ns/op               0 B/op          0 allocs/op
BenchmarkFibComplete     3000000               408 ns/op               0 B/op          0 allocs/op
ok      github.com/abhi-bit/fib 13.074s
```

It would store following blob in Couchbase:

```
Key: 1, value: {“id”:1,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”Fib1”,”n”:1000000000,”ns_op”:2.12,”allocated_bytes_op”:0,”allocs_op”:0}
Key: 2, value: {“id”:2,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”Fib2”,”n”:200000000,”ns_op”:6.82,”allocated_bytes_op”:0,”allocs_op”:0}
Key: 3, value: {“id”:3,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”Fib3”,”n”:100000000,”ns_op”:11.7,”allocated_bytes_op”:0,”allocs_op”:0}
Key: 4, value: {“id”:4,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”Fib10”,”n”:3000000,”ns_op”:403,”allocated_bytes_op”:0,”allocs_op”:0}
Key: 5, value: {“id”:5,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”Fib20”,”n”:30000,”ns_op”:50122,”allocated_bytes_op”:0,”allocs_op”:0}
Key: 6, value: {“id”:6,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”Fib40”,”n”:2,”ns_op”:7.71609019e+08,”allocated_bytes_op”:0,”allocs_op”:0}
Key: 7, value: {“id”:7,”batch_id”:”4fb1b27d31bbc4b34aced01f438fd20b”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-09T13:57:54+05:30”,”name”:”FibComplete”,”n”:3000000,”ns_op”:404,”allocated_bytes_op”:0,”allocs_op”:0}
```
