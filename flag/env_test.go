package flag

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type stringValue string

func (s *stringValue) Set(val string) error { *s = stringValue(val); return nil }
func (*stringValue) Type() string           { return "string" }
func (s *stringValue) String() string       { return string(*s) }

func TestSetFlagsFromEnv(t *testing.T) {
	assert.NoError(t, os.Setenv("FOO", "a"))
	assert.NoError(t, os.Setenv("L1_FOO", "b"))
	assert.NoError(t, os.Setenv("L1_L2_FOO", "c"))

	testCases := []struct {
		path   string
		search bool
		want   string
	}{
		{"", false, "a"},
		{"L1", false, "b"},
		{"L1 L2", false, "c"},
		{"Lx", false, "default"},
		{"L1 Lx", false, "default"},
		{"L1 L2 Lx", false, "default"},
		{"", true, "a"},
		{"L1", true, "b"},
		{"L1 L2", true, "c"},
		{"Lx", true, "default"},
		{"L1 Lx", true, "b"},
		{"L1 L2 Lx", true, "c"},
	}
	for _, tc := range testCases {
		fs := &pflag.FlagSet{}
		v := stringValue("default")
		fs.AddFlag(&pflag.Flag{Name: "FOO", Value: &v})
		assert.NoError(t, SetFlagsFromEnv(fs, tc.search, func(s string) string { return s }, strings.Split(tc.path, " ")...))
		assert.EqualValues(t, tc.want, v)
	}
}
