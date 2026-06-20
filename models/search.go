package models

// GlobalSearchResponse is the response for GET /search?q=.
type GlobalSearchResponse struct {
	Families []FamilyListItem `json:"families"`
	Clients  []Client         `json:"clients"`
	Policies []PolicyListItem `json:"policies"`
}
