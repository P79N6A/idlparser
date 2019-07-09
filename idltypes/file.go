package idltypes

import (
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"path/filepath"
	"strings"
)

// one idl file
type File struct{
	Name string
	NameSpace string
	pkg *Package
	Content []byte
	Version int
	refResult *Result
}

//
func (f *File)NewStruct() *IDLStruct{
	return newStruct(f.pkg)
}

func (f *File)NewService() *IDLService{
	return newService(f.pkg)
}

func (f *File)NewTypedef(name string, refType IDLTypeI) *Typedef{
	return &Typedef{
		typeName: name,
		realType_: refType,
	}
}

func (f *File)NewFunction(rspT IDLTypeI, name string, reqTs *IDLStruct, exT IDLTypeI, oneway bool) *IDLFunction{
	//log.Println("NewFunction:%s, respT:%+v", name, rspT)
	ret := &IDLFunction{
		ReqType:reqTs,
		RspType:rspT,
		IDLNamespace:f.IDLNamespace,
	}

	ret.name_ = name

	return ret
}

func (f *File)NewConst(t IDLTypeI, k string, v interface{}, comment string) *IDLConst{
	return NewIDLConst(k, v, t, comment)
}

func (f *File)NewEnumValue(k string, v int64) *IDLEnumValue{
	return NewIDLEnumValue(k,v, f.packageName)
}

func (f *File)NewEnum() *IDLEnum {
	return NewIDLEnum(f.IDLNamespace, f.FileName)
}

func (f *File)NewField(tag int64, t IDLTypeI, name string) *IDLField {
	ret := &IDLField{
		FieldTag:  tag,
		fieldType: t,
		FieldName: name,
		mod_:f.mod_,
		idlFile: f,
	}
	//log.Printf("[Field]%s", ret.ShortDebugString())

	return ret
}

func (f *File)NewList(vType IDLTypeI) *IDLList{
	ret := &IDLList{}
	ret.SetContainer(containerList, nil, vType)
	ret.mod_ = f.mod_
	return ret
}

func (f *File)NewSet(kType IDLTypeI) *IDLSet{
	ret := &IDLSet{}
	ret.SetContainer(containerSet, kType, nil)
	ret.mod_ = f.mod_
	return ret
}


func (f *File)NewMap(kType IDLTypeI, vType IDLTypeI) *IDLMap{
	ret := &IDLMap{}
	ret.SetContainer(containerMap, kType, vType)
	ret.mod_ = f.mod_
	return ret
}

func (f *File)CreateAttrList() *IDLLobAttrs{
	return NewIDLAttrs(f.IDLNamespace)
}

func (f *File)NewAnnotations() *IDLAnnotations{
	return NewIDLAnnotations()
}

func (f *File)CreateAttr(attrStr string) *IDLLobAttr{
	var ret IDLLobAttr
	firstEq := strings.IndexByte(attrStr, '=')
	if firstEq > 0{
		ret.key = attrStr[:firstEq]
		if firstEq + 1 < len(attrStr) {
			ret.val = attrStr[firstEq+1:]
		}
	}else{
		ret.val = attrStr
	}

	return &ret
}


type FileCollection struct{
	files map[string]File

	protoType *prototype.ProtoType
}

func (fc *FileCollection)AddFile(fileName string, content []byte) error{
	fileExt := filepath.Ext(fileName)
	protoType := prototype.NewPrototype(fileExt)
	if protoType == prototype.Invalid{
		return NewParseError(ErrorCategoryParse, ErrorNoFileTypeInvalid, fmt.Errorf("invalid file ext:%s", fileExt))
	}

	if fc.protoType == nil{
		fc.protoType = &protoType
	}else if *fc.protoType != protoType{
		return NewParseError(ErrorCategoryParse, ErrorNoFileTypeInvalid, fmt.Errorf("invalid file ext:%s", fileExt))
	}

	file := newFile()
	if fc.files == nil{
		fc.files = make(map[string]File)
	}
}

func NewFileCollection() *FileCollection{
	return &FileCollection{
	}
}
