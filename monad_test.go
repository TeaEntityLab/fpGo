// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fpGo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	var m MonadDef

	m = Monad.JustVal(1)
	assert.Equal(t, true, m.IsPresent())
	assert.Equal(t, false, m.IsNil())

	m = Monad.Just(nil)
	assert.Equal(t, false, m.IsPresent())
	assert.Equal(t, true, m.IsNil())
}

func TestOr(t *testing.T) {
	var m MonadDef

	m = Monad.JustVal(1)
	assert.Equal(t, 1, *m.OrVal(3).Val())
	m = Monad.Just(nil)
	assert.Equal(t, 3, *m.OrVal(3).Val())
}

func TestLet(t *testing.T) {
	var m MonadDef

	var letVal int

	letVal = 1
	m = Monad.JustVal(1)
	m.Let(func() {
		letVal = 2
	})
	assert.Equal(t, 2, letVal)

	letVal = 1
	m = Monad.Just(nil)
	m.Let(func() {
		letVal = 3
	})
	assert.Equal(t, 1, letVal)
}
