package jaess

import (
	// "fmt"
	"io"
	"strconv"
	"strings"
)

// Parser instance, consumes a TokenScanner
type Parser struct {
	scanner *TokenScanner
}

// parses a string into an AstNode{type:Program,...}
func Parse(source string) (*Program, error) {
	parser := NewParser(strings.NewReader(source))
	return parser.Parse()
}

// create a new parser
func NewParser(input io.RuneScanner) *Parser {
	parser := new(Parser)
	parser.scanner = NewTokenScanner(input)
	return parser
}

// parses into a full Ast
func (self *Parser) Parse() (*Program, error) {
	node := new(Program)
	node.Type = PROGRAM
	for {
		n, err := self.Next()
		if err != nil {
			return nil, err
		}
		if n != nil {
			node.Body = append(node.Body, n)
		} else {
			break
		}
	}
	return node, nil
}

// parses forward and returns the next statement
func (self *Parser) Next() (AstNode, error) {
	return self.parseStatement()
}

// parses the next statement
func (self *Parser) parseStatement() (AstNode, error) {
	var token *Token
	var err error

	var node AstNode

	for {
		token, err = self.scanner.Next()
		if token == nil || err != nil {
			return nil, err
		}

		switch token.Type {
		case COMMENT:
			continue
		case NEWLINE:
			continue
		}
		if token.Value == ";" {
			self.scanner.UnNext()
			node, err = self.parseEmptyStatement()
			break
		}
		if token.Type == NUMBER || token.Type == STRING || token.Type == DELIMITER {
			self.scanner.UnNext()
			node, err = self.parseExpressionStatement()
			break
		}
		if token.Type == ATOM {
			self.scanner.UnNext()
			if token.Value == "function" {
				node, err = self.parseFunctionDeclaration()
			} else if token.Value == "return" {
				node, err = self.parseReturnStatement()
			} else if token.Value == "var" {
				node, err = self.parseVariableDeclaration()
			} else {
				node, err = self.parseExpressionStatement()
			}
			if err != nil {
				return nil, err
			}
			break
		}

		perr := NewParseError("parser error \"%s\"(%s)", token.Value, token.Type).SetLocation(token.Location)
		panic(perr) // return nil,
	}

	for {
		token, terr := self.scanner.Peek()
		if terr != nil {
			return nil, terr
		}

		if token == nil {
			break
		} else if token.Type == COMMENT {
			_, _ = self.scanner.Next()
			continue
		} else if token.Value == ";" || token.Value == "\n" {
			_, _ = self.scanner.Next()
			break
		} else {
			perr := NewParseError("parser error: STATEMENT...\"%s\"(%s)", token.Value, token.Type)
			// return nil, perr.SetLocation(token.Location)
			panic(perr.SetLocation(token.Location))
		}
	}

	return node, err

}

func (self *Parser) parseEmptyStatement() (AstNode, error) {
	node := new(EmptyStatement)
	node.Type = EMPTY_STATEMENT
	return node, nil
}

// parses a BlockStatement from start
func (self *Parser) parseBlockStatement() (AstNode, error) {
	token, err := self.scanner.Next()
	if err != nil {
		return nil, err
	}
	if token.Value != "{" {
		err := NewParseError("cannot parse BLOCK_STATEMENT<<\"%s\"(%s)", token.Value, token.Type)
		return nil, err.SetLocation(self.scanner.Location)
	}

	node := new(BlockStatement)
	node.Type = BLOCK_STATEMENT

	for {
		nextToken, err := self.scanner.Next()

		if err != nil {
			return nil, err
		}
		if nextToken == nil {
			break
		}
		if "}" == nextToken.Value {
			break
		}
		if "\n" == nextToken.Value {
			continue
		}
		if COMMENT == nextToken.Type {
			continue
		}

		self.scanner.UnNext()
		innerStatement, err := self.parseStatement()
		if err != nil {
			return nil, err
		}
		node.Body = append(node.Body, innerStatement)

	}

	for {
		token, terr := self.scanner.Peek()
		if terr != nil {
			return nil, terr
		}

		if token == nil || token.Value == ";" || token.Value == ")" {
			break
		} else if token.Type == COMMENT {
			_, _ = self.scanner.Next()
			continue
		} else if token.Value == "\n" {
			_, _ = self.scanner.Next()
			break
		} else {
			perr := NewParseError("parser error: BLOCK_STATEMENT...\"%s\"(%s)", token.Value, token.Type)
			panic(perr.SetLocation(token.Location))
			// return nil, perr.SetLocation(token.Location)
		}
	}

	return node, nil
}

// parses and wraps an expression into a statement node
func (self *Parser) parseExpressionStatement() (AstNode, error) {
	node := new(ExpressionStatement)
	node.Type = EXPRESSION_STATEMENT

	exprNode, err := self.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Expression = exprNode

	nextNode, err := self.scanner.Peek()
	if err != nil {
		return nil, err
	}
	if nextNode != nil && nextNode.Value == ";" {
		self.scanner.Next()
	}

	return node, nil
}

// parses and wraps an expression into a statement node
func (self *Parser) parseReturnStatement() (AstNode, error) {
	var err error

	token, err := self.scanner.Next()
	if token == nil || err != nil {
		return nil, err
	}
	if token.Value != "return" {
		err := NewParseError("cannot parse RETURN_STATEMENT<<%s:%s", token.Type, token.Value)
		return nil, err.SetLocation(token.Location)
	}

	node := new(ReturnStatement)
	node.Type = RETURN_STATEMENT
	node.Argument, err = self.parseExpression()
	if err != nil {
		return nil, err
	}
	return node, nil
}

// parses a function declaration into a statement node
func (self *Parser) parseFunctionDeclaration() (AstNode, error) {
	node := new(FunctionDeclaration)
	node.Type = FUNCTION_DECLARATION

	token, err := self.scanner.Next()
	if token == nil || err != nil {
		return nil, err
	}
	if token.Value != "function" {
		err := NewParseError("cannot parse FUNCTION_DECLARATION<<%s:%s", token.Type, token.Value)
		return nil, err.SetLocation(token.Location)
	}

	token, err = self.scanner.Next()
	if token == nil {
		err := NewParseError("cannot parse FUNCTION_DECLARATION<<EOF")
		err.SetLocation(self.scanner.Location)
	}
	if err != nil {
		return nil, err
	}

	if token.Type == ATOM {
		node.Id, err = self.parseIdentifier(token)
		if err != nil {
			return nil, err
		}
		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
	}

	if token.Value != "(" {
		err := NewParseError("cannot parse FUNCTION_DECLARATION<<function %s:%s", token.Type, token.Value)
		return nil, err.SetLocation(token.Location)
	}
	node.Params, err = self.parseParamList()
	node.Defaults = []AstNode{}

	self.scanner.BeginCapture()
	node.Body, err = self.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	capture := self.scanner.FinishCapture()
	node.Source = _TrimFunctionSource(capture.String())

	return node, nil
}

// parses a variable declaration
func (self *Parser) parseVariableDeclaration() (AstNode, error) {
	node := new(VariableDeclaration)
	node.Type = VARIABLE_DECLARATION

	token, err := self.scanner.Next()
	if token == nil || err != nil {
		return nil, err
	}
	if token.Value != "var" {
		err := NewParseError("cannot parse VARIABLE_DECLARATION<<%s:%s", token.Type, token.Value)
		return nil, err.SetLocation(token.Location)
	}

	node.Kind = token.Value

	for {
		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			if len(node.Declarations) == 0 {
				return nil, NewParseError("cannot parse VARIABLE_DECLARATION<<EOF").SetLocation(self.scanner.Location)
			}
			break
		}
		if token.Value == ";" {
			self.scanner.UnNext()
			break
		}

		if token.Type == ATOM {
			declNode := new(VariableDeclarator)
			declNode.Type = VARIABLE_DECLARATOR
			declNode.Id, err = self.parseIdentifier(token)
			if err != nil {
				return nil, err
			}
			token, err = self.scanner.Peek()
			if err != nil {
				return nil, err
			}
			if token != nil && token.Value == "=" {
				_, _ = self.scanner.Next()
				initNode, err := self.parseExpression()
				if err != nil {
					return nil, err
				}
				declNode.Init = initNode
				node.Declarations = append(node.Declarations, declNode)
			}
			continue
		}

		return nil, NewParseError("cannot parse VARIABLE_DECLARATOR<<\"%s\"(%s)", token.Value, token.Type).SetLocation(token.Location)
	}

	return node, nil
}

// parses a list of param patterns
func (self *Parser) parseParamList() ([]AstNode, error) {
	paramList := []AstNode{}

	for {
		token, err := self.scanner.Next()
		if err != nil {
			return nil, err
		}
		if token.Value == ")" {
			break
		}

		if token.Type == ATOM {
			pNode, err := self.parseIdentifier(token)
			if err != nil {
				return nil, err
			}
			paramList = append(paramList, pNode)
		} else {
			err := NewParseError("cannot parse PARAM_LIST<<%s:%s", token.Type, token.Value)
			return nil, err.SetLocation(token.Location)
		}

		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
		if token.Value != "," {
			if token.Value != ")" {
				self.scanner.UnNext()
			}
			break
		}
	}

	return paramList, nil
}

// parses a list of expressions seperated by commas
func (self *Parser) parseArgumentList() ([]AstNode, error) {
	nodeList := []AstNode{}
	for {
		token, err := self.scanner.Peek()
		if err != nil {
			return nil, err
		}
		if token.Value == ")" {
			_, _ = self.scanner.Next()
			break
		}

		nextNode, err := self.parseExpression()
		if err != nil {
			return nil, err
		}
		if nextNode != nil {
			nodeList = append(nodeList, nextNode)
		}

		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
		if token.Value != "," {
			if token.Value != ")" {
				self.scanner.UnNext()
			}
			break
		}
	}
	return nodeList, nil
}

// parses from the start of an expression
func (self *Parser) parseExpression() (AstNode, error) {
  return self.parseExpressionUntil([]string{})
}

// parses an expression until a match in an exclude list, or normal end of expression
func (self *Parser) parseExpressionUntil(excludeList []string) (AstNode, error) {

	var node AstNode

	for {
		token, err := self.scanner.Next()

		// fmt.Printf("\x1b[90mparsing expr<<%s %+v\x1b[0m\n", token.Value, node)

		if err != nil {
			return nil, err
		}
		if token == nil {
			if node == nil {
				perr := NewParseError("cannot parse EXPRESSION<<EOF")
				return nil, perr.SetLocation(token.Location)
			}
			return node, nil
		}

		switch token.Type {
		case NEWLINE, COMMENT:
			continue
		}

		// hack for bin exprs todo fix
		if node != nil {
			for _, exC := range excludeList {
				if exC == token.Value {
					self.scanner.UnNext()
					return node, nil
				}
			}
		}

		switch token.Value {
		case ";", ")", ",", ":", "}":
			self.scanner.UnNext()
			return node, nil
		case "(":
			if node != nil {
				node, err = self.parseCallExpression(node)
			} else {
				node, err = self.parseExpression()
				token, nerr := self.scanner.Next()
				if nerr != nil {
					return nil, err
				}
				if token.Value != ")" {
					perr := NewParseError("cannot parse (EXPRESSION...<<\"%s\"(%s)", token.Value, token.Type)
					return nil, perr.SetLocation(token.Location)
				}
			}
			if err != nil {
				return nil, err
			}
			continue
		}

		if token.Value == "{" {
			if node != nil {
				perr := NewParseError("cannot parse (EXPRESSION...<<\"%s\"(%s)", token.Value, token.Type)
				return nil, perr.SetLocation(token.Location)
			}
			node, err = self.parseObjectExpression(token)
			if err != nil {
				return nil, err
			}
			continue
		}

		if node != nil {
			switch token.Type {
			case STRING, NUMBER:
				self.scanner.UnNext()
				return node, nil
			case ATOM:
				if token.Value == "instanceof" {
					node, err = self.parseBinaryExpression(node, token)
					continue
				}
				self.scanner.UnNext()
				return node, nil
			case OPERATOR:
				switch token.Value {
				case "=", "+=", "-=":
					self.scanner.UnNext()
					node, err = self.parseAssignmentExpression(node)
				case ".":
					self.scanner.UnNext()
					node, err = self.parseMemberExpression(node)
				default:
					node, err = self.parseBinaryExpression(node, token)
				}
				if err != nil {
					return nil, err
				}
				continue
			}
		} else {
			switch token.Type {
			case NUMBER:
				node, err = self.parseLiteral(token)
				continue
			case STRING:
				node, err = self.parseLiteral(token)
				continue
			case OPERATOR:
				if IsUnaryOperator(token) {
					node, err = self.parseUnaryExpression(token)
					continue
				}
			case ATOM:
				switch token.Value {
				case "this":
					node, err = self.parseThisExpression(token)
				case "new":
					node, err = self.parseNewExpression(token)
				case "function":
					node, err = self.parseFunctionExpression(token)
				case "delete":
					node, err = self.parseUnaryExpression(token)
				case "instanceof":
				default:
					node, err = self.parseIdentifier(token)
				}
				if err != nil {
					return nil, err
				}
				continue
			}
		}

		perr := NewParseError("cannot parse EXPRESSION<<'%s'(%s)", token.Value, token.Type)
		// return nil, perr.SetLocation(token.Location)
		panic(perr.SetLocation(token.Location))
	}

	panic("unreachable")
}

// finishes parsing a function expression
func (self *Parser) parseFunctionExpression(token *Token) (AstNode, error) {
	if token.Value != "function" {
		err := NewParseError("cannot parse FUNCTION_EXPRESSION<<'%s'(%s)", token.Value, token.Type)
		return nil, err.SetLocation(token.Location)
	}

	var err error

	node := new(FunctionExpression)
	node.Type = FUNCTION_EXPRESSION

	token, err = self.scanner.Next()
	if err != nil {
		return nil, err
	}
	if token == nil {
		err := NewParseError("cannot parse FUNCTION_EXPRESSION<<EOF")
		return nil, err.SetLocation(self.scanner.Location)
	}

	if token.Type == ATOM {
		node.Id, err = self.parseIdentifier(token)
		if err != nil {
			return nil, err
		}
		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
	}

	if token.Value != "(" {
		err := NewParseError("cannot parse FUNCTION_EXPRESSION<<function %s:%s", token.Type, token.Value)
		return nil, err.SetLocation(token.Location)
	}
	node.Params, err = self.parseParamList()
	node.Defaults = []AstNode{}

	self.scanner.BeginCapture()
	node.Body, err = self.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	capture := self.scanner.FinishCapture()
	node.Source = _TrimFunctionSource(capture.String())

	return node, nil
}

// finishes parsing an object expression
func (self *Parser) parseObjectExpression(token *Token) (AstNode, error) {
	if token.Value != "{" {
		err := NewParseError("cannot parse OBJECT_EXPRESSION<<'%s'(%s)", token.Value, token.Type)
		return nil, err.SetLocation(token.Location)
	}

	var err error

	node := new(ObjectExpression)
	node.Type = OBJECT_EXPRESSION

	for {
		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
		if token == nil {
			err := NewParseError("cannot parse OBJECT_EXPRESSION<<{EOF")
			return nil, err.SetLocation(self.scanner.Location)
		}

		switch token.Value {
		case ",", "\n":
			continue
		}
		if token.Value == "}" {
			break
		}

		propNode := new(Property)
		propNode.Type = PROPERTY
		propNode.Kind = "init"

		switch token.Type {
		case STRING, NUMBER:
			propNode.Key, err = self.parseLiteral(token)
		case ATOM:
			propNode.Key, err = self.parseIdentifier(token)
		default:
			err = NewParseError("cannot parse OBJECT_EXPRESSION<<{'%s'(%s)", token.Value, token.Type).SetLocation(token.Location)
		}
		if err != nil {
			return nil, err
		}

		token, err = self.scanner.Next()
		if err != nil {
			return nil, err
		}
		if token == nil || token.Value != ":" {
			err := NewParseError("cannot parse OBJECT_EXPRESSION<<{_key_'%s'(%s)", token.Value, token.Type)
			// return nil, err.SetLocation(self.scanner.Location)
			panic(err.SetLocation(self.scanner.Location))
		}

		propNode.Value, err = self.parseExpression()
		if err != nil {
			return nil, err
		}

		node.Properties = append(node.Properties, propNode)
	}

	return node, nil
}

// finishes parsing a call expression
func (self *Parser) parseCallExpression(left AstNode) (AstNode, error) {
	node := new(CallExpression)
	node.Type = CALL_EXPRESSION
	node.Callee = left

	argList, err := self.parseArgumentList()
	if err != nil {
		return nil, err
	}
	node.Arguments = argList

	return node, nil
}

// finishes parsing a assignment expression given a left node
func (self *Parser) parseAssignmentExpression(left AstNode) (AstNode, error) {
	node := new(AssignmentExpression)
	node.Type = ASSIGNMENT_EXPRESSION
	node.Left = left

	token, err := self.scanner.Next()
	if token == nil || err != nil {
		return nil, err
	}

	switch token.Value {
	case "=", "+=", "-=":
		node.Operator = token.Value
		right, err := self.parseExpression()
		if token == nil || err != nil {
			return nil, err
		}
		node.Right = right
		return node, nil
	}

	perr := NewParseError("cannot parse ASSIGNMENT_EXPRESSION<<'%s'(%s)", token.Value, token.Type)
	return nil, perr.SetLocation(token.Location)
}

// finishes parsing a member expression given a left node
func (self *Parser) parseMemberExpression(left AstNode) (AstNode, error) {
	node := new(MemberExpression)
	node.Type = MEMBER_EXPRESSION
	node.Object = left

	token, err := self.scanner.Next()
	if token == nil || err != nil {
		return nil, err
	}

	switch token.Value {
	case ".":
		token, err = self.scanner.Next()
		if token == nil || err != nil {
			return nil, err
		}
		right, err := self.parseIdentifier(token)
		if right == nil || err != nil {
			return nil, err
		}
		node.Property = right
		return node, nil
	}

	perr := NewParseError("cannot parse MEMBER_EXPRESSION<<'%s'(%s)", token.Value, token.Type)
	return nil, perr.SetLocation(token.Location)
}

// finishes parsing a binary expression given a left node
func (self *Parser) parseBinaryExpression(left AstNode, token *Token) (AstNode, error) {
	node := new(BinaryExpression)
	node.Type = BINARY_EXPRESSION
	node.Left = left
	node.Operator = token.Value

	var right AstNode
	var err error
	if token.Value == "+" || token.Value == "-" {
		right, err = self.parseExpressionUntil([]string{"+","-"})
	} else {
		right, err = self.parseExpressionUntil([]string{"+","-","*","/"})
	}
	if token == nil || err != nil {
		return nil, err
	}
	node.Right = right
	return node, nil
}

// finishes parsing a unary expression given an operator token
func (self *Parser) parseUnaryExpression(token *Token) (AstNode, error) {
	node := new(UnaryExpression)
	node.Type = UNARY_EXPRESSION
	node.Operator = token.Value
	node.Prefix = true

	var err error
	node.Argument, err = self.parseExpressionUntil([]string{"+","-","*","/"})
	if err != nil {
		return nil, err
	}

	return node, nil
}

// finishes parsing a new expression
func (self *Parser) parseNewExpression(token *Token) (AstNode, error) {
	node := new(NewExpression)
	node.Type = NEW_EXPRESSION

	nextNode, err := self.parseExpression()
	if err != nil {
		return nil, err
	}

	// todo figure out CallExpression/NewExpression relationship
	switch nextNode := nextNode.(type) {
	case *CallExpression:
		node.Arguments = nextNode.Arguments
		node.Callee = nextNode.Callee
		return node, nil
	}

	// token, err = self.scanner.Peek()
	// if err != nil {
	// 	return nil, err
	// }
	// if token == nil {
	// 	perr := NewParseError("cannot parse NEW_EXPRESSION...<<EOF")
	// 	return nil, perr.SetLocation(self.scanner.Location)
	// }
	// if token.Value != "(" {
	// 	perr := NewParseError("cannot parse NEW_EXPRESSION...<<'%s'(%s)", token.Value, token.Type)
	// 	return nil, perr.SetLocation(token.Location)
	// }

	// argList, err := self.parseArgumentList()
	// if err != nil {
	// 	return nil, err
	// }
	// node.Arguments = argList

	return node, nil
}

// finishes parsing a this expression
func (self *Parser) parseThisExpression(token *Token) (AstNode, error) {
	node := new(ThisExpression)
	node.Type = THIS_EXPRESSION
	return node, nil
}

// finishes parsing an identifier
func (self *Parser) parseIdentifier(token *Token) (AstNode, error) {
	node := new(Identifier)
	node.Type = IDENTIFIER
	node.Name = token.Value
	return node, nil
}

// finishes parsing a literal given a token
func (self *Parser) parseLiteral(token *Token) (AstNode, error) {

	switch token.Type {
	case STRING:
		node := new(LiteralString)
		node.Type = LITERAL
		node.Raw = token.Value
		endPos := len(token.Value) - 1
		node.Value = token.Value[1:endPos]
		return node, nil
	case NUMBER:
		node := new(LiteralNumber)
		node.Type = LITERAL
		node.Raw = token.Value
		f, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return nil, err
		}
		node.Value = f
		return node, nil
	}

	perr := NewParseError("cannot parse LITERAL<<'%s'(%s)", token.Value, token.Type)
	return nil, perr.SetLocation(token.Location)
}
