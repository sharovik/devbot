package dto

//DataKey the key struct
type DataKey struct {
	Key string `json:"key"`
}

//Error the error struct
type Error struct {
	Message string  `json:"message"`
	Data    DataKey `json:"data"`
}

//BitBucketErrorResponseBranchCreate the bad response from the bitbucket api
type BitBucketErrorResponseBranchCreate struct {
	Data  DataKey `json:"data"`
	Type  string  `json:"type"`
	Error Error   `json:"error"`
}
