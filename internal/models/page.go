package models

var (
	// DefaultPage default page
	DefaultPage = 1
	// DefaultSize default size
	DefaultSize = 20
)

// PageInformation page information
type PageInformation struct {
	Page                  int   `json:"page,omitempty"`
	Size                  int   `json:"size,omitempty"`
	TotalNumberOfEntities int64 `json:"total_number_of_entities,omitempty"`
	TotalNumberOfPages    int   `json:"total_number_of_pages,omitempty"`
}

// Page page model
type Page struct {
	PageInformation *PageInformation `json:"page_information,omitempty"`
	Entities        interface{}      `json:"entities,omitempty"`
	Footer          interface{}      `json:"footer,omitempty"`
}

// NewPage new page
func NewPage(pif *PageInformation, es interface{}) *Page {
	return &Page{
		PageInformation: &PageInformation{
			Page:                  pif.Page,
			Size:                  pif.Size,
			TotalNumberOfEntities: pif.TotalNumberOfEntities,
			TotalNumberOfPages:    pif.TotalNumberOfPages,
		},
		Entities: es,
	}
}

// GetEntities get entities
func (p *Page) GetEntities() interface{} {
	return p.Entities
}

// PageForm page form
type PageForm struct {
	Page     int      `json:"page,omitempty" form:"page" query:"page"`
	Size     int      `json:"size,omitempty" form:"size" query:"size"`
	Query    string   `json:"query,omitempty" form:"query" query:"query"`
	Sort     string   `json:"sort,omitempty" form:"sort" query:"sort"`
	Sorts    []string `json:"sorts,omitempty" swaggerignore:"true"`
	Reverse  bool     `json:"reverse,omitempty" form:"reverse" query:"reverse"`
	Reverses []bool   `json:"reverses,omitempty" swaggerignore:"true"`
	OrderBy  string   `json:"-" form:"-"`
}

// GetPage get page
func (f *PageForm) GetPage() int {
	if f.Page == 0 {
		f.Page = DefaultPage
	}
	return f.Page
}

// GetSize get size
func (f *PageForm) GetSize() int {
	if f.Size == 0 {
		f.Size = DefaultSize
	}
	return f.Size
}

// GetQuery get query
func (f *PageForm) GetQuery() string {
	return f.Query
}

// GetSort get sort
func (f *PageForm) GetSort() string {
	return f.Sort
}

// GetReverse get reverse
func (f *PageForm) GetReverse() bool {
	return f.Reverse
}

// GetReverses get reverses
func (f *PageForm) GetReverses() []bool {
	return f.Reverses
}

// GetOrderBy get order by
func (f *PageForm) GetOrderBy() string {
	return f.OrderBy
}

// GetSorts get sorts
func (f *PageForm) GetSorts() []string {
	return f.Sorts
}
