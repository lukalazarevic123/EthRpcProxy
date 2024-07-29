package model

type HolderEntity struct {
	BaseModel
	HolderAddress string `json:"holderAddress"`
	IsHolder      bool   `json:"isHolder"`
	BlockNumber   int    `json:"blockNumber"`
}

func (h HolderEntity) GetName() string {
	return "holder"
}
