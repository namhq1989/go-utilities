package language

type Language string

const (
	Unknown    Language = ""
	English    Language = "en"
	Vietnamese Language = "vi"
	Japanese   Language = "ja"
	Korean     Language = "ko"
	Chinese    Language = "zh"
	Thai       Language = "th"
)

func (l Language) String() string {
	return string(l)
}

func (l Language) IsValid() bool {
	return l != Unknown
}

func (l Language) IsEnglish() bool {
	return l == English
}

func (l Language) IsVietnamese() bool {
	return l == Vietnamese
}

func (l Language) IsJapanese() bool {
	return l == Japanese
}

func (l Language) IsKorean() bool {
	return l == Korean
}

func (l Language) IsChinese() bool {
	return l == Chinese
}

func (l Language) IsThai() bool {
	return l == Thai
}

func (l Language) GetCountry() string {
	switch l {
	case Vietnamese:
		return "Vietnam"
	default:
		return ""
	}
}

func ToLanguage(lang string) Language {
	switch lang {
	case English.String():
		return English
	case Vietnamese.String():
		return Vietnamese
	case Japanese.String():
		return Japanese
	case Korean.String():
		return Korean
	case Chinese.String():
		return Chinese
	case Thai.String():
		return Thai
	default:
		return Unknown
	}
}
