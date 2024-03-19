package canvas

type Qualification struct {
	AccountID int
	Name      string
}

var Qualifications []Qualification = []Qualification{
	{
		AccountID: 119,
		Name:      "Certificate III in Early Childhood Education and Care",
	},
	{
		AccountID: 120,
		Name:      "Diploma of Early Childhood Education and Care",
	},
	{
		AccountID: 133,
		Name:      "Certificate IV in Kitchen Management",
	},
	{
		AccountID: 134,
		Name:      "Diploma of Hospitality Management",
	},
	{
		AccountID: 130,
		Name:      "Certificate III in Automotive",
	},
	{
		AccountID: 129,
		Name:      "Certificate IV in Automotive",
	},
	{
		AccountID: 128,
		Name:      "Diploma of Automotive",
	},
	{
		AccountID: 123,
		Name:      "Bricklaying",
	},
	{
		AccountID: 116,
		Name:      "Certificate III in Carpentry",
	},
	{
		AccountID: 114,
		Name:      "Telecommunication",
	},
	{
		AccountID: 124,
		Name:      "Tiling",
	},
	{
		AccountID: 112,
		Name:      "Advanced Diploma of Business",
	},
	{
		AccountID: 111,
		Name:      "Diploma of Business",
	},
	{
		AccountID: 113,
		Name:      "Graduate Diploma of Management (Learning)",
	},
}

func (c *Controller) GetQualifications() []Qualification {
	return Qualifications
}
