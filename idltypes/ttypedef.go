package idltypes

import (
	"fmt"
)

//
// typedef
type Typedef struct{
	typeName string
	realType_ IDLTypeI

	pkg *Package
}

// typedef的名字: {file}:{type_name}
func (def Typedef)Name() string{
	return fmt.Sprintf("%s:%s", def.packageName, def.name_)
}

func (def Typedef)String() string{
	realType := def.GetReftype()
	return fmt.Sprintf("{type=typedef, name=%s, real=%s}", def.Name(), realType)
}

func (def *Typedef)Package() *Package{

}

func (def *Typedef) GetReftype() IDLTypeI{
	ret := def.realType_
	for i := 0; i < 5; i ++{
		if ret == nil{
			return nil
		}

		switch realValue := ret.(type){
		case *Typedef:
			ret = realValue.realType_
		default:
			return ret
		}
	}

	switch ret.(type){
	case *Typedef:
		return nil
	}

	return ret
}
