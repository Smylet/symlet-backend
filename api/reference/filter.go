package reference


type UniversityQueryParams struct {
	Name string `form:"name,omitempty"`
	Code string `form:"code,omitempty"`
	City string `form:"city,omitempty"`
	State string `form:"state,omitempty"`
	Country string `form:"country,omitempty"`
}
