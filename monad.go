// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fpGo

type MonadDef struct {
	ref *interface{}
}

func (self MonadDef) JustVal(in interface{}) MonadDef {
	return MonadDef{ref: &in}
}
func (self MonadDef) Just(in *interface{}) MonadDef {
	return MonadDef{ref: in}
}
func (self MonadDef) OrVal(or interface{}) MonadDef {
	return self.Or(&or)
}
func (self MonadDef) Or(or *interface{}) MonadDef {
	if self.IsNil() {
		return MonadDef{ref: or}
	}

	return self
}

func (self MonadDef) Val() *interface{} {
	return self.ref
}
func (self MonadDef) IsPresent() bool {
	return !(self.IsNil())
}
func (self MonadDef) IsNil() bool {
	return self.ref == nil
}

var Monad MonadDef
