package uaago

type uaaTokenFetcher struct {
	uaaUrl                string
	username              string
	password              string
	insecureSSLSkipVerify bool
}

//matches github.com.cloudfoundry.noaa.consumer.TokenRefresher
type TokenFetcher interface {
	RefreshAuthToken() (token string, authError error)
}

func NewUAATokenFetcher(url string, username string, password string, sslSkipVerify bool) TokenFetcher {
	return &uaaTokenFetcher{
		uaaUrl:                url,
		username:              username,
		password:              password,
		insecureSSLSkipVerify: sslSkipVerify,
	}
}

func (uaa *uaaTokenFetcher) RefreshAuthToken() (string, error) {
	uaaClient, err := NewClient(uaa.uaaUrl)
	if err != nil {
		return "", err
	}

	authToken, err := uaaClient.GetAuthToken(uaa.username, uaa.password, uaa.insecureSSLSkipVerify)
	if err != nil {
		return "", err
	}
	return authToken, nil
}
