// Code generated by proto-gen-dao. DO NOT EDIT.
// source: github.com/gpm2/gpm/proto/apis/gpm/v1/gpm.proto

package dao

import (
	context "context"
	driver "database/sql/driver"
	errors "errors"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	v1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	json "github.com/json-iterator/go"
	dao "github.com/lack-io/vine/lib/dao"
	clause "github.com/lack-io/vine/lib/dao/clause"
	math "math"
	strings "strings"
	time "time"
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

// ServiceArgs the alias of []string
type ServiceArgs []string

// Value return json value, implement driver.Valuer interface
func (m ServiceArgs) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(m)
	return string(b), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *ServiceArgs) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, &m)
}

func (m *ServiceArgs) DaoDataType() string {
	return dao.DefaultDialect.JSONDataType()
}

// ServiceStat the alias of v1.Stat
type ServiceStat v1.Stat

// Value return json value, implement driver.Valuer interface
func (m *ServiceStat) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := json.Marshal(m)
	return string(b), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *ServiceStat) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	return json.Unmarshal(bytes, &m)
}

func (m *ServiceStat) DaoDataType() string {
	return dao.DefaultDialect.JSONDataType()
}

// ServiceS the Schema for Service
type ServiceS struct {
	tx    *dao.DB             `json:"-" dao:"-"`
	exprs []clause.Expression `json:"-" dao:"-"`

	Id                int64        `json:"id,omitempty" dao:"column:id;autoIncrement;primaryKey"`
	Name              string       `json:"name,omitempty" dao:"column:name"`
	Bin               string       `json:"bin,omitempty" dao:"column:bin"`
	Args              ServiceArgs  `json:"args,omitempty" dao:"column:args"`
	Pid               int64        `json:"pid,omitempty" dao:"column:pid"`
	Chroot            string       `json:"chroot,omitempty" dao:"column:chroot"`
	Uid               int32        `json:"uid,omitempty" dao:"column:uid"`
	User              string       `json:"user,omitempty" dao:"column:user"`
	Gid               int32        `json:"gid,omitempty" dao:"column:gid"`
	Group             string       `json:"group,omitempty" dao:"column:group"`
	Version           string       `json:"version,omitempty" dao:"column:version"`
	CreationTimestamp int64        `json:"creationTimestamp,omitempty" dao:"column:creation_timestamp"`
	UpdateTimestamp   int64        `json:"updateTimestamp,omitempty" dao:"column:update_timestamp"`
	Status            string       `json:"status,omitempty" dao:"column:status"`
	Msg               string       `json:"msg,omitempty" dao:"column:msg"`
	Stat              *ServiceStat `json:"stat,omitempty" dao:"column:stat"`
	DeletionTimestamp int64        `json:"deletionTimestamp,omitempty" dao:"column:deletion_timestamp"`
}

func RegistryService() error {
	return dao.DefaultDialect.Migrator().AutoMigrate(&ServiceS{})
}

func ServiceSBuilder() *ServiceS {
	return &ServiceS{tx: dao.DefaultDialect.NewTx()}
}

func (m *ServiceS) SetId(in int64) *ServiceS {
	m.Id = in
	return m
}

func (m *ServiceS) SetName(in string) *ServiceS {
	m.Name = in
	return m
}

func (m *ServiceS) SetBin(in string) *ServiceS {
	m.Bin = in
	return m
}

func (m *ServiceS) SetArgs(in []string) *ServiceS {
	m.Args = in
	return m
}

func (m *ServiceS) AddArgs(in string) *ServiceS {
	if len(m.Args) == 0 {
		m.Args = []string{}
	}
	m.Args = append(m.Args, in)
	return m
}

func (m *ServiceS) FilterArgs(fn func(item string) bool) *ServiceS {
	out := make([]string, 0)
	for _, item := range m.Args {
		if !fn(item) {
			out = append(out, item)
		}
	}
	m.Args = out
	return m
}

func (m *ServiceS) RemoveArgs(index int) *ServiceS {
	if index < 0 || index >= len(m.Args) {
		return m
	}
	if index == len(m.Args)-1 {
		m.Args = m.Args[:index-1]
	} else {
		m.Args = append(m.Args[0:index], m.Args[index+1:]...)
	}
	return m
}

func (m *ServiceS) SetPid(in int64) *ServiceS {
	m.Pid = in
	return m
}

func (m *ServiceS) SetChroot(in string) *ServiceS {
	m.Chroot = in
	return m
}

func (m *ServiceS) SetUid(in int32) *ServiceS {
	m.Uid = in
	return m
}

func (m *ServiceS) SetUser(in string) *ServiceS {
	m.User = in
	return m
}

func (m *ServiceS) SetGid(in int32) *ServiceS {
	m.Gid = in
	return m
}

func (m *ServiceS) SetGroup(in string) *ServiceS {
	m.Group = in
	return m
}

func (m *ServiceS) SetVersion(in string) *ServiceS {
	m.Version = in
	return m
}

func (m *ServiceS) SetCreationTimestamp(in int64) *ServiceS {
	m.CreationTimestamp = in
	return m
}

func (m *ServiceS) SetUpdateTimestamp(in int64) *ServiceS {
	m.UpdateTimestamp = in
	return m
}

func (m *ServiceS) SetStatus(in string) *ServiceS {
	m.Status = in
	return m
}

func (m *ServiceS) SetMsg(in string) *ServiceS {
	m.Msg = in
	return m
}

func (m *ServiceS) SetStat(in *v1.Stat) *ServiceS {
	m.Stat = (*ServiceStat)(in)
	return m
}

func FromService(in *v1.Service) *ServiceS {
	out := new(ServiceS)
	out.tx = dao.DefaultDialect.NewTx()
	if in.Id != 0 {
		out.Id = in.Id
	}
	if in.Name != "" {
		out.Name = in.Name
	}
	if in.Bin != "" {
		out.Bin = in.Bin
	}
	if len(in.Args) != 0 {
		out.Args = in.Args
	}
	if in.Pid != 0 {
		out.Pid = in.Pid
	}
	if in.Chroot != "" {
		out.Chroot = in.Chroot
	}
	if in.Uid != 0 {
		out.Uid = in.Uid
	}
	if in.User != "" {
		out.User = in.User
	}
	if in.Gid != 0 {
		out.Gid = in.Gid
	}
	if in.Group != "" {
		out.Group = in.Group
	}
	if in.Version != "" {
		out.Version = in.Version
	}
	if in.CreationTimestamp != 0 {
		out.CreationTimestamp = in.CreationTimestamp
	}
	if in.UpdateTimestamp != 0 {
		out.UpdateTimestamp = in.UpdateTimestamp
	}
	if in.Status != "" {
		out.Status = in.Status
	}
	if in.Msg != "" {
		out.Msg = in.Msg
	}
	if in.Stat != nil {
		out.Stat = (*ServiceStat)(in.Stat)
	}
	return out
}

func (m *ServiceS) ToService() *v1.Service {
	out := new(v1.Service)
	out.Id = m.Id
	out.Name = m.Name
	out.Bin = m.Bin
	out.Args = m.Args
	out.Pid = m.Pid
	out.Chroot = m.Chroot
	out.Uid = m.Uid
	out.User = m.User
	out.Gid = m.Gid
	out.Group = m.Group
	out.Version = m.Version
	out.CreationTimestamp = m.CreationTimestamp
	out.UpdateTimestamp = m.UpdateTimestamp
	out.Status = m.Status
	out.Msg = m.Msg
	out.Stat = (*v1.Stat)(m.Stat)
	return out
}

func (ServiceS) TableName() string {
	return "services"
}

func (m ServiceS) PrimaryKey() (string, interface{}, bool) {
	return "id", m.Id, m.Id == 0
}

func (m *ServiceS) FindPage(ctx context.Context, page, size int) ([]*v1.Service, int64, error) {
	pk, _, _ := m.PrimaryKey()

	m.exprs = append(m.exprs,
		clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Table: m.TableName(), Name: pk}, Desc: true}}},
		clause.Limit{Offset: int((page - 1) * size), Limit: int(size)},
		clause.Cond().Build("deletion_timestamp", 0),
	)

	data, err := m.findAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	total, _ := m.Count(ctx)
	return data, total, nil
}

func (m *ServiceS) FindAll(ctx context.Context) ([]*v1.Service, error) {
	m.exprs = append(m.exprs, clause.Cond().Build("deletion_timestamp", 0))
	return m.findAll(ctx)
}

func (m *ServiceS) FindPureAll(ctx context.Context) ([]*v1.Service, error) {
	return m.findAll(ctx)
}

func (m *ServiceS) findAll(ctx context.Context) ([]*v1.Service, error) {
	dest := make([]*ServiceS, 0)
	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), m.exprs...)
	if err := tx.Clauses(clauses...).Find(&dest).Error; err != nil {
		return nil, err
	}

	outs := make([]*v1.Service, len(dest))
	for i := range dest {
		outs[i] = dest[i].ToService()
	}

	return outs, nil
}

func (m *ServiceS) Count(ctx context.Context) (total int64, err error) {
	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), clause.Cond().Build("deletion_timestamp", 0))
	clauses = append(clauses, m.exprs...)

	err = tx.Clauses(clauses...).Count(&total).Error
	return
}

func (m *ServiceS) FindOne(ctx context.Context) (*v1.Service, error) {
	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), clause.Cond().Build("deletion_timestamp", 0))
	clauses = append(clauses, m.exprs...)

	if err := tx.Clauses(clauses...).First(&m).Error; err != nil {
		return nil, err
	}

	return m.ToService(), nil
}

func (m *ServiceS) FindPureOne(ctx context.Context) (*v1.Service, error) {
	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	if err := tx.Clauses(m.exprs...).First(&m).Error; err != nil {
		return nil, err
	}

	return m.ToService(), nil
}

func (m *ServiceS) Cond(exprs ...clause.Expression) *ServiceS {
	m.exprs = append(m.exprs, exprs...)
	return m
}

func (m *ServiceS) extractClauses(tx *dao.DB) []clause.Expression {
	exprs := make([]clause.Expression, 0)
	if m.Id != 0 {
		exprs = append(exprs, clause.Cond().Build("id", m.Id))
	}
	if m.Name != "" {
		exprs = append(exprs, clause.Cond().Build("name", m.Name))
	}
	if m.Bin != "" {
		exprs = append(exprs, clause.Cond().Build("bin", m.Bin))
	}
	if len(m.Args) != 0 {
		for _, item := range m.Args {
			exprs = append(exprs, dao.DefaultDialect.JSONBuild("args").Tx(tx).Contains(dao.JSONEq, item))
		}
	}
	if m.Pid != 0 {
		exprs = append(exprs, clause.Cond().Build("pid", m.Pid))
	}
	if m.Chroot != "" {
		exprs = append(exprs, clause.Cond().Build("chroot", m.Chroot))
	}
	if m.Uid != 0 {
		exprs = append(exprs, clause.Cond().Build("uid", m.Uid))
	}
	if m.User != "" {
		exprs = append(exprs, clause.Cond().Build("user", m.User))
	}
	if m.Gid != 0 {
		exprs = append(exprs, clause.Cond().Build("gid", m.Gid))
	}
	if m.Group != "" {
		exprs = append(exprs, clause.Cond().Build("group", m.Group))
	}
	if m.Version != "" {
		exprs = append(exprs, clause.Cond().Build("version", m.Version))
	}
	if m.CreationTimestamp != 0 {
		exprs = append(exprs, clause.Cond().Build("creation_timestamp", m.CreationTimestamp))
	}
	if m.UpdateTimestamp != 0 {
		exprs = append(exprs, clause.Cond().Build("update_timestamp", m.UpdateTimestamp))
	}
	if m.Status != "" {
		exprs = append(exprs, clause.Cond().Build("status", m.Status))
	}
	if m.Msg != "" {
		exprs = append(exprs, clause.Cond().Build("msg", m.Msg))
	}
	if m.Stat != nil {
		for k, v := range dao.FieldPatch(m.Stat) {
			if v == nil {
				exprs = append(exprs, dao.DefaultDialect.JSONBuild("stat").Tx(tx).Op(dao.JSONHasKey, "", strings.Split(k, ".")...))
			} else {
				exprs = append(exprs, dao.DefaultDialect.JSONBuild("stat").Tx(tx).Op(dao.JSONEq, v, strings.Split(k, ".")...))
			}
		}
	}

	return exprs
}

func (m *ServiceS) Create(ctx context.Context) (*v1.Service, error) {
	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	if err := tx.Create(m).Error; err != nil {
		return nil, err
	}

	return m.ToService(), nil
}

func (m *ServiceS) BatchUpdates(ctx context.Context) error {
	if len(m.exprs) == 0 {
		return errors.New("missing conditions")
	}

	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	values := make(map[string]interface{}, 0)
	if m.Name != "" {
		values["name"] = m.Name
	}
	if m.Bin != "" {
		values["bin"] = m.Bin
	}
	if len(m.Args) != 0 {
		values["args"] = m.Args
	}
	if m.Pid != 0 {
		values["pid"] = m.Pid
	}
	if m.Chroot != "" {
		values["chroot"] = m.Chroot
	}
	if m.Uid != 0 {
		values["uid"] = m.Uid
	}
	if m.User != "" {
		values["user"] = m.User
	}
	if m.Gid != 0 {
		values["gid"] = m.Gid
	}
	if m.Group != "" {
		values["group"] = m.Group
	}
	if m.Version != "" {
		values["version"] = m.Version
	}
	if m.CreationTimestamp != 0 {
		values["creation_timestamp"] = m.CreationTimestamp
	}
	if m.UpdateTimestamp != 0 {
		values["update_timestamp"] = m.UpdateTimestamp
	}
	if m.Status != "" {
		values["status"] = m.Status
	}
	if m.Msg != "" {
		values["msg"] = m.Msg
	}
	if m.Stat != nil {
		values["stat"] = m.Stat
	}

	return tx.Clauses(m.exprs...).Updates(values).Error
}

func (m *ServiceS) Updates(ctx context.Context) (*v1.Service, error) {
	pk, pkv, isNil := m.PrimaryKey()
	if isNil {
		return nil, errors.New("missing primary key")
	}

	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx).Where(pk+" = ?", pkv)

	values := make(map[string]interface{}, 0)
	if m.Name != "" {
		values["name"] = m.Name
	}
	if m.Bin != "" {
		values["bin"] = m.Bin
	}
	if len(m.Args) != 0 {
		values["args"] = m.Args
	}
	if m.Pid != 0 {
		values["pid"] = m.Pid
	}
	if m.Chroot != "" {
		values["chroot"] = m.Chroot
	}
	if m.Uid != 0 {
		values["uid"] = m.Uid
	}
	if m.User != "" {
		values["user"] = m.User
	}
	if m.Gid != 0 {
		values["gid"] = m.Gid
	}
	if m.Group != "" {
		values["group"] = m.Group
	}
	if m.Version != "" {
		values["version"] = m.Version
	}
	if m.CreationTimestamp != 0 {
		values["creation_timestamp"] = m.CreationTimestamp
	}
	if m.UpdateTimestamp != 0 {
		values["update_timestamp"] = m.UpdateTimestamp
	}
	if m.Status != "" {
		values["status"] = m.Status
	}
	if m.Msg != "" {
		values["msg"] = m.Msg
	}
	if m.Stat != nil {
		values["stat"] = m.Stat
	}

	if err := tx.Updates(values).Error; err != nil {
		return nil, err
	}

	err := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx).Where(pk+" = ?", pkv).First(m).Error
	if err != nil {
		return nil, err
	}
	return m.ToService(), nil
}

func (m *ServiceS) BatchDelete(ctx context.Context, soft bool) error {
	if len(m.exprs) == 0 {
		return errors.New("missing conditions")
	}

	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	if soft {
		return tx.Clauses(m.exprs...).Updates(map[string]interface{}{"deletion_timestamp": time.Now().UnixNano()}).Error
	}
	return tx.Clauses(m.exprs...).Delete(&ServiceS{}).Error
}

func (m *ServiceS) Delete(ctx context.Context, soft bool) error {
	pk, pkv, isNil := m.PrimaryKey()
	if isNil {
		return errors.New("missing primary key")
	}

	tx := m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx)

	if soft {
		return tx.Where(pk+" = ?", pkv).Updates(map[string]interface{}{"deletion_timestamp": time.Now().UnixNano()}).Error
	}
	return tx.Where(pk+" = ?", pkv).Delete(&ServiceS{}).Error
}

func (m *ServiceS) Tx(ctx context.Context) *dao.DB {
	return m.tx.Session(&dao.Session{}).Table(m.TableName()).WithContext(ctx).Clauses(m.exprs...)
}
