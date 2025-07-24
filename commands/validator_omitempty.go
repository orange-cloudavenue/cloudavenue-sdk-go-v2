package commands

type validatorOmitempty struct{}

func ValidatorOmitempty() Validator {
	return &validatorOmitempty{}
}

func (v *validatorOmitempty) GetKey() string {
	return "omitempty"
}

func (v *validatorOmitempty) GetDescription() string {
	return ""
}

func (v *validatorOmitempty) GetMarkdownDescription() string {
	return ""
}
