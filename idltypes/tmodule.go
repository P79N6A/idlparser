package idltypes

import (
	"fmt"
	"github.com/afLnk/idltranslator/prototype"
	"sort"
	"strings"
	//"log"
)

//
// 定义一个微服务
type MicroModule struct{
	name_   string  // 微服务的Psm
	idls_   FileCollection
	kvs_    map[string]interface{}
	consts_ map[string]interface{}
	errors_ []string

	pkgs_ map[string]*Package
	functions_ map[string]*IDLFunction
}

func (mod *MicroModule)AddError(format string, a ...interface{}){
	mod.errors_ = append(mod.errors_, fmt.Sprintf(format, a...))
}

func (mod *MicroModule)GetErrors() []string{
	return mod.errors_
}

func (mod *MicroModule)GetName() string{
	if mod == nil{
		return "unknown_mod"
	}
	return mod.name_
}

func (mod *MicroModule) FileMarshal(protoType prototype.ProtoType) map[string]string{
	marshalFiles := make(map[string]string)
	if protoType == ProtoType_THRIFT{
		for _, idlFile := range mod.idls_.GetSortedIDLFiles(){
			marshalFiles[idlFile.GetFileName()] = idlFile.Marshal(protoType)
		}
	}else if protoType == ProtoType_PROTOBUF{

	}else{ // 默认mashal成json

	}

	return marshalFiles
}

func (mod *MicroModule) PackageMarshal(protoType prototype.ProtoType) map[string]string{
	marshalPkgs := make(map[string]string)
	if protoType == ProtoType_THRIFT{
		for _, idlPkg := range mod.pkgs_{
			marshalPkgs[idlPkg.pkgName] = idlPkg.Marshal(protoType)
		}
	}else if protoType == ProtoType_PROTOBUF{

	}else{ // 默认mashal成json

	}

	return marshalPkgs
}


func (mod *MicroModule) GetIDLFileCollection() FileCollection {
	return mod.idls_
}

func (mod *MicroModule)LnkTypedef(){
	for _, idlPkg := range mod.pkgs_ {
		for _, def := range idlPkg.defs_ {
			def.realType_ = idlPkg.GetRealType(def.GetName())
			if def.realType_ == nil{
				mod.AddError("type %s.%s is invalid.", def.GetPackage(), def.GetName())
			}
		}
	}

	if mod.functions_ == nil{
		mod.functions_ = make(map[string]*IDLFunction)
	}

	for _, idlPkg := range mod.pkgs_ {
		for _, idlSvc := range idlPkg.svcs_ {
			for _, fuc := range idlSvc.Functions {
				mod.functions_[fuc.GetName()] = fuc
				mod.functions_[strings.ToLower(fuc.GetName())] = fuc
			}
		}
	}
}

func (mod *MicroModule)AddInclude(k, v interface{}){

}

func (mod *MicroModule)AddCppInclude(k interface{}){
}

func (mod *MicroModule)addConst(name string, val interface{}) {
	if mod.consts_ == nil{
		mod.consts_ = make(map[string]interface{})
	}

	_, ok := mod.consts_[name]
	if ok{
		return
	}

	mod.consts_[name] = val
}

func (mod *MicroModule)AddConst(kv *IDLConst){
	mod.addConst(kv.key, kv.val)
}

func (mod *MicroModule) GetStruct(pkgName, name string) *IDLStruct{
	idlPkg := mod.GetIdlPkg(pkgName)
	return idlPkg.strcts_[name]
}

func (mod* MicroModule)GetFunction(name string) *IDLFunction{
	ret, ok := mod.functions_[name]
	if !ok{
		ret, _ = mod.functions_[strings.ToLower(name)]
	}

	return ret
}

type IdlFuncSlice []*IDLFunction
func (c IdlFuncSlice) Len() int {
	return len(c)
}

func (c IdlFuncSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c IdlFuncSlice) Less(i, j int) bool {
	return c[i].GetName() < c[j].GetName()
}

func (mod *MicroModule) GetAllServices() []*IDLService {
	var result []*IDLService
	for _, idlPkg := range mod.pkgs_ {
		for _, idlSvc := range idlPkg.svcs_ {
			result = append(result, idlSvc)
		}
	}

	return result
}

func (mod *MicroModule) GetAllFunction() []*IDLFunction {
	var result IdlFuncSlice
	for _, idlPkg := range mod.pkgs_ {
		for _, idlSvc := range idlPkg.svcs_ {
			for _, fuc := range idlSvc.Functions {
				result = append(result, fuc)
			}
		}
	}

	sort.Sort(result)

	return result
}


func (mod *MicroModule) GetAllStructs() []*IDLStruct {
	var result []*IDLStruct
	for _, idlPkg := range mod.pkgs_ {
		for _, idlStrct := range idlPkg.strcts_ {
			result = append(result, idlStrct)
		}
	}

	return result
}

func (mod* MicroModule) GetOrNewPkg(pkgName string) *Package {
	existPackage := mod.GetIdlPkg(pkgName)
	if existPackage == nil{
		mod.pkgs_[pkgName] = NewIdlPackage(pkgName)
	}

	return mod.GetIdlPkg(pkgName)
}

func (mod* MicroModule) GetIdlPkg(pkgName string) *Package {
	existPkg, ok := mod.pkgs_[pkgName]
	if ok {
		return existPkg
	}

	return nil
}

func (mod *MicroModule)AddType(typ IDLTypeI){
	//log.Printf("AddType:%v\n", typ)
}

func (mod* MicroModule)AddStruct (strct *IDLStruct){
	idlFile := mod.GetOrNewPkg(strct.GetPackage())
	_, ok := idlFile.strcts_[strct.GetName()]
	if ok{
		mod.AddError("dumplicated struct:%s.%s", strct.GetPackage(), strct.GetName())
		return
	}
	//log.Printf("[AddStruct]%s.%s\n", idlFile.GetName(), strct.GetName())
	idlFile.strcts_[strct.GetName()] = strct
}

func (mod* MicroModule)AddPlaceholderTypedef (typeName string, f *File) *IDLTypedef{
	pkgName := f.GetFileName()
	pos := strings.LastIndexByte(typeName, '.')
	if pos != -1{ // 在当前文件里面找
		pkgName = typeName[:pos]
		typeName = typeName[pos+1:]
	}

	idlPkg := mod.GetOrNewPkg(pkgName)
	placeHolderTypedef := &IDLTypedef{
		realType_:nil,
		mod_:f.mod_,
	}

	placeHolderTypedef.packageName = pkgName
	placeHolderTypedef.name_ = typeName
	existDef, ok := idlPkg.defs_[typeName]
	if ok{
		return existDef
	}

	idlPkg.defs_[typeName] = placeHolderTypedef

	return placeHolderTypedef
}

func (mod* MicroModule)AddTypedef (def *IDLTypedef){
	idlPkg := mod.GetOrNewPkg(def.GetPackage())
	_, ok := idlPkg.strcts_[def.GetName()]
	if ok{
		mod.AddError("dumplicated typedef:%s.%s", def.GetPackage(), def.GetName())
		return
	}

	//log.Printf("AddTypedef:%s:%s, realType:%v", def.GetPackage(), def.GetName(), def.realType_)

	existDef, ok := idlPkg.defs_[def.GetName()]
	if ok{
		if def.realType_ != nil{
			existDef.realType_ = def.realType_
		}else{
			def.realType_ = existDef.realType_
		}
	}

	idlPkg.defs_[def.GetName()] = def
}

func (mod* MicroModule)AddService (svc* IDLService){
	idlPkg := mod.GetOrNewPkg(svc.GetPackage())
	_, ok := idlPkg.svcs_[svc.GetName()]
	if ok{
		mod.AddError("[MicroMod][AddService]Fail:%s already exist", svc.GetName())
		return
	}

	idlPkg.svcs_[svc.GetName()] = svc
}

func (mod* MicroModule)AddEnum (enum* IDLEnum){
	idlPkg := mod.GetOrNewPkg(enum.packageName)
	_, ok := idlPkg.enums_[enum.name]
	if ok{
		mod.AddError("dumplicated typedef:%s.%s", enum.packageName, enum.name)
		return
	}

	idlPkg.enums_[enum.name] = enum

	for _, ev := range enum.vals{
		idlPkg.addConst(ev.getConstName(), ev.val)
	}
}

func (scope* MicroModule)AddXception (k interface{}){
}

func (mod* MicroModule)SearchType(n string, f *File) IDLTypeI{
	pkgName := f.GetPackage()
	typeName := n
	pos := strings.LastIndexByte(n, '.')
	if pos != -1{ // 在当前文件里面找
		pkgName = n[:pos]
		typeName = n[pos+1:]
	}

	idlPkg := mod.GetOrNewPkg(pkgName)

	return idlPkg.GetType(typeName)
}

// 这个解析文件时，查找类型使用
// 如果找不到，就创建一个typedef
func (mod* MicroModule)SearchService(n string, f *File) *IDLService{
	pkgName := f.GetPackage()
	typeName := n
	pos := strings.LastIndexByte(n, '.')
	if pos != -1{ // 在当前文件里面找
		pkgName = n[:pos]
		typeName = n[pos+1:]
	}

	idlPkg := mod.GetOrNewPkg(pkgName)

	return idlPkg.svcs_[typeName]
}

func NewMicroModule(name string, idlFiles FileCollection) *MicroModule{
	return &MicroModule{
		name_: name,
		idls_: idlFiles,
		pkgs_: make(map[string]*Package),
	}
}