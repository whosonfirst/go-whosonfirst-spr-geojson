package geojson

import (
	"bufio"
	"bytes"
	"context"
	go_geojson "github.com/paulmach/go.geojson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"io"
)

func ToFeatureCollection(ctx context.Context, rsp spr.StandardPlacesResults, r reader.Reader) (*go_geojson.FeatureCollection, error) {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	err := AsFeatureCollection(ctx, rsp, r, wr)

	if err != nil {
		return nil, err
	}

	wr.Flush()

	return go_geojson.UnmarshalFeatureCollection(buf.Bytes())
}

func AsFeatureCollection(ctx context.Context, rsp spr.StandardPlacesResults, r reader.Reader, wr io.Writer) error {

	_, err := io.WriteString(wr, `{"type":"FeatureCollection", "features": [`)

	if err != nil {
		return err
	}

	for i, pl := range rsp.Results() {

		if i > 0 {

			_, err := io.WriteString(wr, `,`)

			if err != nil {
				return err
			}
		}

		path := pl.Path()
		fh, err := r.Read(ctx, path)

		if err != nil {
			return err
		}

		defer fh.Close()

		_, err = io.Copy(wr, fh)

		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(wr, `]}`)

	if err != nil {
		return err
	}

	return nil
}
