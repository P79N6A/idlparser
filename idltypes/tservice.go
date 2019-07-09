package idltypes

import (
	"fmt"
)

type IDLService struct{
	pkg *Package
	name string
	methods map[string]*IDLFunction
}

func newService(pkg *Package) *IDLService{
	return &IDLService{
		pkg: pkg,
	}
}

func (svc IDLService)GetName() string{
	return svc.name
}

func (svc IDLService)Package() *Package{
	return svc.pkg
}

func (svc *IDLService)AddMethod(f *IDLFunction) error {
	if svc.methods == nil{
		svc.methods = make(map[string]*IDLFunction)
	}

	_, ok := svc.methods[f.GetName()]
	if ok{
		return NewParseError(ErrorCategoryParse, ErrorNoDumplicated, fmt.Errorf("function:%s already exist in service:%s", f.GetName(), svc.name))
	}

	svc.methods[f.GetName()] = f

	return nil
}

func (it *IDLService)SetName(name string){
	it.name = name
}

func (it *IDLService)SetExtends(svc *IDLService){
}

func (it *IDLService)SetAttrs(attrs *IDLLobAttrs){
}
