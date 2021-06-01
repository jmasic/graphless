package main

import (
	"compress/zlib"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/snappy"
)

func main() {
	file, _ := os.Open("/home/chef/Downloads/test_ungzipped")
	info, _ := file.Stat()
	log.Printf("File info is: Size %f, name %s\n", float64(info.Size())/1e6, info.Name())
	file.Close()

	data, _ := ioutil.ReadFile("/home/chef/Downloads/test_ungzipped")

	encoded := make([]byte, 0, 0)
	compressed := snappy.Encode(encoded, data)

	out, _ := os.Create("snappy_out")
	out.Write(compressed)

	out, _ = os.Create("zlib_out")
	w := zlib.NewWriter(out)
	w.Write(data)
}
