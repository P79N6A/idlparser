package idltypes

import (
	//	"log"
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"github.com/sirupsen/logrus"
	"strconv"
)

const( //枚举 EReq
	T_REQUIRED = iota
	T_OPTIONAL
	T_OPT_IN_REQ_OUT
)

var fieldRequirednessMap = map[int]string{
	T_REQUIRED:"required",
	T_OPTIONAL:"optional",
	T_OPT_IN_REQ_OUT:"",
}

// 定义一个字段
type IDLField struct {
	IDLDefine
	fieldType IDLTypeI // 字段的类型

	FieldTag  int64  // 字段的Tag
	FieldName string // 字段的名字

	FieldDefault  interface{} // 字段的默认值
	fieldComment_ string

	mod_    *MicroModule
	idlFile *File // IDL文件

	Required     int // EReq
	XsdOptional_ bool
	XsdNillable_ bool
	XsdAttrs_    *IDLStruct
	Reference_   bool

	annotationsAble

	// 强制类型
	annotationsParsed bool
	forceType       *BaseType
	maxFileSize     int64
}

func (f IDLField)GetAnnotation(key string)string{
	if f.annos == nil{
		return ""
	}

	return ""//f.annotations.GetValue(key)
}

func (f IDLField)GetFileName() string{
	if f.idlFile != nil{
		return f.idlFile.GetFileName()
	}

	return "<nil>"
}

func (f IDLField)String()string{
	return fmt.Sprintf("{file=%s,tag=%d,name=%s,type=%s}", f.GetFileName(), f.FieldTag, f.FieldName, f.GetFieldType())
}

func (f *IDLField)GetFieldType() IDLTypeI{
	if f.fieldType == nil{
		return nil
	}

	tdef, ok := f.fieldType.(*IDLTypedef)
	if ok{
		f.fieldType = tdef.GetRealType()
	}

	return f.fieldType
}

// ShortDebugString
func (f *IDLField)DebugString() string{
	return fmt.Sprintf("%s %s %s = %d", fieldRequirednessMap[f.Required], f.fieldType.GetName(), f.FieldName, f.FieldTag)
}

func (f *IDLField)Marshal(crtFileName string, protoType prototype.ProtoType) string{
	var ret string
	if protoType == ProtoType_THRIFT{
		fileTypeName := f.fieldType.GetName()
		if crtFileName != f.fieldType.GetPackage() && f.fieldType.GetPackage() != ""{
			fileTypeName = fmt.Sprintf("%s.%s", f.fieldType.GetPackage(), f.fieldType.GetName())
		}

		ret = fmt.Sprintf("%4d: %s %s %s, //$ %s", f.FieldTag, fieldRequirednessMap[f.Required], fileTypeName, f.FieldName, f.GetComment())
	}
	return ret
}

func (m *IDLField)SetReq(ereq int){
	m.Required = ereq
}

func (m *IDLField)SetReference(b bool){
}

func (m *IDLField)SetValue(val interface{}){

}

const(
	supportedAnnotationKeyForceType = "forceType"
	supportedAnnotationKeyMaxSize = "maxSize"

	supporedForceTypeBool = "bool"
	supporedForceTypeI8 = "i8"
	supporedForceTypeI16 = "i16"
	supporedForceTypeI32 = "i32"
	supporedForceTypeI64 = "i64"
	supporedForceTypeDouble = "double"
	supporedForceTypeString = "string"
	supporedForceTypeBinary = "binary"
)

func (m *IDLField) initForceType(respType *IDLAnnotation){
	_, ok := m.fieldType.(*BaseType)
	if !ok{
		logrus.Warnf("invalid_force_type: not base")
		return
	}

	if respType == nil{
		logrus.Warnf("invalid_force_type: respType is nil")
		return
	}

	switch respType.val {
	case supporedForceTypeString:
		m.forceType = BaseTypeString
	case supporedForceTypeBinary:
		m.forceType = BaseTypeBinary
	case supporedForceTypeBool:
		m.forceType = BaseTypeBool
	case supporedForceTypeI8:
		m.forceType = BaseTypeI8
	case supporedForceTypeI16:
		m.forceType = BaseTypeI16
	case supporedForceTypeI32:
		m.forceType = BaseTypeI32
	case supporedForceTypeI64:
		m.forceType = BaseTypeI64
	case supporedForceTypeDouble:
		m.forceType = BaseTypeDouble
	default:
		logrus.Warnf("invalid_force_type: unsupported ", respType)
	}
}

// 格式为 maxSize=x.yM/K
func (m *IDLField) initMaxSize(maxSize *IDLAnnotation){
	if maxSize == nil{
		logrus.Warnf("invalid_max_size: nil parameter")
		return
	}

	realType, ok := m.fieldType.(*IDLList)
	if !ok{
		logrus.Warnf("target type is not file list. actual:%s", m.fieldType.GetName())
		return
	}

	valType := realType.GetValueType()
	if valType == nil || valType.GetName() != "File"{
		logrus.Warnf("target type is not file list. actual value type is:%s", valType)
	}

	valLen := len(maxSize.val)
	if valLen == 0{
		logrus.Warnf("invalid maxSize, empty value.")
		return
	}

	var unit float64
	var num float64
	var err error
	switch maxSize.val[valLen-1] {
	case 'm':
		unit = 1024*1024 // 1M
		num, err = strconv.ParseFloat(maxSize.val[:valLen-1], 10)
	case 'M':
		unit = 1024*1024 // 1M
		num, err = strconv.ParseFloat(maxSize.val[:valLen-1], 10)
	case 'k':
		unit = 1024 // 1k
		num, err = strconv.ParseFloat(maxSize.val[:valLen-1], 10)
	case 'K':
		unit = 1024 // 1k
		num, err = strconv.ParseFloat(maxSize.val[:valLen-1], 10)
	default:
		unit = 1 // 1
		num, err = strconv.ParseFloat(maxSize.val, 10)
	}

	if err != nil{
		m.maxFileSize = 1*1024*1024 // 默认最大1M
		return
	}

	m.maxFileSize = int64(unit*num)
}

func (m *IDLField)parseAnnotations(){
	if !m.annotationsParsed {
		m.annotationsParsed = true
		if m.annos != nil {
			for k, v := range m.annos.vals {
				switch k {
				case supportedAnnotationKeyForceType:
					m.initForceType(v)
				case supportedAnnotationKeyMaxSize:
					m.initMaxSize(v)
				}
			}
		}
	}
}

func (m *IDLField) GetForceType()(*BaseType){
	m.parseAnnotations()
	return m.forceType
}

func (m *IDLField) GetMaxFileSize() int64{
	m.parseAnnotations()
	return m.maxFileSize
}

func (m *IDLField)SetXsdOptional(b bool){
}

func (m *IDLField)SetXsdNillable(b bool){
}

func (m *IDLField)SetXsdAttrs(attrs *IDLStruct) {
}

func (m *IDLField)SetComment(comment string) {
	if comment != ""{
		m.fieldComment_ = comment
	}
}

func (m *IDLField)GetComment()string{
	if m.fieldComment_ == ""{
		return "<未设置>"
	}

	return m.fieldComment_
}