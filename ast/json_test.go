package ast_test

import (
	"encoding/json"
	"math"
	"regexp"
	"testing"
	"time"

	"github.com/InfluxCommunity/flux/ast"
	"github.com/InfluxCommunity/flux/ast/asttest"
	"github.com/google/go-cmp/cmp"
)

func TestJSONMarshal(t *testing.T) {
	testCases := []struct {
		name string
		node ast.Node
		want string
	}{
		{
			name: "string interpolation",
			node: &ast.StringExpression{
				Parts: []ast.StringExpressionPart{
					&ast.TextPart{
						Value: "a = ",
					},
					&ast.InterpolatedPart{
						Expression: &ast.Identifier{
							Name: "a",
						},
					},
				},
			},
			want: `{"type":"StringExpression","parts":[{"type":"TextPart","value":"a = "},{"type":"InterpolatedPart","expression":{"type":"Identifier","name":"a"}}]}`,
		},
		{
			name: "paren expression",
			node: &ast.ParenExpression{
				Expression: &ast.StringExpression{
					Parts: []ast.StringExpressionPart{
						&ast.TextPart{
							Value: "a = ",
						},
						&ast.InterpolatedPart{
							Expression: &ast.Identifier{
								Name: "a",
							},
						},
					},
				},
			},
			want: `{"type":"ParenExpression","expression":{"type":"StringExpression","parts":[{"type":"TextPart","value":"a = "},{"type":"InterpolatedPart","expression":{"type":"Identifier","name":"a"}}]}}`,
		},
		{
			name: "simple package",
			node: &ast.Package{
				Package: "foo",
			},
			want: `{"type":"Package","package":"foo","files":null}`,
		},
		{
			name: "package path",
			node: &ast.Package{
				Path:    "bar/foo",
				Package: "foo",
			},
			want: `{"type":"Package","path":"bar/foo","package":"foo","files":null}`,
		},
		{
			name: "simple file",
			node: &ast.File{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.StringLiteral{Value: "hello"},
					},
				},
			},
			want: `{"type":"File","package":null,"imports":null,"body":[{"type":"ExpressionStatement","expression":{"type":"StringLiteral","value":"hello"}}]}`,
		},
		{
			name: "file",
			node: &ast.File{
				Metadata: "parser-type=none",
				Package: &ast.PackageClause{
					Name: &ast.Identifier{Name: "foo"},
				},
				Imports: []*ast.ImportDeclaration{{
					As:   &ast.Identifier{Name: "b"},
					Path: &ast.StringLiteral{Value: "path/bar"},
				}},
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.StringLiteral{Value: "hello"},
					},
				},
			},
			want: `{"type":"File","metadata":"parser-type=none","package":{"type":"PackageClause","name":{"type":"Identifier","name":"foo"}},"imports":[{"type":"ImportDeclaration","as":{"type":"Identifier","name":"b"},"path":{"type":"StringLiteral","value":"path/bar"}}],"body":[{"type":"ExpressionStatement","expression":{"type":"StringLiteral","value":"hello"}}]}`,
		},
		{
			name: "block",
			node: &ast.Block{
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: &ast.StringLiteral{Value: "hello"},
					},
				},
			},
			want: `{"type":"Block","body":[{"type":"ExpressionStatement","expression":{"type":"StringLiteral","value":"hello"}}]}`,
		},
		{
			name: "expression statement",
			node: &ast.ExpressionStatement{
				Expression: &ast.StringLiteral{Value: "hello"},
			},
			want: `{"type":"ExpressionStatement","expression":{"type":"StringLiteral","value":"hello"}}`,
		},
		{
			name: "return statement",
			node: &ast.ReturnStatement{
				Argument: &ast.StringLiteral{Value: "hello"},
			},
			want: `{"type":"ReturnStatement","argument":{"type":"StringLiteral","value":"hello"}}`,
		},
		{
			name: "option statement",
			node: &ast.OptionStatement{
				Assignment: &ast.VariableAssignment{
					ID: &ast.Identifier{Name: "task"},
					Init: &ast.ObjectExpression{
						Properties: []*ast.Property{
							{
								Key:   &ast.Identifier{Name: "name"},
								Value: &ast.StringLiteral{Value: "foo"},
							},
							{
								Key: &ast.Identifier{Name: "every"},
								Value: &ast.DurationLiteral{
									Values: []ast.Duration{
										{
											Magnitude: 1,
											Unit:      "h",
										},
									},
								},
							},
						},
					},
				},
			},
			want: `{"type":"OptionStatement","assignment":{"type":"VariableAssignment","id":{"type":"Identifier","name":"task"},"init":{"type":"ObjectExpression","properties":[{"type":"Property","key":{"type":"Identifier","name":"name"},"value":{"type":"StringLiteral","value":"foo"}},{"type":"Property","key":{"type":"Identifier","name":"every"},"value":{"type":"DurationLiteral","values":[{"magnitude":1,"unit":"h"}]}}]}}}`,
		},
		{
			name: "builtin statement",
			node: &ast.BuiltinStatement{
				ID: &ast.Identifier{Name: "task"},
				Ty: ast.TypeExpression{
					Ty: &ast.NamedType{
						BaseNode: ast.BaseNode{},
						ID: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "int",
						},
					},
				},
			},
			want: `{"type":"BuiltinStatement","id":{"type":"Identifier","name":"task"},"ty":{"type":"TypeExpression","monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}},"constraints":null}}`,
		},
		{
			name: "NamedType",
			node: &ast.NamedType{
				BaseNode: ast.BaseNode{},
				ID: &ast.Identifier{
					BaseNode: ast.BaseNode{},
					Name:     "int",
				},
			},
			want: `{"type":"NamedType","name":{"type":"Identifier","name":"int"}}`,
		},
		{
			name: "TvarType",
			node: &ast.TvarType{
				BaseNode: ast.BaseNode{},
				ID: &ast.Identifier{
					BaseNode: ast.BaseNode{},
					Name:     "A",
				},
			},
			want: `{"type":"TvarType","name":{"type":"Identifier","name":"A"}}`,
		},
		{
			name: "ArrayType",
			node: &ast.ArrayType{
				BaseNode: ast.BaseNode{},
				ElementType: &ast.ArrayType{
					BaseNode: ast.BaseNode{},
					ElementType: &ast.NamedType{
						BaseNode: ast.BaseNode{},
						ID: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "int",
						},
					},
				},
			},
			// r#"{"type":"ArrayType","element":{"type":"ArrayType","element":{"type":"NamedType","name":{"type":"Identifier","name":"A"}}}}"#
			want: `{"type":"ArrayType","element":{"type":"ArrayType","element":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}}`,
		},
		{
			name: "RecordType",
			node: &ast.RecordType{
				BaseNode: ast.BaseNode{},
				Properties: []*ast.PropertyType{
					{
						BaseNode: ast.BaseNode{},
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "A",
						},
						Ty: &ast.NamedType{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "int",
							},
						},
					},
				},
				Tvar: &ast.Identifier{
					BaseNode: ast.BaseNode{},
					Name:     "A",
				},
			},
			want: `{"type":"RecordType","properties":[{"type":"PropertyType","name":{"type":"Identifier","name":"A"},"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}],"tvar":{"type":"Identifier","name":"A"}}`,
		},
		{
			name: "RecordType_NoTvar_NoProp",
			node: &ast.RecordType{
				BaseNode:   ast.BaseNode{},
				Properties: []*ast.PropertyType{},
				Tvar:       nil,
			},
			//{"type":"RecordType","properties":[]}"#);
			want: `{"type":"RecordType","properties":[]}`,
		},
		{
			name: "RecordType_NoTvar",
			node: &ast.RecordType{
				BaseNode: ast.BaseNode{},
				Properties: []*ast.PropertyType{
					{
						BaseNode: ast.BaseNode{},
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "A",
						},
						Ty: &ast.NamedType{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "int",
							},
						},
					},
				},
				Tvar: nil,
			},
			// {"type":"RecordType","properties":[{"name":{"name":"A"},"monotype":{"type":"NamedType","name":{"name":"int"}}}]}"#
			want: `{"type":"RecordType","properties":[{"type":"PropertyType","name":{"type":"Identifier","name":"A"},"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}]}`,
		},
		{
			name: "FunctionType_NoParams",
			node: &ast.FunctionType{
				BaseNode:   ast.BaseNode{},
				Parameters: []*ast.ParameterType{},
				Return: &ast.NamedType{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "int",
					},
				},
			},
			want: `{"type":"FunctionType","parameters":[],"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}`,
		},
		{
			name: "FunctionType_Required",
			node: &ast.FunctionType{
				BaseNode: ast.BaseNode{},
				Parameters: []*ast.ParameterType{
					{
						BaseNode: ast.BaseNode{},
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "B",
						},
						Ty: &ast.NamedType{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "string",
							},
						},
						Kind: "Required"},
				},
				Return: &ast.NamedType{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "uint",
					},
				},
			},
			want: `{"type":"FunctionType","parameters":[{"type":"Required","name":{"type":"Identifier","name":"B"},"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"string"}}}],"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"uint"}}}`,
		},
		{
			name: "FunctionType_Optional",
			node: &ast.FunctionType{
				BaseNode: ast.BaseNode{},
				Parameters: []*ast.ParameterType{
					{
						BaseNode: ast.BaseNode{},
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "A",
						},
						Ty: &ast.NamedType{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "int",
							},
						},
						Kind: "Optional"},
				},
				Return: &ast.NamedType{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "int",
					},
				},
			},
			//r#"{"type":"FunctionType","parameters":[{"type":"Optional","name":{"name":"A"},"monotype":{"type":"NamedType","name":{"name":"int"}}}],"monotype":{"type":"NamedType","name":{"name":"int"}}}"#
			want: `{"type":"FunctionType","parameters":[{"type":"Optional","name":{"type":"Identifier","name":"A"},"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}],"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}`,
		},
		{
			name: "FunctionType_Named_Pipe",
			node: &ast.FunctionType{
				BaseNode: ast.BaseNode{},
				Parameters: []*ast.ParameterType{
					{
						BaseNode: ast.BaseNode{},
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "A",
						},
						Ty: &ast.NamedType{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "int",
							},
						},
						Kind: "Pipe"},
				},
				Return: &ast.NamedType{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "int",
					},
				},
			},
			want: `{"type":"FunctionType","parameters":[{"type":"Pipe","name":{"type":"Identifier","name":"A"},"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}],"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}`,
		},
		{
			name: "FunctionType_UnNamed_Pipe",
			node: &ast.FunctionType{
				BaseNode: ast.BaseNode{},
				Parameters: []*ast.ParameterType{
					{
						BaseNode: ast.BaseNode{},
						Name:     nil,
						Ty: &ast.NamedType{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "int",
							},
						},
						Kind: "Pipe"},
				},
				Return: &ast.NamedType{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "int",
					},
				},
			},
			want: `{"type":"FunctionType","parameters":[{"type":"Pipe","monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}],"monotype":{"type":"NamedType","name":{"type":"Identifier","name":"int"}}}`,
		},
		{
			name: "TypeExpression Test",
			node: &ast.TypeExpression{
				BaseNode: ast.BaseNode{},
				Ty: &ast.FunctionType{
					BaseNode: ast.BaseNode{},
					Parameters: []*ast.ParameterType{
						{
							BaseNode: ast.BaseNode{},
							Name: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "a",
							},
							Ty: &ast.TvarType{
								BaseNode: ast.BaseNode{},
								ID: &ast.Identifier{
									BaseNode: ast.BaseNode{},
									Name:     "T",
								},
							},
							Kind: "Required",
						},
						{
							BaseNode: ast.BaseNode{},
							Name: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "b",
							},
							Ty: &ast.TvarType{
								BaseNode: ast.BaseNode{},
								ID: &ast.Identifier{
									BaseNode: ast.BaseNode{},
									Name:     "T",
								},
							},
							Kind: "Required",
						},
					},
					Return: &ast.TvarType{
						BaseNode: ast.BaseNode{},
						ID: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "T",
						},
					},
				},
				Constraints: []*ast.TypeConstraint{
					{
						BaseNode: ast.BaseNode{},
						Tvar: &ast.Identifier{
							BaseNode: ast.BaseNode{},
							Name:     "T",
						},
						Kinds: []*ast.Identifier{
							{
								BaseNode: ast.BaseNode{},
								Name:     "Addable",
							},
							{
								BaseNode: ast.BaseNode{},
								Name:     "Divisible",
							},
						},
					},
				},
			},
			want: `{"type":"TypeExpression","monotype":{"type":"FunctionType","parameters":[{"type":"Required","name":{"type":"Identifier","name":"a"},"monotype":{"type":"TvarType","name":{"type":"Identifier","name":"T"}}},{"type":"Required","name":{"type":"Identifier","name":"b"},"monotype":{"type":"TvarType","name":{"type":"Identifier","name":"T"}}}],"monotype":{"type":"TvarType","name":{"type":"Identifier","name":"T"}}},"constraints":[{"type":"TypeConstraint","tvar":{"type":"Identifier","name":"T"},"kinds":[{"type":"Identifier","name":"Addable"},{"type":"Identifier","name":"Divisible"}]}]}`,
		},
		{
			name: "test statement",
			node: &ast.TestStatement{
				Assignment: &ast.VariableAssignment{
					ID: &ast.Identifier{Name: "mean"},
					Init: &ast.ObjectExpression{
						Properties: []*ast.Property{
							{
								Key: &ast.Identifier{
									Name: "want",
								},
								Value: &ast.IntegerLiteral{
									Value: 0,
								},
							},
							{
								Key: &ast.Identifier{
									Name: "got",
								},
								Value: &ast.IntegerLiteral{
									Value: 0,
								},
							},
						},
					},
				},
			},
			want: `{"type":"TestStatement","assignment":{"type":"VariableAssignment","id":{"type":"Identifier","name":"mean"},"init":{"type":"ObjectExpression","properties":[{"type":"Property","key":{"type":"Identifier","name":"want"},"value":{"type":"IntegerLiteral","value":"0"}},{"type":"Property","key":{"type":"Identifier","name":"got"},"value":{"type":"IntegerLiteral","value":"0"}}]}}}`,
		},
		{
			name: "test case statement",
			node: &ast.TestCaseStatement{

				ID: &ast.Identifier{Name: "test case statement id"},
				Extends: &ast.StringLiteral{
					Value: "extended test case",
				},
				Block: &ast.Block{
					Body: []ast.Statement{
						&ast.ExpressionStatement{
							Expression: &ast.StringLiteral{Value: "expression statement"},
						},
					},
				},
			},
			want: `{"type":"TestCaseStatement","id":{"type":"Identifier","name":"test case statement id"},"extends":{"type":"StringLiteral","value":"extended test case"},"block":{"type":"Block","body":[{"type":"ExpressionStatement","expression":{"type":"StringLiteral","value":"expression statement"}}]}}`,
		},
		{
			name: "qualified option statement",
			node: &ast.OptionStatement{
				Assignment: &ast.MemberAssignment{
					Member: &ast.MemberExpression{
						Object: &ast.Identifier{
							Name: "alert",
						},
						Property: &ast.Identifier{
							Name: "state",
						},
					},
					Init: &ast.StringLiteral{
						Value: "Warning",
					},
				},
			},
			want: `{"type":"OptionStatement","assignment":{"type":"MemberAssignment","member":{"type":"MemberExpression","object":{"type":"Identifier","name":"alert"},"property":{"type":"Identifier","name":"state"}},"init":{"type":"StringLiteral","value":"Warning"}}}`,
		},
		{
			name: "variable assignment",
			node: &ast.VariableAssignment{
				ID:   &ast.Identifier{Name: "a"},
				Init: &ast.StringLiteral{Value: "hello"},
			},
			want: `{"type":"VariableAssignment","id":{"type":"Identifier","name":"a"},"init":{"type":"StringLiteral","value":"hello"}}`,
		},
		{
			name: "call expression",
			node: &ast.CallExpression{
				Callee:    &ast.Identifier{Name: "a"},
				Arguments: []ast.Expression{&ast.StringLiteral{Value: "hello"}},
			},
			want: `{"type":"CallExpression","callee":{"type":"Identifier","name":"a"},"arguments":[{"type":"StringLiteral","value":"hello"}]}`,
		},
		{
			name: "pipe expression",
			node: &ast.PipeExpression{
				Argument: &ast.Identifier{Name: "a"},
				Call: &ast.CallExpression{
					Callee:    &ast.Identifier{Name: "a"},
					Arguments: []ast.Expression{&ast.StringLiteral{Value: "hello"}},
				},
			},
			want: `{"type":"PipeExpression","argument":{"type":"Identifier","name":"a"},"call":{"type":"CallExpression","callee":{"type":"Identifier","name":"a"},"arguments":[{"type":"StringLiteral","value":"hello"}]}}`,
		},
		{
			name: "member expression with identifier",
			node: &ast.MemberExpression{
				Object:   &ast.Identifier{Name: "a"},
				Property: &ast.Identifier{Name: "b"},
			},
			want: `{"type":"MemberExpression","object":{"type":"Identifier","name":"a"},"property":{"type":"Identifier","name":"b"}}`,
		},
		{
			name: "member expression with string literal",
			node: &ast.MemberExpression{
				Object:   &ast.Identifier{Name: "a"},
				Property: &ast.StringLiteral{Value: "b"},
			},
			want: `{"type":"MemberExpression","object":{"type":"Identifier","name":"a"},"property":{"type":"StringLiteral","value":"b"}}`,
		},
		{
			name: "index expression",
			node: &ast.IndexExpression{
				Array: &ast.Identifier{Name: "a"},
				Index: &ast.IntegerLiteral{Value: 3},
			},
			want: `{"type":"IndexExpression","array":{"type":"Identifier","name":"a"},"index":{"type":"IntegerLiteral","value":"3"}}`,
		},
		{
			name: "arrow function expression",
			node: &ast.FunctionExpression{
				Params: []*ast.Property{{Key: &ast.Identifier{Name: "a"}}},
				Body:   &ast.StringLiteral{Value: "hello"},
			},
			want: `{"type":"FunctionExpression","params":[{"type":"Property","key":{"type":"Identifier","name":"a"},"value":null}],"body":{"type":"StringLiteral","value":"hello"}}`,
		},
		{
			name: "binary expression",
			node: &ast.BinaryExpression{
				Operator: ast.AdditionOperator,
				Left:     &ast.StringLiteral{Value: "hello"},
				Right:    &ast.StringLiteral{Value: "world"},
			},
			want: `{"type":"BinaryExpression","operator":"+","left":{"type":"StringLiteral","value":"hello"},"right":{"type":"StringLiteral","value":"world"}}`,
		},
		{
			name: "unary expression",
			node: &ast.UnaryExpression{
				Operator: ast.NotOperator,
				Argument: &ast.BooleanLiteral{Value: true},
			},
			want: `{"type":"UnaryExpression","operator":"not","argument":{"type":"BooleanLiteral","value":true}}`,
		},
		{
			name: "logical expression",
			node: &ast.LogicalExpression{
				Operator: ast.OrOperator,
				Left:     &ast.BooleanLiteral{Value: false},
				Right:    &ast.BooleanLiteral{Value: true},
			},
			want: `{"type":"LogicalExpression","operator":"or","left":{"type":"BooleanLiteral","value":false},"right":{"type":"BooleanLiteral","value":true}}`,
		},
		{
			name: "array expression",
			node: &ast.ArrayExpression{
				Elements: []ast.Expression{&ast.StringLiteral{Value: "hello"}},
			},
			want: `{"type":"ArrayExpression","elements":[{"type":"StringLiteral","value":"hello"}]}`,
		},
		{
			name: "dict expression",
			node: &ast.DictExpression{
				Elements: []*ast.DictItem{{Key: &ast.StringLiteral{Value: "a"}, Val: &ast.IntegerLiteral{Value: 0}}, {Key: &ast.StringLiteral{Value: "b"}, Val: &ast.IntegerLiteral{Value: 1}}, {Key: &ast.StringLiteral{Value: "c"}, Val: &ast.IntegerLiteral{Value: 2}}},
			},
			want: `{"type":"DictExpression","elements":[{"type":"DictItem","key":{"type":"StringLiteral","value":"a"},"val":{"type":"IntegerLiteral","value":"0"}},{"type":"DictItem","key":{"type":"StringLiteral","value":"b"},"val":{"type":"IntegerLiteral","value":"1"}},{"type":"DictItem","key":{"type":"StringLiteral","value":"c"},"val":{"type":"IntegerLiteral","value":"2"}}]}`,
		},
		{
			name: "object expression",
			node: &ast.ObjectExpression{
				Properties: []*ast.Property{{
					Key:   &ast.Identifier{Name: "a"},
					Value: &ast.StringLiteral{Value: "hello"},
				}},
			},
			want: `{"type":"ObjectExpression","properties":[{"type":"Property","key":{"type":"Identifier","name":"a"},"value":{"type":"StringLiteral","value":"hello"}}]}`,
		},
		{
			name: "object expression with string literal key",
			node: &ast.ObjectExpression{
				Properties: []*ast.Property{{
					Key:   &ast.StringLiteral{Value: "a"},
					Value: &ast.StringLiteral{Value: "hello"},
				}},
			},
			want: `{"type":"ObjectExpression","properties":[{"type":"Property","key":{"type":"StringLiteral","value":"a"},"value":{"type":"StringLiteral","value":"hello"}}]}`,
		},
		{
			name: "object expression implicit keys",
			node: &ast.ObjectExpression{
				Properties: []*ast.Property{{
					Key: &ast.Identifier{Name: "a"},
				}},
			},
			want: `{"type":"ObjectExpression","properties":[{"type":"Property","key":{"type":"Identifier","name":"a"},"value":null}]}`,
		},
		{
			name: "conditional expression",
			node: &ast.ConditionalExpression{
				Test:       &ast.BooleanLiteral{Value: true},
				Alternate:  &ast.StringLiteral{Value: "false"},
				Consequent: &ast.StringLiteral{Value: "true"},
			},
			want: `{"type":"ConditionalExpression","test":{"type":"BooleanLiteral","value":true},"consequent":{"type":"StringLiteral","value":"true"},"alternate":{"type":"StringLiteral","value":"false"}}`,
		},
		{
			name: "property",
			node: &ast.Property{
				Key:   &ast.Identifier{Name: "a"},
				Value: &ast.StringLiteral{Value: "hello"},
			},
			want: `{"type":"Property","key":{"type":"Identifier","name":"a"},"value":{"type":"StringLiteral","value":"hello"}}`,
		},
		{
			name: "identifier",
			node: &ast.Identifier{
				Name: "a",
			},
			want: `{"type":"Identifier","name":"a"}`,
		},
		{
			name: "string literal",
			node: &ast.StringLiteral{
				Value: "hello",
			},
			want: `{"type":"StringLiteral","value":"hello"}`,
		},
		{
			name: "boolean literal",
			node: &ast.BooleanLiteral{
				Value: true,
			},
			want: `{"type":"BooleanLiteral","value":true}`,
		},
		{
			name: "float literal",
			node: &ast.FloatLiteral{
				Value: 42.1,
			},
			want: `{"type":"FloatLiteral","value":42.1}`,
		},
		{
			name: "integer literal",
			node: &ast.IntegerLiteral{
				Value: math.MaxInt64,
			},
			want: `{"type":"IntegerLiteral","value":"9223372036854775807"}`,
		},
		{
			name: "unsigned integer literal",
			node: &ast.UnsignedIntegerLiteral{
				Value: math.MaxUint64,
			},
			want: `{"type":"UnsignedIntegerLiteral","value":"18446744073709551615"}`,
		},
		{
			name: "regexp literal",
			node: &ast.RegexpLiteral{
				Value: regexp.MustCompile(`.*`),
			},
			want: `{"type":"RegexpLiteral","value":".*"}`,
		},
		{
			name: "duration literal",
			node: &ast.DurationLiteral{
				Values: []ast.Duration{
					{
						Magnitude: 1,
						Unit:      "h",
					},
					{
						Magnitude: 1,
						Unit:      "h",
					},
				},
			},
			want: `{"type":"DurationLiteral","values":[{"magnitude":1,"unit":"h"},{"magnitude":1,"unit":"h"}]}`,
		},
		{
			name: "datetime literal",
			node: &ast.DateTimeLiteral{
				Value: time.Date(2017, 8, 8, 8, 8, 8, 8, time.UTC),
			},
			want: `{"type":"DateTimeLiteral","value":"2017-08-08T08:08:08.000000008Z"}`,
		},
		{
			name: "object expression with source locations and errors",
			node: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{
					Loc: &ast.SourceLocation{
						File:   "foo.flux",
						Start:  ast.Position{Line: 1, Column: 1},
						End:    ast.Position{Line: 1, Column: 13},
						Source: "{a: \"hello\"}",
					},
				},
				Properties: []*ast.Property{{
					BaseNode: ast.BaseNode{
						Loc: &ast.SourceLocation{
							File:   "foo.flux",
							Start:  ast.Position{Line: 1, Column: 2},
							End:    ast.Position{Line: 1, Column: 12},
							Source: "a: \"hello\"",
						},
						Errors: []ast.Error{{Msg: "an error"}},
					},
					Key: &ast.Identifier{BaseNode: ast.BaseNode{
						Loc: &ast.SourceLocation{
							File:   "foo.flux",
							Start:  ast.Position{Line: 1, Column: 2},
							End:    ast.Position{Line: 1, Column: 3},
							Source: "a",
						},
					},
						Name: "a",
					},
					Value: &ast.StringLiteral{BaseNode: ast.BaseNode{
						Loc: &ast.SourceLocation{
							File:   "foo.flux",
							Start:  ast.Position{Line: 1, Column: 5},
							End:    ast.Position{Line: 1, Column: 12},
							Source: "\"hello\"",
						},
						Errors: []ast.Error{{Msg: "an error"}, {Msg: "another error"}},
					},
						Value: "hello",
					},
				}},
			},
			want: `{"type":"ObjectExpression","location":{"file":"foo.flux","start":{"line":1,"column":1},"end":{"line":1,"column":13},"source":"{a: \"hello\"}"},"properties":[{"type":"Property","location":{"file":"foo.flux","start":{"line":1,"column":2},"end":{"line":1,"column":12},"source":"a: \"hello\""},"errors":[{"msg":"an error"}],"key":{"type":"Identifier","location":{"file":"foo.flux","start":{"line":1,"column":2},"end":{"line":1,"column":3},"source":"a"},"name":"a"},"value":{"type":"StringLiteral","location":{"file":"foo.flux","start":{"line":1,"column":5},"end":{"line":1,"column":12},"source":"\"hello\""},"errors":[{"msg":"an error"},{"msg":"another error"}],"value":"hello"}}]}`,
		},
		{
			name: "Comments in BaseNode",
			node: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Comments: []ast.Comment{{Text: "This is a comment"}},
				},
				Name: "A",
			},
			want: `{"type":"Identifier","comments":[{"text":"This is a comment"}],"name":"A"}`,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.node)
			if err != nil {
				t.Fatal(err)
			}
			if got := string(data); got != tc.want {
				t.Errorf("unexpected json data:\nwant:%s\ngot: %s\n", tc.want, got)
			}
			node, err := ast.UnmarshalNode(data)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tc.node, node, asttest.CompareOptions...) {
				t.Errorf("unexpected node after unmarshalling: -want/+got:\n%s", cmp.Diff(tc.node, node, asttest.CompareOptions...))
			}
		})
	}
}
