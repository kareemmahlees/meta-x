package lib

type TablePropsValidator struct {
	Type     string      `json:"type" validate:"required"`
	Nullable interface{} `json:"nullable" validate:"omitempty,boolean"`
	Default  interface{} `json:"default" validate:"omitempty,alpha" `
	Unique   interface{} `json:"unique" validate:"omitempty,boolean"`
}
