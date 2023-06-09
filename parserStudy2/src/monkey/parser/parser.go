package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// 우선순위 테이블
// 하단의 연산자 우선순위에 따라 token.PLUS와 token.MINUS는 우선순위가 같다.
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// 연산자 우선순위
const (
	_ int = iota //iota를 이용해 뒤에 나오는 상수에게 1씩 증가하는 숫자를 값으로제공한다.  _는 0이되고 이후에 나오는 상수는 1부터 7까지 할당받는다.
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// 프랫파서 구현의 핵심아이디어는 파싱함수를 토큰타입과 연관짓는 것이다.
// 파서가 토큰 타입을 만날 때마다 파싱 함수가 적절한 표현식을 파싱하고 그 표현식을 나타내는 AST노드를 하나 반환한다. 각각의 토큰 타입은 토큰이 전위연산자인지 중위 연산자인지에 따라 최대 두 개의 파싱 함수와
// 연관지을 수 있다.
type (
	prefixParseFn func() ast.Expression               //전위 파싱 함수, 전위 연산자와 연관된 토큰 타입을 만나면 이 함수가 호출된다.
	infixParseFn  func(ast.Expression) ast.Expression //중위 파싱 함수,  여기서 받는 인수 ast.Expression는 중위연산자의 좌측에 위치한다. / 중위 연산자와 연관된 토큰 타입을 만나면 이 함수가 호출된다.
)

// 파서는 3개 필드를 가지고 있다. l, curToken, peekToken이다. l은 input에서 다음 토큰을 얻기위해 반복해서 NextToken을 호출하는 데 쓰인다.
// curToken과 peekToken은 lexer에서 사용한 포인터인 position, readPosition과 매우 유사하게 동작한다.
type Parser struct {
	l         *lexer.Lexer //현재의 렉서 인스턴스를 가리키는 포인터
	curToken  token.Token  //현재 토큰
	peekToken token.Token  //그 다음 토큰
	errors    []string     //에러를  처리하기 위한 선언

	//파서가 토큰 타입에 맞게 prefixParseFn이나 infixParseFn을 선택하도록 map을 두 개 추가한다.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// 자기설명적이고 nextToken메서드는 curToken과 peekToken을 다음 위치로 보내는 짧은 도움 메서드다.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	//전위 연산자
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression) //token.BANG과 token.MINUS는 연관된 파싱 함수가 같다.
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	//중위 연산자
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	//bool 타입
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	//그룹 표현식
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	//if
	p.registerPrefix(token.IF, p.parseIfExpression)

	return p
}

// parsePrefixExpression과 다른 점은 left를 인수로 받는 점
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	//defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	//현재 토큰의 우선순위를 precedence에 넣어둠
	precedence := p.curPrecedence()
	//토큰을 진행시킴
	p.nextToken()
	//Right필드에 parseExpression 함수를 호출한 결과값을 담는다.
	expression.Right = p.parseExpression(precedence)
	return expression
}

// p.peekToken이 갖는 토큰타입과 연관된 우선순위를 반환한다. 타입을 못찾으면 LOWEST를 반환한다.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST //연산자가 가질 수 있는 우선순위 중 가장 낮은 값
}

// p.curToken이 갖는 토큰타입과 연관된 우선순위를 반환한다. 타입을 못찾으면 LOWEST를 반환한다.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	//defer untrace(trace("parsePrefixExpression"))
	//ast.PrefixExpression를 만든다.
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	//현재 위치에 있는 토큰 token.LET 토큰으로 *ast.LetStatement 노드를 만든다.
	stmt := &ast.LetStatement{Token: p.curToken}

	//이후 다음에 원하는 토큰이 오는지 확인하기 위해 expectPeek을 호출한다.
	//우선 token.IDENT가 오기를 기대한다.
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	//token.IDENT는 *ast.Identifier 노드를 만드는데 사용된다.
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	//그러고나서 등호가 오기를 기대한다.
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO : 세미콜론을 만날 때까지 표현식을 건너뛴다.
	//세미콜론을 만나기 전까지 등호 이후의 표현식(expression)을 건너뛴다.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt

}

// return 문 statement를 파싱하는 함수
// 현재 위치에 있는 토큰으로 ast.ReturnStatement를 만들어냄
// 그리고 nextToken을 호출해서 파서를 다음에 올 표현식이 있는 곳에 위치시킨다.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	//TODO: 세미콜론을 만날 때까지 표현식을 건너뛴다.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 그룹 표현식을 파싱하기 위한 함수
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// if문 파싱하기 위한 함수
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// if 와 else에 있는 블록 스테이츠먼츠를 파싱하기 위한 함수다.
// p.curToken과 p.peekToken을 필요한 만큼만 진행시켰기 때문에,parseBlockStatements가 호출된 시점에 p.curToken은 { 을 보고 있을 것이고 토큰 타입은 token.LBRACE가 될 것이다.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram은 가장 먼저 AST의 루트 노드인 *ast.Program을 만든다.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	program.Statements = []ast.Statement{}

	//token.EOF 토큰을 만날 때까지 input에 있는 모든 토큰을 대상으로 for-loop를 반복한다.
	for p.curToken.Type != token.EOF {
		//반복할 때마다 parseStatement를 호출해 명령문 statement를 파싱한다.
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		//p.curToken과 p.peekToken 둘 다 진행시키는 nextToken 메서드를 반복적으로 호출한다.
		p.nextToken()
	}

	return program
}

// parseStatement는 curToken의 타입을 파싱해서 ast.Statement로 반환한다.
// 이렇게 반환된 ast.Statement는 AST 루트 노드의 Statements 슬라이스에 추가된다.
// 더는 반환할 것이 없으면 *ast.Program 루트노드를 반환한다.
// 표현식을 파악하기위한 메서드
// monkey언어에는 let문과 return문만 있어서 나머지는 전부 표현식문으로 파싱한다.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	//let문일 경우
	case token.LET:
		return p.parseLetStatement()
		//let문일 경우
	case token.RETURN:
		return p.parseReturnStatement()
		//let문일 경우
	default:
		return p.parseExpressionStatement()
	}
}

// 표현식을 파싱하는 함수
// 반환값은 *ast.ExpressionStatement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	//defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parseExpression은 p.curToken.Type이 전위로 연관된 파싱함수가 있는지 검사한다. 만약 그런 파싱함수가 있으면 호출하고 없다면 nil을 반환한다.
// 프랫파서의 요체
func (p *Parser) parseExpression(precedence int) ast.Expression {
	//defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	//반복문 몸체(body)에서 parseExpression 메서드는 다음 토큰에 맞는 infixParseFn을 찾는다.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	//defer untrace(trace("parseExpression"))
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// expectPeek 메서드내에서 nextToken을 호출해 토큰을 진행시킨다.
// 원하는 토큰 타입이 오는 지 확인한다.
// 다음 함수는 모든 파서가 공유하는 단정(assert) 함수다.
// 다음 토큰 타입을 검사해 토큰 간의 순서를 올바르게 강제할 용도로 사용한다.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		//nextToken을 통해 curToken, peekToken을 다음 데이터로 파싱을 할 수 있게 진행한다.
		return true
	} else {
		//에러체크하는 case
		p.peekError(t)
		return false
	}
}

// parser map에 파싱함수를 추가하는 도움 메서드 두개 정의
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
