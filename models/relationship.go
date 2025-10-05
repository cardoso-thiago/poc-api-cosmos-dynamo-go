package models

type RelationshipType struct {
	Relationship string `json:"relationship"`
}

type Relationship struct {
	Relationship      string             `json:"relationship"`
	RelationshipTypes []RelationshipType `json:"relationship_type,omitempty"`
}

type Response struct {
	Data []Relationship `json:"data"`
}
