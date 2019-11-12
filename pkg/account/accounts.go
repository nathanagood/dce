package account

type Accounts []Account

// GetAccountsByStatus - Returns the accounts by status
func GetAccountsByStatus(status Status, input Reader) (Accounts, error) {
	accounts := Accounts{}
	err := input.GetAccountsByStatus(string(status), accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountsByPrincipalID Get a list of accounts based on Principal ID
func GetAccountsByPrincipalID(principalID string, input Reader) (Accounts, error) {
	accounts := Accounts{}
	err := input.GetAccountsByPrincipalID(principalID, accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
