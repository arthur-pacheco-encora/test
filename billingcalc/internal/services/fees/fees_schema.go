package fees

import "github.com/shopspring/decimal"

// Contains the info about the Asset for Staking Calculations
type CalcTableEntry struct {
	Asset      string   `json:"asset"`
	Validator  string   `json:"validator"`
	Operations []string `json:"operations"`
	On_chain   bool     `json:"on_chain"`
	Claimable  bool     `json:"claimable"`
}

type CalcTable []CalcTableEntry

type StakingOutput struct {
	ServiceType             string          `json:"serviceType"`
	Asset                   string          `json:"asset"`
	Amount                  decimal.Decimal `json:"amount"`
	CollectedOnChainAlready bool            `json:"collectedOnChainAlready"`
	EarnedRewards           decimal.Decimal `json:"earnedRewards"`
	FeeRates                decimal.Decimal `json:"feeRates"`
	ItemCategory            string          `json:"itemCategory"`
	ItemDescription         string          `json:"itemDescription"`
	ItemQuantity            string          `json:"itemQuantity"`
	Memo                    string          `json:"memo"`
	MonthlyRate             string          `json:"monthlyRate"`
}

type OrgResult struct {
	OrgName  string          `json:"orgName"`
	Accounts []AccountResult `json:"accounts"`
}

type AccountResult struct {
	AccName       string          `json:"clientName"`
	BillingTerms  string          `json:"billingTerms"`
	CustomerID    string          `json:"customerID"`
	DisplayName   string          `json:"displayName"`
	EntityID      string          `json:"entityID"`
	InvoiceNumber string          `json:"invoiceNumber"`
	ExternalID    string          `json:"externalID"`
	InvoiceDate   string          `json:"invoiceDate"`
	DueDate       string          `json:"dueDate"`
	Assets        []StakingOutput `json:"assets"`
}

type Warning struct {
	OrgName     string `json:"orgName"`
	AccName     string `json:"accName"`
	Asset       string `json:"asset"`
	Description string `json:"description"`
}

type StakingSummary []OrgResult

type CalculatedFees struct {
	Summary StakingSummary
	Warns   []Warning
}
