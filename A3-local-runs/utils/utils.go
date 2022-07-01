package utils

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func ParseUint(valueString string, bitSize int) uint64 {
	value, err := strconv.ParseUint(valueString, 10, bitSize)
	if err != nil {
		panic(err)
	}

	return value
}

func ParseInt(valueString string) int {
	value, err := strconv.Atoi(valueString)
	if err != nil {
		panic(err)
	}

	return value
}

func ParseTimeWithNano(valueString string) time.Time {
	value, err := time.Parse(time.RFC3339Nano, valueString)
	if err != nil {
		panic(err)
	}

	return value
}

func ParseBool(valueString string) bool {
	value, err := strconv.ParseBool(valueString)
	if err != nil {
		panic(err)
	}

	return value
}

func MeasureDurationMs(start time.Time) int64 {
	elapsed := time.Since(start)
	return elapsed.Nanoseconds() / 1e6
	//log.Printf("Executed %s; Duration: %v s; Additional info: %s\n", functionName, elapsed.Nanoseconds()/1e6, logMessage)
}

var EPSILON float64 = 0.00000001

func FloatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func Monitor(runId string, tag string, formattedMessage string) {
	log.Printf("MONITORING:%s:%s:%s\n", runId, tag, formattedMessage)
}

func GetLogger() *log.Logger {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.9999Z07:00"})
	return log.StandardLogger()
}

func ZLibCompress(uncompressed []byte) []byte {
	compressed := new(bytes.Buffer)

	w := zlib.NewWriter(compressed)
	defer w.Close()
	_, err := w.Write(uncompressed)
	w.Flush()
	if err != nil {
		panic(err)
	}

	return compressed.Bytes()
}

func ZLibDecompress(compressed []byte) []byte {
	b := bytes.NewBuffer(compressed)

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	decompressed, _ := ioutil.ReadAll(r)

	if err != nil {
		panic(err)
	}

	return decompressed
}

func IsLocal() bool {
	argv := os.Args[1:]
	return len(argv) > 0 && argv[0] == "local"
}
