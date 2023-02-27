package pojo

import "fmt"

type Job struct {
	Name string
	Update string
	Salary string
	Postion string
	Experience string
	Degree string
	Tags []string
	Describe string
	Url string
	Company
}

type Company struct {
	CName string
	CType string
	CSize string
	MainBusiness []string
	CDescribe string
}

func (j Job) String() string {
	return fmt.Sprintf(
		"{\n\tName: %s\n\tUpdate: %s\n\tSalary: %s\n\tPostion: %s\n\tExperience: %s\n\tDegree: %s\n\tTags: %s\n\tDescribe: %s\n\tUrl: %s\n\tCompany: %s\n}",
		j.Name, j.Update, j.Salary, j.Postion, j.Experience, j.Degree, j.Tags, j.Describe, j.Url, j.Company,
	)
}

func (c Company) String() string {
	return fmt.Sprintf(
		"{\n\t\tCname: %s\n\t\tCtype: %s\n\t\tCsize: %s\n\t\tMainBusiness: %s\n\t\tCdescribe: %s\n\t}",
		c.CName, c.CType, c.CSize, c.MainBusiness, c.CDescribe,
	)
}