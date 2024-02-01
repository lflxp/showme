package consted

// 表单字段类型
type FieldType string

const (
	StringField FieldType = "string"
	IntField    FieldType = "int"
	Int16Field  FieldType = "int16"
	Int64Field  FieldType = "int64"
	Text        FieldType = "text"
	Textarea    FieldType = "textarea"
	Radio       FieldType = "radio"
	Select      FieldType = "select"
	MultiSelect FieldType = "multiselect"
	File        FieldType = "file"
	Time        FieldType = "time"
	OneToOne    FieldType = "o2o"
	OneToMany   FieldType = "o2m"
	ManyToMany  FieldType = "m2m"
	Password    FieldType = "password"
)

// html组件类型
type HtmlComponentType int

const (
	FormComponent HtmlComponentType = iota
	EditFormComponent
	AdminColumns
	StrToHtml
	BeegoLi
	Upload
)
