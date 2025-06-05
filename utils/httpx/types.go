package httpclient

type Config struct {
	ActName string // 操作名称
}

type FormValue struct {
	Value     interface{}
	IsFile    bool
	FieldName string
	FileName  string
}

type FormFile struct {
	Path     string // 文件路径
	Field    string // form字段名（可选，默认键名）
	FileName string // 文件名（可选）
}
