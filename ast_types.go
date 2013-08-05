package jsparser

// a node in the ast
type AstNode struct {
  Type AstKind
  meta map[string] interface{}
}


func NewAstNode(kind AstKind) *AstNode {
  return &AstNode{kind, map[string]interface{}{"type": kind.String() }}
}

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
  AST_UNARY_EXPRESSION
  AST_FUNCTION_EXPRESSION
  AST_CALL_EXPRESSION
  AST_VARIABLE_DECLARATION
  AST_VARIABLE_DECLARATOR
  AST_NEW_EXPRESSION
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
  case AST_UNARY_EXPRESSION:
    return "UnaryExpression"
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
  case AST_BINARY_EXPRESSION:
    return "BinaryExpression"
  }
  return "<#error: bad value>"
}



