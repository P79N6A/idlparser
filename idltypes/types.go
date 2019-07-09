package idltypes

type IDLTypeI interface {
	Name() string     // 类型的名字
	Package() Package // 类型所属的包
	File() File       // 类型定义的文件
}

type annotationsAble struct{
	annos *IDLAnnotations
}

func (def *annotationsAble)SetAnnotations(annos *IDLAnnotations){
	def.annos = annos
}

type IDLDefine struct{
	doc_ string
	hasDoc_ bool
	name_ string
}

func (doc *IDLDefine)SetName(name string){

}

func (doc *IDLDefine)GetPackage() string{
	return ""
}

func (doc *IDLDefine)GetFileName() string{
	return ""
}

func (doc *IDLDefine)GetName() string{
	if doc == nil{
		return "unknown_doc"
	}
	return doc.name_
}

func NewDefine(doc string) *IDLDefine{
	return &IDLDefine{
		doc_:doc,
		hasDoc_:true,
	}
}

type IDLNamespace struct{
	packageName   string
	FileName      string
	namespace     string
	language      string
	LanguageAnons map[string]*IDLAnnotations
}

type IDLDocI interface {
	SetDoc(doc string)
	GetDoc() string
	GetName() string
}
