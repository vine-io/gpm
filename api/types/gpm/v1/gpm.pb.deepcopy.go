// Code generated by proto-gen-deepcopy. DO NOT EDIT.
// source: github.com/vine-io/gpm/api/types/gpm/v1/gpm.proto

package gpmv1

import ()

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Service) DeepCopyInto(out *Service) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SysProcAttr != nil {
		in, out := &in.SysProcAttr, &out.SysProcAttr
		*out = new(SysProcAttr)
		(*in).DeepCopyInto(*out)
	}
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = new(ProcLog)
		(*in).DeepCopyInto(*out)
	}
	if in.Stat != nil {
		in, out := &in.Stat, &out.Stat
		*out = new(Stat)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *SysProcAttr) DeepCopyInto(out *SysProcAttr) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *ServiceSpec) DeepCopyInto(out *ServiceSpec) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SysProcAttr != nil {
		in, out := &in.SysProcAttr, &out.SysProcAttr
		*out = new(SysProcAttr)
		(*in).DeepCopyInto(*out)
	}
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = new(ProcLog)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *UpgradeSpec) DeepCopyInto(out *UpgradeSpec) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *EditServiceSpec) DeepCopyInto(out *EditServiceSpec) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SysProcAttr != nil {
		in, out := &in.SysProcAttr, &out.SysProcAttr
		*out = new(SysProcAttr)
		(*in).DeepCopyInto(*out)
	}
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = new(ProcLog)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *ProcLog) DeepCopyInto(out *ProcLog) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Stat) DeepCopyInto(out *Stat) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *GpmInfo) DeepCopyInto(out *GpmInfo) {
	*out = *in
	if in.Stat != nil {
		in, out := &in.Stat, &out.Stat
		*out = new(Stat)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *Package) DeepCopyInto(out *Package) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *InstallServiceIn) DeepCopyInto(out *InstallServiceIn) {
	*out = *in
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(ServiceSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Pack != nil {
		in, out := &in.Pack, &out.Pack
		*out = new(Package)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *InstallServiceResult) DeepCopyInto(out *InstallServiceResult) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *UpgradeServiceIn) DeepCopyInto(out *UpgradeServiceIn) {
	*out = *in
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		*out = new(UpgradeSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Pack != nil {
		in, out := &in.Pack, &out.Pack
		*out = new(Package)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *UpgradeServiceResult) DeepCopyInto(out *UpgradeServiceResult) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *ServiceLog) DeepCopyInto(out *ServiceLog) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *ServiceVersion) DeepCopyInto(out *ServiceVersion) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *FileInfo) DeepCopyInto(out *FileInfo) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *UpdateIn) DeepCopyInto(out *UpdateIn) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *UpdateResult) DeepCopyInto(out *UpdateResult) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *ExecIn) DeepCopyInto(out *ExecIn) {
	*out = *in
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *ExecResult) DeepCopyInto(out *ExecResult) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *PullResult) DeepCopyInto(out *PullResult) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *PushIn) DeepCopyInto(out *PushIn) {
	*out = *in
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *TerminalIn) DeepCopyInto(out *TerminalIn) {
	*out = *in
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopyInto is an auto-generated deepcopy function, coping the receiver, writing into out. in must be no-nil.
func (in *TerminalResult) DeepCopyInto(out *TerminalResult) {
	*out = *in
}
