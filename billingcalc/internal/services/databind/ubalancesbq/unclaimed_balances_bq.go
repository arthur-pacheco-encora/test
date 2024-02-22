package ubalancesbq

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/bigqueryutils"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/slices"
)

type UnclaimedBalances struct {
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BUSINESS_DAY_DATE    bigquery.NullDate   `bigquery:"historic_daily_blockchain_balances_business_day_date"`
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_ADDRESS              bigquery.NullString `bigquery:"historic_daily_blockchain_balances_address"`
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_ASSET_TYPE_ID        bigquery.NullString `bigquery:"historic_daily_blockchain_balances_asset_type_id"`
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_ORGANIZATION_KEY_ID  bigquery.NullString `bigquery:"historic_daily_blockchain_balances_organization_key_id"`
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BALANCE_TYPE         bigquery.NullString `bigquery:"historic_daily_blockchain_balances_balance_type"`
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BALANCE_STR          bigquery.NullString `bigquery:"historic_daily_blockchain_balances_balance_str"`
	HISTORIC_DAILY_BLOCKCHAIN_BALANCES_LAST_UPDATED_AT_TIME bigquery.NullString `bigquery:"historic_daily_blockchain_balances_last_updated_at_time"`
	ORGANIZATIONS_ORG_NAME                                  bigquery.NullString `bigquery:"organizations_org_name"`
	CUSTODY_ACCOUNTS_ACCOUNT_NAME                           bigquery.NullString `bigquery:"custody_accounts_account_name"`
	USD_PRICE                                               bigquery.NullString `bigquery:"usd_price"`
	USD_VALUE                                               bigquery.NullString `bigquery:"usd_value"`
}

func GetDataQuery(business_day_start time.Time, business_day_end time.Time) string {
	business_day_start_str := business_day_start.Format(time.RFC3339)
	business_day_end_str := business_day_end.Format(time.RFC3339)

	query := fmt.Sprintf(`
	WITH historic_daily_blockchain_balances AS (SELECT
		address,
		asset_type_id,
		organization_key_id,
		balance_str,
		balance_type,
		business_day,
		last_updated_at
	FROM
		adb_operations.historic_daily_blockchain_balances
	WHERE
		business_day BETWEEN DATE('%s') AND DATE('%s')
		AND (balance_type = 'DELEGATION_REWARDS')
	)
	,  vaults_metadata AS (SELECT
		vm.deposit_required,
		vm.system_account_id,
		a.anchorage_entity,
		vm.vault_unique_id,
		vm.withdraw_required,
		COALESCE(dp.parent_affiliate_id, a.affiliate_id) AS client_internal_id,
		vm.name,
		vm.organization_key_id,
		vm.vault_sub_id
		FROM
		client_operations.vaults_metadata vm
		LEFT JOIN
		client_operations.accounts_confidential a
		ON
		vm.system_account_id = a.account_id
		LEFT JOIN
		kyc_operations.duplicate_affiliates dp
		ON
		a.affiliate_id = dp.affiliate_id
	)
	, daily_balances AS (SELECT
		db.*,
		ydb.unit_price_usd AS prev_day_unit_price_usd,
		ydb.quantity AS prev_day_quantity,
		COALESCE(dp.parent_affiliate_id, a.affiliate_id) AS client_internal_id
	  FROM
		client_operations.daily_balances_confidential AS db
	  LEFT JOIN
		client_operations.accounts_confidential AS a
	  ON
		db.account_id = a.account_id
	  LEFT JOIN
		kyc_operations.duplicate_affiliates AS dp
	  ON
		a.affiliate_id = dp.affiliate_id
	  LEFT JOIN
		  client_operations.daily_balances_confidential AS ydb
	ON
	  db.asset_unique_id = ydb.asset_unique_id AND db.date = DATE_ADD(ydb.date, INTERVAL 1 DAY) AND COALESCE(db.token_id, "") = COALESCE(ydb.token_id, "")
	  )
	SELECT
		(historic_daily_blockchain_balances.business_day ) AS historic_daily_blockchain_balances_business_day_date,
		historic_daily_blockchain_balances.address  AS historic_daily_blockchain_balances_address,
		shared.ASSET_SYMBOL_MAP(historic_daily_blockchain_balances.asset_type_id)  AS historic_daily_blockchain_balances_asset_type_id,
		historic_daily_blockchain_balances.organization_key_id  AS historic_daily_blockchain_balances_organization_key_id,
		historic_daily_blockchain_balances.balance_type  AS historic_daily_blockchain_balances_balance_type,
		historic_daily_blockchain_balances.balance_str  AS historic_daily_blockchain_balances_balance_str,
			(FORMAT_TIMESTAMP('%%F %%T', historic_daily_blockchain_balances.last_updated_at )) AS historic_daily_blockchain_balances_last_updated_at_time,
		organizations.org_name  AS organizations_org_name,
		custody_accounts.name  AS custody_accounts_account_name
		,CAST(AVG(daily_balances.unit_price_usd) as STRING) AS usd_price
		,CAST(ROUND((CAST(historic_daily_blockchain_balances.balance_str AS FLOAT64) * AVG(daily_balances.unit_price_usd)), 8)  as STRING) AS usd_value
	FROM historic_daily_blockchain_balances
	LEFT JOIN client_operations.address_details
		AS addresses ON historic_daily_blockchain_balances.address = addresses.address
	LEFT JOIN vaults_metadata ON (
		addresses.organization_key_id = vaults_metadata.organization_key_id
		AND addresses.vault_sub_id = vaults_metadata.vault_sub_id
	)
	LEFT JOIN client_operations.custody_accounts  AS custody_accounts ON (CONCAT(custody_accounts.affiliate_id, "-", LOWER		(custody_accounts.anchorage_entity))) = (CONCAT(vaults_metadata.client_internal_id, "-", LOWER(vaults_metadata.anchorage_entity)))
	LEFT JOIN client_operations.organizations_confidential  AS organizations ON historic_daily_blockchain_balances.organization_key_id = organizations.org_id
	-- new
	LEFT JOIN daily_balances ON historic_daily_blockchain_balances.business_day = daily_balances.date
      AND shared.ASSET_SYMBOL_MAP(historic_daily_blockchain_balances.asset_type_id) = shared.ASSET_SYMBOL_MAP(daily_balances.asset_type)
	-- new
	WHERE
	((historic_daily_blockchain_balances.balance_str ) <> '0' OR (historic_daily_blockchain_balances.balance_str ) IS NULL) AND
	(historic_daily_blockchain_balances.balance_type ) = 'DELEGATION_REWARDS' AND
	((( historic_daily_blockchain_balances.business_day  ) >= (DATE('%s')) AND ( historic_daily_blockchain_balances.business_day  ) < (DATE('%s'))))
	GROUP BY
	1
	,2
	,3
	,4
	,5
	,6
	,7
	,8
	,9
	ORDER BY
    1 DESC
	`, business_day_start_str, business_day_end_str, business_day_start_str, business_day_end_str)

	return query
}

func StructToSlice(iter bigqueryutils.ResultIterator) ([][]string, error) {
	table := [][]string{}

	for {
		bigqueryRow := UnclaimedBalances{}
		err := iter.Next(&bigqueryRow)
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			errorMessage := "Error iterating over. " + err.Error()
			return nil, errors.New(errorMessage)
		}

		newRow := []string{}
		newRow = slices.Insert(newRow, int(ubalances.ColAsset), bigqueryRow.HISTORIC_DAILY_BLOCKCHAIN_BALANCES_ASSET_TYPE_ID.StringVal)
		newRow = slices.Insert(newRow, int(ubalances.ColBalanceType), bigqueryRow.HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BALANCE_TYPE.StringVal)
		newRow = slices.Insert(newRow, int(ubalances.ColDailyBalanceStr), bigqueryRow.HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BALANCE_STR.StringVal)
		newRow = slices.Insert(newRow, int(ubalances.ColDailyBalanceDate), bigqueryRow.HISTORIC_DAILY_BLOCKCHAIN_BALANCES_LAST_UPDATED_AT_TIME.String())
		newRow = slices.Insert(newRow, int(ubalances.ColAccountName), bigqueryRow.CUSTODY_ACCOUNTS_ACCOUNT_NAME.String())
		newRow = slices.Insert(newRow, int(ubalances.ColUsdPrice), bigqueryRow.USD_PRICE.String())
		newRow = slices.Insert(newRow, int(ubalances.ColUsdValue), bigqueryRow.USD_VALUE.String())

		table = append(table, newRow)
	}

	return table, nil
}
