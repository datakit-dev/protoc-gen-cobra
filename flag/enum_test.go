package flag

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func TestEnumVar(t *testing.T) {
	fs := &pflag.FlagSet{}
	var p Enum
	EnumVar(fs, &p, "foo", "")
	assert.NoError(t, fs.Set("foo", "TUE"))
	assert.Equal(t, Enum_TUE, p)
	assert.NoError(t, fs.Set("foo", "wed"))
	assert.Equal(t, Enum_WED, p)
	v := fs.Lookup("foo").Value
	assert.Equal(t, "string", v.Type())
}

func TestEnumPointerVar(t *testing.T) {
	fs := &pflag.FlagSet{}
	var p *Enum
	EnumPointerVar(fs, &p, "foo", "")
	assert.NoError(t, fs.Set("foo", "TUE"))
	want := Enum_TUE
	assert.Equal(t, &want, p)
	assert.NoError(t, fs.Set("foo", "2"))
	want = Enum_WED
	assert.Equal(t, &want, p)
	v := fs.Lookup("foo").Value
	assert.Equal(t, "string", v.Type())
}

func TestEnumSliceVar(t *testing.T) {
	fs := &pflag.FlagSet{}
	var p []Enum
	EnumSliceVar(fs, &p, "foo", "")
	assert.NoError(t, fs.Set("foo", "MON,1"))
	assert.Equal(t, []Enum{Enum_MON, Enum_TUE}, p)
	assert.NoError(t, fs.Set("foo", "WED"))
	assert.Equal(t, []Enum{Enum_MON, Enum_TUE, Enum_WED}, p)
	v := fs.Lookup("foo").Value
	assert.Equal(t, "slice", v.Type())
}

type Enum int32

const (
	Enum_MON Enum = 0
	Enum_TUE Enum = 1
	Enum_WED Enum = 2
)

func (x Enum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Enum) Descriptor() protoreflect.EnumDescriptor {
	return enumTypes[0].Descriptor()
}

func (Enum) Type() protoreflect.EnumType {
	return &enumTypes[0]
}

func (x Enum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

var enumTypes = make([]protoimpl.EnumInfo, 1)

func init() {
	protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			RawDescriptor: []byte{
				0x0a, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x65,
				0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2a, 0x21, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x07,
				0x0a, 0x03, 0x4d, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x55, 0x45, 0x10, 0x01,
				0x12, 0x07, 0x0a, 0x03, 0x57, 0x45, 0x44, 0x10, 0x02, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70,
				0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
			},
			NumEnums: 1,
		},
		GoTypes:   []interface{}{(Enum)(0)},
		EnumInfos: enumTypes,
	}.Build()
}
