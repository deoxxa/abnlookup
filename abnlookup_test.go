package abnlookup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type abnTest struct {
	abn  string
	name string
	err1 error
	err2 error
}

var (
	tests = []abnTest{
		{"87 007 382 031", "", ErrWrongLength, nil},
		{"87007382xxx", "", ErrInvalidFormat, nil},
		{"87007382031", "", ErrInvalidChecksum, nil},
		{"87007382032", "MOVING DATA PTY LTD", nil, nil},
		{"87007382033", "", ErrInvalidChecksum, nil},
		{"51824753556", "AUSTRALIAN TAXATION OFFICE", nil, nil},
		{"51824999396", "", nil, ErrRecordNotFound},
	}
)

func TestValidate(t *testing.T) {
	for _, e := range tests {
		t.Run(e.abn, func(t *testing.T) {
			if e.err1 != nil {
				assert.Contains(t, Validate(e.abn).Error(), e.err1.Error())
			} else {
				assert.NoError(t, Validate(e.abn))
			}
		})
	}
}

func TestFetch(t *testing.T) {
	for _, e := range tests {
		t.Run(e.abn, func(t *testing.T) {
			if e.err1 != nil {
				assert.Contains(t, Lookup(e.abn).Error(), e.err1.Error())
			} else if e.err2 != nil {
				assert.Contains(t, Lookup(e.abn).Error(), e.err2.Error())
			} else {
				res, err := Fetch(e.abn)
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, e.name, res.Name)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	for _, e := range tests {
		t.Run(e.abn, func(t *testing.T) {
			if e.err1 != nil {
				assert.Contains(t, Lookup(e.abn).Error(), e.err1.Error())
			} else if e.err2 != nil {
				assert.Contains(t, Lookup(e.abn).Error(), e.err2.Error())
			} else {
				assert.NoError(t, Lookup(e.abn))
			}
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	for _, e := range tests {
		b.Run(e.abn, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Validate(e.abn)
			}
		})
	}
}
