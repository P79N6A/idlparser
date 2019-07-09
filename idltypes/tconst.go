package idltypes

type IDLConst struct{
	key string
	val interface{}
	valType IDLTypeI
	comment string
}

func (cst IDLConst)GetName() string{
	return cst.key
}

func NewIDLConst(k string, v interface{}, t IDLTypeI, comment string) *IDLConst{
	return &IDLConst{
		key:k,
		val:v,
		valType: t,
		comment: comment,
	}
}
