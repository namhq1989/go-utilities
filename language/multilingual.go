package language

type Multilingual struct {
	English    string `json:"en,omitempty"`
	Vietnamese string `json:"vi,omitempty"` // Vietnamese
}

func (m Multilingual) GetLocalized(lang string) Multilingual {
	// English by default
	result := Multilingual{
		English: m.English,
	}

	switch lang {
	case Vietnamese.String():
		result.Vietnamese = m.Vietnamese
	}

	return result
}

func (m Multilingual) IsEmpty() bool {
	return m.English == "" && m.Vietnamese == ""
}
