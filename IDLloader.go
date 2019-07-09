package idlparser

import (
	"code.byted.org/ee/lobster-idlloader/types"
	log_p "code.byted.org/ee/lobster-pliers/util/log-p"
	"context"
	"fmt"
	"strings"
	"sync"
)

var	g_MicroServiceMap sync.Map

func GetMicroService(ctx context.Context, name string) *types.MicroModule{
	m, ok := g_MicroServiceMap.Load(name)
	if ok{
		mod, ok1 := m.(*types.MicroModule)
		if ok1{
			metrics.EmitCounter("get.succ", map[string]string{"name":name})
			metrics.EmitCounter("get.succ_exact", map[string]string{"name":name})
			return mod
		}
	}

	m, ok = g_MicroServiceMap.Load(strings.ToLower(name))
	if ok{
		mod, ok2 := m.(*types.MicroModule)
		if ok2{
			metrics.EmitCounter("get.succ", map[string]string{"name":name})
			metrics.EmitCounter("get.succ_lower", map[string]string{"name":name})
			return mod
		}
	}

	metrics.EmitCounter("get.fail", map[string]string{"name":name})

	return nil
}

func reloadMicroService(ctx context.Context, name string, idlFiles types.IDLFileCollection)(*types.MicroModule, error) {
	ret := types.NewMicroModule(name, idlFiles)
	for _, idlFile := range idlFiles.GetSortedIDLFiles(){
		if !idlFile.ProtoType().IsValid(){
			metrics.EmitCounter("parse.invalid_file_type", map[string]string{"name":name})
			log_p.Errorf(ctx,"%s No parser for proto:%s, file:%s", name, idlFile.ProtoType().Str(), idlFile.GetFileName())
			continue
		}

		parser := NewParser(idlFile.ProtoType())
		if nil == parser {
			metrics.EmitCounter("parse.unsupported_proto_type", map[string]string{"name":name})
			log_p.Errorf(ctx,"%s No parser for proto:%s, file:%s", name, idlFile.ProtoType().Str(), idlFile.GetFileName())
			continue
		}

		idlFile.SetMod(ret)
		err := parser.Parse(ctx, ret, idlFile)
		if err != nil{
			metrics.EmitCounter("parse.file_error", map[string]string{"name":name})
			log_p.Errorf(ctx,"%s parser file:%s error:%s", name, idlFile.FileName, err.Error())
			continue
		}
	}

	ret.LnkTypedef()
	errs := ret.GetErrors()
	if nil != errs && len(errs) != 0{
		for _, err := range errs{
			log_p.Errorf(ctx,"%s parser error:%s", name, err)
		}

		metrics.EmitCounter("parse.error", map[string]string{"name":name})
		return nil, fmt.Errorf("the following errors when loading [%s]:\n%s", name, strings.Join(errs, "\n"))
	}

	g_MicroServiceMap.Store(name, ret)
	g_MicroServiceMap.Store(strings.ToLower(name), ret)
	metrics.EmitCounter("parse.succ", map[string]string{"name":name})

	return ret, nil
}

func LoadMicroService(ctx context.Context, name string, fileCollection types.IDLFileCollection, force bool)(*types.MicroModule, error){
	if name == ""{
		metrics.EmitCounter("request.invalid.name_is_nil", nil)
		return nil, fmt.Errorf("invalid_load_request, name is empty.")
	}

	if fileCollection.IsEmpty() {
		metrics.EmitCounter("request.invalid.idl_empty", map[string]string{"name":name})
		return nil, fmt.Errorf("invalid_load_request, idl files is empty.")
	}

	metrics.EmitCounter("request.vaild", map[string]string{"name":name})

	if force{
		metrics.EmitCounter("reload.force", map[string]string{"name":name})
		return reloadMicroService(ctx, name, fileCollection)
	}

	existMod := GetMicroService(ctx, name)
	if existMod == nil{
		metrics.EmitCounter("reload.new_micro_service", map[string]string{"name":name})
		return reloadMicroService(ctx, name, fileCollection)
	}

	existIdlFileCollection := existMod.GetIDLFileCollection()
	if !existIdlFileCollection.IsEqual(fileCollection){
		metrics.EmitCounter("reload.new_file", map[string]string{"name":name})
		return reloadMicroService(ctx, name, fileCollection)
	}

	metrics.EmitCounter("reload.ignore", map[string]string{"name":name})
	return existMod, nil
}


func BatchLoadMicroServices(ctx context.Context, idlFileCollectionMap map[string]types.IDLFileCollection, force bool)(map[string]*types.MicroModule, error){
	var errs []string
	ret := make(map[string]*types.MicroModule)
	for modName, idlFiles := range idlFileCollectionMap {
		mod, err := LoadMicroService(ctx, modName, idlFiles, force)
		if err != nil{
			errs = append(errs, err.Error())
		}else{
			ret[modName] = mod
		}
	}

	if errs != nil || len(errs) > 0 {
		return nil, fmt.Errorf("load error:\n%s", strings.Join(errs, "\n"))
	}

	return ret, nil
}