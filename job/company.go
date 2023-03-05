package job

import "fmt"

type Company struct {
	CId          int    `db:"cid"`
	CName        string `db:"cname"`
	CType        string `db:"type"`
	CSize        string `db:"size"`
	MainBusiness string `db:"main_business"`
	CDescribe    string `db:"cdescribe"`
}

func (c Company) String() string {
	return fmt.Sprintf(
		"{\n\t\tCId: %d\n\t\tCname: %s\n\t\tCtype: %s\n\t\tCsize: %s\n\t\tMainBusiness: %s\n\t\tCdescribe: %s\n\t}",
		c.CId, c.CName, c.CType, c.CSize, c.MainBusiness, c.CDescribe,
	)
}
