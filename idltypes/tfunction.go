package idltypes

import (
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"log"

	// "log"
)

const(
	EmptyJson = "{\n}"
)

func getJsonSpaceHolder(level int) (string){
	var ret string
	for tmp := 0; tmp < level; tmp ++{
		ret += "\t\t"
	}
	return  ret
}

func genTypeJsonObj(idlType IDLTypeI, defaultVal interface{}, level int, isReq bool, ets *ExpandedTypes) string {
	var ret string
	if idlType == nil{
		return "---haha---"
	}

	if ets.IsExpanded(idlType){
		return "{}"
	}

	newEts := NewExpandedTypes(ets, idlType)
	switch realType := idlType.(type){
	case *IDLStruct:
		ret = "{"
		is1st := true
		for _, x := range realType.GetOrderTagMembers() {
			var fieldType = x.GetFieldType()
			if nil == fieldType {
				continue
			}

			if idlType.GetName() == "Head" && x.FieldTag > 30{
				continue
			}

			if !is1st {
				ret += fmt.Sprintf(",\n%s%q:%s", getJsonSpaceHolder(level), x.FieldName, genTypeJsonObj(x.GetFieldType(), x.FieldDefault, level+1,true, newEts))
			}else{
				ret += fmt.Sprintf("\n%s%q:%s",getJsonSpaceHolder(level), x.FieldName, genTypeJsonObj(x.GetFieldType(), x.FieldDefault, level+1,true, newEts))
			}

			is1st = false
		}
		ret += "\n" + getJsonSpaceHolder(level - 1) + "}"
		return ret
	case *IDLSet:
		oneObjStr := genTypeJsonObj(realType.GetKeyType(), nil, level + 1,isReq, newEts)
		ret = "["
		ret += "\n" + getJsonSpaceHolder(level) + oneObjStr
		ret += ",\n" + getJsonSpaceHolder(level) + oneObjStr
		ret += "\n" + getJsonSpaceHolder(level-1) + "]"
	case *IDLList:
		oneObjStr := genTypeJsonObj(realType.GetValueType(), nil, level + 1,isReq, newEts)
		ret = "["
		ret += "\n" + getJsonSpaceHolder(level) + oneObjStr
		ret += ",\n" + getJsonSpaceHolder(level) + oneObjStr
		ret += "\n" + getJsonSpaceHolder(level-1) + "]"
	case *IDLMap:
		oneObjStr := genTypeJsonObj(realType.GetValueType(), nil, level + 1, isReq, newEts)
		ret = "{"
		ret += "\n" + getJsonSpaceHolder(level) + "\"key1\":" + oneObjStr
		ret += ",\n" + getJsonSpaceHolder(level) + "\"key2\":"+ oneObjStr
		ret += "\n" + getJsonSpaceHolder(level-1) + "}"
		return ret

	case *BaseType:
		if realType == BaseTypeString{
			return "\"\""
		}
		if realType == BaseTypeDouble{
			return "0.0"
		}
		return "0"
	case *IDLEnum:
		for _, v := range realType.vals{
			if ret == ""{
				ret += fmt.Sprintf("\"%d(%s)", v.val, v.name)
			}else{
				ret += fmt.Sprintf("/%d(%s)", v.val, v.name)
			}
		}
		ret += "\""
	default:
		log.Printf("2unsupported type:%T, %s", realType, realType.GetName())
	}

	return ret
}

// 定义一个接口
type IDLFunction struct{
	IDLDefine
	IDLNamespace
	RspType  IDLTypeI
	ReqType  *IDLStruct
	comment_ string
	attrs_  *IDLLobAttrs

	annotationsAble
}

func (fnc *IDLFunction)GetName() string{
	if fnc == nil{
		return "unknown_func"
	}
	return fnc.name_
}

func (fnc *IDLFunction)DebugString() string{
	return fmt.Sprintf("%s %s(%s)", fnc.RspType.GetName(), fnc.GetName(), fnc.ReqType.GetName())
}

func (fnc *IDLFunction)SetComment(comment string) {
	if comment != ""{
		//log.Println("comment:", comment)
		fnc.comment_ = comment
	}
}

func (fnc *IDLFunction)GetComment() string{
	if fnc.comment_ == ""{
		return "<未设置>"
	}

	return fnc.comment_
}

func (fnc *IDLFunction)GetJsonReq() string {
	for _, tagMember := range fnc.ReqType.TagMembers{
		return tagMember.GetFieldType().Json()
	}
	return "{}"
}

func (fnc *IDLFunction)GetJsonReqV2(reqType int64)(string, error){
	if fnc == nil{
		return EmptyJson, fmt.Errorf("function is nil")
	}

	var realReq *IDLStruct
	var ok bool
	for _, tagMember := range fnc.ReqType.TagMembers{
		realReq, ok = tagMember.GetFieldType().(*IDLStruct)
		if !ok{
			return EmptyJson, fmt.Errorf("LobMethod's func member:%d is not struct", tagMember.FieldTag)
		}
	}

	if realReq == nil{
		return EmptyJson, fmt.Errorf("func  is invalid.")
	}

	ret := "{"
	is1st := true
	for _, x := range realReq.GetOrderTagMembers() {
		var fieldType = x.GetFieldType()
		if nil == fieldType {
			continue
		}

		if x.FieldName == "Base" && x.FieldTag == 255{
			continue
		}

		if x.FieldName == "Head" && x.FieldTag == 1{
			// 透传模式之外的，全过滤
			if reqType != RespType_Penetrate{
				continue
			}
		}

		if !is1st {
			ret += fmt.Sprintf(",\n\t\t%q:%s", x.FieldName, genTypeJsonObj(x.GetFieldType(), x.FieldDefault, 2,true, nil))
		}else{
			ret += fmt.Sprintf("\n\t\t%q:%s", x.FieldName, genTypeJsonObj(x.GetFieldType(), x.FieldDefault, 2,true, nil))
		}

		is1st = false
	}
	ret += "\n}"

	return string(ret), nil
}

func (fnc *IDLFunction)GetJsonRespV2(respType int64) (string, error){
	if fnc == nil{
		return EmptyJson, fmt.Errorf("LobMethod's func is nil")
	}

	realResp, ok := fnc.RspType.(*IDLStruct)
	if !ok{
		return EmptyJson, fmt.Errorf("func<%s> resp type<%s> is not struct", fnc.GetName(), fnc.RspType.GetName())
	}

	if realResp == nil{
		return EmptyJson, fmt.Errorf("LobMethod's request is invalid.")
	}

	ret := "{"
	level := 1
	is1st := true
	if respType == RespType_Code{
		ret = "{\n" + "\t\t\"code\": 0,\n" + "\t\t\"msg\":\"success\",\n"+"\t\t\"data\":{"
		level ++
	}

	for _, x := range realResp.GetOrderTagMembers() {
		var fieldType = x.GetFieldType()
		if nil == fieldType {
			continue
		}

		if x.FieldName == "Base" && x.FieldTag == 255{
			continue
		}

		// 除非是透传，不然Base和BaseResp全过滤
		if x.FieldName == "BaseResp" && x.FieldTag == 255 {
			if respType != RespType_Penetrate {
				continue
			}
		}

		if x.FieldName == "Head" && x.FieldTag == 1{
			// 透传模式
			if respType != RespType_Penetrate{
				continue
			}
		}

		if !is1st {
			ret += fmt.Sprintf(",\n%s%q:%s", getJsonSpaceHolder(level), x.FieldName, genTypeJsonObj(x.GetFieldType(), x.FieldDefault, level+1,false, nil))
		}else{
			ret += fmt.Sprintf("\n%s%q:%s",getJsonSpaceHolder(level), x.FieldName, genTypeJsonObj(x.GetFieldType(), x.FieldDefault, level+1,false, nil))
		}
		is1st = false
	}
	if respType == RespType_Code {
		ret += "\n\t\t}"
	}
	ret += "\n}"

	return string(ret), nil
}

func (fnc *IDLFunction)GetJsonResp() string {
	return fnc.RspType.Json()
}

func (fnc *IDLFunction)Mashal(crtFileName string, protoType prototype.ProtoType) string{
	var ret string
	if protoType == ProtoType_THRIFT{
		if fnc.RspType.GetPackage() != crtFileName{
			ret += fnc.RspType.GetPackage()
			ret += "."
		}
		ret += fnc.RspType.GetName()
		ret += " "
		ret += fnc.GetName()
		ret += "("
		for _, tagMember := range fnc.ReqType.TagMembers{
			if tagMember.fieldType.GetPackage() != crtFileName{
				ret += fmt.Sprintf("%d:", tagMember.FieldTag)
				ret += fnc.ReqType.GetPackage()
				ret += "."
			}
			ret += tagMember.fieldType.GetName()
			ret += " "
			ret += tagMember.FieldName
		}
		ret += "); //$"
		ret += fnc.GetComment()
	}

	return ret
}

func (it *IDLFunction)SetAttrs(attrs *IDLLobAttrs){
	it.attrs_ = attrs
}
