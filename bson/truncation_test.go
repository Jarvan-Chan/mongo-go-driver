// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"
	"testing"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/internal/assert"
)

type inputArgs struct {
	Name string
	Val  *float64
}

type outputArgs struct {
	Name string
	Val  *int64
}

func TestTruncation(t *testing.T) {
	t.Run("truncation", func(t *testing.T) {
		inputName := "truncation"
		inputVal := 4.7892

		input := inputArgs{Name: inputName, Val: &inputVal}

		buf := new(bytes.Buffer)
		vw, err := bsonrw.NewBSONValueWriter(buf)
		assert.Nil(t, err)
		enc := NewEncoder(vw)
		enc.IntMinSize()
		enc.SetRegistry(DefaultRegistry)
		err = enc.Encode(&input)
		assert.Nil(t, err)

		var output outputArgs
		dc := bsoncodec.DecodeContext{
			Registry: DefaultRegistry,
			Truncate: true,
		}

		err = UnmarshalWithContext(dc, buf.Bytes(), &output)
		assert.Nil(t, err)

		assert.Equal(t, inputName, output.Name)
		assert.Equal(t, int64(inputVal), *output.Val)
	})
	t.Run("no truncation", func(t *testing.T) {
		inputName := "no truncation"
		inputVal := 7.382

		input := inputArgs{Name: inputName, Val: &inputVal}

		buf := new(bytes.Buffer)
		vw, err := bsonrw.NewBSONValueWriter(buf)
		assert.Nil(t, err)
		enc := NewEncoder(vw)
		enc.IntMinSize()
		enc.SetRegistry(DefaultRegistry)
		err = enc.Encode(&input)
		assert.Nil(t, err)

		var output outputArgs
		dc := bsoncodec.DecodeContext{
			Registry: DefaultRegistry,
			Truncate: false,
		}

		// case throws an error when truncation is disabled
		err = UnmarshalWithContext(dc, buf.Bytes(), &output)
		assert.NotNil(t, err)
	})
}
