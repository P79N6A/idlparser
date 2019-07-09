package idltypes

type IDLAnnotation struct{
	key string
	val string
	comment string
}

func NewAnnotation(k, v string, comment string) *IDLAnnotation{
	return &IDLAnnotation{
		key:k,
		val:v,
		comment: comment,
	}
}

type IDLAnnotations struct{
	vals map[string]*IDLAnnotation
}

func NewIDLAnnotations() *IDLAnnotations{
	return &IDLAnnotations{
		vals: make(map[string]*IDLAnnotation),
	}
}

func (annos *IDLAnnotations)Add(anno *IDLAnnotation){
	if anno != nil {
		annos.vals[anno.key] = anno
	}
}