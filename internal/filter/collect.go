package filter

import "github.com/getkin/kin-openapi/openapi3"


func (o *OpenAPISpecFilter) collectOperationRefs(op *openapi3.Operation, refs map[string]struct{}) {
	for _, param := range op.Parameters {
		if param.Ref != "" {
			refs[param.Ref] = struct{}{}
		}
	}

	if op.RequestBody != nil {
		if op.RequestBody.Ref != "" {
			refs[op.RequestBody.Ref] = struct{}{}
		}

		if rb := op.RequestBody.Value; rb != nil {
			o.collectRequestBodyRefs(rb, refs)
		}
	}

	for _, resp := range op.Responses.Map() {
		if resp.Ref != "" {
			refs[resp.Ref] = struct{}{}
		}
		if r := resp.Value; r != nil {
			for _, media := range r.Content {
				o.collectSchemaRefRecursively(media.Schema, refs)
			}
		}
	}

	// TODO: add callbacks collection
}

func (o *OpenAPISpecFilter) collectRequestBodyRefs(rb *openapi3.RequestBody, refs map[string]struct{}) {
	for _, media := range rb.Content {
		o.collectSchemaRefRecursively(media.Schema, refs)
	}
}

func (o *OpenAPISpecFilter) collectSchemaRefRecursively(s *openapi3.SchemaRef, refs map[string]struct{}) {
	if s == nil {
		return
	}

	if s.Ref != "" {
		refs[s.Ref] = struct{}{}
	}

	if s.Value != nil {
		o.collectSchemaRecursively(s.Value, refs)
	}
}

func (o *OpenAPISpecFilter) collectSchemaRefsRecursively(s openapi3.SchemaRefs, refs map[string]struct{}) {
	for _, sr := range s {
		o.collectSchemaRefRecursively(sr, refs)
	}
}

func (o *OpenAPISpecFilter) collectSchemasRecursively(s openapi3.Schemas, refs map[string]struct{}) {
	for _, sr := range s {
		o.collectSchemaRefRecursively(sr, refs)
	}
}

func (o *OpenAPISpecFilter) collectSchemaRecursively(s *openapi3.Schema, refs map[string]struct{}) {
	o.collectSchemaRefsRecursively(s.OneOf, refs)
	o.collectSchemaRefsRecursively(s.AnyOf, refs)
	o.collectSchemaRefsRecursively(s.AllOf, refs)
	o.collectSchemaRefRecursively(s.Not, refs)
	
	o.collectSchemaRefRecursively(s.Items, refs)

	o.collectSchemasRecursively(s.Properties, refs)
	o.collectSchemaRefRecursively(s.AdditionalProperties.Schema, refs)	
}
