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
	rc.CollectParameters(op.Parameters)
	if op.RequestBody != nil {
		rc.CollectRequestBodyRef(op.RequestBody)
	}
	rc.CollectResponses(op.Responses)
	rc.CollectCallbacks(op.Callbacks)
}

func (rc *RefsCollector) CollectParameters(params openapi3.Parameters) {
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

func (rc *RefsCollector) collectSchemaRef(sr *openapi3.SchemaRef) {
	if sr == nil {
		return
	}
	if sr.Ref != "" {
		rc.AddRef(sr.Ref)
	}
	if sr.Value != nil {
		rc.collectSchema(sr.Value)
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

func (rc *RefsCollector) collectHeaders(headers openapi3.Headers) {
	for _, header := range headers {
		if header.Ref != "" {
			rc.AddRef(header.Ref)
		}
		if h := header.Value; h != nil {
			rc.collectParameter(&h.Parameter) // Header type embeds the Parameter type
		}
	}
}

func (rc *RefsCollector) CollectRequestBodyRef(rbr *openapi3.RequestBodyRef) {
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

func (rc *RefsCollector) CollectResponses(resps *openapi3.Responses) {
	for _, resp := range resps.Map() {
		if resp.Ref != "" {
			rc.AddRef(resp.Ref)
		}
		if r := resp.Value; r != nil {
			rc.collectHeaders(r.Headers)
			rc.collectContent(r.Content)
			rc.collectLinks(r.Links)
		}
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

func (rc *RefsCollector) CollectCallbacks(callbacks openapi3.Callbacks) {
	for _, callback := range callbacks {
		if callback.Ref != "" {
			rc.AddRef(callback.Ref)
		}
		if c := callback.Value; c != nil {
			for _, path := range c.Map() {
				rc.collectPathItem(path)
			}
		}
	}
}

func (rc *RefsCollector) collectPathItem(path *openapi3.PathItem) {
	if path.Ref != "" {
		rc.AddRef(path.Ref)
	}
	for _, op := range path.Operations() {
		rc.CollectOperation(op)
	}
	rc.CollectParameters(path.Parameters)
}

func (rc *RefsCollector) collectSchemaRefs(s openapi3.SchemaRefs) {
	for _, sr := range s {
		rc.collectSchemaRef(sr)
	}
}

func (rc *RefsCollector) collectSchemas(s openapi3.Schemas) {
	for _, sr := range s {
		rc.collectSchemaRef(sr)
	}
}

func (rc *RefsCollector) collectSchema(s *openapi3.Schema) {
	rc.collectSchemaRefs(s.OneOf)
	rc.collectSchemaRefs(s.AnyOf)
	rc.collectSchemaRefs(s.AllOf)
	rc.collectSchemaRef(s.Not)

	rc.collectSchemaRef(s.Items)

	rc.collectSchemas(s.Properties)
	rc.collectSchemaRef(s.AdditionalProperties.Schema)
}
