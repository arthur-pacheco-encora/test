package operationsstatusesbq

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/bigqueryutils"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/slices"
)

type OperationsStatuses struct {
	DELEGATION_STATUSES_ACCOUNT_NAME                bigquery.NullString `bigquery:"delegation_statuses_account_name"`
	DELEGATION_STATUSES_ACTIVE_DELEGATED_QUANTITY   bigquery.NullString `bigquery:"delegation_statuses_active_delegated_quantity"`
	DELEGATION_STATUSES_ACTIVE_DELEGATED_VALUE_USD  bigquery.NullString `bigquery:"delegation_statuses_active_delegated_value_usd"`
	DELEGATION_STATUSES_ANCHORAGE_ENTITY            bigquery.NullString `bigquery:"delegation_statuses_anchorage_entity"`
	DELEGATION_STATUSES_ASSET_SUB_ID                bigquery.NullString `bigquery:"delegation_statuses_asset_sub_id"`
	DELEGATION_STATUSES_ASSET_TYPE                  bigquery.NullString `bigquery:"delegation_statuses_asset_type"`
	DELEGATION_STATUSES_ASSET_UNIQUE_ID             bigquery.NullString `bigquery:"delegation_statuses_asset_unique_id"`
	DELEGATION_STATUSES_DATE_DATE                   bigquery.NullString `bigquery:"delegation_statuses_date_date"`
	DELEGATION_STATUSES_IS_ANCHORAGE_VALIDATOR      bigquery.NullString `bigquery:"delegation_statuses_is_anchorage_validator"`
	DELEGATION_STATUSES_LAST_UPDATED_DATE           bigquery.NullString `bigquery:"delegation_statuses_last_updated_date"`
	DELEGATION_STATUSES_ORG_ID                      bigquery.NullString `bigquery:"delegation_statuses_org_id"`
	DELEGATION_STATUSES_ORG_NAME                    bigquery.NullString `bigquery:"delegation_statuses_org_name"`
	DELEGATION_STATUSES_PENDING_DELEGATED_QUANTITY  bigquery.NullString `bigquery:"delegation_statuses_pending_delegated_quantity"`
	DELEGATION_STATUSES_PENDING_DELEGATED_VALUE_USD bigquery.NullString `bigquery:"delegation_statuses_pending_delegated_value_usd"`
	DELEGATION_STATUSES_ACCOUNT_ID                  bigquery.NullString `bigquery:"delegation_statuses_account_id"`
	DELEGATION_STATUSES_UNBONDED_QUANTITY           bigquery.NullString `bigquery:"delegation_statuses_unbonded_quantity"`
	DELEGATION_STATUSES_UNBONDED_VALUE_USD          bigquery.NullString `bigquery:"delegation_statuses_unbonded_value_usd"`
	DELEGATION_STATUSES_UNIT_PRICE_USD              bigquery.NullString `bigquery:"delegation_statuses_unit_price_usd"`
	DELEGATION_STATUSES_VALIDATOR_ADDRESS           bigquery.NullString `bigquery:"delegation_statuses_validator_address"`
	DELEGATION_STATUSES_VALIDATOR_ID                bigquery.NullString `bigquery:"delegation_statuses_validator_id"`
	VAULTS_METADATA_VAULT_NAME                      bigquery.NullString `bigquery:"vaults_metadata_vault_name"`
	DELEGATION_STATUSES_VAULT_SUB_ID                bigquery.NullString `bigquery:"edit"`
	DELEGATION_STATUSES_VAULT_UNIQUE_ID             bigquery.NullString `bigquery:"delegation_statuses_vault_unique_id"`
	DELEGATION_STATUSES_WALLET_UNIQUE_ID            bigquery.NullString `bigquery:"delegation_statuses_wallet_unique_id"`
}

func GetDataQuery(business_day_start time.Time, business_day_end time.Time) string {
	business_day_start_str := business_day_start.Format(time.RFC3339)
	business_day_end_str := business_day_end.Format(time.RFC3339)

	query := fmt.Sprintf(`
	WITH delegation_statuses AS (SELECT
        hds.account_id,
        active_delegated_quantity,
        active_delegated_value_usd,
        hds.anchorage_entity,
        asset_sub_id,
        asset_type,
        asset_unique_id,
        date,
        is_anchorage_validator,
        last_updated_at,
        hds.org_id,
        pending_delegated_quantity,
        pending_delegated_value_usd,
        unbonded_quantity,
        unbonded_value_usd,
        unit_price_usd,
        validator_address,
        validator_id,
        hds.vault_sub_id,
        hds.vault_unique_id,
        wallet_unique_id,
        accounts.name as account_name,
        orgs.org_name
      FROM
        client_operations.historic_delegation_statuses_confidential AS hds
      LEFT JOIN
        client_operations.accounts_confidential AS accounts
        ON
          hds.account_id = accounts.account_id
      LEFT JOIN
        client_operations.organizations_confidential as orgs
        ON
          hds.org_id = orgs.org_id
    )
  	,vaults_metadata AS (SELECT
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
	SELECT
		delegation_statuses.account_name AS delegation_statuses_account_name,
		CAST(delegation_statuses.active_delegated_quantity AS STRING) AS delegation_statuses_active_delegated_quantity,
		CAST(delegation_statuses.active_delegated_value_usd AS STRING) AS delegation_statuses_active_delegated_value_usd,
		delegation_statuses.anchorage_entity AS delegation_statuses_anchorage_entity,
		delegation_statuses.asset_sub_id AS delegation_statuses_asset_sub_id,
		shared.ASSET_SYMBOL_MAP(delegation_statuses.asset_type) AS delegation_statuses_asset_type,
		delegation_statuses.asset_unique_id AS delegation_statuses_asset_unique_id,
			CAST(delegation_statuses.date as STRING ) AS delegation_statuses_date_date,
			(CASE WHEN delegation_statuses.is_anchorage_validator  THEN 'Yes' ELSE 'No' END) AS delegation_statuses_is_anchorage_validator,
			CAST(DATE(delegation_statuses.last_updated_at) as STRING) AS delegation_statuses_last_updated_date,
		delegation_statuses.org_id AS delegation_statuses_org_id,
		delegation_statuses.org_name AS delegation_statuses_org_name,
		CAST(delegation_statuses.pending_delegated_quantity as STRING) AS delegation_statuses_pending_delegated_quantity,
		CAST(delegation_statuses.pending_delegated_value_usd as STRING) AS delegation_statuses_pending_delegated_value_usd,
		delegation_statuses.account_id AS delegation_statuses_account_id,
		CAST(delegation_statuses.unbonded_quantity as STRING) AS delegation_statuses_unbonded_quantity,
		CAST(delegation_statuses.unbonded_value_usd as STRING) AS delegation_statuses_unbonded_value_usd,
		CAST(delegation_statuses.unit_price_usd as STRING) AS delegation_statuses_unit_price_usd,
		delegation_statuses.validator_address AS delegation_statuses_validator_address,
		delegation_statuses.validator_id AS delegation_statuses_validator_id,
		vaults_metadata.name AS vaults_metadata_vault_name,
		delegation_statuses.vault_sub_id AS delegation_statuses_vault_sub_id,
		delegation_statuses.vault_unique_id AS delegation_statuses_vault_unique_id,
		delegation_statuses.wallet_unique_id AS delegation_statuses_wallet_unique_id
	FROM delegation_statuses
	LEFT JOIN vaults_metadata ON delegation_statuses.vault_unique_id = vaults_metadata.vault_unique_id
	WHERE ((( delegation_statuses.date  ) >= (DATE('%s')) AND ( delegation_statuses.date  ) < (DATE('%s'))))
	GROUP BY
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
		13,
		14,
		15,
		16,
		17,
		18,
		19,
		20,
		21,
		22,
		23,
		24
	ORDER BY
		8 DESC
	`, business_day_start_str, business_day_end_str)

	return query
}

func StructToSlice(iter bigqueryutils.ResultIterator) ([][]string, error) {
	table := [][]string{}

	for {
		bigqueryRow := OperationsStatuses{}
		err := iter.Next(&bigqueryRow)
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			errorMessage := "Error iterating over. " + err.Error()
			return nil, errors.New(errorMessage)
		}

		newRow := []string{}
		newRow = slices.Insert(newRow, int(operationsstatuses.ColStatusesAccountName), bigqueryRow.DELEGATION_STATUSES_ACCOUNT_NAME.StringVal)
		newRow = slices.Insert(newRow, int(operationsstatuses.ColStatusesActiveDelegatedValue), bigqueryRow.DELEGATION_STATUSES_ACTIVE_DELEGATED_VALUE_USD.StringVal)
		newRow = slices.Insert(newRow, int(operationsstatuses.ColStatusesAssetType), bigqueryRow.DELEGATION_STATUSES_ASSET_TYPE.StringVal)
		newRow = slices.Insert(newRow, int(operationsstatuses.ColStatusesDate), bigqueryRow.DELEGATION_STATUSES_DATE_DATE.StringVal)
		newRow = slices.Insert(newRow, int(operationsstatuses.ColCosmosValidatorsRate), getExternalValidatorPercentage(bigqueryRow.DELEGATION_STATUSES_IS_ANCHORAGE_VALIDATOR.StringVal))
		table = append(table, newRow)
	}

	return table, nil
}

func getExternalValidatorPercentage(isAnchorageValidator string) string {
	var cosmosValidatorsRate string

	if isAnchorageValidator == "Yes" {
		cosmosValidatorsRate = "0"
	} else {
		cosmosValidatorsRate = "100.00%"
	}
	return cosmosValidatorsRate
}
