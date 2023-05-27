package lexer

import (
	'go/token'
	'monkey/token'
)

// position과 readPosition 모두 입력문자열에 있는 문자에 인덱스로 접근하기 위해 사용된다.
// 입력문자열을 가리키는 포인터가 두 개인 이유는 다음 처리 대상을 알아내려면 입력 문자열에서 다음 문자를 미리 살펴봄과 동시에 현재 문자를 보존할 수 있어야 하기 때문이다.
type Lexer struct {
	input        string
	position     int  //입력해서 현재 위치(현재 문자를 가리킴)
	readPosition int  //입력에서 현재 읽는 위치 (현재 문자의 다음을 가리킴)
	ch           byte //현재 조사하고 있는 문자, 현재문자가 곧 byte 타입을 갖는 ch다.
}

// 문자열 input에서 렉서가 현재 보고 있는 위치를 다음으로 이동하기위한 메서드
func (l *Lexer) readChar() {

	//문자열 input의 끝에 도달했는지 확인함
	if l.readPosition >= len(l.input) {
		//끝에 도달했다면 l.ch에 아스키 코드 문자 NUL에 해당하는 0을 넣는다.
		l.ch = 0
	} else {
		//끝에 도달 못했다면 l.input[l.readPosition]으로 접근해서 l.ch에 다음 문자를 저장한다.
		l.ch = l.input[l.readPosition]
	}
	//l.position을 l.readPosition으로 업데이트하고 readPosition은 1 증가시킴
	l.position = l.readPosition
	//l.readPosition은 항상 다음에 읽어야할 위치, l.position은 마지막으로 읽은 위치
	l.readPosition += 1
	//렉서는 유니코드가 아니라 아스키 문자만 지원함, 렉서를 단순화하기 위함
	//유니코드 지원은 독자가 알아서 개선하기..
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ''
		tok.Type = token.EOF
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
