package main

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

	cb := func(ctx context.Context, path string) (string, error) {
		return geojson.WhosOnFirstPathWithString(path)
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
