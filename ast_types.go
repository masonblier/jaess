package jsparser

import (
	"fmt"
	// "encoding/json"
)

// an AstNode
type AstNode interface {
}

// types of ast nodes
type AstKind int

const (
	AST_UNKNOWN AstKind = iota
	AST_LITERAL
	AST_IDENTIFIER
	AST_PROGRAM
	AST_FUNCTION_DECLARATION
	AST_BLOCK_STATEMENT
	AST_EXPRESSION_STATEMENT
	AST_ASSIGNMENT_EXPRESSION
	AST_MEMBER_EXPRESSION
	AST_THIS_EXPRESSION
  AST_FUNCTION_EXPRESSION
  AST_CALL_EXPRESSION
  AST_VARIABLE_DECLARATION
  AST_VARIABLE_DECLARATOR
  AST_NEW_EXPRESSION
	AST_UNARY_EXPRESSION
	AST_BINARY_EXPRESSION
)

func (self AstKind) String() string {
	switch self {

	case AST_UNKNOWN:
		return "AstUnknown"
	case AST_LITERAL:
		return "Literal"
	case AST_IDENTIFIER:
		return "Identifier"
	case AST_PROGRAM:
		return "Program"
	case AST_FUNCTION_DECLARATION:
		return "FunctionDeclaration"
	case AST_BLOCK_STATEMENT:
		return "BlockStatement"
	case AST_EXPRESSION_STATEMENT:
		return "ExpressionStatement"
	case AST_ASSIGNMENT_EXPRESSION:
		return "AssignmentExpression"
	case AST_MEMBER_EXPRESSION:
		return "MemberExpression"
	case AST_THIS_EXPRESSION:
		return "ThisExpression"
  case AST_FUNCTION_EXPRESSION:
    return "FunctionExpression"
  case AST_CALL_EXPRESSION:
    return "CallExpression"
  case AST_VARIABLE_DECLARATION:
    return "VariableDeclaration"
  case AST_VARIABLE_DECLARATOR:
    return "VariableDeclarator"
  case AST_NEW_EXPRESSION:
    return "NewExpression"
  case AST_UNARY_EXPRESSION:
    return "UnaryExpression"
	case AST_BINARY_EXPRESSION:
		return "BinaryExpression"
	}
	return "<#error: bad value>"
}

func (self AstKind) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", self)
	return []byte(str), nil
}

type AstNodeLiteralString struct {
	Kind  AstKind `json:"type"`
	Value string  `json:"value"`
	Raw   string  `json:"raw"`
}

type AstNodeLiteralNumber struct {
	Kind  AstKind `json:"type"`
	Value float64 `json:"value"`
	Raw   string  `json:"raw"`
}

type AstNodeIdentifier struct {
	Kind AstKind `json:"type"`
	Name string  `json:"name"`
}

type AstNodeProgram struct {
  Kind AstKind `json:"type"`
  Body []AstNode `json:"body"`
}

type AstNodeFunctionDeclaration struct {
  Kind AstKind `json:"type"`
  Id AstNode `json:"id"`
  Params []AstNode `json:"params"`
  Defaults []AstNode `json:"defaults"`
  Body AstNode `json:"body"`
  Rest AstNode `json:"rest"`
  Generator bool `json:"generator"`
  Expression bool `json:"expression"`
}

type AstNodeBlockStatement struct {
  Kind AstKind `json:"type"`
  Body []AstNode `json:"body"`
}

type AstNodeExpressionStatement struct {
	Kind       AstKind `json:"type"`
	Expression AstNode `json:"expression"`
}

type AstNodeAssignmentExpression struct {
  Kind     AstKind `json:"type"`
  Operator string  `json:"operator"`
  Left     AstNode `json:"left"`
  Right    AstNode `json:"right"`
}

type AstNodeMemberExpression struct {
  Kind AstKind `json:"type"`
  Computed bool `json:"computed"`
  Object AstNode `json:"object"`
  Property AstNode `json:"property"`
}

type AstNodeThisExpression struct {
  Kind AstKind `json:"type"`
}

type AstNodeFunctionExpression struct {
  Kind AstKind `json:"type"`
  Id AstNode `json:"id"`
  Params []AstNode `json:"params"`
  Defaults []AstNode `json:"defaults"`
  Body AstNode `json:"body"`
  Rest AstNode `json:"rest"`
  Generator bool `json:"generator"`
  Expression bool `json:"expression"`
}

type AstNodeCallExpression struct {
  Kind AstKind `json:"type"`
  Callee AstNode `json:"callee"`
  Arguments []AstNode `json:"arguments"`
}

type AstNodeVariableDeclaration struct {
  Kind AstKind `json:"type"`
  Declarations []AstNode `json:"declarations"`
  Keyword string `json:"kind"`
}

type AstNodeVariableDeclarator struct {
  Kind AstKind `json:"type"`
  Id AstNode `json:"id"`
  Init AstNode `json:"init"`
}

type AstNodeNewExpression struct {
  Kind AstKind `json:"type"`
  Callee AstNode `json:"callee"`
  Arguments []AstNode `json:"arguments"`
}

type AstNodeUnaryExpression struct {
  Kind       AstKind `json:"type"`
  Operator string `json:"operator"`
  Argument AstNode `json:"argument"`
  Prefix bool  `json:"prefix"`
}

type AstNodeBinaryExpression struct {
	Kind     AstKind `json:"type"`
	Operator string  `json:"operator"`
	Left     AstNode `json:"left"`
	Right    AstNode `json:"right"`
}
