package idltypes

import (
	"fmt"
	"strings"
)

type NameSpace struct {
	full string
	names []string
}

func (ns NameSpace)String() string{
	return ns.full
}

func (ns NameSpace)LastName() string{
	if len(ns.names) == 0{
		return ""
	}

	return ns.names[len(ns.names) - 1]
}

func NewNameSpace(nameString string) NameSpace{
	return NameSpace{
		full: nameString,
		names: strings.Split(nameString, "."),
	}
}

//
// 表示一个go的包
type Package struct{
	name string
	nameSpace NameSpace
	services map[string]IDLService
	structs map[string]*IDLStruct
	typedefs map[string]*Typedef
	types map[string]IDLTypeI
	enums map[string]*IDLEnum
	consts map[string]interface{}
}

func NewPackage(nameSpace string) *Package {
	return &Package{
		nameSpace: NewNameSpace(nameSpace),
	}
}

func (pkg *Package)addConst(name string, value interface{}) error {
	if pkg.consts == nil{
		pkg.consts = make(map[string]interface{})
	}

	_, exist := pkg.consts[name]
	if exist{
		return NewParseError(ErrorCategoryParse, ErrorNoDumplicated, fmt.Errorf("const %s already exist in package:%s", name, pkg.nameSpace))
	}

	pkg.consts[name] = value
	return nil
}

// Query a type
func (pkg *Package)queryType(name string) IDLTypeI{
	ret, ok := pkg.types[name]
	if !ok{
		return nil // not exist.
	}

	switch realRet := ret.(type){
	case *Typedef:
		return realRet.GetReftype()
	}

	return ret
}
