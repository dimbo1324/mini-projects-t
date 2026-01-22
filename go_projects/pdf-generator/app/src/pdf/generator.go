package pdf

type Generator struct {
	templatePath string
}

func NewGenerator(templatePath string) *Generator {
	return &Generator{templatePath: templatePath}
}
