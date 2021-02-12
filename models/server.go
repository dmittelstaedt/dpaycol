package models

import "database/sql"

// ServerEnvironment enum
type ServerEnvironment string

// ServerRole enum
type ServerRole string

// Constants for ServerEnvironment
const (
	Development ServerEnvironment = "development"
	Test        ServerEnvironment = "test"
	Stage       ServerEnvironment = "stage"
	Producation ServerEnvironment = "production"
)

// Constants for ServerRole
const (
	App ServerRole = "app"
	Web ServerRole = "web"
	Db  ServerRole = "db"
)

// Server holds server information
type Server struct {
	ID          int               `json:"id"`
	VMName      string            `json:"vmName"`
	OS          string            `json:"os"`
	IPAddress   string            `json:"ipAddress"`
	Hostname    string            `json:"hostname"`
	Domain      string            `json:"domain"`
	CPU         int               `json:"cpu"`
	Memory      int               `json:"memory"`
	Environment ServerEnvironment `json:"environment"`
	Role        ServerRole        `json:"role"`
	VSwitch     string            `json:"vswitch"`
	VLan        string            `json:"vlan"`
	OSA         string            `json:"osa"`
	Comment     sql.NullString    `json:"comment"`
	Auto        int               `json:"auto"`
	LparID      int               `json:"lparId"`
	ServiceID   sql.NullInt64     `json:"serviceId"`
}
