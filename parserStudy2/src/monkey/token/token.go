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
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	//구분자
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

//토큰 리터럴에 맞는 TokenType을 반환할 함수를 정의함

// 식별자가 예약어인지 정의한 함수
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// 키워드 테이블을 검사해서 주어진 식별자가 예약어인지 아닌지 살펴본다.
// 만약 예약어라면 TokenType 상수를 반환한다.
// 만약 예약어가 아니라면 그냥 token.IDENT를 반환한다.  token.IDENT는 사용자 정의 식별자를 나타내는 TokenType이다.
// 다음 함수를 통해 식별자와 예약어 렉싱을 마무리할 수 있다.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
