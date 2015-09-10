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
key: 6e2d0fe23f09a753d74da034466756a4 value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”Fib40”,”n”:2,”ns_op”:7.52662411e+08,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
key: dacf590ae84f0cd2053af8e649f8cec7 value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”FibComplete”,”n”:3000000,”ns_op”:404,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
key: 5ff124ca1fe828c651cd0977a82ea115 value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”Fib1”,”n”:1000000000,”ns_op”:2.14,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
key: bb73de238bfa8f83d8a03e2396b6600f value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”Fib2”,”n”:200000000,”ns_op”:6.71,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
key: bf531c6aba2830021deb879707bcc1ad value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”Fib3”,”n”:100000000,”ns_op”:11.7,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
key: 00be4ca6248e45939391b630e286a3a3 value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”Fib10”,”n”:3000000,”ns_op”:400,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
key: f5d47b8275784503790ad87ca42433da value: {“batch_id”:”a592505ebf75679d7083fa098b28eab2”,”latest_sha”:”6bcdb74bf1”,”datetime”:”2015-09-10T10:49:17+05:30”,”name”:”Fib20”,”n”:30000,”ns_op”:49884,”allocated_bytes_op”:0,”allocs_op”:0,”mb_per_s”:0}
```
