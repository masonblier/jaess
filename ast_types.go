package jsparser

import (
	"fmt"
	// "encoding/json"
)

// an AstNode
type AstNode interface {
}

// types of ast nodes
type AstType int

const (
	AST_UNKNOWN AstType = iota
	LITERAL
	IDENTIFIER
	PROPERTY
	PROGRAM
	FUNCTION_DECLARATION
	BLOCK_STATEMENT
	EXPRESSION_STATEMENT
	ASSIGNMENT_EXPRESSION
	MEMBER_EXPRESSION
	THIS_EXPRESSION
	FUNCTION_EXPRESSION
	CALL_EXPRESSION
	VARIABLE_DECLARATION
	VARIABLE_DECLARATOR
	NEW_EXPRESSION
	OBJECT_EXPRESSION
	UNARY_EXPRESSION
	BINARY_EXPRESSION
)

func (self AstType) String() string {
	switch self {

	case AST_UNKNOWN:
		return "AstUnknown"
	case LITERAL:
		return "Literal"
	case IDENTIFIER:
		return "Identifier"
	case PROPERTY:
		return "Property"
	case PROGRAM:
		return "Program"
	case FUNCTION_DECLARATION:
		return "FunctionDeclaration"
	case BLOCK_STATEMENT:
		return "BlockStatement"
	case EXPRESSION_STATEMENT:
		return "ExpressionStatement"
	case ASSIGNMENT_EXPRESSION:
		return "AssignmentExpression"
	case MEMBER_EXPRESSION:
		return "MemberExpression"
	case THIS_EXPRESSION:
		return "ThisExpression"
	case FUNCTION_EXPRESSION:
		return "FunctionExpression"
	case CALL_EXPRESSION:
		return "CallExpression"
	case VARIABLE_DECLARATION:
		return "VariableDeclaration"
	case VARIABLE_DECLARATOR:
		return "VariableDeclarator"
	case NEW_EXPRESSION:
		return "NewExpression"
	case OBJECT_EXPRESSION:
		return "ObjectExpression"
	case UNARY_EXPRESSION:
		return "UnaryExpression"
	case BINARY_EXPRESSION:
		return "BinaryExpression"

	}
	return "<#error: bad value>"
}

func (self AstType) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", self)
	return []byte(str), nil
}

type LiteralString struct {
	Type  AstType `json:"type"`
	Value string  `json:"value"`
	Raw   string  `json:"raw"`
}

type LiteralNumber struct {
	Type  AstType `json:"type"`
	Value float64 `json:"value"`
	Raw   string  `json:"raw"`
}

type Identifier struct {
	Type AstType `json:"type"`
	Name string  `json:"name"`
}

type Property struct {
	Type  AstType `json:"type"`
	Key   AstNode `json:"key"`
	Value AstNode `json:"value"`
	Kind  string  `json:"kind"`
}

type Program struct {
	Type AstType   `json:"type"`
	Body []AstNode `json:"body"`
}

type FunctionDeclaration struct {
	Type       AstType   `json:"type"`
	Id         AstNode   `json:"id"`
	Params     []AstNode `json:"params"`
	Defaults   []AstNode `json:"defaults"`
	Body       AstNode   `json:"body"`
	Rest       AstNode   `json:"rest"`
	Generator  bool      `json:"generator"`
	Expression bool      `json:"expression"`
}

type BlockStatement struct {
	Type AstType   `json:"type"`
	Body []AstNode `json:"body"`
}

type ExpressionStatement struct {
	Type       AstType `json:"type"`
	Expression AstNode `json:"expression"`
}

type AssignmentExpression struct {
	Type     AstType `json:"type"`
	Operator string  `json:"operator"`
	Left     AstNode `json:"left"`
	Right    AstNode `json:"right"`
}

type MemberExpression struct {
	Type     AstType `json:"type"`
	Computed bool    `json:"computed"`
	Object   AstNode `json:"object"`
	Property AstNode `json:"property"`
}

type ThisExpression struct {
	Type AstType `json:"type"`
}

type FunctionExpression struct {
	Type       AstType   `json:"type"`
	Id         AstNode   `json:"id"`
	Params     []AstNode `json:"params"`
	Defaults   []AstNode `json:"defaults"`
	Body       AstNode   `json:"body"`
	Rest       AstNode   `json:"rest"`
	Generator  bool      `json:"generator"`
	Expression bool      `json:"expression"`
}

type CallExpression struct {
	Type      AstType   `json:"type"`
	Callee    AstNode   `json:"callee"`
	Arguments []AstNode `json:"arguments"`
}

type VariableDeclaration struct {
	Type         AstType   `json:"type"`
	Declarations []AstNode `json:"declarations"`
	Kind         string    `json:"kind"`
}

type VariableDeclarator struct {
	Type AstType `json:"type"`
	Id   AstNode `json:"id"`
	Init AstNode `json:"init"`
}

type NewExpression struct {
	Type      AstType   `json:"type"`
	Callee    AstNode   `json:"callee"`
	Arguments []AstNode `json:"arguments"`
}

type ObjectExpression struct {
	Type       AstType   `json:"type"`
	Properties []AstNode `json:"properties"`
}

type UnaryExpression struct {
	Type     AstType `json:"type"`
	Operator string  `json:"operator"`
	Argument AstNode `json:"argument"`
	Prefix   bool    `json:"prefix"`
}

type BinaryExpression struct {
	Type     AstType `json:"type"`
	Operator string  `json:"operator"`
	Left     AstNode `json:"left"`
	Right    AstNode `json:"right"`
}
