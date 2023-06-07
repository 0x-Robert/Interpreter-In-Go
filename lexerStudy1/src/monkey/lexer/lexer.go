package lexer

import (
	"monkey/token"
)

// position과 readPosition 모두 입력문자열에 있는 문자에 인덱스로 접근하기 위해 사용된다.
// 입력문자열을 가리키는 포인터가 두 개인 이유는 다음 처리 대상을 알아내려면 입력 문자열에서 다음 문자를 미리 살펴봄과 동시에 현재 문자를 보존할 수 있어야 하기 때문이다.
type Lexer struct {
	input        string
	position     int  //입력에서 현재 위치(현재 문자를 가리킴)
	readPosition int  //입력에서 현재 읽는 위치 (현재 문자의 다음을 가리킴)
	ch           byte //현재 조사하고 있는 문자, 현재문자가 곧 byte 타입을 갖는 ch다.
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
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

// 렉서의 핵심 함수다.  렉서는 소스코드를 입력받아서 토큰열로 출력해주는 기능이 핵심이다.
func (l *Lexer) NextToken() token.Token {
	//token.go에서 정의한 token
	var tok token.Token

	//공백처리를 위한 함수
	l.skipWhitespace()

	switch l.ch {
	case '=':
		//렉서가 입력에서 ==을 만나면 렉서는 token.EQ 하나를 만드는 게 아니라 token.ASSIGN을 두 개 생성한다.
		//그래서 peekChar 메서드를 사용해야한다.
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

//readChar 함수와 비슷한 기능
//다음에 나올 입력을 미리 살펴본다 = peek
//l.position이나 l.readPosition을 증가시키지 않는다.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// 토큰타입과 현재 조사하고 있는 문자를 입력으로 받는다. 반환은 토큰타입으로 반환한다.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	//토큰 타입과 토큰스트링 리터럴로 반환
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// 현재 토큰의 Literal 필드를 채운다.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 이 책에서 글자란 isLetter 함수가 참으로 판별하는 문자(character)을 뜻한다.
func isLetter(ch byte) bool {
	//_ 문자를 글자로 다루겠다는 뜻이고 식별자와 예약어에 사용하겠다는 뜻
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 공백문자를 통째로 지나가는 함수
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// 전달받은 byte가 0부터 9사이의 라틴 숫자인지 아닌지 여부만 반환한다.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
