package filter

import "github.com/getkin/kin-openapi/openapi3"

type RefsCollector struct {
	refs map[string]struct{}
}

func NewRefsCollector() *RefsCollector {
	return &RefsCollector{
		refs: make(map[string]struct{}),
	}
}

func (rc *RefsCollector) AddRef(ref string) {
	rc.refs[ref] = struct{}{}
}

func (rc *RefsCollector) Refs() map[string]struct{} {
	return rc.refs
}

func (rc *RefsCollector) CollectOperation(op *openapi3.Operation) {
	rc.collectParameters(op.Parameters)
	if op.RequestBody != nil {
		rc.collectRequestBodyRef(op.RequestBody)
	}
	rc.collectResponses(op.Responses)
	rc.collectCallbacks(op.Callbacks)
}

func (rc *RefsCollector) collectParameters(params openapi3.Parameters) {
	for _, param := range params {
		if param.Ref != "" {
			rc.AddRef(param.Ref)
		}
		if p := param.Value; p != nil {
			rc.collectParameter(p)
		}
	}
}

func (rc *RefsCollector) collectParameter(param *openapi3.Parameter) {
	rc.collectSchemaRef(param.Schema)
	rc.collectExamples(param.Examples)
	rc.collectContent(param.Content)
}

func (rc *RefsCollector) collectParameterRef(paramr *openapi3.ParameterRef) {
	if paramr.Ref != "" {
		rc.AddRef(paramr.Ref)
	}
	if p := paramr.Value; p != nil {
		rc.collectParameter(p)
	}
}

func (rc *RefsCollector) collectSchemaRef(scr *openapi3.SchemaRef) {
	if scr == nil {
		return
	}
	if scr.Ref != "" {
		rc.AddRef(scr.Ref)
	}
	if scr.Value != nil {
		rc.collectSchema(scr.Value)
	}
}

func (rc *RefsCollector) collectExamples(examples openapi3.Examples) {
	for _, example := range examples {
		if example.Ref != "" {
			rc.AddRef(example.Ref)
		}
	}
}

func (rc *RefsCollector) collectContent(content openapi3.Content) {
	for _, media := range content {
		rc.collectSchemaRef(media.Schema)
		rc.collectExamples(media.Examples)
		for _, enc := range media.Encoding {
			rc.collectEncoding(enc)
		}
	}
}

func (rc *RefsCollector) collectEncoding(enc *openapi3.Encoding) {
	rc.collectHeaders(enc.Headers)
}

func (rc *RefsCollector) collectHeaderRef(hr *openapi3.HeaderRef) {
	if hr.Ref != "" {
		rc.AddRef(hr.Ref)
	}
	if h := hr.Value; h != nil {
		rc.collectParameter(&h.Parameter) // Header type embeds the Parameter type
	}
}

func (rc *RefsCollector) collectHeaders(headers openapi3.Headers) {
	for _, hr := range headers {
		rc.collectHeaderRef(hr)
	}
}

func (rc *RefsCollector) collectRequestBodyRef(rbr *openapi3.RequestBodyRef) {
	if rbr.Ref != "" {
		rc.AddRef(rbr.Ref)
	}
	if rb := rbr.Value; rb != nil {
		rc.collectRequestBodyRefs(rb)
	}
}

func (rc *RefsCollector) collectRequestBodyRefs(rb *openapi3.RequestBody) {
	rc.collectContent(rb.Content)
}

func (rc *RefsCollector) collectResponseRef(respr *openapi3.ResponseRef) {
	if respr.Ref != "" {
		rc.AddRef(respr.Ref)
	}
	if r := respr.Value; r != nil {
		rc.collectHeaders(r.Headers)
		rc.collectContent(r.Content)
		rc.collectLinks(r.Links)
	}
}

func (rc *RefsCollector) collectResponses(resps *openapi3.Responses) {
	for _, respr := range resps.Map() {
		rc.collectResponseRef(respr)
	}
}

func (rc *RefsCollector) collectLinks(links openapi3.Links) {
	for _, link := range links {
		rc.collectLinkRef(link)
	}
}

func (rc *RefsCollector) collectLinkRef(lr *openapi3.LinkRef) {
	if lr.Ref != "" {
		rc.AddRef(lr.Ref)
	}
	// TODO: add OperationRefs handling
}

func (rc *RefsCollector) collectCallbackRef(cbr *openapi3.CallbackRef) {
	if cbr.Ref != "" {
		rc.AddRef(cbr.Ref)
	}
	if c := cbr.Value; c != nil {
		for _, path := range c.Map() {
			rc.collectPathItem(path)
		}
	}
}

func (rc *RefsCollector) collectCallbacks(callbacks openapi3.Callbacks) {
	for _, cbr := range callbacks {
		rc.collectCallbackRef(cbr)
	}
}

func (rc *RefsCollector) collectPathItem(path *openapi3.PathItem) {
	if path.Ref != "" {
		rc.AddRef(path.Ref)
	}
	for _, op := range path.Operations() {
		rc.CollectOperation(op)
	}
	rc.collectParameters(path.Parameters)
}

func (rc *RefsCollector) collectSchemaRefs(scrs openapi3.SchemaRefs) {
	for _, sr := range scrs {
		rc.collectSchemaRef(sr)
	}
}

func (rc *RefsCollector) collectSchemas(scs openapi3.Schemas) {
	for _, sr := range scs {
		rc.collectSchemaRef(sr)
	}
}

func (rc *RefsCollector) collectSchema(sc *openapi3.Schema) {
	rc.collectSchemaRefs(sc.OneOf)
	rc.collectSchemaRefs(sc.AnyOf)
	rc.collectSchemaRefs(sc.AllOf)
	rc.collectSchemaRef(sc.Not)

	rc.collectSchemaRef(sc.Items)

	rc.collectSchemas(sc.Properties)
	rc.collectSchemaRef(sc.AdditionalProperties.Schema)
}

func (rc *RefsCollector) collectSecurityScheme(secsc *openapi3.SecuritySchemeRef) {
	if secsc.Ref != "" {
		rc.AddRef(secsc.Ref)
	}
}

func (rc *RefsCollector) collectExampleRef(exr *openapi3.ExampleRef) {
	if exr.Ref != "" {
		rc.AddRef(exr.Ref)
	}
}
