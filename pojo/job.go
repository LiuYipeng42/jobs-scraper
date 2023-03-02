package pojo

import "fmt"

type Job struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	Update     string `db:"update"`
	Salary     string `db:"salary"`
	Position   string `db:"position"`
	Experience string `db:"experience"`
	Degree     string `db:"degree"`
	Tags       string `db:"tags"`
	Describe   string `db:"jdescribe"`
	Url        string `db:"url"`
	Company
}

type Company struct {
	CId          int    `db:"cid"`
	CName        string `db:"cname"`
	CType        string `db:"type"`
	CSize        string `db:"size"`
	MainBusiness string `db:"main_business"`
	CDescribe    string `db:"cdescribe"`
}

func (j Job) String() string {
	return fmt.Sprintf(
		"{\n\tId: %d\n\tName: %s\n\tUpdate: %s\n\tSalary: %s\n\tPostion: %s\n\tExperience: %s\n\tDegree: %s\n\tTags: %s\n\tDescribe: %s\n\tUrl: %s\n\tCompany: %s\n}",
		j.Id, j.Name, j.Update, j.Salary, j.Position, j.Experience, j.Degree, j.Tags, j.Describe, j.Url, j.Company,
	)
}

func (c Company) String() string {
	return fmt.Sprintf(
		"{\n\t\tCId: %d\n\t\tCname: %s\n\t\tCtype: %s\n\t\tCsize: %s\n\t\tMainBusiness: %s\n\t\tCdescribe: %s\n\t}",
		c.CId, c.CName, c.CType, c.CSize, c.MainBusiness, c.CDescribe,
	)
}
