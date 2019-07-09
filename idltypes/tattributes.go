package idltypes

type IDLLobAttr struct{
	key string
	val string
}

func NewIDLAttribute(k, v string) *IDLLobAttr {
	return &IDLLobAttr{
		key:k,
		val:v,
	}
}

type IDLLobAttrs struct{
	kvs []*IDLLobAttr
	name string
}

func NewIDLAttrs(ns IDLNamespace) *IDLLobAttrs {
	return &IDLLobAttrs{}
}

func (as *IDLLobAttrs)SetName(name string){
	as.name = name
}

func (as *IDLLobAttrs)Append(attr *IDLLobAttr){
	as.kvs = append(as.kvs, attr)
}
