// Package docs 在运行时构造 OpenAPI 3.0 规范，对外提供 /swagger/doc.json。
//
// 设计上不在 handlers 内部塞反射，而采用"路由注册处旁路调用 Register(...)"
// 的方式集中维护。比手工维护 openapi.json 文件更安全（一次只改一处，避免漏标）。
package docs

import (
	"encoding/json"
	"sync"
)

// ============================================================
// OpenAPI 3.0 基础结构
// ============================================================

// Spec OpenAPI 3.0 文档主对象。
type Spec struct {
	OpenAPI    string           `json:"openapi"`
	Info       Info             `json:"info"`
	Servers    []Server         `json:"servers"`
	Tags       []Tag            `json:"tags"`
	Components *Components      `json:"components"`
	Paths      map[string]*Path `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

type Server struct {
	URL string `json:"url"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Components struct {
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
	Schemas          map[string]*Schema         `json:"schemas,omitempty"`
}

type SecurityScheme struct {
	Type        string `json:"type"`
	Scheme      string `json:"scheme,omitempty"`
	Description string `json:"description,omitempty"`
}

// Schema 是 OpenAPI schema 的最小实现（覆盖普通对象/引用/数组/基本类型）。
type Schema struct {
	Ref                  string             `json:"$ref,omitempty"`
	Type                 string             `json:"type,omitempty"`
	Format               string             `json:"format,omitempty"`
	Description          string             `json:"description,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	Enum                 []any             `json:"enum,omitempty"`
	AdditionalProperties any               `json:"additionalProperties,omitempty"`
}

// Path 一个 URL 路径下若干 HTTP 方法的集合。
type Path struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

// Operation 单个 HTTP 方法的元信息。
type Operation struct {
	Tags        []string             `json:"tags,omitempty"`
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	Security    []map[string]any     `json:"security,omitempty"`
	Parameters  []*Parameter         `json:"parameters,omitempty"`
	RequestBody *RefBody             `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // "query" / "header" / "path"
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type RefBody struct {
	Required bool                  `json:"required,omitempty"`
	Content  map[string]*MediaType `json:"content"`
}

type MediaType struct {
	Schema *Schema `json:"schema,omitempty"`
}

type Response struct {
	Description string                `json:"description"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

// ============================================================
// 注册中心
// ============================================================

var (
	once sync.Once
	spec *Spec
)

// Get 返回当前进程的 OpenAPI 文档（首次访问时构造）。
func Get() *Spec {
	once.Do(build)
	return spec
}

// Register 注册一个端点。
//
// 多次调用按 path 聚合到同一个 Path 对象；同一 path+method 重复注册 panic。
func Register(method, path string, op *Operation) {
	once.Do(build) // 保证 schemas 已初始化
	p, ok := spec.Paths[path]
	if !ok {
		p = &Path{}
		spec.Paths[path] = p
	}
	switch method {
	case "GET":
		if p.Get != nil {
			panic("docs: duplicate GET " + path)
		}
		p.Get = op
	case "POST":
		if p.Post != nil {
			panic("docs: duplicate POST " + path)
		}
		p.Post = op
	case "PUT":
		if p.Put != nil {
			panic("docs: duplicate PUT " + path)
		}
		p.Put = op
	case "DELETE":
		if p.Delete != nil {
			panic("docs: duplicate DELETE " + path)
		}
		p.Delete = op
	default:
		panic("docs: unsupported method " + method)
	}
}

// MustJSON 将 spec 序列化为 JSON。
func MustJSON() []byte {
	b, err := json.MarshalIndent(Get(), "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}

// ============================================================
// Schema 注册助手（用于声明可复用对象）
// ============================================================

// AddSchema 注册/更新一个可复用 schema（key 必须全局唯一）。
func AddSchema(name string, s *Schema) {
	if spec.Components == nil {
		spec.Components = &Components{}
	}
	if spec.Components.Schemas == nil {
		spec.Components.Schemas = map[string]*Schema{}
	}
	spec.Components.Schemas[name] = s
}

// Ref 创建一个引用形式的 schema。
func Ref(name string) *Schema { return &Schema{Ref: "#/components/schemas/" + name} }

// Obj 创建一个 object schema。
func Obj(props map[string]*Schema, required ...string) *Schema {
	return &Schema{Type: "object", Properties: props, Required: required}
}

// Str 创建一个 string schema。
func Str(desc, format string, enum ...any) *Schema {
	s := &Schema{Type: "string", Description: desc}
	if format != "" {
		s.Format = format
	}
	if len(enum) > 0 {
		s.Enum = enum
	}
	return s
}

// Int 创建一个 integer schema。
func Int(desc, format string) *Schema {
	s := &Schema{Type: "integer", Description: desc}
	if format != "" {
		s.Format = format
	}
	return s
}

// Arr 创建一个 array schema。
func Arr(item *Schema, desc string) *Schema {
	return &Schema{Type: "array", Description: desc, Items: item}
}

// ============================================================
// Operation 构造助手
// ============================================================

// Op 构造一个 Operation，传入 callback 自定义响应/参数等。
func Op(summary string, tags ...string) *Operation {
	return &Operation{
		Summary:   summary,
		Tags:      tags,
		Responses: map[string]*Response{},
	}
}

// Sec 加一个安全要求。
func (o *Operation) Sec(name string) *Operation {
	o.Security = append(o.Security, map[string]any{name: []any{}})
	return o
}

// Param 加一个 query/header/path 参数。
func (o *Operation) Param(name, in, desc string, required bool, schema *Schema) *Operation {
	o.Parameters = append(o.Parameters, &Parameter{
		Name:        name,
		In:          in,
		Description: desc,
		Required:    required,
		Schema:      schema,
	})
	return o
}

// JSONBody 声明 JSON 请求体。
func (o *Operation) JSONBody(schema *Schema, required bool, mediaDescs ...string) *Operation {
	desc := "application/json"
	if len(mediaDescs) > 0 {
		desc = mediaDescs[0]
	}
	o.RequestBody = &RefBody{
		Required: required,
		Content:  map[string]*MediaType{desc: {Schema: schema}},
	}
	return o
}

// Res 声明一个响应。
func (o *Operation) Res(status, description string, schema *Schema) *Operation {
	r := &Response{Description: description}
	if schema != nil {
		r.Content = map[string]*MediaType{"application/json": {Schema: schema}}
	}
	o.Responses[status] = r
	return o
}

// ============================================================
// 可复用 Schema 定义
// ============================================================

var (
	// User 用户模型
	User = Obj(map[string]*Schema{
		"id":    Int("用户ID", "int64"),
		"name":  Str("用户名", ""),
		"email": Str("邮箱", "email"),
		"role":  Str("角色", "", "admin", "user"),
	}, "id", "name", "email")

	// UserCreate 创建用户请求体
	UserCreate = Obj(map[string]*Schema{
		"name":  Str("用户名", ""),
		"email": Str("邮箱", "email"),
		"role":  Str("角色", "", "admin", "user"),
	}, "name", "email")

	// Message 通用消息响应
	Message = Obj(map[string]*Schema{
		"message": Str("提示信息", ""),
	})
)

// ============================================================
// 一次性构建：注册所有 schema + 全部端点
// ============================================================

func build() {
	spec = &Spec{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       "Go 标准库 API",
			Version:     "1.0.0",
			Description: "Go 原生 net/http + Knife4j Vue3 示例，无需任何第三方包",
		},
		Servers: []Server{{URL: "/"}},
		Tags: []Tag{
			{Name: "用户管理", Description: "用户 CRUD 操作"},
			{Name: "系统", Description: "系统相关接口"},
		},
		Paths: map[string]*Path{},
	}

	// 注册可复用 schemas
	AddSchema("User", User)
	AddSchema("UserCreate", UserCreate)
	AddSchema("Message", Message)

	// 端点清单（与 main.go 中的路由注册保持同步）
	registerEndpoints()
}

// registerEndpoints 集中声明所有 HTTP 端点的 OpenAPI 元信息。
func registerEndpoints() {

	// -------- 用户管理 --------

	Register("GET", "/api/users",
		Op("获取用户列表", "用户管理").
			Res("200", "成功", Arr(Ref("User"), "")).
			Res("500", "服务端错误", Ref("Message")),
	)

	Register("GET", "/api/users/{id}",
		Op("根据ID获取用户", "用户管理").
			Param("id", "path", "用户ID", true, Int("", "int64")).
			Res("200", "成功", Ref("User")).
			Res("404", "用户不存在", Ref("Message")),
	)

	Register("POST", "/api/users",
		Op("创建用户", "用户管理").
			JSONBody(Ref("UserCreate"), true).
			Res("200", "创建成功", Ref("User")).
			Res("400", "参数错误", Ref("Message")),
	)

	Register("PUT", "/api/users/{id}",
		Op("更新用户", "用户管理").
			Param("id", "path", "用户ID", true, Int("", "int64")).
			JSONBody(Ref("UserCreate"), true).
			Res("200", "更新成功", Ref("User")).
			Res("404", "用户不存在", Ref("Message")),
	)

	Register("DELETE", "/api/users/{id}",
		Op("删除用户", "用户管理").
			Param("id", "path", "用户ID", true, Int("", "int64")).
			Res("200", "删除成功", Ref("Message")).
			Res("404", "用户不存在", Ref("Message")),
	)

	// -------- 系统 --------

	Register("GET", "/api/health",
		Op("健康检查", "系统").
			Res("200", "成功", Ref("Message")),
	)
}
