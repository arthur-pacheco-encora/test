package rewards

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

const (
	ColOrganization  databind.Column = 0
	ColAccount       databind.Column = 1
	ColOpeType       databind.Column = 5
	ColAsset         databind.Column = 6
	ColAnchAssetQty  databind.Column = 11
	ColAnchValue     databind.Column = 12
	ColThirdPtQty    databind.Column = 13
	ColThirdPtValue  databind.Column = 14
	ColBizDay        databind.Column = 17
	ColAccInternalID databind.Column = 18
)

type Rewards struct {
	organizations map[string]Organization
}

type Organization struct {
	Name        string
	DisplayName string
	accounts    map[string]Account
}

type Account struct {
	Id          databind.AccountID
	Name        string
	DisplayName string
	assets      map[string]Asset
}

type Asset struct {
	Name           string
	claimedRewards []ClaimedReward
}

type ClaimedReward struct {
	AnchorageAssetQty  decimal.Decimal
	AnchorageUsdValue  decimal.Decimal
	BusinessDay        time.Time
	OperationType      string
	ThirdPartyAssetQty decimal.Decimal
	ThirdPartyUsdValue decimal.Decimal
}

func (r *Rewards) GetOrganizationNames() map[string]Organization {
	return r.organizations
}

func (r *Rewards) GetAccountNames(organization string) map[string]Account {
	org := r.organizations[organization]
	return org.accounts
}

func (r *Rewards) GetAccountById(accountId databind.AccountID) Account {
	for _, o := range r.organizations {
		for _, a := range o.accounts {
			if a.Id == accountId {
				return a
			}
		}
	}
	return Account{}
}

func (a *Account) GetAssets() map[string]Asset {
	return a.assets
}

func (a *Asset) GetClaimedRewards() []ClaimedReward {
	return a.claimedRewards
}

func (r *Rewards) IsEmpty() bool {
	return len(r.organizations) == 0
}

func NewRewards(table [][]string) *Rewards {
	rewards := &Rewards{
		organizations: make(map[string]Organization),
	}

	for _, row := range table {
		organizationName := sanitization.SanitizeName(row[ColOrganization])
		org, exists := rewards.organizations[organizationName]
		if !exists {
			org = Organization{
				Name:        organizationName,
				DisplayName: row[ColOrganization],
				accounts:    make(map[string]Account),
			}
		}

		accountId := row[ColAccInternalID]
		acc, accExists := org.accounts[accountId]
		if !accExists {
			acc = Account{
				Id:          databind.AccountID(accountId),
				Name:        sanitization.SanitizeName(row[ColAccount]),
				DisplayName: row[ColAccount],
				assets:      make(map[string]Asset),
			}
		}

		assetName := row[ColAsset]
		asset, exists := acc.assets[assetName]
		if !exists {
			asset = Asset{
				Name:           assetName,
				claimedRewards: make([]ClaimedReward, 0),
			}
		}

		claimedReward := ClaimedReward{
			AnchorageAssetQty:  converter.FromStrToDecimal(row[ColAnchAssetQty]),
			AnchorageUsdValue:  converter.FromStrToDecimal(row[ColAnchValue]),
			BusinessDay:        converter.FromStrToTimeStamp(row[ColBizDay]),
			OperationType:      row[ColOpeType],
			ThirdPartyAssetQty: converter.FromStrToDecimal(row[ColThirdPtQty]),
			ThirdPartyUsdValue: converter.FromStrToDecimal(row[ColThirdPtValue]),
		}

		asset.claimedRewards = append(asset.claimedRewards, claimedReward)
		acc.assets[assetName] = asset
		org.accounts[accountId] = acc
		rewards.organizations[organizationName] = org
	}

	return rewards
}
