package genconfig

type GenerationConfig struct {
	TargetDir    string
	PackageName  string
	TemplateName string
	Interfaces   []InterfaceType
}
