package db

import (
	"fmt"
	"strings"
)

type Template struct {
	Id             int
	Name           string
	Alias          string
	Subject        string
	Content        string
	Description    string
	SiteId         int
	IsSmsTemplate  int
	Lang           string
	RenderTemplate bool
}

type MailRecord struct {
	Subject string
	Content string
}

func GetMailContent(templateAlias string, siteID int, customVariables map[string]interface{}) (*MailRecord, error) {
	template, err := getMailTemplateByTemplateId(templateAlias, siteID)
	parseMailContentToTemplate(&template.Content, customVariables)
	if err != nil {
		// eğer nil döneceksen diğer dönüşte referans dönmeli
		return nil, err
	}

	// anonim bir struct ile sadece senderin ihtiyacı olanları dönelim, template structu bu classta lazım
	// Anonim struct düşündüğüm gibi olmadı, dönüş parametresi olarak struct tanımlamak lazım
	return &MailRecord{Subject: template.Subject, Content: template.Content}, nil
}

func getMailTemplateByTemplateId(templateAlias string, siteID int) (*Template, error) {
	var template Template
	// Burayı gorma taşıyabilirim
	err := DB.QueryRow("SELECT * FROM cms_email_templates WHERE alias=? and site_id=?", templateAlias, siteID).Scan(
		&template.Id,
		&template.Name,
		&template.Alias,
		&template.Subject,
		&template.Content,
		&template.Description,
		&template.SiteId,
		&template.IsSmsTemplate,
		&template.Lang,
		&template.RenderTemplate,
	)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func parseMailContentToTemplate(templateContent *string, customVariables map[string]interface{}) {
	for key, value := range customVariables {
		stringValue := fmt.Sprintf("%v", value)
		*templateContent = strings.Replace(*templateContent, "{{"+key+"}}", stringValue, -1)
	}
}
