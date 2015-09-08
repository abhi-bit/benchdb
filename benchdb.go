package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/abhi-bit/benchdb/bench"
)

var (
	conn = flag.String("conn", "", `Couchbase connection string,
                                    http://admin:pass@localhost:8091/default`)
	testBench = flag.String("test.bench", ".", "regex to match benchmarks to run")
)

var nsha = 10

func main() {
	flag.Parse()
	c := *conn
	bregex := *testBench

	err := (&bench.BenchKVStore{
		id: 0,
		Config: &bench.BenchDbConfig{
			Regex:  bregex,
			ShaLen: nsha,
		},
		Driver:  "couchbase",
		Connstr: c,
	}).Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
