package cmd

import (
	"fmt"

	"github.com/flwyd/adif-multitool/adif"
)

func readers(rs ...adif.ReadWriter) map[adif.Format]adif.Reader {
	res := make(map[adif.Format]adif.Reader)
	for _, r := range rs {
		f, err := adif.ParseFormat(r.String())
		if err != nil {
			panic(fmt.Sprintf("Unknown adif.Format %q: %v", r, err))
		}
		res[f] = r.(adif.Reader)
	}
	return res
}

func writers(rs ...adif.ReadWriter) map[adif.Format]adif.Writer {
	res := make(map[adif.Format]adif.Writer)
	for _, r := range rs {
		f, err := adif.ParseFormat(r.String())
		if err != nil {
			panic(fmt.Sprintf("Unknown adif.Format %q: %v", r, err))
		}
		res[f] = r.(adif.Writer)
	}
	return res
}
