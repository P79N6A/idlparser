package thrift

import "github.com/sirupsen/logrus"

type Lexer struct {
	HasErrors bool
	ErrMsg    string

	Src      []byte   //当前文档的指针
	CurChar    byte   // 当前读的字符
	In_stack   []byte // 当前已经在栈内
	Un_stack   []byte // 已读的吐回
	Mark int
	CurLineNo    int    // 当前所在行数
}

func NewThriftLexer(content []byte) *Lexer {
	cl := &Lexer{
		Src: content,
	}
	cl.Next()
	return cl
}

func (cl *Lexer)Next() {
	if len(cl.Un_stack) != 0 {
		cl.CurChar = cl.Un_stack[len(cl.Un_stack)-1]
		cl.Un_stack = cl.Un_stack[:len(cl.Un_stack)-1]
		return
	}

	cl.In_stack = append(cl.In_stack, cl.CurChar)
	if len(cl.Src) == 0 {
		cl.CurChar = 0 // EOF
		return
	}

	cl.CurChar = cl.Src[0]
	if cl.CurChar == '\n'{
		cl.CurLineNo ++
	}

	cl.Src = cl.Src[1:]
}

func (cl *Lexer)Unget(b byte) {
	cl.Un_stack = append(cl.Un_stack, b)
}

func (cl *Lexer) Error(s string) {
	cl.ErrMsg = s
	cl.HasErrors = true
	logrus.Warnf("err:%s, line:%d", s, cl.CurLineNo)
}

func (cl *Lexer)Reduced(rule, state int, lval *yySymType) bool{
	return false
}

func (cl *Lexer)Match(){
	cl.Mark = len(cl.In_stack)
}