package thrift

import (
	"log"
	idl "code.byted.org/ee/lobster-idlloader/types"
	"fmt"
	"runtime/debug"
	"context"
)

const(
	PROGRAM = 1
	INCLUDES = 2
	ALREADY_PROCESSED = 3
)

var(
	G_parent_prefix = "ffd"
	G_parse_mode = PROGRAM

	crtMicroMod *idl.MicroModule
	crtFile *idl.IDLFile // 当前正在处理的文件

	hasFileInProcess = make(chan bool, 1) // 是否有文件正在处理中，一定要带缓冲

	G_doctext string
)

type Parser struct{
}

func (p *Parser)Type() idl.ProtoType{
	return idl.ProtoType_THRIFT
}

func (p *Parser)Parse(ctx context.Context, mod *idl.MicroModule, file *idl.IDLFile)(err error){
	defer func() {
		<- hasFileInProcess // 当前文件处理，释放
		if e := recover(); e != nil {
			err = fmt.Errorf("recover error:%+v, stack:\n%s", e, debug.Stack())
		}
	}()

	hasFileInProcess <- true // 当前有文件在处理中

	crtMicroMod = mod
	crtFile = file
	lex := NewThriftLexer(file.GetContent())
	ret := yyParse(lex)

	if 0 != ret{
		fmt.Errorf("yyParse fail:%d", ret)
	}

	return err
}

func NotSupported(errMsg string, params...interface{}){
	for idx, v := range params{
		errMsg += fmt.Sprintf(" param_%d:,%+v", idx, v)
	}

	log.Println(errMsg)
}

func DestoryDocText(){
	//log.Printf("clear doc")
}

func CaptureDocText() string{
	ret := G_doctext
	G_doctext = ""
	return ret
}

//
// 这个不处理
func IncludeFile(file string){
	//log.Printf("[Mod:%s][File:%s]include:%s", crtMicroMod.GetName(), crtFile.GetName(), file)
}

func Program(){
}