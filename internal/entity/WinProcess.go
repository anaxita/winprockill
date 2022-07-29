package entity

type WinProcess struct {
	ID       int64  `json:"Id"`
	Name     string `json:"Name"`
	UserName string `json:"UserName"`
}
