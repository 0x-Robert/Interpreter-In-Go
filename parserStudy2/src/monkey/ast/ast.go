package ast

import (
	"bytes"
	"monkey/token"
)

//생성할 AST는 노드로만 구성될 것이고 각각의 노드는 서로 연결될 것이다. AST도 결국 트리다.

// AST를 구성하는 모든 노드는 Node 인터페이스를 구현해야한다. 따라서 모든 노드는 TokenLiteral메서드를 제공해야하고 TokenLiteral 메서드는 토큰에 대응하는 리터럴값을 반환해야 한다.
type Node interface {
	//토큰에 대응하는 리터럴 값을 반환해야한다.
	//디버깅과 테스트용도로만 사용된다.
	TokenLiteral() string
	String() string
}

// 어떤 노드는 Statement 인터페이스를 구현한다.
type Statement interface {
	Node
	//더미 메서드 포함 꼭 필요하지 않지만 Go 컴파일러가 에러를 처리하는 데 도움을 줄 수 있음
	statementNode()
}

// 어떤 노드는 Expression 인터페이스를 구현한다.
type Expression interface {
	Node
	//더미 메서드 포함 꼭 필요하지 않지만 Go 컴파일러가 에러를 처리하는 데 도움을 줄 수 있음
	expressionNode()
}

// 첫 번째 Node 인터페이스 구현체인 Program 노드다.
// Program 노드는 파서가 생산하는 모든 AST의 루트 노드가 된다.
// 모든 몽키 프로그램은 일련의 명령문으로 이루어진다.
type Program struct {
	//위 명령문들은 Program.Statements에 들어 있고, Statement 인터페이스를 구현하는 AST노드로 구성된 슬라이스(slice)일 뿐이다.
	Statements []Statement
}

// ast.Expression 인터페이스를 충족한다.
type IntegerLiteral struct {
	Token token.Token
	Value int64 //ast.Identifier와 구조체 자체에서 눈에 띄게 다른 점은 Value가 문자열이 아니라 int64다. 소스코드에서 정수 리터럴이 표현하는 문자의 실제값을 담았다.
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// 전위 연산자
type PrefixExpression struct {
	Token    token.Token //전위 연산자 토큰, 예를 들면 !, -
	Operator string      // - or !를 담을 필드
	Right    Expression  //연산자의 오른쪽에 나올 표현식을 담을 필드
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// 중위 연산자
type InfixExpression struct {
	Token    token.Token //연산자 토큰 예를 들어 5 + 5에서 +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// 다음 메서드가 있으면 디버깅 목적으로 AST 노드를 출력해볼 수 있고 또 다른 AST 노드와 비교도 할 수 있다.
// 다음 String 메서드는 작업 대부분을 *ast.Program.Statements에 위임(delegate)한다.
func (p *Program) String() string {
	//버퍼를 선언
	var out bytes.Buffer

	//각 명령문의 String 메서드를 호출해서 반환값을 버퍼에 쓴다. 이후 버퍼를 문자열로 반환한다.
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token //token.IDENT 토큰
	Value string
}

func (i *Identifier) expressionNode()      {} //Expression 인터페이스
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// 변수 바인딩에 사용할 노드가 어떤 모습이어야 좋을지 생각해보자
// 예를 들어 let x = 5를 생각해보자. 이걸 바인딩하려면 변수이름과 등호 오른쪽 표현식 필드, 그리고 AST 노드와 연관된 토큰도 추적할 수 있어야한다.
// 3개 필드는 바꿔 말하면 식별자 필드,  값을 내는 표현식 필드, 나머지하나는 토큰 필드가 필요하다.
type LetStatement struct {
	Token token.Token //토큰 필드
	Name  *Identifier //변수 바인딩 식별자 필드
	Value Expression  //값을 생성하는 표현식 필드
}

// 다음은 각각 Statement와 Node 인터페이스를 각각 구현하고 있다.
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier 구조체는 Expression 인터페이스를 구현한다.왜 Expression이냐 하면 파서 프로그램을 단순하게 만들기 위해서다.
// 몽키프로그램의 다른 부분에서는 식별자가 값을 생성하기도 한다. 예를 들면 다음과 같다. let x = 값_생성 식별자에서는 값을 생성한다.
// 노드 타입의 수를 가능한 작게 만들기  위해 Identifier 노드를 사용하려 한다.

// return 문 파싱
// return 문은 return 키워드 하나와 표현식 하나로 구성된다.
type ReturnStatement struct {
	Token       token.Token // 'return' 토큰, 토큰타입이다.
	ReturnValue Expression  // 반환될 표현식을 담는다.
}

func (rs *ReturnStatement) statementNode()       {}                          //Node 인터페이스를 충족한다. *ast.LetStatement에 정의된 메서드와 동일해 보인다.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal } //Node 인터페이스를 충족한다. *ast.LetStatement에 정의된 메서드와 동일해 보인다.

// 표현식을 지원하는 언어는 하나의 행을 하나의 표현식으로 구성할 수 있다.
// 이제 이런 유형의 노드를 AST에 추가하자!
// 표현식 AST
// 다음은 모든 노드 타입이 가지는 Token필드와 표현식이 가지는 Expression 필드다.
// ast.ExpressionStatement는 ast.Statement 인터페이스를 충족한다.
// ast.ExpressionStatement를 ast.Program의 필드인 Statements슬라이스에 추가할 수 있다.
type ExpressionStatement struct {
	Token      token.Token //표현식의 첫 번째 토큰
	Expression Expression  //표현식
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
