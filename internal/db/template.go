package db

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

func GetMailContent(templateID int, customVariables map[string]interface{}) (*Template, error) {
	template, err := getMailTemplateByTemplateId(templateID)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func getMailTemplateByTemplateId(templateID int) (*Template, error) {
	var template Template
	// Burayı gorma taşıyabilirim
	err := DB.QueryRow("SELECT * FROM cms_email_templates WHERE id=?", templateID).Scan(
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
