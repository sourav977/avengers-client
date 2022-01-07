package avengersclient

type Avenger struct {
	ID     int    `json:"_id,omitempty"`
	Name   string `json:"name,omitempty"`
	Alias  string `json:"alias,omitempty"`
	Weapon string `json:"weapon,omitempty"`
}

type UpdateResult struct {
	MatchedCount  int `json:"MatchedCount,omitempty"`
	ModifiedCount int `json:"ModifiedCount,omitempty"`
	UpsertedCount int `json:"UpsertedCount,omitempty"`
}

type DeleteResult struct {
	DeletedCount int `json:"DeletedCount,omitempty"`
}
