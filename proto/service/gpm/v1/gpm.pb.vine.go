// Code generated by proto-gen-vine. DO NOT EDIT.
// source: github.com/gpm2/gpm/proto/service/gpm/v1/gpm.proto

package gpmv1

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gpm2/gpm/proto/apis/gpm/v1"
	client "github.com/lack-io/vine/core/client"
	server "github.com/lack-io/vine/core/server"
	api "github.com/lack-io/vine/lib/api"
	apipb "github.com/lack-io/vine/proto/apis/api"
	openapi "github.com/lack-io/vine/proto/apis/openapi"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// API Endpoints for GpmService service
func NewGpmServiceEndpoints() []*apipb.Endpoint {
	return []*apipb.Endpoint{
		&apipb.Endpoint{
			Name:        "GpmService.Healthz",
			Description: "GpmService.Healthz",
			Path:        []string{"/api/healthz"},
			Method:      []string{"GET"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.ListService",
			Description: "GpmService.ListService",
			Path:        []string{"/api/v1/Service/"},
			Method:      []string{"GET"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.GetService",
			Description: "GpmService.GetService",
			Path:        []string{"/api/v1/Service/{id}"},
			Method:      []string{"GET"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.GetServiceByName",
			Description: "GpmService.GetServiceByName",
			Path:        []string{"/api/v1/Service/{name}"},
			Method:      []string{"GET"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.CreateService",
			Description: "GpmService.CreateService",
			Path:        []string{"/api/v1/Service/"},
			Method:      []string{"POST"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.StartService",
			Description: "GpmService.StartService",
			Path:        []string{"/api/v1/Service/{id}/action/start"},
			Method:      []string{"PATCH"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.StopService",
			Description: "GpmService.StopService",
			Path:        []string{"/api/v1/Service/{id}/action/stop"},
			Method:      []string{"PATCH"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.RebootService",
			Description: "GpmService.RebootService",
			Path:        []string{"/api/v1/Service/{id}/action/reboot"},
			Method:      []string{"PATCH"},
			Body:        "*",
			Handler:     "rpc",
		},
		&apipb.Endpoint{
			Name:        "GpmService.DeleteService",
			Description: "GpmService.DeleteService",
			Path:        []string{"/api/v1/Service/{id}"},
			Method:      []string{"DELETE"},
			Body:        "*",
			Handler:     "rpc",
		},
	}
}

// Swagger OpenAPI 3.0 for GpmService service
func NewGpmServiceOpenAPI() *openapi.OpenAPI {
	return &openapi.OpenAPI{
		Openapi: "3.0.1",
		Info: &openapi.OpenAPIInfo{
			Title:       "GpmServiceService",
			Description: "OpenAPI3.0 for GpmService",
			Version:     "v1.0.0",
		},
		Servers: []*openapi.OpenAPIServer{},
		Tags: []*openapi.OpenAPITag{
			&openapi.OpenAPITag{
				Name:        "GpmService",
				Description: "OpenAPI3.0 for GpmService",
			},
		},
		Paths: map[string]*openapi.OpenAPIPath{
			"/api/v1/Service/{id}/action/start": &openapi.OpenAPIPath{
				Patch: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "启动服务",
					Description: "GpmService StartService",
					OperationId: "GpmServiceStartService",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "id",
							In:          "path",
							Description: "StartServiceReq field id",
							Required:    true,
							Explode:     true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					RequestBody: &openapi.PathRequestBody{
						Description: "StartService StartServiceReq",
						Content: &openapi.PathRequestBodyContent{
							ApplicationJson: &openapi.ApplicationContent{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/v1.StartServiceReq",
								},
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.StartServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
			"/api/v1/Service/{id}/action/stop": &openapi.OpenAPIPath{
				Patch: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "停止服务",
					Description: "GpmService StopService",
					OperationId: "GpmServiceStopService",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "id",
							In:          "path",
							Description: "StopServiceReq field id",
							Required:    true,
							Explode:     true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					RequestBody: &openapi.PathRequestBody{
						Description: "StopService StopServiceReq",
						Content: &openapi.PathRequestBodyContent{
							ApplicationJson: &openapi.ApplicationContent{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/v1.StopServiceReq",
								},
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.StopServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
			"/api/v1/Service/{id}/action/reboot": &openapi.OpenAPIPath{
				Patch: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "重启服务",
					Description: "GpmService RebootService",
					OperationId: "GpmServiceRebootService",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "id",
							In:          "path",
							Description: "RebootServiceReq field id",
							Required:    true,
							Explode:     true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					RequestBody: &openapi.PathRequestBody{
						Description: "RebootService RebootServiceReq",
						Content: &openapi.PathRequestBodyContent{
							ApplicationJson: &openapi.ApplicationContent{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/v1.RebootServiceReq",
								},
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.RebootServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
			"/api/healthz": &openapi.OpenAPIPath{
				Get: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Description: "GpmService Healthz",
					OperationId: "GpmServiceHealthz",
					Parameters:  []*openapi.PathParameters{},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.Empty"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
			"/api/v1/Service/": &openapi.OpenAPIPath{
				Get: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "查询所有服务",
					Description: "GpmService ListService",
					OperationId: "GpmServiceListService",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "page",
							In:          "query",
							Description: "页码",
							Style:       "form",
							Explode:     true,
							Schema: &openapi.Schema{
								Type:             "integer",
								Format:           "int32",
								ExclusiveMinimum: true,
								Minimum:          0,
								Default:          "1",
							},
						},
						&openapi.PathParameters{
							Name:        "size",
							In:          "query",
							Description: "每页数量",
							Style:       "form",
							Explode:     true,
							Schema: &openapi.Schema{
								Type:             "integer",
								Format:           "int32",
								ExclusiveMinimum: true,
								Minimum:          0,
								Default:          "10",
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.ListServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
				Post: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "新建服务",
					Description: "GpmService CreateService",
					OperationId: "GpmServiceCreateService",
					RequestBody: &openapi.PathRequestBody{
						Description: "CreateService CreateServiceReq",
						Content: &openapi.PathRequestBodyContent{
							ApplicationJson: &openapi.ApplicationContent{
								Schema: &openapi.Schema{
									Ref: "#/components/schemas/v1.CreateServiceReq",
								},
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.CreateServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
			"/api/v1/Service/{id}": &openapi.OpenAPIPath{
				Get: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "查询单个服务",
					Description: "GpmService GetService",
					OperationId: "GpmServiceGetService",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "id",
							In:          "path",
							Description: "服务 id",
							Required:    true,
							Explode:     true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.GetServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
				Delete: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "删除服务",
					Description: "GpmService DeleteService",
					OperationId: "GpmServiceDeleteService",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "id",
							In:          "path",
							Description: "DeleteServiceReq field id",
							Required:    true,
							Explode:     true,
							Schema: &openapi.Schema{
								Type:   "integer",
								Format: "int64",
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.DeleteServiceRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
			"/api/v1/Service/{name}": &openapi.OpenAPIPath{
				Get: &openapi.OpenAPIPathDocs{
					Tags:        []string{"GpmService"},
					Summary:     "通过名称查询单个服务",
					Description: "GpmService GetServiceByName",
					OperationId: "GpmServiceGetServiceByName",
					Parameters: []*openapi.PathParameters{
						&openapi.PathParameters{
							Name:        "name",
							In:          "path",
							Description: "服务 id",
							Required:    true,
							Explode:     true,
							Schema: &openapi.Schema{
								Type: "string",
							},
						},
					},
					Responses: map[string]*openapi.PathResponse{
						"200": &openapi.PathResponse{
							Description: "successful response (stream response)",
							Content: &openapi.PathRequestBodyContent{
								ApplicationJson: &openapi.ApplicationContent{
									Schema: &openapi.Schema{Ref: "#/components/schemas/v1.GetServiceByNameRsp"},
								},
							},
						},
					},
					Security: []*openapi.PathSecurity{},
				},
			},
		},
		Components: &openapi.OpenAPIComponents{
			SecuritySchemes: &openapi.SecuritySchemes{},
			Schemas: map[string]*openapi.Model{
				"v1.StartServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
					},
					Required: []string{"id"},
				},
				"v1.StartServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.StopServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
					},
					Required: []string{"id"},
				},
				"v1.StopServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.RebootServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
					},
					Required: []string{"id"},
				},
				"v1.RebootServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.Empty": &openapi.Model{
					Type:       "object",
					Properties: map[string]*openapi.Schema{},
				},
				"v1.ListServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"meta": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.PageMeta",
						},
					},
				},
				"v1.ListServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"services": &openapi.Schema{
							Type:  "array",
							Items: &openapi.Schema{Ref: "#/components/schemas/v1.Service"},
						},
						"total": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
					},
				},
				"v1.CreateServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"name": &openapi.Schema{
							Type: "string",
						},
						"bin": &openapi.Schema{
							Type: "string",
						},
						"args": &openapi.Schema{
							Type:  "array",
							Items: &openapi.Schema{Type: "string"},
						},
						"dir": &openapi.Schema{
							Type: "string",
						},
						"env": &openapi.Schema{
							AdditionalProperties: &openapi.Schema{},
						},
						"sysProcAttr": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.SysProcAttr",
						},
						"log": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.ProcLog",
						},
						"version": &openapi.Schema{
							Type: "string",
						},
						"autoRestart": &openapi.Schema{
							Type:   "integer",
							Format: "int32",
						},
					},
					Required: []string{"name", "bin"},
				},
				"v1.CreateServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.GetServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
					},
					Required: []string{"id"},
				},
				"v1.GetServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.DeleteServiceReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
					},
					Required: []string{"id"},
				},
				"v1.DeleteServiceRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.GetServiceByNameReq": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"name": &openapi.Schema{
							Type: "string",
						},
					},
					Required: []string{"name"},
				},
				"v1.GetServiceByNameRsp": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"service": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Service",
						},
					},
				},
				"v1.Service": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"id": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
						"name": &openapi.Schema{
							Type: "string",
						},
						"bin": &openapi.Schema{
							Type: "string",
						},
						"args": &openapi.Schema{
							Type:  "array",
							Items: &openapi.Schema{Type: "string"},
						},
						"pid": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
						"dir": &openapi.Schema{
							Type: "string",
						},
						"env": &openapi.Schema{
							AdditionalProperties: &openapi.Schema{},
						},
						"sysProcAttr": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.SysProcAttr",
						},
						"log": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.ProcLog",
						},
						"version": &openapi.Schema{
							Type: "string",
						},
						"autoRestart": &openapi.Schema{
							Type:   "integer",
							Format: "int32",
						},
						"creationTimestamp": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
						"updateTimestamp": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
						"startTimestamp": &openapi.Schema{
							Type:   "integer",
							Format: "int64",
						},
						"status": &openapi.Schema{
							Type: "string",
							Enum: []string{"init", "running", "stopped", "failed", "upgrading"},
						},
						"msg": &openapi.Schema{
							Type: "string",
						},
						"stat": &openapi.Schema{
							Type: "object",
							Ref:  "#/components/schemas/v1.Stat",
						},
					},
					Required: []string{"name", "bin"},
				},
				"v1.PageMeta": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"page": &openapi.Schema{
							Type:             "integer",
							Format:           "int32",
							ExclusiveMinimum: true,
							Minimum:          0,
							Default:          "1",
						},
						"size": &openapi.Schema{
							Type:             "integer",
							Format:           "int32",
							ExclusiveMinimum: true,
							Minimum:          0,
							Default:          "10",
						},
					},
				},
				"v1.SysProcAttr": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"chroot": &openapi.Schema{
							Type: "string",
						},
						"uid": &openapi.Schema{
							Type:   "integer",
							Format: "int32",
						},
						"user": &openapi.Schema{
							Type: "string",
						},
						"gid": &openapi.Schema{
							Type:   "integer",
							Format: "int32",
						},
						"group": &openapi.Schema{
							Type: "string",
						},
					},
				},
				"v1.ProcLog": &openapi.Model{
					Type:       "object",
					Properties: map[string]*openapi.Schema{},
				},
				"v1.Stat": &openapi.Model{
					Type: "object",
					Properties: map[string]*openapi.Schema{
						"cpuPercent": &openapi.Schema{
							Type:   "number",
							Format: "double",
						},
						"memory": &openapi.Schema{},
						"memPercent": &openapi.Schema{
							Type:   "number",
							Format: "float",
						},
					},
				},
			},
		},
	}
}

// Client API for GpmService service
// +gen:openapi
type GpmService interface {
	// +gne:summary=gpm 检测 gpm 服务状态
	// +gen:get=/api/healthz
	Healthz(ctx context.Context, in *Empty, opts ...client.CallOption) (*Empty, error)
	// +gen:summary=查询所有服务
	// +gen:get=/api/v1/Service/
	ListService(ctx context.Context, in *ListServiceReq, opts ...client.CallOption) (*ListServiceRsp, error)
	// +gen:summary=查询单个服务
	// +gen:get=/api/v1/Service/{id}
	GetService(ctx context.Context, in *GetServiceReq, opts ...client.CallOption) (*GetServiceRsp, error)
	// +gen:summary=通过名称查询单个服务
	// +gen:get=/api/v1/Service/{name}
	GetServiceByName(ctx context.Context, in *GetServiceByNameReq, opts ...client.CallOption) (*GetServiceByNameRsp, error)
	// +gen:summary=新建服务
	// +gen:post=/api/v1/Service/
	CreateService(ctx context.Context, in *CreateServiceReq, opts ...client.CallOption) (*CreateServiceRsp, error)
	// +gen:summary=启动服务
	// +gen:patch=/api/v1/Service/{id}/action/start
	StartService(ctx context.Context, in *StartServiceReq, opts ...client.CallOption) (*StartServiceRsp, error)
	// +gen:summary=停止服务
	// +gen:patch=/api/v1/Service/{id}/action/stop
	StopService(ctx context.Context, in *StopServiceReq, opts ...client.CallOption) (*StopServiceRsp, error)
	// +gen:summary=重启服务
	// +gen:patch=/api/v1/Service/{id}/action/reboot
	RebootService(ctx context.Context, in *RebootServiceReq, opts ...client.CallOption) (*RebootServiceRsp, error)
	// +gen:summary=删除服务
	// +gen:delete=/api/v1/Service/{id}
	DeleteService(ctx context.Context, in *DeleteServiceReq, opts ...client.CallOption) (*DeleteServiceRsp, error)
}

type gpmService struct {
	c    client.Client
	name string
}

func NewGpmService(name string, c client.Client) GpmService {
	return &gpmService{
		c:    c,
		name: name,
	}
}

func (c *gpmService) Healthz(ctx context.Context, in *Empty, opts ...client.CallOption) (*Empty, error) {
	req := c.c.NewRequest(c.name, "GpmService.Healthz", in)
	out := new(Empty)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) ListService(ctx context.Context, in *ListServiceReq, opts ...client.CallOption) (*ListServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.ListService", in)
	out := new(ListServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) GetService(ctx context.Context, in *GetServiceReq, opts ...client.CallOption) (*GetServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.GetService", in)
	out := new(GetServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) GetServiceByName(ctx context.Context, in *GetServiceByNameReq, opts ...client.CallOption) (*GetServiceByNameRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.GetServiceByName", in)
	out := new(GetServiceByNameRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) CreateService(ctx context.Context, in *CreateServiceReq, opts ...client.CallOption) (*CreateServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.CreateService", in)
	out := new(CreateServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) StartService(ctx context.Context, in *StartServiceReq, opts ...client.CallOption) (*StartServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.StartService", in)
	out := new(StartServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) StopService(ctx context.Context, in *StopServiceReq, opts ...client.CallOption) (*StopServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.StopService", in)
	out := new(StopServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) RebootService(ctx context.Context, in *RebootServiceReq, opts ...client.CallOption) (*RebootServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.RebootService", in)
	out := new(RebootServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gpmService) DeleteService(ctx context.Context, in *DeleteServiceReq, opts ...client.CallOption) (*DeleteServiceRsp, error) {
	req := c.c.NewRequest(c.name, "GpmService.DeleteService", in)
	out := new(DeleteServiceRsp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GpmService service
// +gen:openapi
type GpmServiceHandler interface {
	// +gne:summary=gpm 检测 gpm 服务状态
	// +gen:get=/api/healthz
	Healthz(context.Context, *Empty, *Empty) error
	// +gen:summary=查询所有服务
	// +gen:get=/api/v1/Service/
	ListService(context.Context, *ListServiceReq, *ListServiceRsp) error
	// +gen:summary=查询单个服务
	// +gen:get=/api/v1/Service/{id}
	GetService(context.Context, *GetServiceReq, *GetServiceRsp) error
	// +gen:summary=通过名称查询单个服务
	// +gen:get=/api/v1/Service/{name}
	GetServiceByName(context.Context, *GetServiceByNameReq, *GetServiceByNameRsp) error
	// +gen:summary=新建服务
	// +gen:post=/api/v1/Service/
	CreateService(context.Context, *CreateServiceReq, *CreateServiceRsp) error
	// +gen:summary=启动服务
	// +gen:patch=/api/v1/Service/{id}/action/start
	StartService(context.Context, *StartServiceReq, *StartServiceRsp) error
	// +gen:summary=停止服务
	// +gen:patch=/api/v1/Service/{id}/action/stop
	StopService(context.Context, *StopServiceReq, *StopServiceRsp) error
	// +gen:summary=重启服务
	// +gen:patch=/api/v1/Service/{id}/action/reboot
	RebootService(context.Context, *RebootServiceReq, *RebootServiceRsp) error
	// +gen:summary=删除服务
	// +gen:delete=/api/v1/Service/{id}
	DeleteService(context.Context, *DeleteServiceReq, *DeleteServiceRsp) error
}

func RegisterGpmServiceHandler(s server.Server, hdlr GpmServiceHandler, opts ...server.HandlerOption) error {
	type gpmServiceImpl interface {
		Healthz(ctx context.Context, in *Empty, out *Empty) error
		ListService(ctx context.Context, in *ListServiceReq, out *ListServiceRsp) error
		GetService(ctx context.Context, in *GetServiceReq, out *GetServiceRsp) error
		GetServiceByName(ctx context.Context, in *GetServiceByNameReq, out *GetServiceByNameRsp) error
		CreateService(ctx context.Context, in *CreateServiceReq, out *CreateServiceRsp) error
		StartService(ctx context.Context, in *StartServiceReq, out *StartServiceRsp) error
		StopService(ctx context.Context, in *StopServiceReq, out *StopServiceRsp) error
		RebootService(ctx context.Context, in *RebootServiceReq, out *RebootServiceRsp) error
		DeleteService(ctx context.Context, in *DeleteServiceReq, out *DeleteServiceRsp) error
	}
	type GpmService struct {
		gpmServiceImpl
	}
	h := &gpmServiceHandler{hdlr}
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.Healthz",
		Description: "GpmService.Healthz",
		Path:        []string{"/api/healthz"},
		Method:      []string{"GET"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.ListService",
		Description: "GpmService.ListService",
		Path:        []string{"/api/v1/Service/"},
		Method:      []string{"GET"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.GetService",
		Description: "GpmService.GetService",
		Path:        []string{"/api/v1/Service/{id}"},
		Method:      []string{"GET"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.GetServiceByName",
		Description: "GpmService.GetServiceByName",
		Path:        []string{"/api/v1/Service/{name}"},
		Method:      []string{"GET"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.CreateService",
		Description: "GpmService.CreateService",
		Path:        []string{"/api/v1/Service/"},
		Method:      []string{"POST"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.StartService",
		Description: "GpmService.StartService",
		Path:        []string{"/api/v1/Service/{id}/action/start"},
		Method:      []string{"PATCH"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.StopService",
		Description: "GpmService.StopService",
		Path:        []string{"/api/v1/Service/{id}/action/stop"},
		Method:      []string{"PATCH"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.RebootService",
		Description: "GpmService.RebootService",
		Path:        []string{"/api/v1/Service/{id}/action/reboot"},
		Method:      []string{"PATCH"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, api.WithEndpoint(&apipb.Endpoint{
		Name:        "GpmService.DeleteService",
		Description: "GpmService.DeleteService",
		Path:        []string{"/api/v1/Service/{id}"},
		Method:      []string{"DELETE"},
		Body:        "*",
		Handler:     "rpc",
	}))
	opts = append(opts, server.OpenAPIHandler(NewGpmServiceOpenAPI()))
	return s.Handle(s.NewHandler(&GpmService{h}, opts...))
}

type gpmServiceHandler struct {
	GpmServiceHandler
}

func (h *gpmServiceHandler) Healthz(ctx context.Context, in *Empty, out *Empty) error {
	return h.GpmServiceHandler.Healthz(ctx, in, out)
}

func (h *gpmServiceHandler) ListService(ctx context.Context, in *ListServiceReq, out *ListServiceRsp) error {
	return h.GpmServiceHandler.ListService(ctx, in, out)
}

func (h *gpmServiceHandler) GetService(ctx context.Context, in *GetServiceReq, out *GetServiceRsp) error {
	return h.GpmServiceHandler.GetService(ctx, in, out)
}

func (h *gpmServiceHandler) GetServiceByName(ctx context.Context, in *GetServiceByNameReq, out *GetServiceByNameRsp) error {
	return h.GpmServiceHandler.GetServiceByName(ctx, in, out)
}

func (h *gpmServiceHandler) CreateService(ctx context.Context, in *CreateServiceReq, out *CreateServiceRsp) error {
	return h.GpmServiceHandler.CreateService(ctx, in, out)
}

func (h *gpmServiceHandler) StartService(ctx context.Context, in *StartServiceReq, out *StartServiceRsp) error {
	return h.GpmServiceHandler.StartService(ctx, in, out)
}

func (h *gpmServiceHandler) StopService(ctx context.Context, in *StopServiceReq, out *StopServiceRsp) error {
	return h.GpmServiceHandler.StopService(ctx, in, out)
}

func (h *gpmServiceHandler) RebootService(ctx context.Context, in *RebootServiceReq, out *RebootServiceRsp) error {
	return h.GpmServiceHandler.RebootService(ctx, in, out)
}

func (h *gpmServiceHandler) DeleteService(ctx context.Context, in *DeleteServiceReq, out *DeleteServiceRsp) error {
	return h.GpmServiceHandler.DeleteService(ctx, in, out)
}
