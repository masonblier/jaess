package jaess

import (
	"fmt"
	// "io"
	"bytes"
	"encoding/json"
)

// an AstNode
type AstNode interface {
	AstType() AstType
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
	EMPTY_STATEMENT
	BLOCK_STATEMENT
	EXPRESSION_STATEMENT
	IF_STATEMENT
	FOR_STATEMENT
	// LABELED_STATEMENT
	// BREAK_STATEMENT
	// CONTINUE_STATEMENT
	// WITH_STATEMENT
	// SWITCH_STATEMENT
	RETURN_STATEMENT
	// THROW_STATEMENT
	ASSIGNMENT_EXPRESSION
	MEMBER_EXPRESSION
	THIS_EXPRESSION
	FUNCTION_EXPRESSION
	CALL_EXPRESSION
	VARIABLE_DECLARATION
	VARIABLE_DECLARATOR
	NEW_EXPRESSION
	ARRAY_EXPRESSION
	OBJECT_EXPRESSION
	UNARY_EXPRESSION
	BINARY_EXPRESSION
	UPDATE_EXPRESSION
)

type AstNodeMeta struct {
	Type  AstType `json:"type"`
}

type LiteralNull struct {
	AstNodeMeta
	Value interface{}  `json:"value"`
	Raw   string  `json:"raw"`
}

type LiteralString struct {
	AstNodeMeta
	Value string  `json:"value"`
	Raw   string  `json:"raw"`
}

type LiteralNumber struct {
	AstNodeMeta
	Value float64 `json:"value"`
	Raw   string  `json:"raw"`
}

type Identifier struct {
	AstNodeMeta
	Name string  `json:"name"`
}

type Property struct {
	AstNodeMeta
	Key   AstNode `json:"key"`
	Value AstNode `json:"value"`
	Kind  string  `json:"kind"`
}

type Program struct {
	AstNodeMeta
	Body []AstNode `json:"body"`
}

type FunctionDeclaration struct {
	AstNodeMeta
	Id         AstNode   `json:"id"`
	Params     []AstNode `json:"params"`
	Defaults   []AstNode `json:"defaults"`
	Body       AstNode   `json:"body"`
	Rest       AstNode   `json:"rest"`
	Generator  bool      `json:"generator"`
	Expression bool      `json:"expression"`
	Source     string 	 `json:"-"`
}

type EmptyStatement struct {
	AstNodeMeta
}

type BlockStatement struct {
	AstNodeMeta
	Body []AstNode `json:"body"`
}

type ExpressionStatement struct {
	AstNodeMeta
	Expression AstNode `json:"expression"`
}

type IfStatement struct {
	AstNodeMeta
	Test AstNode `json:"test"`
	Consequent AstNode `json:"consequent"`
	Alternate AstNode `json:"alternate"`
}

type ForStatement struct {
	AstNodeMeta
	Init AstNode `json:"init"`
	Test AstNode `json:"test"`
	Update AstNode `json:"update"`
	Body AstNode `json:"body"`
}

type ReturnStatement struct {
	AstNodeMeta
	Argument AstNode `json:"argument"`
}

type AssignmentExpression struct {
	AstNodeMeta
	Operator string  `json:"operator"`
	Left     AstNode `json:"left"`
	Right    AstNode `json:"right"`
}

type MemberExpression struct {
	AstNodeMeta
	Computed bool    `json:"computed"`
	Object   AstNode `json:"object"`
	Property AstNode `json:"property"`
}

type ThisExpression struct {
	AstNodeMeta
}

type FunctionExpression struct {
	AstNodeMeta
	Id         AstNode   `json:"id"`
	Params     []AstNode `json:"params"`
	Defaults   []AstNode `json:"defaults"`
	Body       AstNode   `json:"body"`
	Rest       AstNode   `json:"rest"`
	Generator  bool      `json:"generator"`
	Expression bool      `json:"expression"`
	Source     string 	 `json:"-"`
}

type CallExpression struct {
	AstNodeMeta
	Callee    AstNode   `json:"callee"`
	Arguments []AstNode `json:"arguments"`
}

type VariableDeclaration struct {
	AstNodeMeta
	Declarations []AstNode `json:"declarations"`
	Kind         string    `json:"kind"`
}

type VariableDeclarator struct {
	AstNodeMeta
	Id   AstNode `json:"id"`
	Init AstNode `json:"init"`
}

type NewExpression struct {
	AstNodeMeta
	Callee    AstNode   `json:"callee"`
	Arguments []AstNode `json:"arguments"`
}

type ArrayExpression struct {
	AstNodeMeta
	Elements []AstNode `json:"elements"`
}

type ObjectExpression struct {
	AstNodeMeta
	Properties []AstNode `json:"properties"`
}

type UnaryExpression struct {
	AstNodeMeta
	Operator string  `json:"operator"`
	Argument AstNode `json:"argument"`
	Prefix   bool    `json:"prefix"`
}

type BinaryExpression struct {
	AstNodeMeta
	Operator string  `json:"operator"`
	Left     AstNode `json:"left"`
	Right    AstNode `json:"right"`
}

type UpdateExpression struct {
	AstNodeMeta
	Operator string `json:"operator"`
	Argument AstNode `json:"argument"`
	Prefix bool `json:"prefix"`
}

func (self AstNodeMeta) AstType() AstType {
	return self.Type
}

func FormattedAstBuffer(ast AstNode) (*bytes.Buffer, error) {
	jsonStr, err := json.Marshal(ast)
	if err != nil {
		return nil, err
	}

	var dst bytes.Buffer
	err = json.Indent(&dst, jsonStr, "", "    ")
	if err != nil {
		return nil, err
	}

	return &dst, nil
}

func FormattedAstString(ast AstNode) string {
	buf, err := FormattedAstBuffer(ast)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (self AstType) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", self)
	return []byte(str), nil
}

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
	case EMPTY_STATEMENT:
		return "EmptyStatement"
	case BLOCK_STATEMENT:
		return "BlockStatement"
	case EXPRESSION_STATEMENT:
		return "ExpressionStatement"
	case IF_STATEMENT:
		return "IfStatement"
	case FOR_STATEMENT:
		return "ForStatement"
	// case LABELED_STATEMENT:
	// 	return "LabeledStatement"
	// case BREAK_STATEMENT:
	// 	return "BreakStatement"
	// case CONTINUE_STATEMENT:
	// 	return "ContinueStatement"
	// case WITH_STATEMENT:
	// 	return "WithStatement"
	// case SWITCH_STATEMENT:
	// 	return "SwitchStatement"
	case RETURN_STATEMENT:
		return "ReturnStatement"
	// case THROW_STATEMENT:
	// 	return "ThrowStatement"
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
	case ARRAY_EXPRESSION:
		return "ArrayExpression"
	case OBJECT_EXPRESSION:
		return "ObjectExpression"
	case UNARY_EXPRESSION:
		return "UnaryExpression"
	case BINARY_EXPRESSION:
		return "BinaryExpression"
	case UPDATE_EXPRESSION:
		return "UpdateExpression"

	}
	return "<#error: bad value>"
}
