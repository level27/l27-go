package l27

// POST /login
func (c *Client) Login(username string, password string) (Login, error) {
	return c.Login2FA(&LoginRequest{Username: username, Password: password})
}

// Alternative to Login() that accepts full requests data, for 2FA etc.
func (c *Client) Login2FA(request *LoginRequest) (Login, error) {
	var login Login

	err := c.invokeAPI("POST", "login", request, &login)

	return login, err
}

// POST /login
// Fetches login information for the current API key.
func (c *Client) LoginInfo() (Login, error) {
	var login Login

	err := c.invokeAPI("POST", "login", nil, &login)

	return login, err
}

type Login struct {
	Success bool `json:"success"`
	User    struct {
		ID           IntID    `json:"id"`
		Username     string   `json:"username"`
		Email        string   `json:"email"`
		Firstname    string   `json:"firstName"`
		Lastname     string   `json:"lastName"`
		Roles        []string `json:"roles"`
		Status       string   `json:"status"`
		Language     string   `json:"language"`
		Organisation struct {
			ID          IntID  `json:"id"`
			Name        string `json:"name"`
			Street      string `json:"street"`
			Housenumber string `json:"houseNumber"`
			Zip         string `json:"zip"`
			City        string `json:"city"`
		} `json:"organisation"`
		Country struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"country"`
		Fullname string `json:"fullname"`
	} `json:"user"`
	Hash    string `json:"hash"`
	Hash2FA string `json:"hash2fa"`
}

type LoginRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	SixDigitCode    string `json:"6digitCode"`
	TrustThisDevice bool   `json:"trustThisDevice"`
	TwoFAToken      string `json:"2faToken"`
}
