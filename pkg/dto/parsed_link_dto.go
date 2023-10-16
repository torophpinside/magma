package dto

type ParsedLinkDTO struct {
	Links  []string `json:"links"`
	Target string   `json:"target"`
	Dork   string   `json:"dork"`
}
