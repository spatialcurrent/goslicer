// =================================================================
//
// Copyright (C) 2020 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slicer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlicerStrictZero(t *testing.T) {
	in := "Hello World"
	s, err := New(true)
	assert.NoError(t, err)
	out := s.MustSlice(in)
	assert.Equal(t, "Hello World", out)
}

func TestSlicerStrictOne(t *testing.T) {
	in := "Hello World"
	s, err := New(true, 6)
	assert.NoError(t, err)
	out := s.MustSlice(in)
	assert.Equal(t, "World", out)
}

func TestSlicerLooseOne(t *testing.T) {
	in := "Hello World"
	s, err := New(false, 100)
	assert.NoError(t, err)
	out := s.MustSlice(in)
	assert.Equal(t, "", out)
	//
	s, err = New(false, -100)
	assert.NoError(t, err)
	out = s.MustSlice(in)
	assert.Equal(t, "Hello World", out)
}

func TestSlicerStrictTwo(t *testing.T) {
	in := "Hello World"
	s, err := New(true, 0, 5)
	assert.NoError(t, err)
	out := s.MustSlice(in)
	assert.Equal(t, "Hello", out)
}

func TestSlicerStrictTwoNegative(t *testing.T) {
	in := "Hello World"
	s, err := New(true, 0, -6)
	assert.NoError(t, err)
	out := s.MustSlice(in)
	assert.Equal(t, "Hello", out)
}

func TestSlicerLooseTwoNegative(t *testing.T) {
	in := "Hello World"
	s, err := New(false, -100, -6)
	assert.NoError(t, err)
	out := s.MustSlice(in)
	assert.Equal(t, "Hello", out)
}
