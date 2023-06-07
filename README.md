# Interpreter-In-Go

## 렉싱

- 소스코드 > 토큰 > 추상구문트리
- 소스코드를 토큰열로 바꾸는 작업을 렉싱이라고 한다. 또는 어휘분석(lexical analysis)라고 부른다.
- 토큰은 자체로 쉽게 분류할 수 있는 작은 자료구조다.
- 파서는 전달받은 토큰열을 추상구문트리(Abstract Syntax Tree)로 바꾼다.

## 토큰 정의하기

- 변수이름은 식별자로 부른다.
- 보일러플레이트란 매번 변경 없이 반복되는 코드를 뜻한다.

## 렉서

- 렉서는 소스코드를 입력받고 소스코드를 표현하는 토큰열을 결과로 출력한다.
- 렉서는 입력받은 코드를 훑어가면서 토큰을 인식할 때마다 결과를 출력한다.
- 버퍼도 필요없고 토큰 저장도 필요없다. 다음 토큰을 출력하는 NextToken메서드 하나면 충분하다.

## REPL(Read Eval Print Loop)

- 대다수의 인터프리터 언어는 이 REPL을 가지고 있다. (JavaScript, Python, Ruby, Lisp류의 언어 등)
- REPL은 콘솔 혹은 대화형 모드라고 부른다.
- REPL은 입력을 읽고(Read), 인터프리터에 보내 평가하고(Eval), 인터프리터의 결과물을 출력하고(Print) 이런 동작을 반복하기(Loop) 때문에 REPL(Read, Eval, Print, Loop)이라고 부른다.

## 파서(Parser)

위키피디아의 설명이지만 개념을 참고해보자

- 파서는 (주로 문자로 된) 입력 데이터를 받아 자료구조를 만들어 내는 소프트웨어 컴포넌트다. 자료구조 형태는 파스 트리(parse tree), 추상구문트리일 수 있고, 그렇지 않으면 다른 계층 구조일 수도 있다.
  파서는 자료구조를 만들면서 입력에 대응하는 구조화된 표현을 더하기도 구문이 올바른지 검사하기도 한다. (중략) 보통은 파서 앞에 어휘 분석기(lexical analyzer)를 따로 두기도 한다.

- 파서는 입력을 표현하는 자료구조로 변환한다.
- 다음과 같은 자바스크립트 Json 파서가 있다. 개념수준에서는 프로그래밍 언어의 파서와 Json 파서는 개념이 같다.

```
> var input = '{"name": "Thorsten", "age":28 }'
> var output = JSON.parse(input)
> output
{name : 'Thorsten', age:28}
> output.name
'Thorsten'
>output.age
28
>
```
- 인터프리터와 컴파일러 관련 주제에서 소스코드를 내부적으로 표현할 때 쓰는 자료구조를 "구문트리(syntax tree)" 내지 추상구문트리(Abstract Syntax Tree, 이하 AST)라고 부른다. 
- 파싱 프로세스를 구문 분석(syntactic analysis)라고도 부른다.
- 프로그래밍 언어를 파싱할 때 크게 두 가지 전략이 있다. 하향식(top-down)과 상향식(bottom-up)이다. 
- 하향식의 예로는 재귀적 하향 파싱(recursive decent parsing), 얼리 파싱(Earley parsing), 예측적 파싱(predictive parsing)등은 모두 하향식 파싱을 변형한 것들이다. 
- 여기서는 재귀적 하향 파서를 만들 것이다.  이 파서는 프랫 파서(Pratt parser)라고 불리기도 한다. 왜냐하면 최초로 만든 사람의 이름이 본 프랫(Vaughan Pratt)이다. 
- 하향식 파서는 AST의 루트노드를 생성하는 것으로 시작해서 점차 아래쪽으로 파싱해 나간다.


