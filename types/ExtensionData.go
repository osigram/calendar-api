package types

type ExtensionData struct {
	Id             int64  `json:"id" db:"id"`
	AdditionalData string `json:"additionalData" db:"additional_data"`
}
