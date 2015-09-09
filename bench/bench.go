package bench

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/couchbase/go-couchbase"
	"golang.org/x/tools/benchmark/parse"
)

var mutex = &sync.Mutex{}

type BenchDB interface {
	Run() error
	WriteSet(parse.Set) (int, error)
}

type BenchDBConfig struct {
	Regex  string
	ShaLen int
}

type BenchKVStore struct {
	Id      int64
	Config  *BenchDBConfig
	Driver  string
	Connstr string

    bucketObj *couchbase.Bucket
}

type Doc struct {
	Id               int64   `json:"id"`
	BatchId          string  `json:"batch_id"`
	LatestSha        string  `json:"latest_sha"`
	DateTime         string  `json:"datetime"`
	Name             string  `json:"name"`
	N                int     `json:"n"`
	NsOp             float64 `json:"ns_op"`
	AllocatedBytesOp uint64  `json:"allocated_bytes_op"`
	AllocsOp         uint64  `json:"allocs_op"`
}

func (benchdb *BenchKVStore) Run() error {
	cmd := exec.Command("go", "test", "-bench", benchdb.Config.Regex,
		"-test.run", "XXX", "-benchmem")
	var out bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &out)
	cmd.Stderr = io.Writer(os.Stderr)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to run command, err: %v\n", err)
	}

	benchSet, err := parse.ParseSet(&out)
	if err != nil {
		return fmt.Errorf("Failed to parse benchmark data, err: %\n", err)
	}

	_, err = benchdb.WriteSet(benchSet)
	if err != nil {
		fmt.Errorf("Failed to write benchSet to CB, err: %v\n", err)
	}
	return nil
}

func (benchdb *BenchKVStore) WriteSet(benchSet parse.Set) (int, error) {
	conn, err := couchbase.Connect(benchdb.Connstr)
	mf("bucket connect", err)

	pool, err := conn.GetPool("default")
	mf("Pool fetch", err)

	u, err := url.Parse(benchdb.Connstr)
	bName := strings.Split(u.Path, "/")[1]
	bucket, err := pool.GetBucket(bName)
	benchdb.bucketObj = bucket

	batchId, err := uuid()
	if err != nil {
		return 0, fmt.Errorf("Could not generate batch id, err: %v", err)
	}

	for _, b := range benchSet {
		n := len(b)
        fmt.Printf("LENGTH n: %d b: %v\n", n, b)
		for i := 0; i < n; i++ {
			val := b[i]
            fmt.Printf("val: %v\n", val)
			err := benchdb.saveBenchmark(batchId, *val)
			if err != nil {
				return 0, fmt.Errorf("Failed to save benchmark, err: %v", err)
			}
		}
	}
	return 0, nil
}

func (benchdb *BenchKVStore) saveBenchmark(batchId string, b parse.Benchmark) error {
	sha, err := latestGitSha(benchdb.Config.ShaLen)
	if err != nil {
		return err
	}

	sName := strings.TrimPrefix(strings.TrimSpace(b.Name), "Benchmark")

	t := time.Now()
	ts := t.Format(time.RFC3339)

	//Increment the counter
	mutex.Lock()
	benchdb.Id = benchdb.Id + 1
	mutex.Unlock()

	bStats := &Doc{
		Id:               benchdb.Id,
		BatchId:          batchId,
		LatestSha:        sha,
		DateTime:         ts,
		Name:             sName,
		N:                b.N,
		NsOp:             b.NsPerOp,
		AllocatedBytesOp: b.AllocedBytesPerOp,
		AllocsOp:         b.AllocsPerOp}

	data, err := json.Marshal(bStats)
	if err != nil {
		fmt.Println(err)
		return err
	}

	key := fmt.Sprintf("%d", benchdb.Id)
    fmt.Printf("Key: %s, value: %v\n", key, string(data))
	err = benchdb.bucketObj.SetRaw(key, 0, data)
	mf("Bucket set failed", err)
    return err
}

func latestGitSha(n int) (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("Failed to capture git sha, err: %v\n", err)
	}
	return string(out[:n]), nil
}

func uuid() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("Failed to generate uuid, err: %v\n", err)
	}
	return hex.EncodeToString(b), nil
}

func mf(errString string, err error) {
	if err != nil {
		log.Printf("%s err: %v\n", errString, err)
	}
}
