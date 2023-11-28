package db

import "fmt"

const (
	Template1 = `Mail Şablonu 1`
	Template2 = `Mail Şablonu 2`
)

func GetMailTemplate(templateID int) (string, error) {
	// Şablonları döndür
	switch templateID {
	case 1:
		return Template1, nil
	case 2:
		return Template2, nil
	default:
		return "", fmt.Errorf("unknown template ID: %d", templateID)
	}
}
