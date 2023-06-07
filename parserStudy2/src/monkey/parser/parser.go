package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
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

// 프랫파서 구현의 핵심아이디어는 파싱함수를 토큰타입과 연관짓는 것이다.
// 파서가 토큰 타입을 만날 때마다 파싱 함수가 적절한 표현식을 파싱하고 그 표현식을 나타내는 AST노드를 하나 반환한다. 각각의 토큰 타입은 토큰이 전위연산자인지 중위 연산자인지에 따라 최대 두 개의 파싱 함수와
// 연관지을 수 있다.
type (
	prefixParseFn func() ast.Expression               //전위 파싱 함수, 전위 연산자와 연관된 토큰 타입을 만나면 이 함수가 호출된다.
	infixParseFun func(ast.Expression) ast.Expression //중위 파싱 함수,  여기서 받는 인수 ast.Expression는 중위연산자의 좌측에 위치한다. / 중위 연산자와 연관된 토큰 타입을 만나면 이 함수가 호출된다.
)

// parser map에 파싱함수를 추가하는 도움 메서드 두개 정의
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 자기설명적이고 nextToken메서드는 curToken과 peekToken을 다음 위치로 보내는 짧은 도움 메서드다.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	return p
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
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
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
