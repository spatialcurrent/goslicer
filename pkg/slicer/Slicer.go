// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package slicer provides an easy API to slice a stream of bytes using a custom slice of indicies.
package slicer

import (
	"fmt"
)

// Slicer provides an easy API to slice a stream of bytes.
type Slicer struct {
	indicies []int
	strict   bool // strict mode
}

// New returns a new slicer that slices an input stream of bytes using the given indicies.
// If strict is false, then any data contained entirely within the given bounds is returned, rather than an error.
func New(strict bool, indicies ...int) (*Slicer, error) {
	if len(indicies) > 2 {
		return nil, fmt.Errorf("invalid number of indicies %d", len(indicies))
	}
	return &Slicer{strict: strict, indicies: indicies}, nil
}

func (s Slicer) strictSliceString(in string) (string, error) {
	if len(s.indicies) == 2 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		if start < 0 {
			return "", fmt.Errorf("when in strict mode, start index %d (%d) must be greater than or equal to zero", s.indicies[0], start)
		}
		end := s.indicies[1]
		if end < 0 {
			end = len(in) + end
		}
		if end < start {
			return "", fmt.Errorf("start index %d (%d) must be before end index %d (%d)", s.indicies[0], start, s.indicies[1], end)
		}
		if end > len(in) {
			return "", fmt.Errorf("when in strict mode, end index %d (%d) cannot be greater than the length of the input string %q (%d)", s.indicies[1], end, in, len(in))
		}
		return in[start:end], nil
	}
	if len(s.indicies) == 1 {
		start := s.indicies[0]
		if start < 0 {
			return "", fmt.Errorf("when in strict mode, start index %d (%d) must be greater than or equal to zero", s.indicies[0], start)
		}
		if start > len(in)-1 {
			return "", fmt.Errorf("when in strict mode, start index %d (%d) cannot be greater than the length of the input string %q (%d) minus 1", s.indicies[0], start, in, len(in))
		}
		return in[start:], nil
	}
	return in[:], nil
}

func (s Slicer) looseSliceString(in string) (string, error) {
	if len(s.indicies) == 2 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		end := s.indicies[1]
		if end < 0 {
			end = len(in) + end
		}
		if end < start {
			return "", fmt.Errorf("start index %d (%d) must be before end index %d (%d)", s.indicies[0], start, s.indicies[1], end)
		}
		if start < 0 {
			start = 0
		}
		if start > len(in) {
			return "", nil
		}
		if end < 0 {
			return "", nil
		}
		if end > len(in) {
			return in[start:], nil
		}
		return in[start:end], nil
	}
	if len(s.indicies) == 1 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		if start < 0 {
			start = 0
		}
		if start > len(in)-1 {
			return "", nil
		}
		return in[start:], nil
	}
	return in[:], nil
}

// SliceString returns the sliced version of the given string.
func (s Slicer) SliceString(in string) (string, error) {
	if s.strict {
		return s.strictSliceString(in)
	}
	return s.looseSliceString(in)
}

// MustSliceString returns the sliced version of the given string and panics if there is an error.
func (s Slicer) MustSliceString(in string) string {
	out, err := s.SliceString(in)
	if err != nil {
		panic(err)
	}
	return out
}

func (s Slicer) strictSliceBytes(in []byte) ([]byte, error) {
	if len(s.indicies) == 2 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		if start < 0 {
			return make([]byte, 0), fmt.Errorf("when in strict mode, start index %d (%d) must be greater than or equal to zero", s.indicies[0], start)
		}
		end := s.indicies[1]
		if end < 0 {
			end = len(in) + end
		}
		if end < start {
			return make([]byte, 0), fmt.Errorf("start index %d (%d) must be before end index %d (%d)", s.indicies[0], start, s.indicies[1], end)
		}
		if end > len(in) {
			return make([]byte, 0), fmt.Errorf("when in strict mode, end index %d (%d) cannot be greater than the length of the input bytes %q (%d)", s.indicies[1], end, in, len(in))
		}
		return in[start:end], nil
	}
	if len(s.indicies) == 1 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		if start < 0 {
			return make([]byte, 0), fmt.Errorf("when in strict mode, start index %d (%d) must be greater than or equal to zero", s.indicies[0], start)
		}
		if start > len(in)-1 {
			return make([]byte, 0), fmt.Errorf("when in strict mode, start index %d (%d) cannot be greater than the length of the input bytes %q (%d) minus 1", s.indicies[0], start, in, len(in))
		}
		return in[start:], nil
	}
	return in[:], nil
}

func (s Slicer) looseSliceBytes(in []byte) ([]byte, error) {
	if len(s.indicies) == 2 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		end := s.indicies[1]
		if end < 0 {
			end = len(in) + end
		}
		if end < start {
			return make([]byte, 0), fmt.Errorf("start index %d (%d) must be before end index %d (%d)", s.indicies[0], start, s.indicies[1], end)
		}
		if start > len(in) {
			return make([]byte, 0), nil
		}
		if end < 0 {
			return make([]byte, 0), nil
		}
		if start < 0 {
			start = 0
		}
		if end > len(in) {
			return in[start:], nil
		}
		return in[start:end], nil
	}
	if len(s.indicies) == 1 {
		start := s.indicies[0]
		if start < 0 {
			start = len(in) + start
		}
		if start < 0 {
			start = 0
		}
		if start > len(in)-1 {
			return make([]byte, 0), nil
		}
		return in[start:], nil
	}
	return in[:], nil
}

// SliceBytes returns the sliced version of the given slice of bytes.
func (s Slicer) SliceBytes(in []byte) ([]byte, error) {
	if s.strict {
		return s.strictSliceBytes(in)
	}
	return s.looseSliceBytes(in)
}

// MustSliceBytes returns the sliced version of the given slice of bytes and panics if there is an error.
func (s Slicer) MustSliceBytes(in []byte) []byte {
	out, err := s.SliceBytes(in)
	if err != nil {
		panic(err)
	}
	return out
}

// Slice returns the sliced version of the given slice of bytes.
func (s Slicer) Slice(in interface{}) (interface{}, error) {
	switch x := in.(type) {
	case string:
		return s.SliceString(x)
	case *string:
		return s.SliceString(*x)
	case []byte:
		return s.SliceBytes(x)
	case *[]byte:
		return s.SliceBytes(*x)
	}
	return nil, fmt.Errorf("unknown type %T", in)
}

// MustSlice returns the sliced version of the given slice of bytes and panics if there is an error.
func (s Slicer) MustSlice(in interface{}) interface{} {
	switch x := in.(type) {
	case string:
		return s.MustSliceString(x)
	case *string:
		return s.MustSliceString(*x)
	case []byte:
		return s.MustSliceBytes(x)
	case *[]byte:
		return s.MustSliceBytes(*x)
	}
	panic(fmt.Errorf("unknown type %T", in))
}
