package docx

import (
	"custom-resume-builder/prompt"
	"os"
	"strings"

	"github.com/fumiama/go-docx"
)

const (
	fontName = "Arial"

	nameSize        = "80"
	designationSize = "25"
	contactSize     = "23"
	sectionSize     = "25"
	bodySize        = "23"
	smallSize       = "20"
)

func addContactInfo(
	doc *docx.Docx,
	contact prompt.Contact,
) {
	p := doc.AddParagraph().
		Justification("center")

	first := true

	addSeparator := func() {
		if !first {
			p.AddText("\u00A0·\u00A0").
				Size(contactSize)
		}

		first = false
	}

	if strings.TrimSpace(contact.Phone) != "" {
		addSeparator()

		p.AddText(strings.TrimSpace(contact.Phone)).
			Size(contactSize)
	}

	if strings.TrimSpace(contact.Email) != "" {
		addSeparator()

		p.AddText(strings.TrimSpace(contact.Email)).
			Size(contactSize)
	}

	if strings.TrimSpace(contact.Github) != "" {
		addSeparator()

		link := p.AddLink(
			"GitHub",
			normalizeURL(contact.Github),
		)

		link.Run.Size(contactSize)
	}

	if strings.TrimSpace(contact.Linkedin) != "" {
		addSeparator()

		link := p.AddLink(
			"LinkedIn",
			normalizeURL(contact.Linkedin),
		)

		link.Run.Size(contactSize)
	}

	if strings.TrimSpace(contact.Website) != "" {
		addSeparator()

		link := p.AddLink(
			"Website",
			normalizeURL(contact.Website),
		)

		link.Run.Size(contactSize)
	}

	if strings.TrimSpace(contact.Location) != "" {
		addSeparator()

		p.AddText(strings.TrimSpace(contact.Location)).
			Size(contactSize)
	}
}

func normalizeURL(url string) string {
	url = strings.TrimSpace(url)

	if url == "" {
		return ""
	}

	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		return "https://" + url
	}

	return url
}


func GenerateResumeDocx(resume prompt.Resume, outputPath string) error {
	// A4 document
	doc := docx.New().WithDefaultTheme().WithA4Page()

	setPageMargins(doc)
	// =========================================================
	// HEADER
	// =========================================================

	addCenteredText(
		doc,
		resume.Profile.Name,
		nameSize,
		true,
	)

	if resume.Profile.Designation != "" {
		addCenteredText(
			doc,
			resume.Profile.Designation,
			designationSize,
			false,
		)
	}

	// Contact information
	addContactInfo(doc, resume.Contact)

	// =========================================================
	// SUMMARY
	// =========================================================

	if strings.TrimSpace(resume.Summary) != "" {
		addSectionTitle(doc, "SUMMARY")
		addBodyParagraph(doc, resume.Summary)
	}

	// =========================================================
	// SKILLS
	// =========================================================

	if hasSkills(resume.Skills) {
		addSectionTitle(doc, "SKILLS")

		addSkillRow(
			doc,
			"Languages",
			resume.Skills.Languages,
		)

		addSkillRow(
			doc,
			"Frameworks & Libraries",
			resume.Skills.FrameworksLibraries,
		)

		addSkillRow(
			doc,
			"Databases",
			resume.Skills.DatabasesCaching,
		)

		addSkillRow(
			doc,
			"DevOps & Infrastructure",
			resume.Skills.DevOpsInfrastructure,
		)
	}

	// =========================================================
	// EXPERIENCE
	// =========================================================

	if len(resume.Experience) > 0 {
		addSectionTitle(doc, "EXPERIENCE")

		for _, experience := range resume.Experience {
			addExperience(doc, experience)
		}
	}

	// =========================================================
	// PROJECTS
	// =========================================================

	if len(resume.Projects) > 0 {
		addSectionTitle(doc, "PROJECTS AND WORKS")

		for _, project := range resume.Projects {
			addProject(doc, project)
		}
	}

	// =========================================================
	// EDUCATION
	// =========================================================

	if len(resume.Education) > 0 {
		addSectionTitle(doc, "EDUCATION")

		for _, education := range resume.Education {
			addEducation(doc, education)
		}
	}

	// =========================================================
	// WRITE
	// =========================================================

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = doc.WriteTo(file)

	return err
}

// =============================================================
// HEADER
// =============================================================

func addCenteredText(
	doc *docx.Docx,
	text string,
	size string,
	bold bool,
) {
	if strings.TrimSpace(text) == "" {
		return
	}

	p := doc.AddParagraph().
		Justification("center")

	t := p.AddText(strings.TrimSpace(text)).
		Size(size).
		Color("#002147")

	if bold {
		t.Bold()
	}
}

// =============================================================
// SECTION TITLE
// =============================================================

func addSectionTitle(doc *docx.Docx, title string) {
	p := doc.AddParagraph()

	if p.Properties == nil {
		p.Properties = &docx.ParagraphProperties{}
	}

	p.Properties.Spacing = &docx.Spacing{
		Before: 160,
	}

	p.AddText(strings.ToUpper(title)).
		Size(sectionSize).
		Bold().
		Color("#002147")
}

func setPageMargins(doc *docx.Docx) {
	for _, item := range doc.Document.Body.Items {
		if sectPr, ok := item.(*docx.SectPr); ok {
			sectPr.PgMar = &docx.PgMar{
				Top:    650,
				Bottom: 650,
				Left:   650,
				Right:  650,
				Header: 400,
				Footer: 400,
				Gutter: 0,
			}
			return
		}
	}
}


// =============================================================
// SUMMARY / BODY
// =============================================================

func addBodyParagraph(
	doc *docx.Docx,
	text string,
) {
	if strings.TrimSpace(text) == "" {
		return
	}

	p := doc.AddParagraph()

	p.AddText(strings.TrimSpace(text)).
		Size(bodySize)
}

// =============================================================
// SKILLS
// =============================================================

func addSkillRow(
	doc *docx.Docx,
	label string,
	skills []string,
) {
	if len(skills) == 0 {
		return
	}

	filtered := make([]string, 0, len(skills))

	for _, skill := range skills {
		skill = strings.TrimSpace(skill)

		if skill != "" {
			filtered = append(filtered, skill)
		}
	}

	if len(filtered) == 0 {
		return
	}

	p := doc.AddParagraph()

	p.AddText(label+":").
		Size(bodySize).
		Bold()

	p.AddText(" "+strings.Join(filtered, ", ")).
		Size(bodySize)
}

// =============================================================
// EXPERIENCE
// =============================================================

func addExperience(
	doc *docx.Docx,
	exp prompt.Experience,
) {
	if strings.TrimSpace(exp.Company) == "" &&
		strings.TrimSpace(exp.Role) == "" &&
		len(exp.Highlights) == 0 {
		return
	}

	// ---------------------------------------------------------
	// Company / Type / Date
	// ---------------------------------------------------------

	p := doc.AddParagraph()
	addEntrySpacing(p)

	if exp.Company != "" {
		p.AddText(strings.TrimSpace(exp.Company)).
			Size(bodySize).
			Bold()
	}

	if exp.Type != "" {
		p.AddText(" — "+strings.TrimSpace(exp.Type)).
			Size(smallSize)
	}

	// Tab separates date from company information.
	if exp.Date != "" {
		p.AddText("\u00A0|\u00A0").
    		Size(smallSize)

		p.AddText(strings.TrimSpace(exp.Date)).
			Size(smallSize)
	}

	// ---------------------------------------------------------
	// Role
	// ---------------------------------------------------------

	if strings.TrimSpace(exp.Role) != "" {
		role := doc.AddParagraph()

		role.AddText(strings.TrimSpace(exp.Role)).
			Size(smallSize).
			Bold()
	}

	// ---------------------------------------------------------
	// Highlights
	// ---------------------------------------------------------

	for _, highlight := range exp.Highlights {
		addBullet(doc, highlight)
	}
}

// =============================================================
// PROJECT
// =============================================================

func addProject(
	doc *docx.Docx,
	project prompt.Project,
) {
	if strings.TrimSpace(project.Name) == "" &&
		len(project.Highlights) == 0 {
		return
	}

	p := doc.AddParagraph()
	addEntrySpacing(p)

	if project.Name != "" {
		p.AddText(strings.TrimSpace(project.Name)).
			Size(bodySize).
			Bold()
	}

	if project.Type != "" {
		p.AddText(" — " + strings.TrimSpace(project.Type)).
			Size(smallSize)
	}

	for _, highlight := range project.Highlights {
		addBullet(doc, highlight)
	}
}

// =============================================================
// EDUCATION
// =============================================================

func addEducation(
	doc *docx.Docx,
	education prompt.Education,
) {
	if strings.TrimSpace(education.Degree) == "" &&
		strings.TrimSpace(education.Institution) == "" {
		return
	}

	// Degree + period
	p := doc.AddParagraph()
	addEntrySpacing(p)

	if education.Degree != "" {
		p.AddText(strings.TrimSpace(education.Degree)).
			Size(bodySize).
			Bold()
	}

	if education.Period != "" {
		p.AddText("\u00A0|\u00A0").
    		Size(smallSize)

		p.AddText(strings.TrimSpace(education.Period)).
			Size(smallSize)
	}

	// Institution + university
	var parts []string

	if education.Institution != "" {
		parts = append(
			parts,
			strings.TrimSpace(education.Institution),
		)
	}

	if education.University != "" {
		parts = append(
			parts,
			strings.TrimSpace(education.University),
		)
	}

	if len(parts) > 0 {
		p = doc.AddParagraph()

		p.AddText(strings.Join(parts, " — ")).
			Size(smallSize)
	}
}

// =============================================================
// BULLET
// =============================================================

func addBullet(
	doc *docx.Docx,
	text string,
) {
	text = strings.TrimSpace(text)

	if text == "" {
		return
	}

	p := doc.AddParagraph()

	if p.Properties == nil {
		p.Properties = &docx.ParagraphProperties{}
	}

	p.Properties.Ind = &docx.Ind{
		Left: 200,
	}

	p.AddText("•\u00A0"+text).
		Size(bodySize)
}

// =============================================================
// SKILL CHECK
// =============================================================

func hasSkills(skills prompt.Skills) bool {
	return len(skills.Languages) > 0 ||
		len(skills.FrameworksLibraries) > 0 ||
		len(skills.DatabasesCaching) > 0 ||
		len(skills.DevOpsInfrastructure) > 0
}

func addEntrySpacing(p *docx.Paragraph) {
	if p.Properties == nil {
		p.Properties = &docx.ParagraphProperties{}
	}

	p.Properties.Spacing = &docx.Spacing{
		Before: 80,
	}
}