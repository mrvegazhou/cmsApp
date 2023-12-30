package models

func GetModels() []interface{} {
	return []interface{}{
		&AppUser{}, &SiteConfig{},
	}
}
