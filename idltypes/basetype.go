package idltypes

// base idltypes: int8, int16....
type BaseType struct {
	name string
}

func newSimplpeBaseType(name string) BaseType{
	return BaseType{ name : name}
}

var(
	BaseTypeStop = newSimplpeBaseType("STOP")
	BaseTypeVoid = newSimplpeBaseType("VOID")
	BaseTypeBool = newSimplpeBaseType("BOOL")
	BaseTypeByte = newSimplpeBaseType("BYTE")
	BaseTypeI8 = newSimplpeBaseType("I8")
	BaseTypeDouble = newSimplpeBaseType("DOUBLE")
	BaseTypeI16 = newSimplpeBaseType("I16")
	BaseTypeI32 = newSimplpeBaseType("I32")
	BaseTypeI64 = newSimplpeBaseType("I64")
	BaseTypeString = newSimplpeBaseType("STRING")
	BaseTypeUTF7 = newSimplpeBaseType("UTF7")
	BaseTypeUTF8 = newSimplpeBaseType("UTF8")
	BaseTypeUTF16 = newSimplpeBaseType("UTF16")
	BaseTypeBinary = newSimplpeBaseType("BINARY")
	BaseTypeSlist = newSimplpeBaseType("SLIST")
)
