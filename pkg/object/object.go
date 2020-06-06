package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ripta/mwnci/pkg/ast"
)

type Type string

const (
	NullObj  = "NULL"
	ErrorObj = "ERROR"

	FloatObj   = "FLOAT"
	IntegerObj = "INTEGER"
	BooleanObj = "BOOLEAN"
	StringObj  = "STRING"

	ReturnValueObj = "RETURN_VALUE"

	FunctionObj = "FUNCTION"
	BuiltinObj  = "BUILTIN"

	ArrayObj = "ARRAY"
	HashObj  = "HASH"
)

type HashKey struct {
	Type  Type
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() Type
	Inspect() string
	Location() ast.Location
}

type Integer struct {
	Value int64
	Loc   ast.Location
}

func (i *Integer) Type() Type             { return IntegerObj }
func (i *Integer) Inspect() string        { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Location() ast.Location { return i.Loc }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Float struct {
	Value float64
	Loc   ast.Location
}

func (f *Float) Type() Type             { return FloatObj }
func (f *Float) Inspect() string        { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Location() ast.Location { return f.Loc }
func (f *Float) HashKey() HashKey {
	return HashKey{Type: f.Type(), Value: uint64(f.Value)}
}

type Boolean struct {
	Value bool
	Loc   ast.Location
}

func (b *Boolean) Type() Type             { return BooleanObj }
func (b *Boolean) Location() ast.Location { return b.Loc }
func (b *Boolean) Inspect() string        { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Null struct {
	Loc ast.Location
}

func (n *Null) Type() Type             { return NullObj }
func (n *Null) Location() ast.Location { return n.Loc }
func (n *Null) Inspect() string        { return "null" }

type ReturnValue struct {
	Value Object
	Loc   ast.Location
}

func (rv *ReturnValue) Type() Type             { return ReturnValueObj }
func (rv *ReturnValue) Location() ast.Location { return rv.Loc }
func (rv *ReturnValue) Inspect() string        { return rv.Value.Inspect() }

type Error struct {
	Message string
	Loc     ast.Location
}

func (e *Error) Type() Type             { return ErrorObj }
func (e *Error) Location() ast.Location { return e.Loc }
func (e *Error) Inspect() string        { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	Loc        ast.Location
}

func (f *Function) Type() Type             { return FunctionObj }
func (f *Function) Location() ast.Location { return f.Loc }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
	Loc   ast.Location
}

func (s *String) Type() Type             { return StringObj }
func (s *String) Location() ast.Location { return s.Loc }
func (s *String) Inspect() string        { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Builtin func(args ...Object) Object

func (b Builtin) Type() Type             { return BuiltinObj }
func (b Builtin) Location() ast.Location { return ast.Location{} }
func (b Builtin) Inspect() string        { return "builtin function" }

type Array struct {
	Elements []Object
	Loc      ast.Location
}

func (ao *Array) Type() Type             { return ArrayObj }
func (ao *Array) Location() ast.Location { return ao.Loc }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

var _ Object = (*Array)(nil)

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
	Loc   ast.Location
}

func (h *Hash) Type() Type             { return HashObj }
func (h *Hash) Location() ast.Location { return h.Loc }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		p := fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect())
		pairs = append(pairs, p)
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
