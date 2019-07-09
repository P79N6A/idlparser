package idlparser

import (
	"code.byted.org/ee/lobster-idlloader/thrift"
	"code.byted.org/ee/lobster-idlloader/types"
	log_p "code.byted.org/ee/lobster-pliers/util/log-p"
	"context"
	"github.com/afLnk/idltranslator/idltypes"
	"github.com/afLnk/idltranslator/prototype"
	"strings"
)

type ParserI interface {
	Type() types.ProtoType
	Parse(ctx context.Context, mod *types.MicroModule, file *types.IDLFile) error
}

func NewParser(protoType prototype.ProtoType) ParserI{
	switch protoType {
	case prototype.Thrift:
		return thrift.Parser{}

	}
	ret, ok := g_lexers[protoType]
	if ok{
		return ret
	}

	return nil
}

//
// parser idl files.
func Parse(fileCollection idltypes.FileCollection) idltypes.Result {
	parser := NewParser().Parse()
	ret := idltypes.NewResult(fileCollection)
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