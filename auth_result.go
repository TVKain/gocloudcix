package gocloudcix

/*
AuthResult is the result from the request that was used to obtain a provider
client's token. It is returned from ProviderClient.GetAuthResult().
*/
type AuthResult interface {
	ExtractTokenID() (string, error)
}
