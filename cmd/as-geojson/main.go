package main

/*

curl -s 'http://localhost:8080/api/point-in-polygon?latitude=37.618582632478834&longitude=-122.38769531250001&format=' \

|

/usr/local/whosonfirst/go-whosonfirst-spr-geojson/bin/as-geojson \
	-reader-uri fs:///usr/local/data/woeplanet-state-us/data \
	-path 'places.#.spr:id' \
	-as-wof

*/

import (
	"bufio"
	"context"
	"flag"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI.")
	path := flag.String("path", "places.#.wof:path", "A valid tidwall/gjson query path for finding the path for each element in your SPR response.")

	wof_cb := flag.Bool("as-wof", false, "Parse each item returned by -path as though it were a Who's On First style ID/path.")

	flag.Parse()

	ctx := context.Background()

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}

	wr := io.MultiWriter(os.Stdout)

	var cb geojson.JSONPathResolverCallback

	if *wof_cb {

		cb = func(ctx context.Context, path string) (string, error) {
			return geojson.WhosOnFirstPathWithString(path)
		}
	}

	resolver_func := geojson.JSONPathResolverFuncWithCallback(*path, cb)

	as_opts := &geojson.AsFeatureCollectionOptions{
		Reader:           r,
		Writer:           wr,
		JSONPathResolver: resolver_func,
	}

	err = geojson.AsFeatureCollectionWithJSON(ctx, body, as_opts)

	if err != nil {
		log.Fatal(err)
	}

}
