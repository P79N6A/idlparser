%{
package thrift

import(
    //"fmt"
    "log"
    idl "code.byted.org/ee/lobster-idlloader/types"
)
%}

%{
/**
 * This global variable is used for automatic numbering of field indices etc.
 * when parsing the members of a struct. Field values are automatically
 * assigned starting from -1 and working their way down.
 */
    var y_field_val = -1
/**
 * This global variable is used for automatic numbering of enum values.
 * y_enum_val is the last value assigned; the next auto-assigned value will be
 * y_enum_val+1, and then it continues working upwards.  Explicitly specified
 * enum values reset y_enum_val to that value.
 */
    var y_enum_val int64 = -1
    var G_arglist = 0
    var struct_is_struct int64 = 0
    var struct_is_union int64 = 1
    var g_strict = 255
    var G_program_doctext_candidate string
    var G_program_doctext_status int
%}

/**
 * This structure is used by the parser to hold the data types associated with
 * various parse nodes.
 */
%union {
  id string
  iconst int64
  dconst float64
  bconst bool
  tbool bool
  tdoc *idl.IDLDefine
  idltype idl.IDLTypeI
  tbase *idl.IDLBaseType
  ttypedef *idl.IDLTypedef
  tanno *idl.IDLAnnotation
  tannos *idl.IDLAnnotations
  tlobattr *idl.IDLLobAttr
  tlobattrs *idl.IDLLobAttrs
  tenum *idl.IDLEnum
  tenumv *idl.IDLEnumValue
  senum *idl.IDLEnum
  senumv *idl.IDLEnumValue
  tconst *idl.IDLConst
  tconstv interface{}
  idlconstlist []interface{}
  idlconstmap map[interface{}]interface{}
  tstruct *idl.IDLStruct
  tservice *idl.IDLService
  tfunction *idl.IDLFunction
  tfield *idl.IDLField
  dtext string
  ereq int // TField:EReq
  idldoci idl.IDLTypeI
  tlist *idl.IDLList
  tmap *idl.IDLMap
  tset *idl.IDLSet
  tcontainer idl.IDLContainerI
  idlfieldid int64
}

/**
 * Strings identifier
 */
%token<id>     tok_identifier
%token<id>     tok_literal
%token<dtext>  tok_doctext
%token<dtext>  tok_fieldcomment
%token<dtext>  tok_lobattr

/**
 * Constant values
 */
%token<iconst> tok_int_constant
%token<dconst> tok_dub_constant

/**
 * Header keywords
 */
%token tok_include
%token tok_namespace
%token tok_cpp_include
%token tok_cpp_type
%token tok_xsd_all
%token tok_xsd_optional
%token tok_xsd_nillable
%token tok_xsd_attrs

/**
 * Base datatype keywords
 */
%token tok_void
%token tok_bool
%token tok_string
%token tok_binary
%token tok_slist
%token tok_senum
%token tok_i8
%token tok_i16
%token tok_i32
%token tok_i64
%token tok_double

/**
 * Complex type keywords
 */
%token tok_map
%token tok_list
%token tok_set

/**
 * Function modifiers
 */
%token tok_oneway

/**
 * Thrift language keywords
 */
%token tok_typedef
%token tok_struct
%token tok_xception
%token tok_throws
%token tok_extends
%token tok_service
%token tok_enum
%token tok_const
%token tok_required
%token tok_optional
%token tok_union
%token tok_reference

/**
 * Grammar nodes
 */

%type<tbase>     BaseType
%type<tbase>     SimpleBaseType
%type<tcontainer>   ContainerType
%type<tcontainer>   SimpleContainerType
%type<tmap>   MapType
%type<tset>   SetType
%type<tlist>   ListType

%type<tdoc>      Definition
%type<idldoci>   TypeDefinition

%type<ttypedef>  Typedef

%type<tannos>   TypeAnnotations
%type<tannos>   TypeAnnotationList
%type<tanno>     TypeAnnotation
%type<id>        TypeAnnotationValue

%type<tfield>    Field
%type<idlfieldid> FieldIdentifier
%type<ereq>      FieldRequiredness
%type<idltype>   FieldType
%type<tconstv>   FieldValue
%type<tstruct>   FieldList
%type<tbool>     FieldReference

%type<tenum>     Enum
%type<tenum>     EnumDefList
%type<tenumv>    EnumDef
%type<tenumv>    EnumValue

%type<senum>     Senum
%type<senum>     SenumDefList
%type<id>        SenumDef

%type<tconst>    Const
%type<tconstv>   ConstValue
%type<idlconstlist>  ConstList
%type<idlconstlist>  ConstListContents
%type<idlconstmap>   ConstMap
%type<idlconstmap>   ConstMapContents

%type<iconst>    StructHead
%type<tstruct>   Struct
%type<tstruct>   Xception
%type<tservice>  Service

%type<tlobattr>       LobAttribute
%type<tlobattrs>      LobAttributeList

%type<tfunction> Function
%type<idltype>   FunctionType
%type<tservice>  FunctionList

%type<tstruct>   Throws
%type<tservice>  Extends
%type<tbool>     Oneway
%type<tbool>     XsdAll
%type<tbool>     XsdOptional
%type<tbool>     XsdNillable
%type<tstruct>   XsdAttributes
%type<id>        CppType

%type<dtext>     CaptureDocText
%type<dtext>     FieldComment

%%

/**
 * Thrift Grammar Implementation.
 *
 * For the most part this source file works its way top down from what you
 * might expect to find in a typical .thrift file, i.e. type definitions and
 * namespaces up top followed by service definitions using those types.
 */

Program:
  HeaderList DefinitionList
    {
      //log.Println("Program -> HeaderList DefinitionList")
      Program()
    }

CaptureDocText:
    {
      //log.Println("CaptureDocText")
      $$ = CaptureDocText()
    }

DestroyDocText:
    {
      DestoryDocText()
    }

/* We have to DestroyDocText here, otherwise it catches the doctext on the first real element. */
HeaderList:
  HeaderList DestroyDocText Header
    {
      // log.Println("[yacc]HeaderList -> HeaderList Header")
    }
|
    {
      // log.Println("[yacc]HeaderList -> ")
    }

Header:
  Include
    {
      //log.Println("[yacc]Header -> Include")
    }
| tok_namespace tok_identifier tok_identifier TypeAnnotations
    {
      //log.Println("[yacc]Header -> tok_namespace tok_identifier tok_identifier")
      crtFile.SetNamespaceWithLang($2, $3);
      if ($4 != nil) {
        crtFile.SetNamespaceAnnotations($2, $4)
      }
    }
| tok_namespace '*' tok_identifier
    {
      //log.Println("[yacc]Header -> tok_namespace * tok_identifier")
      crtFile.SetNamespaceWithLang("*", $3);
    }
| tok_cpp_include tok_literal
    {
      //log.Println("[yacc]Header -> tok_cpp_include tok_literal");
      crtFile.AddCppInclude($2)
    }

Include:
  tok_include tok_literal
    {
        IncludeFile($2)
    }

DefinitionList:
  DefinitionList CaptureDocText Definition
    {
      //log.Println("[yacc]DefintionList -> DefinitionList CaptureDocText Definition")
    }
|
    {
      //log.Println("[yacc]DefinitionList -> ");
    }

Definition:
  Const
    {
      //log.Println("[yacc]Definition -> Const")
      crtMicroMod.AddConst($1);
      $$ = idl.NewDefine($1.GetName())
    }
| TypeDefinition
    {
      //log.Println("[yacc]Definition -> TypeDefinition", $1.GetName())
      crtMicroMod.AddType($1)
      $$ = idl.NewDefine($1.GetName())
    }
| Service
    {
      //log.Println("[yacc]Definition -> Service")
      crtMicroMod.AddService($1)
      $$ = idl.NewDefine($1.GetName())
    }

TypeDefinition:
  Typedef
    {
      //log.Println("[yacc]TypeDefinition -> Typedef")
      crtMicroMod.AddTypedef($1)
      $$ = $1
    }
| Enum
    {
      //log.Println("[yacc]TypeDefinition -> Enum")
      crtMicroMod.AddEnum($1)
      $$ = $1
    }
| Senum
    {
        NotSupported("Senum")
    }
| Struct
    {
      //log.Println("[yacc]TypeDefinition -> Struct");
      crtMicroMod.AddStruct($1)
      $$ = $1
    }
| Xception
    {
      //log.Println("[yacc]TypeDefinition -> Xception");
      crtMicroMod.AddXception($1)
      $$ = $1
    }

CommaOrSemicolonOptional:
  ','
    {}
| ';'
    {}
|
    {}

Typedef:
  tok_typedef FieldType tok_identifier TypeAnnotations CommaOrSemicolonOptional
    {
      td := crtFile.NewTypedef($2, $3);
      $$ = td
      if $4 != nil {
        $$.SetAnnotations($4)
      }
    }

Enum:
  tok_enum tok_identifier '{' EnumDefList '}' TypeAnnotations
    {
      //log.Println("[yacc]Enum -> tok_enum tok_identifier { EnumDefList }")
      $$ = $4
      $$.SetName($2)
      if $6 != nil {
        $$.SetAnnotations($6)
      }
    }

EnumDefList:
  EnumDefList EnumDef
    {
      //log.Println("[yacc]EnumDefList -> EnumDefList EnumDef")
      $$ = $1
      $$.Append($2)
    }
|
    {
      //log.Println("[yacc]EnumDefList -> ")
      $$ = crtFile.NewEnum()
      y_enum_val = -1
    }

EnumDef:
  CaptureDocText EnumValue TypeAnnotations CommaOrSemicolonOptional FieldComment
    {
      $$ = $2
      $$.SetComment($5)
	  if $3 != nil {
        $$.SetAnnotations($3)
      }
    }

EnumValue:
  tok_identifier '=' tok_int_constant
    {
      //log.Println("[yacc]EnumValue -> tok_identifier = tok_int_constant")
      y_enum_val = $3
      $$ = crtFile.NewEnumValue($1, y_enum_val);
    }
 |
  tok_identifier
    {
      y_enum_val ++
      $$ = crtFile.NewEnumValue($1, y_enum_val)
    }

Senum:
  tok_senum tok_identifier '{' SenumDefList '}' TypeAnnotations
    {
      NotSupported("Senum", $2, $4, $6)
    }

SenumDefList:
  SenumDefList SenumDef
    {
      NotSupported("Senum", $1, $2)
    }
|
    {
      NotSupported("Senum")
    }

SenumDef:
  tok_literal CommaOrSemicolonOptional
    {
      NotSupported("Senum", $1)
    }

Const:
  tok_const FieldType tok_identifier '=' ConstValue CommaOrSemicolonOptional FieldComment
    {
      //log.Println("[yacc]Const -> tok_const FieldType tok_identifier = ConstValue")
      $$ = crtFile.NewConst($2, $3, $5, $7)
      crtMicroMod.AddConst($$)
    }

ConstValue:
  tok_int_constant
    {
      $$ = $1
    }
| tok_dub_constant
    {
      $$ = $1
    }
| tok_literal
    {
      $$ = $1
    }
| tok_identifier
    {
      $$ = $1
    }
| ConstList
    {
      $$ = $1
    }
| ConstMap
    {
      $$ = $1
    }

ConstList:
  '[' ConstListContents ']'
    {
      log.Println("[yacc]ConstList => [ ConstListContents ]")
      $$ = $2;
    }

ConstListContents:
  ConstListContents ConstValue CommaOrSemicolonOptional
    {
      log.Println("[yacc]ConstListContents => ConstListContents ConstValue CommaOrSemicolonOptional");
      $$ = $1;
      $$ = append($$, $2)
    }
|
    {
      $$ = nil
    }

ConstMap:
  '{' ConstMapContents '}'
    {
      $$ = $2;
    }

ConstMapContents:
  ConstMapContents ConstValue ':' ConstValue CommaOrSemicolonOptional
    {
      //log.Println("[yacc]ConstMapContents => ConstMapContents ConstValue CommaOrSemicolonOptional")
      $$ = $1
      $$[$2] = $4
    }
|
    {
      $$ = make(map[interface{}]interface{})
    }

StructHead:
  tok_struct
    {
      $$ = struct_is_struct
    }
| tok_union
    {
      $$ = struct_is_union
    }

Struct:
  StructHead tok_identifier XsdAll '{' FieldList '}' TypeAnnotations
    {
      $5.SetXsdAll($3)
      $5.SetUnion($1 == struct_is_union)
      $$ = $5
      $$.SetName($2);
      if ($7 != nil) {
        $$.SetAnnotations($7)
      }
    }

XsdAll:
  tok_xsd_all
    {
      $$ = true;
    }
|
    {
      $$ = false;
    }

XsdOptional:
  tok_xsd_optional
    {
      $$ = true
    }
|
    {
      $$ = false
    }

XsdNillable:
  tok_xsd_nillable
    {
      $$ = true
    }
|
    {
      $$ = false
    }

XsdAttributes:
  tok_xsd_attrs '{' FieldList '}'
    {
      $$ = $3
    }
|
    {
      $$ = nil
    }

Xception:
  tok_xception tok_identifier '{' FieldList '}' TypeAnnotations
    {
      $4.SetName($2)
      $4.SetXception(true)
      $$ = $4
      if $6 != nil {
        $$.SetAnnotations($6)
      }
    }

Service:
  LobAttributeList tok_service tok_identifier Extends '{' FlagArgs FunctionList UnflagArgs '}' TypeAnnotations
    {
      //log.Println("what's the hell")
      //log.Println("[yacc]Service -> tok_service tok_identifier { FunctionList }")
      $$ = $7
      $$.SetName($3);
      $$.SetExtends($4);
      if ($10 != nil) {
        $$.SetAnnotations($10)
      }
      if $1 != nil{
        $$.SetAttrs($1)
        $1.SetName("service_" + $3)
      }
    }

FlagArgs:
    {
       G_arglist = 1;
    }

UnflagArgs:
    {
       G_arglist = 0;
    }

Extends:
  tok_extends tok_identifier
    {
      log.Println("[yacc]Extends -> tok_extends tok_identifier")
      $$ = crtMicroMod.SearchService($2, crtFile)
    }
|
    {
      $$ = nil
    }

FunctionList:
  FunctionList Function
    {
      //log.Println("[yacc]FunctionList -> FunctionList Function", $1, $2)
      $$ = $1
      $$.AddFunction($2)
    }
|
    {
      //log.Println("[yacc]FunctionList -> ");
      $$ = crtFile.NewService()
    }

Function:
  CaptureDocText LobAttributeList Oneway FunctionType tok_identifier '(' FieldList ')' Throws TypeAnnotations CommaOrSemicolonOptional FieldComment
    {
      $7.SetReqStructName($5)
      $$ = crtFile.NewFunction($4, $5, $7, $9, $3)
      if ($10 != nil) {
        $$.SetAnnotations($10)
      }

      $$.SetComment($12)
      $$.SetAttrs($2)
      $2.SetName("func_" + $5)
    }

Oneway:
  tok_oneway
    {
      $$ = true
      //log.Println("oneway:", $$)
    }
|
    {
      $$ = false
      //log.Println("oneway:", $$)
    }

Throws:
  tok_throws '(' FieldList ')'
    {
      //log.Println("[yacc]Throws -> tok_throws ( FieldList )");
      $$ = $3
    }
|
    {
      //log.Println("Throws ->")
      $$ = crtFile.NewEmptyStruct()
    }

FieldList:
  FieldList Field
    {
      //log.Println("[yacc]FieldList -> FieldList , Field");
      $$ = $1
      err := $$.Append($2)
      if err != nil{
          panic(err.Error())
      }
    }
|
    {
      //log.Println("[yacc]FieldList -> ");
      y_field_val = -1;
      $$ = crtFile.NewStruct()
    }

Field:
  CaptureDocText FieldIdentifier FieldRequiredness FieldType FieldReference tok_identifier FieldValue XsdOptional XsdNillable XsdAttributes TypeAnnotations CommaOrSemicolonOptional FieldComment
    {
      $$ = crtFile.NewField($2, $4, $6)
      $$.SetReference($5);
      $$.SetReq($3);
      if ($7 != nil) {
        $$.SetValue($7)
      }
      $$.SetXsdOptional($8)
      $$.SetXsdNillable($9)
      if ($10 != nil) {
        $$.SetXsdAttrs($10)
      }
      if ($11 != nil) {
        $$.SetAnnotations($11)
      }
      $$.SetComment($13)
    }

FieldIdentifier:
  tok_int_constant ':'
    {
      $$ = $1
    }
|
    {
      $$ = 0
    }

FieldReference:
  tok_reference
    {
      $$ = true
    }
|
   {
     $$ = false
   }

FieldRequiredness:
  tok_required
    {
      $$ = idl.T_REQUIRED
    }
| tok_optional
    {
      if G_arglist > 0 {
        $$ = idl.T_OPT_IN_REQ_OUT
      } else {
        $$ = idl.T_OPTIONAL
      }
    }
|
    {
      $$ = idl.T_OPT_IN_REQ_OUT
    }

FieldValue:
  '=' ConstValue
    {
      if (G_parse_mode == PROGRAM) {
        $$ = $2
      } else {
        $$ = nil
      }
    }
|
    {
      $$ = nil
    }

FunctionType:
  FieldType
    {
      //log.Println("FunctionType -> FieldType", $1)
      $$ = $1
    }
| tok_void
    {
      //log.Println("FunctionType -> tok_void")
      $$ = nil
    }

FieldType:
  tok_identifier
    {
      //log.Println("[yacc]FieldType -> tok_identifier", $1);
      $$ = crtMicroMod.SearchType($1, crtFile)
      if ($$ == nil) {
         $$ = crtMicroMod.AddPlaceholderTypedef($1, crtFile)
      }
    }
| BaseType
    {
      //log.Println("[yacc]FieldType -> BaseType");
      $$ = $1
    }
| ContainerType
    {
      //log.Println("[yacc]FieldType -> ContainerType");
      $$ = $1
    }

BaseType: SimpleBaseType TypeAnnotations
    {
      if ($2 != nil) {
        $$ = $1
        $$.SetAnnotations($2)
      } else {
        $$ = $1
      }
    }

SimpleBaseType:
  tok_string
    {
      $$ = idl.BaseTypeString
    }
| tok_binary
    {
      $$ = idl.BaseTypeBinary
    }
| tok_slist
    {
      $$ = idl.BaseTypeSlist
    }
| tok_bool
    {
      $$ = idl.BaseTypeBool
    }
| tok_i8
    {
      $$ = idl.BaseTypeI8
    }
| tok_i16
    {
      $$ = idl.BaseTypeI16
    }
| tok_i32
    {
      $$ = idl.BaseTypeI32
    }
| tok_i64
    {
      $$ = idl.BaseTypeI64
    }
| tok_double
    {
      $$ = idl.BaseTypeDouble
    }

ContainerType: SimpleContainerType TypeAnnotations
    {
      $$ = $1
      if $2 != nil {
         $$.SetAnnotations($2)
      }
    }

SimpleContainerType:
  MapType
    {
      //log.Println("[yacc]SimpleContainerType -> MapType")
      $$ = $1;
    }
| SetType
    {
      //log.Println("[yacc]SimpleContainerType -> SetType")
      $$ = $1;
    }
| ListType
    {
      //log.Println("[yacc]SimpleContainerType -> ListType")
      $$ = $1
    }

MapType:
  tok_map CppType '<' FieldType ',' FieldType '>'
    {
      //log.Println("[yacc]MapType -> tok_map <FieldType, FieldType>")
      $$ = crtFile.NewMap($4, $6)
      if ($2 != "") {
        $$.SetCppName($2)
      }
    }

SetType:
  tok_set CppType '<' FieldType '>'
    {
      //log.Println("[yacc]SetType -> tok_set<FieldType>")
      $$ = crtFile.NewSet($4)
      if ($2 != "") {
        $$.SetCppName($2)
      }
    }

ListType:
  tok_list '<' FieldType '>' CppType
    {
      $$ = crtFile.NewList($3)
      if ($5 != "") {
        $$.SetCppName($5)
      }
    }

CppType:
  tok_cpp_type tok_literal
    {
      $$ = $2
    }
|
    {
      $$ = ""
    }

TypeAnnotations:
  '(' TypeAnnotationList ')'
    {
      //log.Println("TypeAnno '(' TypeAnnotationList ')'")
      $$ = $2
    }
|
    {
      //log.Println("TypeAnno")
      $$ = nil
    }

TypeAnnotationList:
  TypeAnnotationList TypeAnnotation
    {
      $$ = $1
      $$.Add($2)
    }
|
    {
      //log.Println("TypeAnnotationList")
      $$ = crtFile.NewAnnotations()
    }

TypeAnnotation:
  tok_identifier TypeAnnotationValue CommaOrSemicolonOptional FieldComment
    {
      $$ = idl.NewAnnotation($1, $2, $4)
    }

TypeAnnotationValue:
  '=' tok_literal
    {
      $$ = $2;
    }
|
    {
      $$ = "1";
    }

LobAttributeList:
  LobAttributeList LobAttribute
    {
      $$ = $1
      $$.Append($2)
    }
|
    {
      $$ = crtFile.CreateAttrList()
    }

LobAttribute:
    tok_lobattr
    {
        $$ = crtFile.CreateAttr($1[3:])
    }

FieldComment:
    tok_fieldcomment
    {
        $$ = $1[3:]
        //log.Println($1)
    }
|
    {
        $$ = ""
    }

%%