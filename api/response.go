package api

type ResponseImage struct {
	ID    int32  `json:"id"`
	Thumb string `json:"thumb"`
	Large string `json:"large"`
}

type Response struct {
	Image ResponseImage `json:"image"`
	Count int32         `json:"count"`
}
