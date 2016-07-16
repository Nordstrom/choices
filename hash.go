package choices

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
)

const longScale = float64(0xFFFFFFFFFFFFFFF)

func hash(ctx context.Context) (int64, error) {
	cv := ctxValuer{ctx: ctx}
	nn := cv.string("namespace")
	en := cv.string("experiment")
	pn := cv.string("param")
	units := cv.stringSlice("units")
	if cv.err != nil {
		return 0, cv.err
	}

	key := config.globalSalt + "." + nn + "." + en + "." + pn + ":" + strings.Join(units, ".")
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(key)))
	i, err := strconv.ParseInt(hash[:15], 16, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func getUniform(hash int64, min, max float64) float64 {
	return min + (max-min)*(float64(hash)/longScale)
}

type ctxValuer struct {
	ctx context.Context
	err error
}

func (c *ctxValuer) string(s string) string {
	if c.err != nil {
		return ""
	}

	val, ok := c.ctx.Value(s).(string)
	if !ok {
		c.err = fmt.Errorf("%q is not a string", s)
		return ""
	}
	return val
}

func (c *ctxValuer) stringSlice(s string) []string {
	if c.err != nil {
		return []string{}
	}
	val, ok := c.ctx.Value(s).([]string)
	if !ok {
		c.err = fmt.Errorf("%q is not a slice of string", s)
		return []string{}
	}
	return val
}
