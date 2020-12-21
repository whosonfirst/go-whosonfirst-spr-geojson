package geojson

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"strconv"
)

type SPRPathResolver func(context.Context, spr.StandardPlacesResult) (string, error)

type JSONPathResolver func(context.Context, []byte) ([]string, error)

func WhosOnFirstSPRPathResolverFunc() SPRPathResolver {

	fn := func(ctx context.Context, r spr.StandardPlacesResult) (string, error) {

		id, err := strconv.ParseInt(r.Id(), 10, 64)

		if err != nil {
			return "", err
		}

		return uri.Id2RelPath(id)
	}

	return fn
}

func JSONPathResolverFunc(gjson_path string) JSONPathResolver {

	fn := func(ctx context.Context, body []byte) ([]string, error) {

		path_rsp := gjson.GetBytes(body, gjson_path)

		if !path_rsp.Exists() {
			return nil, errors.New("Missing path")
		}

		paths := make([]string, 0)

		for _, p := range path_rsp.Array() {
			paths = append(paths, p.String())
		}

		return paths, nil
	}

	return fn
}
