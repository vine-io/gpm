// Code generated by proto-gen-validator. DO NOT EDIT.
// source: github.com/gpm2/gpm/proto/service/gpm/v1/gpm.proto

package gpmv1

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gpm2/gpm/proto/apis/gpm/v1"
	is "github.com/lack-io/vine/util/is"
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

func (m *Empty) Validate() error {
	return m.ValidateE("")
}

func (m *Empty) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *UpdateSelfReq) Validate() error {
	return m.ValidateE("")
}

func (m *UpdateSelfReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *UpdateSelfRsp) Validate() error {
	return m.ValidateE("")
}

func (m *UpdateSelfRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *InfoReq) Validate() error {
	return m.ValidateE("")
}

func (m *InfoReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *InfoRsp) Validate() error {
	return m.ValidateE("")
}

func (m *InfoRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *ListServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *ListServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *ListServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *ListServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *GetServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *GetServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *GetServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *GetServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *CreateServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *CreateServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if m.Spec == nil {
		errs = append(errs, fmt.Errorf("field '%sspec' is required", prefix))
	} else {
		errs = append(errs, m.Spec.ValidateE(prefix+"spec."))
	}
	return is.MargeErr(errs...)
}

func (m *CreateServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *CreateServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *EditServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *EditServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if m.Spec == nil {
		errs = append(errs, fmt.Errorf("field '%sspec' is required", prefix))
	} else {
		errs = append(errs, m.Spec.ValidateE(prefix+"spec."))
	}
	return is.MargeErr(errs...)
}

func (m *EditServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *EditServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *StartServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *StartServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *StartServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *StartServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *StopServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *StopServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *StopServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *StopServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *RebootServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *RebootServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *RebootServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *RebootServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *DeleteServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *DeleteServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *DeleteServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *DeleteServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *WatchServiceLogReq) Validate() error {
	return m.ValidateE("")
}

func (m *WatchServiceLogReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *WatchServiceLogRsp) Validate() error {
	return m.ValidateE("")
}

func (m *WatchServiceLogRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *InstallServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *InstallServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if m.Spec == nil {
		errs = append(errs, fmt.Errorf("field '%sspec' is required", prefix))
	} else {
		errs = append(errs, m.Spec.ValidateE(prefix+"spec."))
	}
	if m.Pack == nil {
		errs = append(errs, fmt.Errorf("field '%spack' is required", prefix))
	} else {
		errs = append(errs, m.Pack.ValidateE(prefix+"pack."))
	}
	return is.MargeErr(errs...)
}

func (m *InstallServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *InstallServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *ListServiceVersionsReq) Validate() error {
	return m.ValidateE("")
}

func (m *ListServiceVersionsReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *ListServiceVersionsRsp) Validate() error {
	return m.ValidateE("")
}

func (m *ListServiceVersionsRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *UpgradeServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *UpgradeServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	if len(m.Version) == 0 {
		errs = append(errs, fmt.Errorf("field '%sversion' is required", prefix))
	}
	if m.Pack == nil {
		errs = append(errs, fmt.Errorf("field '%spack' is required", prefix))
	} else {
		errs = append(errs, m.Pack.ValidateE(prefix+"pack."))
	}
	return is.MargeErr(errs...)
}

func (m *UpgradeServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *UpgradeServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *RollbackServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *RollbackServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *RollbackServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *RollbackServiceRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *LsReq) Validate() error {
	return m.ValidateE("")
}

func (m *LsReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Path) == 0 {
		errs = append(errs, fmt.Errorf("field '%spath' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *LsRsp) Validate() error {
	return m.ValidateE("")
}

func (m *LsRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *PullReq) Validate() error {
	return m.ValidateE("")
}

func (m *PullReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *PullRsp) Validate() error {
	return m.ValidateE("")
}

func (m *PullRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *PushReq) Validate() error {
	return m.ValidateE("")
}

func (m *PushReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *PushRsp) Validate() error {
	return m.ValidateE("")
}

func (m *PushRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *ExecReq) Validate() error {
	return m.ValidateE("")
}

func (m *ExecReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if m.In == nil {
		errs = append(errs, fmt.Errorf("field '%sin' is required", prefix))
	} else {
		errs = append(errs, m.In.ValidateE(prefix+"in."))
	}
	return is.MargeErr(errs...)
}

func (m *ExecRsp) Validate() error {
	return m.ValidateE("")
}

func (m *ExecRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}

func (m *TerminalReq) Validate() error {
	return m.ValidateE("")
}

func (m *TerminalReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if m.In == nil {
		errs = append(errs, fmt.Errorf("field '%sin' is required", prefix))
	} else {
		errs = append(errs, m.In.ValidateE(prefix+"in."))
	}
	return is.MargeErr(errs...)
}

func (m *TerminalRsp) Validate() error {
	return m.ValidateE("")
}

func (m *TerminalRsp) ValidateE(prefix string) error {
	errs := make([]error, 0)
	return is.MargeErr(errs...)
}
