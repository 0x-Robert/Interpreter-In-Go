package token

type TokenType string //서로 다른 여러 값을 TokenType으로 필요한만큼 정의해서 사용가능하다.

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" //어떤 토큰이나 문자를 렉서가 알 수없다는 뜻이다.
	EOF     = "EOF"     //파일의 끝을 말한다.

	//식별자 + 리터럴
	IDENT = "IDENT" // add, foobar, x,y, ...
	INT   = "INT"   //1343456

	//연산자
	ASSIGN = "="
	PLUS   = "+"

	//구분자
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//예약어
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
