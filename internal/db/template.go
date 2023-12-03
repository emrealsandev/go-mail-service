package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Template struct {
	// Eğer typeleri sql paketi Nullxxx diye vermezsek değer db den null gelirse patlıyor
	Id             sql.NullInt64  `db:"id"`
	Name           sql.NullString `db:"name"`
	Alias          sql.NullString `db:"alias"`
	Subject        sql.NullString `db:"subject"`
	Content        sql.NullString `db:"content"`
	Description    sql.NullString `db:"description"`
	SiteId         sql.NullInt16  `db:"site_id"`
	IsSmsTemplate  sql.NullBool   `db:"is_sms_template"`
	Lang           sql.NullString `db:"lang"`
	RenderTemplate sql.NullBool   `db:"render_template"`
}

type MailRecord struct {
	Subject string
	Content string
}

func GetMailContent(templateAlias string, siteID int, customVariables map[string]interface{}) (*MailRecord, error) {
	template, err := getMailTemplateByTemplateId(templateAlias, siteID)
	if !(template.Content.Valid || template.Subject.Valid) { // Nullxxx alanları string functinolarında kullanamıyoruz .String ile dönüştürmek lazım
		return nil, errors.New("Mail içeriği geçersiz")
	}
	parseMailContentToTemplate(&template.Content.String, customVariables)
	if err != nil {
		// eğer nil döneceksen diğer dönüşte referans dönmeli
		return nil, err
	}

	// anonim bir struct ile sadece senderin ihtiyacı olanları dönelim, template structu bu classta lazım
	// Anonim struct düşündüğüm gibi olmadı, dönüş parametresi olarak struct tanımlamak lazım
	return &MailRecord{Subject: template.Subject.String, Content: template.Content.String}, nil
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
		stringValue := fmt.Sprintf("%v", value) // value.(string) şeklinde de bir convert var ama bu daha esnek
		*templateContent = strings.Replace(*templateContent, "{{"+key+"}}", stringValue, -1)
	}
}
