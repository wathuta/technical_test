package fieldmask

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wathuta/technical_test/orders/internal/common"
	"google.golang.org/genproto/protobuf/field_mask"
)

// MaxSize is the hard limit on the number of fields in a single field mask.
const MaxSize = 96

// Mask represents a list of fields in a resource.
type Mask struct {
	// Fields included in the mask, in snake case.
	Fields []string
}

// New returns a new Mask created from a protobuf mask.
func New(mask *field_mask.FieldMask) (*Mask, error) {
	if len(mask.Paths) > MaxSize {
		return nil, fmt.Errorf("number of fields is %d, maximum allowed is %d", len(mask.Paths), MaxSize)
	}

	fields := make([]string, len(mask.Paths))
	for i, k := range mask.Paths {
		fields[i] = strcase.ToSnake(k)
	}

	return &Mask{
		Fields: fields,
	}, nil
}

func (f *Mask) RemoveOutputOnly() {
	newFields := make([]string, 0, len(f.Fields))
	for _, k := range f.Fields {
		if !common.IsFieldOutputOnly(k) {
			newFields = append(newFields, k)
		}
	}
	f.Fields = newFields
}
