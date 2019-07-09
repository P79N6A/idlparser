package prototype

type ProtoType string

func NewPrototype(fileSuffix string) ProtoType{
	switch fileSuffix {
	case thriftSuffix:
		return Thrift
	case protobufferSuffix:
		return ProtoBuffer
	}

	return Invalid
}

const(
	Invalid = ProtoType("invalid")
	Thrift  = ProtoType("thrift")
	ProtoBuffer = ProtoType("pb")

	thriftSuffix = ".thrift"
	protobufferSuffix = ".pb"
)

