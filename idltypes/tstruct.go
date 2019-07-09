package idltypes

import (
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"sort"
	//"log"
)

type IDLStruct struct{
	pkg *Package
	name string
	tagMembers map[int64] *IDLField
	nameMebers map[string] *IDLField
}

func newStruct(pkg *Package) *IDLStruct{
	return &IDLStruct{
		pkg:pkg,
		tagMembers:make(map[int64]*IDLField),
		nameMebers:make(map[string]*IDLField),
	}
}

func (it *IDLStruct)Name() string{
	return it.name
}

func (it *IDLStruct)SetName(name string){
	it.name = name
}

func (it *IDLStruct)Package() Package{
	return it.pkg
}

func (it *IDLStruct)FullName() string{
	return fmt.Sprintf("%s.%s", it.pkg.name, it.name)
}

func (it *IDLStruct)SetXSDAll(b bool){
}

func (it *IDLStruct)SetUnion(b bool){
}

func (it *IDLStruct)SetXception(b bool){
}

func (it *IDLStruct)AppendField(f *IDLField) error{
	_, ok := it.tagMembers[f.FieldTag]
	if ok{
		return NewParseError(ErrorCategoryParse, ErrorNoDumplicated, fmt.Errorf("tag %d already exist in struct:%s", it.FullName()))
	}

	_, ok = it.NameMembers[f.FieldName]
	if ok{
		return fmt.Errorf("invalid thrift idl: struct %s has dumplicated tag:%d", it.name_, f.FieldTag)
	}

	it.TagMembers[f.FieldTag] = f
	it.NameMembers[f.FieldName] = f

	return nil
}

func (it *IDLStruct)GetOrderTagMembers() []*IDLField {
	if it.orderTagMembers == nil{
		var tags []int
		for _, tag := range it.TagMembers{
			tags = append(tags, int(tag.FieldTag))
		}
		sort.Ints(tags)

		for _, tag := range tags{
			it.orderTagMembers = append(it.orderTagMembers, it.TagMembers[int64(tag)])
		}
	}
	return it.orderTagMembers
}

func (it *IDLStruct)Marshal(crtFileName string, protoType prototype.ProtoType) string{
	typeName := it.name_
	if crtFileName != it.GetPackage() {
		typeName = it.GetName()
	}
	s := fmt.Sprintf("struct %s{\n", typeName)

	var tags []int
	for _, tag := range it.TagMembers{
		tags = append(tags, int(tag.FieldTag))
	}
	sort.Ints(tags)

	for _, tag := range tags{
		s += "    " + it.TagMembers[int64(tag)].Marshal(crtFileName, protoType)
		s += "\n"
	}
	s += "}"

	return s
}