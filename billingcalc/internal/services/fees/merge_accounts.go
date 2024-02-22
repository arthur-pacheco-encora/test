package fees

func MergeAccounts(accounts1, accounts2 []AccountResult) []AccountResult {
	mergedAccounts := make(map[string]*AccountResult)

	for _, accounts := range [][]AccountResult{accounts1, accounts2} {
		for _, account := range accounts {
			key := account.AccName + account.BillingTerms + account.CustomerID
			if existingAccount, ok := mergedAccounts[key]; ok {
				existingAccount.Assets = append(existingAccount.Assets, account.Assets...)
			} else {
				accountCopy := account
				mergedAccounts[key] = &accountCopy
			}
		}
	}

	var result []AccountResult
	for _, account := range mergedAccounts {
		result = append(result, *account)
	}
	return result
}
