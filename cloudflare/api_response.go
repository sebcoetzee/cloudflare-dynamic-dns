package cloudflare

type APIResponse struct {
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
