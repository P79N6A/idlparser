package idltypes

import (
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"log"
)

const(
	ContainerType_DEFAULT = 0
	ContainerType_MAP = 10
	ContainerType_SET = 20
	ContainerType_LIST = 30
)

const (
	containerMap containerType = 10
	containerSet containerType = 20
	containerList containerType = 30
)

type containerType int

type IDLContainerI interface {
	IDLTypeI
}

type IDLContainer struct{
	IDLDefine
	keyType       IDLTypeI
	valType       IDLTypeI
	containerType containerType
	mod_          *MicroModule

	annotationsAble
}

// TODO
func (c IDLContainer)String() string{
	return ""
}

func (c *IDLContainer)SetCppName(name string){
	log.Printf("[%s]cpp name:%s", c.GetName(), name)
}

func (c *IDLContainer)Marshal(crtFileName string, protoType prototype.ProtoType) string { // 把这个类型转化成字符串
	return ""
}

func (c *IDLContainer)GetKeyType() IDLTypeI{
	if c.keyType == nil{
		return nil
	}

	tdef, ok := c.keyType.(*Typedef)
	if ok{
		c.keyType = tdef.GetReftype()
	}

	return c.keyType
}

func (c *IDLContainer)GetValueType() IDLTypeI{
	if c.valType == nil{
		return nil
	}

	tdef, ok := c.valType.(*IDLTypedef)
	if ok{
		c.valType = tdef.GetRealType()
	}
	return c.valType
}


func (c *IDLContainer)Json() string{
	if c == nil{
		return "\"\""
	}

	var ret string

	switch c.containerType {
	case ContainerType_MAP:
		ret += "{"
		ret += "\"key1\":"
		ret += c.valType.JsonDefault()
		ret += ",\"key2\":"
		ret += c.valType.JsonDefault()
		ret += "}"
	case ContainerType_LIST:
		ret += "["
		ret += c.valType.JsonDefault()
		ret += ","
		ret += c.valType.JsonDefault()
		ret += "]"
	case ContainerType_SET:
		ret += "["
		ret += c.keyType.JsonDefault()
		ret += ","
		ret += c.keyType.JsonDefault()
		ret += "]"
	}

	return ret
}

func (c *IDLContainer)JsonDefault() string{
	return c.Json()
}


// 返回类型名字 map<xx,xx>的方式
func (c IDLContainer)Name() string{
	switch c.containerType {
	case containerMap:
		return fmt.Sprintf("MAP<%s,%s>", c.keyType.GetName(), c.valType.GetName())
	case containerSet:
		return fmt.Sprintf("SET<%s>", c.keyType.GetName())
	case containerList:
		return fmt.Sprintf("LIST<%s>", c.valType.GetName())
	}

	return fmt.Sprintf("UnknownContainer<%s,%s>", c.keyType.GetName(), c.valType.GetName())
}

// 返回类型名字 map<xx,xx>的方式
func (c *IDLContainer)GetName() string{
	if c == nil{
		return "NilContainer"
	}

	switch c.containerType {
	case ContainerType_MAP:
		return fmt.Sprintf("MAP<%s,%s>", c.keyType.GetName(), c.valType.GetName())
	case ContainerType_LIST:
		return fmt.Sprintf("LIST<%s>", c.valType.GetName())
	case ContainerType_SET:
		return fmt.Sprintf("SET<%s>", c.keyType.GetName())
	}

	return fmt.Sprintf("UnknownContainer<%s,%s>", c.keyType.GetName(), c.valType.GetName())
}

func (c *IDLContainer)SetContainer(t containerType, kType IDLTypeI, vType IDLTypeI){
	c.containerType = t
	c.keyType = kType
	c.valType = vType
	//log.Printf("[Container][New]%s", ret.Str())
}

// 返回类型
func (c *IDLContainer)Str() string{
	return c.GetName()
}

type IDLMap struct{
	IDLContainer
}

type IDLList struct{
	IDLContainer
}

type IDLSet struct {
	IDLContainer
}



