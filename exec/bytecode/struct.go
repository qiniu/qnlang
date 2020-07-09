/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package bytecode

import (
	"reflect"
	"strconv"

	"github.com/qiniu/x/log"
)

// -----------------------------------------------------------------------------

// A StructField describes a single field in a struct.
type StructField = reflect.StructField

// A StructInfo represents a struct information.
type StructInfo struct {
	Fields []StructField
	typ    reflect.Type
}

// Struct creates a new struct.
func Struct(fields []StructField) *StructInfo {
	return &StructInfo{
		Fields: fields,
		typ:    reflect.StructOf(fields),
	}
}

// Type returns the struct type.
func (p *StructInfo) Type() reflect.Type {
	return p.typ
}

// -----------------------------------------------------------------------------

type varsContext = reflect.Value

func makeVarList(vars []*Var) []StructField {
	exists := make(map[string]bool, len(vars))
	items := make([]StructField, len(vars))
	for i, v := range vars {
		if _, ok := exists[v.name]; ok {
			v.name = "Q" + strconv.Itoa(i) + v.name[1:]
		} else {
			exists[v.name] = true
		}
		items[i].Type = v.typ
		items[i].Name = v.name
	}
	return items
}

func makeVarsContextType(vars []*Var, ctx *Context) reflect.Type {
	t := Struct(makeVarList(vars)).Type()
	if log.CanOutput(log.Ldebug) {
		nestDepth := ctx.getNestDepth()
		for i, v := range vars {
			if v.nestDepth != nestDepth || v.idx != uint32(i) {
				log.Panicln("makeVarsContext failed: unexpected var -", v.name[1:], v.nestDepth, nestDepth)
			}
		}
	}
	return t
}

func (p *varManager) makeVarsContext(ctx *Context) varsContext {
	if p.tcache == nil {
		p.tcache = makeVarsContextType(p.vlist, ctx)
	}
	return reflect.New(p.tcache).Elem()
}

func (ctx *varScope) addrVar(idx uint32) interface{} {
	return ctx.vars.Field(int(idx)).Addr().Interface()
}

func (ctx *varScope) getVar(idx uint32) interface{} {
	return ctx.vars.Field(int(idx)).Interface()
}

func (ctx *varScope) setVar(idx uint32, v interface{}) {
	x := ctx.vars.Field(int(idx))
	setValue(x, v)
}

func (ctx *varScope) setVarField(idx uint32, v interface{}, index interface{}) {
	x := ctx.vars.Field(int(idx))
	setValue(x.FieldByIndex(index.([]int)), v)
}

func setValue(x reflect.Value, v interface{}) {
	if v != nil {
		x.Set(reflect.ValueOf(v))
		return
	}
	x.Set(reflect.Zero(x.Type()))
}

func getValueOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t)
}

func getKeyOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.Key())
}

func getElementOf(v interface{}, t reflect.Type) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.Elem())
}

func getArgOf(v interface{}, t reflect.Type, i int) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(t.In(i))
}

func getRetOf(v interface{}, fun *FuncInfo, i int) reflect.Value {
	if v != nil {
		return reflect.ValueOf(v)
	}
	return reflect.Zero(fun.vlist[i].typ)
}

// -----------------------------------------------------------------------------
