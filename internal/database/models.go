package database

import (
	"encoding/json"
	"time"
)

// User model
type User struct {
	ID           string    `gorm:"primaryKey"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CompanyID    int       `json:"companyID"`
	AssignedJobs []*Job    `json:"assignedJobs" gorm:"many2many:user_jobs;"`
}

// Company model
type Company struct {
	ID              int         `json:"id" gorm:"primaryKey"`
	CompanyName     string      `json:"companyName"`
	CompanySize     json.Number `json:"companySize"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	Users           []User      `json:"users"`
	CustomMilestone []CustomMilestone
}

// JobStatus enum
type JobStatus string

const (
	//Active enums
	Active JobStatus = "Active"
	//Expired enums
	Expired = "Expired"
	//Fulfilled enums
	Fulfilled = "Fulfilled"
	//Ended enums
	Ended = "Ended"
)

// Job model
type Job struct {
	ID                    int                `gorm:"primaryKey" json:"id"`
	Location              string             `json:"location" gorm:"notNull"`
	PayRate               float32            `json:"payRate" gorm:"notNull"`
	Description           string             `json:"description" gorm:"notNull"`
	Title                 string             `json:"title" gorm:"notNull"`
	City                  string             `json:"city" gorm:"notNull"`
	State                 string             `json:"state" gorm:"notNull"`
	CreatedAt             time.Time          `json:"createdAt" gorm:"notNull"`
	UpdatedAt             time.Time          `json:"updatedAt" gorm:"notNull"`
	JobDuties             string             `json:"jobDuties" gorm:"notNull"`
	Status                string             `json:"status" gorm:"notNull;default:Active"`
	ExpireDate            time.Time          `json:"expireDate" gorm:"notNull"`
	EmploymentType        string             `json:"employmentType" gorm:"notNull"`
	Industry              string             `json:"industry" gorm:"notNull"`
	Openings              int                `json:"openings" gorm:"notNull"`
	CompanyID             int                `json:"companyID" gorm:"notNull"`
	Recruiters            []*User            `json:"recruiters" gorm:"many2many:user_jobs;"`
	Milestones            []*Milestone       `json:"milestones" gorm:"many2many:job_milestones;"`
	CustomMilestones      []*CustomMilestone `json:"customMilestones" gorm:"many2many:job_custom_milestones;"`
	Benefits              []*Benefits        `json:"benefits" gorm:"many2many:job_benefits;"`
	RequiresAuthorization bool
}

// UserJob association table
type UserJob struct {
	UserID    string    `gorm:"primaryKey" json:"userId"`
	JobID     string    `gorm:"primaryKey" json:"jobId"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Application table
type Application struct {
	ID                int `gorm:"primaryKey" json:"id"`
	FirstName         string
	LastName          string
	Status            string
	Email             string
	PhoneNumber       string
	Address           string
	City              string
	State             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Files             []ApplicantFile
	EmploymentHistory []EmployerHistory
	EducationHistory  []EducationHistory
	Milestone         Milestone
	MilestoneID       int
	CustomMilestone   CustomMilestone
	CustomMilestoneID int
	WorkAuthorized    bool
	ConvictedFelon    bool
}

// ApplicantFile table
type ApplicantFile struct {
	ID            int    `gorm:"primaryKey"`
	Type          string `json:"type"`
	Path          string `json:"path"`
	ApplicationID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// EducationHistory  table
type EducationHistory struct {
	ID              int `gorm:"primaryKey"`
	Type            string
	FromDate        time.Time
	ToDate          time.Time
	Graduated       bool
	InstitutionName string
	Major           string
	Degree          string
	ApplicationID   int
}

// EmployerHistory table
type EmployerHistory struct {
	ID              int `gorm:"primaryKey"`
	EmployerName    string
	EmployerAddress string
	Position        string
	Duties          string
	FromDate        time.Time
	ToDate          time.Time
	Current         bool
	LeaveReason     string
	ApplicationID   int
}

// Milestone table
type Milestone struct {
	ID   int `gorm:"primaryKey"`
	Name string
	Jobs []*Job `json:"jobs" gorm:"many2many:job_milestones;"`
}

// CustomMilestone table
type CustomMilestone struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	CompanyID int
	Jobs      []*Job `json:"jobs" gorm:"many2many:job_custom_milestones;"`
}

// Benefits table
type Benefits struct {
	ID   int `gorm:"primaryKey"`
	Type string
	Jobs []*Job `json:"jobs" gorm:"many2many:job_benefits;"`
}
