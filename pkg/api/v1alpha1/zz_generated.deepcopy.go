package v1alpha1

func (in *Application) DeepCopy() *Application                    {}
func (in *Application) DeepCopyInto(out *Application)             {}
func (in *ApplicationSpec) DeepCopy() *ApplicationSpec            {}
func (in *ApplicationSpec) DeepCopyInto(out *ApplicationSpec)     {}
func (in *ApplicationStatus) DeepCopy() *ApplicationStatus        {}
func (in *ApplicationStatus) DeepCopyInto(out *ApplicationStatus) {}
