package idltypes

import (
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"github.com/sirupsen/logrus"
)

type IDLEnumValue struct{
	em *IDLEnum
	name string
	val int64
	comment string
	packageName string

	annotationsAble
}

func (ev *IDLEnumValue)SetComment(comment string){
	ev.comment = comment
}

func (ev IDLEnumValue)getConstName() string{
	if ev.em == nil{
		logrus.Warnf("invalid enum value:%s=%d, no emun", ev.name, ev.val)
		return ""
	}

	return fmt.Sprintf("%s_%s", ev.em.name, ev.name)
}

func NewIDLEnumValue(name string, val int64, packageName string) *IDLEnumValue{
	return &IDLEnumValue{
		name:name,
		val:val,
		packageName: packageName,
	}
}

type IDLEnum struct{
	name string
	packageName string
	vals map[string]*IDLEnumValue
	annotationsAble
}

func (em *IDLEnum)Append(val *IDLEnumValue){
	if val == nil{
		return
	}

	val.em = em
	em.vals[val.name] = val
}

func (em *IDLEnum)SetName(name string){
	em.name = name
}

func (em IDLEnum)GetName() string{
	return em.name
}

func (em IDLEnum)GetFileName()string{
	return ""
}

func (em IDLEnum)GetPackage()string{
	return em.packageName
}

func (em IDLEnum)Marshal(crtFileName string, protoType prototype.ProtoType) string{
	return ""
}

func (em IDLEnum)JsonDefault() string{
	return ""
}

func (em IDLEnum)Json() string{
	return ""
}

func (em IDLEnum)String() string{
	return ""
}

func (em IDLEnum)GetValues() map[string]int64{
	ret := make(map[string]int64)
	for _, v := range em.vals{
		ret[v.name] = v.val
	}
	return ret
}

func (em IDLEnum)GetValue(key string) (int64, bool){
	ret, ok := em.vals[key]
	if !ok{
		return 0, false
	}

	return ret.val, ok
}

func NewIDLEnum(ns IDLNamespace, packageName string) *IDLEnum{
	ret := IDLEnum{
		vals: make(map[string]*IDLEnumValue),
		packageName: packageName,
	}
	return &ret
}
