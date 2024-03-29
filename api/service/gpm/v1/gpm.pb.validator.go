// Code generated by proto-gen-validator
// source: github.com/vine-io/gpm/api/service/gpm/v1/gpm.proto

package gpmv1

import (
	fmt "fmt"
	_ "github.com/vine-io/gpm/api/types/gpm/v1"
	is "github.com/vine-io/vine/util/is"
)

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

func (m *RestartServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *RestartServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *RestartServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *RestartServiceRsp) ValidateE(prefix string) error {
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

func (m *ForgetServiceReq) Validate() error {
	return m.ValidateE("")
}

func (m *ForgetServiceReq) ValidateE(prefix string) error {
	errs := make([]error, 0)
	if len(m.Name) == 0 {
		errs = append(errs, fmt.Errorf("field '%sname' is required", prefix))
	}
	return is.MargeErr(errs...)
}

func (m *ForgetServiceRsp) Validate() error {
	return m.ValidateE("")
}

func (m *ForgetServiceRsp) ValidateE(prefix string) error {
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
