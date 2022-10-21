package l27

func (c *Client) Login(username string, password string) (Login, error) {
	var login Login

	err := c.invokeAPI("POST", "login", &LoginRequest{Username: username, Password: password}, &login)

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
	Hash string `json:"hash"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
