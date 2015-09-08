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
	"sync/atomic"
	"time"

	"github.com/couchbase/go-couchbase"
	"golang.org/x/tools/benchmark/parse"
)

type BenchDB interface {
	Run() error
	WriteSet(parse.Set) (int, error)
}

type BenchDBConfig struct {
	Regex  string
	ShaLen int
}

type counter64 int64

type BenchKVStore struct {
	id      *counter64
	Config  *BenchDBConfig
	Driver  string
	Connstr string
	// sample connstr: http://ops:password@localhost:8091/default
	bucketObj *couchbase.Bucket
}

type JSONTime time.Time

type Doc struct {
	id                 int64    `json:"id"`
	batch_id           string   `json:"batch_id"`
	latest_sha         string   `json:"latest_sha"`
	datetime           JSONTime `json:"datetime"`
	name               string   `json:"name"`
	n                  int      `json:"n"`
	ns_op              float64  `json:"ns_op"`
	allocated_bytes_op uint64   `json:"allocated_bytes_op"`
	allocs_op          uint64   `json:"allocs_op"`
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan 2 15:04:05 +0530 IST 2006"))
	return []byte(stamp), nil
}

func (c *counter64) increment() int64 {
	var next int64
	for {
		next = int64(*c) + 1
		if atomic.CompareAndSwapInt64((*int64)(c), int64(*c), next) {
			return next
		}
	}
}

func (c *counter64) get() int64 {
	return atomic.LoadInt64((*int64)(c))
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
		for i := 0; i < n; i++ {
			fmt.Println("AbhI:", b[i], batchId)
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
	ts := JSONTime(time.Now())
	//Increment the counter
	benchdb.id.increment()
	id := benchdb.id.get()

	bStats := &Doc{
		id:                 id,
		batch_id:           batchId,
		latest_sha:         sha,
		datetime:           ts,
		name:               sName,
		n:                  b.N,
		ns_op:              b.NsPerOp,
		allocated_bytes_op: b.AllocedBytesPerOp,
		allocs_op:          b.AllocsPerOp}

	data, err := json.Marshal(bStats)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%d", id)
	err = benchdb.bucketObj.SetRaw(key, 0, data)
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
		log.Fatalf("%s err: %v\n", errString, err)
	}
}
