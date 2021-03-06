// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mixer/v1/config/client/api_spec.proto

package client

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// HTTPAPISpec defines the canonical configuration for generating
// API-related attributes from HTTP requests based on the method and
// uri templated path matches. It is sufficient for defining the API
// surface of a service for the purposes of API attribute
// generation. It is not intended to represent auth, quota,
// documentation, or other information commonly found in other API
// specifications, e.g. OpenAPI.
//
// Existing standards that define operations (or methods) in terms of
// HTTP methods and paths can be normalized to this format for use in
// Istio. For example, a simple petstore API described by OpenAPIv2
// [here](https://github.com/googleapis/gnostic/blob/master/examples/v2.0/yaml/petstore-simple.yaml)
// can be represented with the following HTTPAPISpec.
//
// ```yaml
// apiVersion: config.istio.io/v1alpha2
// kind: HTTPAPISpec
// metadata:
//   name: petstore
//   namespace: default
// spec:
//   attributes:
//     attributes:
//       api.service:
//         stringValue: petstore.swagger.io
//       api.version:
//         stringValue: 1.0.0
//   patterns:
//   - attributes:
//       attributes:
//         api.operation:
//           stringValue: findPets
//     httpMethod: GET
//     uriTemplate: /api/pets
//   - attributes:
//       attributes:
//         api.operation:
//           stringValue: addPet
//     httpMethod: POST
//     uriTemplate: /api/pets
//   - attributes:
//       attributes:
//         api.operation:
//           stringValue: findPetById
//     httpMethod: GET
//     uriTemplate: /api/pets/{id}
//   - attributes:
//       attributes:
//         api.operation:
//           stringValue: deletePet
//     httpMethod: DELETE
//     uriTemplate: /api/pets/{id}
//   apiKeys:
//   - query: api-key
// ```
type HTTPAPISpec struct {
	// List of attributes that are generated when *any* of the HTTP
	// patterns match. This list typically includes the "api.service"
	// and "api.version" attributes.
	Attributes *Attributes `protobuf:"bytes,1,opt,name=attributes,proto3" json:"attributes,omitempty"`
	// List of HTTP patterns to match.
	Patterns []*HTTPAPISpecPattern `protobuf:"bytes,2,rep,name=patterns,proto3" json:"patterns,omitempty"`
	// List of APIKey that describes how to extract an API-KEY from an
	// HTTP request. The first API-Key match found in the list is used,
	// i.e. 'OR' semantics.
	//
	// The following default policies are used to generate the
	// `request.api_key` attribute if no explicit APIKey is defined.
	//
	//     `query: key, `query: api_key`, and then `header: x-api-key`
	//
	ApiKeys              []*APIKey `protobuf:"bytes,3,rep,name=api_keys,json=apiKeys,proto3" json:"api_keys,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *HTTPAPISpec) Reset()         { *m = HTTPAPISpec{} }
func (m *HTTPAPISpec) String() string { return proto.CompactTextString(m) }
func (*HTTPAPISpec) ProtoMessage()    {}
func (*HTTPAPISpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6b15fd2f44b459, []int{0}
}

func (m *HTTPAPISpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPAPISpec.Unmarshal(m, b)
}
func (m *HTTPAPISpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPAPISpec.Marshal(b, m, deterministic)
}
func (m *HTTPAPISpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPAPISpec.Merge(m, src)
}
func (m *HTTPAPISpec) XXX_Size() int {
	return xxx_messageInfo_HTTPAPISpec.Size(m)
}
func (m *HTTPAPISpec) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPAPISpec.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPAPISpec proto.InternalMessageInfo

func (m *HTTPAPISpec) GetAttributes() *Attributes {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *HTTPAPISpec) GetPatterns() []*HTTPAPISpecPattern {
	if m != nil {
		return m.Patterns
	}
	return nil
}

func (m *HTTPAPISpec) GetApiKeys() []*APIKey {
	if m != nil {
		return m.ApiKeys
	}
	return nil
}

// HTTPAPISpecPattern defines a single pattern to match against
// incoming HTTP requests. The per-pattern list of attributes is
// generated if both the http_method and uri_template match. In
// addition, the top-level list of attributes in the HTTPAPISpec is also
// generated.
//
// ```yaml
// pattern:
// - attributes
//     api.operation: doFooBar
//   httpMethod: GET
//   uriTemplate: /foo/bar
// ```
type HTTPAPISpecPattern struct {
	// List of attributes that are generated if the HTTP request matches
	// the specified http_method and uri_template. This typically
	// includes the "api.operation" attribute.
	Attributes *Attributes `protobuf:"bytes,1,opt,name=attributes,proto3" json:"attributes,omitempty"`
	// HTTP request method to match against as defined by
	// [rfc7231](https://tools.ietf.org/html/rfc7231#page-21). For
	// example: GET, HEAD, POST, PUT, DELETE.
	HttpMethod string `protobuf:"bytes,2,opt,name=http_method,json=httpMethod,proto3" json:"http_method,omitempty"`
	// Types that are valid to be assigned to Pattern:
	//	*HTTPAPISpecPattern_UriTemplate
	//	*HTTPAPISpecPattern_Regex
	Pattern              isHTTPAPISpecPattern_Pattern `protobuf_oneof:"pattern"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *HTTPAPISpecPattern) Reset()         { *m = HTTPAPISpecPattern{} }
func (m *HTTPAPISpecPattern) String() string { return proto.CompactTextString(m) }
func (*HTTPAPISpecPattern) ProtoMessage()    {}
func (*HTTPAPISpecPattern) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6b15fd2f44b459, []int{1}
}

func (m *HTTPAPISpecPattern) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPAPISpecPattern.Unmarshal(m, b)
}
func (m *HTTPAPISpecPattern) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPAPISpecPattern.Marshal(b, m, deterministic)
}
func (m *HTTPAPISpecPattern) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPAPISpecPattern.Merge(m, src)
}
func (m *HTTPAPISpecPattern) XXX_Size() int {
	return xxx_messageInfo_HTTPAPISpecPattern.Size(m)
}
func (m *HTTPAPISpecPattern) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPAPISpecPattern.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPAPISpecPattern proto.InternalMessageInfo

func (m *HTTPAPISpecPattern) GetAttributes() *Attributes {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *HTTPAPISpecPattern) GetHttpMethod() string {
	if m != nil {
		return m.HttpMethod
	}
	return ""
}

type isHTTPAPISpecPattern_Pattern interface {
	isHTTPAPISpecPattern_Pattern()
}

type HTTPAPISpecPattern_UriTemplate struct {
	UriTemplate string `protobuf:"bytes,3,opt,name=uri_template,json=uriTemplate,proto3,oneof"`
}

type HTTPAPISpecPattern_Regex struct {
	Regex string `protobuf:"bytes,4,opt,name=regex,proto3,oneof"`
}

func (*HTTPAPISpecPattern_UriTemplate) isHTTPAPISpecPattern_Pattern() {}

func (*HTTPAPISpecPattern_Regex) isHTTPAPISpecPattern_Pattern() {}

func (m *HTTPAPISpecPattern) GetPattern() isHTTPAPISpecPattern_Pattern {
	if m != nil {
		return m.Pattern
	}
	return nil
}

func (m *HTTPAPISpecPattern) GetUriTemplate() string {
	if x, ok := m.GetPattern().(*HTTPAPISpecPattern_UriTemplate); ok {
		return x.UriTemplate
	}
	return ""
}

func (m *HTTPAPISpecPattern) GetRegex() string {
	if x, ok := m.GetPattern().(*HTTPAPISpecPattern_Regex); ok {
		return x.Regex
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*HTTPAPISpecPattern) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*HTTPAPISpecPattern_UriTemplate)(nil),
		(*HTTPAPISpecPattern_Regex)(nil),
	}
}

// APIKey defines the explicit configuration for generating the
// `request.api_key` attribute from HTTP requests.
//
// See [API Keys](https://swagger.io/docs/specification/authentication/api-keys)
// for a general overview of API keys as defined by OpenAPI.
type APIKey struct {
	// Types that are valid to be assigned to Key:
	//	*APIKey_Query
	//	*APIKey_Header
	//	*APIKey_Cookie
	Key                  isAPIKey_Key `protobuf_oneof:"key"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *APIKey) Reset()         { *m = APIKey{} }
func (m *APIKey) String() string { return proto.CompactTextString(m) }
func (*APIKey) ProtoMessage()    {}
func (*APIKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6b15fd2f44b459, []int{2}
}

func (m *APIKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_APIKey.Unmarshal(m, b)
}
func (m *APIKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_APIKey.Marshal(b, m, deterministic)
}
func (m *APIKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_APIKey.Merge(m, src)
}
func (m *APIKey) XXX_Size() int {
	return xxx_messageInfo_APIKey.Size(m)
}
func (m *APIKey) XXX_DiscardUnknown() {
	xxx_messageInfo_APIKey.DiscardUnknown(m)
}

var xxx_messageInfo_APIKey proto.InternalMessageInfo

type isAPIKey_Key interface {
	isAPIKey_Key()
}

type APIKey_Query struct {
	Query string `protobuf:"bytes,1,opt,name=query,proto3,oneof"`
}

type APIKey_Header struct {
	Header string `protobuf:"bytes,2,opt,name=header,proto3,oneof"`
}

type APIKey_Cookie struct {
	Cookie string `protobuf:"bytes,3,opt,name=cookie,proto3,oneof"`
}

func (*APIKey_Query) isAPIKey_Key() {}

func (*APIKey_Header) isAPIKey_Key() {}

func (*APIKey_Cookie) isAPIKey_Key() {}

func (m *APIKey) GetKey() isAPIKey_Key {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *APIKey) GetQuery() string {
	if x, ok := m.GetKey().(*APIKey_Query); ok {
		return x.Query
	}
	return ""
}

func (m *APIKey) GetHeader() string {
	if x, ok := m.GetKey().(*APIKey_Header); ok {
		return x.Header
	}
	return ""
}

func (m *APIKey) GetCookie() string {
	if x, ok := m.GetKey().(*APIKey_Cookie); ok {
		return x.Cookie
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*APIKey) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*APIKey_Query)(nil),
		(*APIKey_Header)(nil),
		(*APIKey_Cookie)(nil),
	}
}

// HTTPAPISpecReference defines a reference to an HTTPAPISpec. This is
// typically used for establishing bindings between an HTTPAPISpec and an
// IstioService. For example, the following defines an
// HTTPAPISpecReference for service `foo` in namespace `bar`.
//
// ```yaml
// - name: foo
//   namespace: bar
// ```
type HTTPAPISpecReference struct {
	// REQUIRED. The short name of the HTTPAPISpec. This is the resource
	// name defined by the metadata name field.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Optional namespace of the HTTPAPISpec. Defaults to the encompassing
	// HTTPAPISpecBinding's metadata namespace field.
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HTTPAPISpecReference) Reset()         { *m = HTTPAPISpecReference{} }
func (m *HTTPAPISpecReference) String() string { return proto.CompactTextString(m) }
func (*HTTPAPISpecReference) ProtoMessage()    {}
func (*HTTPAPISpecReference) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6b15fd2f44b459, []int{3}
}

func (m *HTTPAPISpecReference) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPAPISpecReference.Unmarshal(m, b)
}
func (m *HTTPAPISpecReference) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPAPISpecReference.Marshal(b, m, deterministic)
}
func (m *HTTPAPISpecReference) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPAPISpecReference.Merge(m, src)
}
func (m *HTTPAPISpecReference) XXX_Size() int {
	return xxx_messageInfo_HTTPAPISpecReference.Size(m)
}
func (m *HTTPAPISpecReference) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPAPISpecReference.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPAPISpecReference proto.InternalMessageInfo

func (m *HTTPAPISpecReference) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *HTTPAPISpecReference) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

// HTTPAPISpecBinding defines the binding between HTTPAPISpecs and one or more
// IstioService. For example, the following establishes a binding
// between the HTTPAPISpec `petstore` and service `foo` in namespace `bar`.
//
// ```yaml
// apiVersion: config.istio.io/v1alpha2
// kind: HTTPAPISpecBinding
// metadata:
//   name: my-binding
//   namespace: default
// spec:
//   services:
//   - name: foo
//     namespace: bar
//   apiSpecs:
//   - name: petstore
//     namespace: default
// ```
type HTTPAPISpecBinding struct {
	// REQUIRED. One or more services to map the listed HTTPAPISpec onto.
	Services []*IstioService `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	// REQUIRED. One or more HTTPAPISpec references that should be mapped to
	// the specified service(s). The aggregate collection of match
	// conditions defined in the HTTPAPISpecs should not overlap.
	ApiSpecs             []*HTTPAPISpecReference `protobuf:"bytes,2,rep,name=api_specs,json=apiSpecs,proto3" json:"api_specs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *HTTPAPISpecBinding) Reset()         { *m = HTTPAPISpecBinding{} }
func (m *HTTPAPISpecBinding) String() string { return proto.CompactTextString(m) }
func (*HTTPAPISpecBinding) ProtoMessage()    {}
func (*HTTPAPISpecBinding) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb6b15fd2f44b459, []int{4}
}

func (m *HTTPAPISpecBinding) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPAPISpecBinding.Unmarshal(m, b)
}
func (m *HTTPAPISpecBinding) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPAPISpecBinding.Marshal(b, m, deterministic)
}
func (m *HTTPAPISpecBinding) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPAPISpecBinding.Merge(m, src)
}
func (m *HTTPAPISpecBinding) XXX_Size() int {
	return xxx_messageInfo_HTTPAPISpecBinding.Size(m)
}
func (m *HTTPAPISpecBinding) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPAPISpecBinding.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPAPISpecBinding proto.InternalMessageInfo

func (m *HTTPAPISpecBinding) GetServices() []*IstioService {
	if m != nil {
		return m.Services
	}
	return nil
}

func (m *HTTPAPISpecBinding) GetApiSpecs() []*HTTPAPISpecReference {
	if m != nil {
		return m.ApiSpecs
	}
	return nil
}

func init() {
	proto.RegisterType((*HTTPAPISpec)(nil), "istio.mixer.v1.config.client.HTTPAPISpec")
	proto.RegisterType((*HTTPAPISpecPattern)(nil), "istio.mixer.v1.config.client.HTTPAPISpecPattern")
	proto.RegisterType((*APIKey)(nil), "istio.mixer.v1.config.client.APIKey")
	proto.RegisterType((*HTTPAPISpecReference)(nil), "istio.mixer.v1.config.client.HTTPAPISpecReference")
	proto.RegisterType((*HTTPAPISpecBinding)(nil), "istio.mixer.v1.config.client.HTTPAPISpecBinding")
}

func init() {
	proto.RegisterFile("mixer/v1/config/client/api_spec.proto", fileDescriptor_fb6b15fd2f44b459)
}

var fileDescriptor_fb6b15fd2f44b459 = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0xb6, 0x9b, 0x36, 0x4d, 0xc6, 0xbf, 0xd3, 0xaa, 0xfa, 0xc9, 0x44, 0x55, 0x88, 0xdc, 0x22,
	0x45, 0x1c, 0x6c, 0x1a, 0xc4, 0x85, 0x0b, 0x4a, 0x0e, 0x28, 0x51, 0x41, 0x44, 0xdb, 0x9c, 0xe0,
	0x10, 0xb9, 0xce, 0xd4, 0x59, 0xa5, 0xf1, 0x9a, 0xdd, 0x4d, 0x54, 0xbf, 0x11, 0x17, 0x1e, 0x80,
	0x37, 0xe0, 0x11, 0x7a, 0x6d, 0x9f, 0x80, 0x47, 0x40, 0xbb, 0xeb, 0x3a, 0xe1, 0x5f, 0x84, 0xc4,
	0xc9, 0x9e, 0x99, 0x6f, 0xbe, 0x99, 0x6f, 0x76, 0x06, 0x9e, 0x2c, 0xd9, 0x0d, 0x8a, 0x68, 0x7d,
	0x16, 0x25, 0x3c, 0xbb, 0x62, 0x69, 0x94, 0x5c, 0x33, 0xcc, 0x54, 0x14, 0xe7, 0x6c, 0x2a, 0x73,
	0x4c, 0xc2, 0x5c, 0x70, 0xc5, 0xc9, 0x31, 0x93, 0x8a, 0xf1, 0xd0, 0x80, 0xc3, 0xf5, 0x59, 0x68,
	0xc1, 0xa1, 0x05, 0xb7, 0x8e, 0x52, 0x9e, 0x72, 0x03, 0x8c, 0xf4, 0x9f, 0xcd, 0x69, 0x3d, 0xaa,
	0xa8, 0x63, 0xa5, 0x04, 0xbb, 0x5c, 0x29, 0x94, 0x65, 0xe8, 0xf4, 0x0f, 0x55, 0x25, 0x8a, 0x35,
	0x4b, 0xd0, 0xa2, 0x82, 0x5b, 0x17, 0xbc, 0xe1, 0x64, 0x32, 0xee, 0x8f, 0x47, 0x17, 0x39, 0x26,
	0xe4, 0x25, 0xc0, 0x86, 0xc9, 0x77, 0x3b, 0x6e, 0xd7, 0xeb, 0xb5, 0xc2, 0x9f, 0x3a, 0xeb, 0x57,
	0x08, 0xba, 0x85, 0x26, 0x6f, 0xa0, 0x91, 0xc7, 0x4a, 0xa1, 0xc8, 0xa4, 0xbf, 0xd7, 0xa9, 0x75,
	0xbd, 0xde, 0xb3, 0x70, 0x97, 0xa6, 0x70, 0xab, 0xf0, 0xd8, 0x26, 0xd2, 0x8a, 0x81, 0xbc, 0x82,
	0x86, 0x1e, 0xd0, 0x02, 0x0b, 0xe9, 0xd7, 0x0c, 0xdb, 0xe9, 0x6e, 0xb6, 0xfe, 0x78, 0x74, 0x8e,
	0x05, 0x3d, 0x8c, 0x73, 0x76, 0x8e, 0x85, 0x0c, 0xbe, 0xb8, 0x40, 0x7e, 0xad, 0xf0, 0x4f, 0x0a,
	0x1f, 0x83, 0x37, 0x57, 0x2a, 0x9f, 0x2e, 0x51, 0xcd, 0xf9, 0xcc, 0xdf, 0xeb, 0xb8, 0xdd, 0x26,
	0x05, 0xed, 0x7a, 0x6b, 0x3c, 0xe4, 0x04, 0xfe, 0x5b, 0x09, 0x36, 0x55, 0xb8, 0xcc, 0xaf, 0x63,
	0x85, 0x7e, 0x4d, 0x23, 0x86, 0x0e, 0xf5, 0x56, 0x82, 0x4d, 0x4a, 0x27, 0xf9, 0x1f, 0x0e, 0x04,
	0xa6, 0x78, 0xe3, 0xef, 0x97, 0x51, 0x6b, 0x0e, 0x9a, 0x70, 0x58, 0xaa, 0x0f, 0x3e, 0x40, 0xdd,
	0xca, 0xd1, 0xe0, 0x8f, 0x2b, 0x14, 0x85, 0xe9, 0xd4, 0x80, 0x8d, 0x49, 0x7c, 0xa8, 0xcf, 0x31,
	0x9e, 0xa1, 0xb0, 0x5d, 0x0c, 0x1d, 0x5a, 0xda, 0x3a, 0x92, 0x70, 0xbe, 0x60, 0x9b, 0xea, 0xa5,
	0x3d, 0x38, 0x80, 0xda, 0x02, 0x8b, 0x60, 0x08, 0x47, 0x5b, 0x73, 0xa1, 0x78, 0x85, 0x02, 0xb3,
	0x04, 0x09, 0x81, 0xfd, 0x2c, 0x5e, 0xa2, 0xad, 0x44, 0xcd, 0x3f, 0x39, 0x86, 0xa6, 0xfe, 0xca,
	0x3c, 0x4e, 0xb0, 0xd4, 0xbb, 0x71, 0x04, 0x9f, 0x7f, 0x1c, 0xf1, 0x80, 0x65, 0x33, 0x96, 0xa5,
	0xe4, 0x35, 0x34, 0xca, 0x2d, 0xd3, 0x03, 0xd6, 0x4f, 0xf7, 0x74, 0xf7, 0xd3, 0x8d, 0x74, 0xf0,
	0xc2, 0xa6, 0xd0, 0x2a, 0x97, 0xbc, 0x83, 0xe6, 0xc3, 0x8d, 0x3c, 0x6c, 0x54, 0xef, 0xaf, 0x37,
	0xaa, 0xd2, 0x45, 0xf5, 0x1e, 0x69, 0x8f, 0x1c, 0xbc, 0xf8, 0x7a, 0xd7, 0x76, 0xbe, 0xdd, 0xb5,
	0x9d, 0x4f, 0xf7, 0x6d, 0xe7, 0xf6, 0xbe, 0xed, 0xbe, 0x3f, 0xb1, 0x74, 0x8c, 0xeb, 0x6b, 0x8c,
	0x7e, 0x7f, 0x32, 0x97, 0x75, 0x73, 0x2b, 0xcf, 0xbf, 0x07, 0x00, 0x00, 0xff, 0xff, 0x4c, 0x93,
	0xf6, 0x84, 0xc9, 0x03, 0x00, 0x00,
}
