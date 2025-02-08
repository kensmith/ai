package model

type Model struct {
	Name     string
	Provider string
}

type Models []Model

var ModelMap = map[string]Models{}

func Register(model Model) {
	_, ok := ModelMap[model.Provider]
	if !ok {
		ModelMap[model.Provider] = Models{}
	}

	ModelMap[model.Provider] = append(ModelMap[model.Provider], model)
}
