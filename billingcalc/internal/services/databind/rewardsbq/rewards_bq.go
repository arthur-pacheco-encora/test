package rewardsbq

import (
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/bigqueryutils"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/slices"
)

type Rewards struct {
	OPERATIONS_ORGANIZATION_NAME                   bigquery.NullString `bigquery:"operations_organization_name"`
	OPERATIONS_ACCOUNT_NAME                        bigquery.NullString `bigquery:"operations_account_name"`
	OPERATIONS_VAULT_NAME                          bigquery.NullString `bigquery:"operations_vault_name"`
	OPERATIONS_CREATED_TIME                        bigquery.NullString `bigquery:"operations_created_time"`
	OPERATIONS_END_TIME                            bigquery.NullString `bigquery:"operations_end_time"`
	OPERATIONS_TYPE                                bigquery.NullString `bigquery:"operations_type"`
	OPERATIONS_ASSET_TYPE                          bigquery.NullString `bigquery:"operations_asset_type"`
	OPERATIONS_TRANSACTION_ID                      bigquery.NullString `bigquery:"operations_transaction_id"`
	OPERATIONS_DESTINATION_ADDRESS                 bigquery.NullString `bigquery:"operations_destination_address"`
	OPERATIONS_SOURCE_ADDRESSES                    bigquery.NullString `bigquery:"operations_source_addresses"`
	OPERATIONS_STATUS                              bigquery.NullString `bigquery:"operations_status"`
	OPERATIONS_TOTAL_ANCHORAGE_REWARD_PART         bigquery.NullString `bigquery:"operations_total_anchorage_reward_part"`
	OPERATIONS_TOTAL_ANCHORAGE_USD_REWARD_PART     bigquery.NullString `bigquery:"operations_total_anchorage_usd_reward_part"`
	OPERATIONS_TOTAL_NON_ANCHORAGE_REWARD_PART     bigquery.NullString `bigquery:"operations_total_non_anchorage_reward_part"`
	OPERATIONS_TOTAL_NON_ANCHORAGE_USD_REWARD_PART bigquery.NullString `bigquery:"operations_total_non_anchorage_usd_reward_part"`
	OPERATIONS_TOTAL_USD_VALUE                     bigquery.NullString `bigquery:"operations_total_usd_value"`
	OPERATIONS_TOTAL_ASSET_QUANTITY                bigquery.NullString `bigquery:"operations_total_asset_quantity"`
	BUSINESS_DAY                                   bigquery.NullString `bigquery:"business_day"`
	OPERATIONS_ACCOUNT_INTERNAL_ID                 bigquery.NullString `bigquery:"operations_client_internal_id"`
}

func GetDataQuery(day_start time.Time, day_end time.Time) string {
	day_start_str := day_start.Format(time.RFC3339)
	day_end_str := day_end.Format(time.RFC3339)

	query := fmt.Sprintf(`
	WITH operations AS (SELECT
		ops.* EXCEPT (account_name, organization_name, vault_name),
		COALESCE(dp.parent_affiliate_id, a.affiliate_id) AS client_internal_id,
		vm.name AS vault_name,
		dest_orgs.org_name AS destination_org_name,
		orgs.org_name AS organization_name,
		a.name AS account_name
	  FROM
		client_operations.operations_confidential ops
	  LEFT JOIN
		client_operations.accounts_confidential a
	  ON
		ops.account_id = a.account_id
	  LEFT JOIN
		kyc_operations.duplicate_affiliates dp
	  ON
		a.affiliate_id = dp.affiliate_id
	  LEFT JOIN
		client_operations.vaults_metadata vm
	  ON
		ops.vault_unique_id = vm.vault_unique_id
	  LEFT JOIN
		client_operations.organizations_confidential orgs
	  ON
		ops.org_id = orgs.org_id
	  LEFT JOIN
		client_operations.organizations_confidential dest_orgs
	  ON
		ops.destination_org_id = dest_orgs.org_id
	  WHERE
		((( ops.end_time ) >= (TIMESTAMP('%s')) AND ( ops.end_time ) < (TIMESTAMP('%s'))))
	)
	SELECT
		operations.organization_name  AS operations_organization_name,
		operations.account_name  AS operations_account_name,
		operations.vault_name  AS operations_vault_name,
			(FORMAT_TIMESTAMP('%%F %%T', operations.created_at )) AS operations_created_time,
			(FORMAT_TIMESTAMP('%%F %%T', operations.end_time )) AS operations_end_time,
		operations.type  AS operations_type,
		shared.ASSET_SYMBOL_MAP(operations.asset_type) AS operations_asset_type,
		operations.transaction_id  AS operations_transaction_id,
		operations.destination_address  AS operations_destination_address,
		operations.source_addresses  AS operations_source_addresses,
		operations.status  AS operations_status,
		operations.client_internal_id AS operations_client_internal_id,
		(FORMAT_TIMESTAMP('%%F', operations.end_time )) AS business_day,
		CAST(COALESCE(SUM(operations.anchorage_reward_part ), 0) as STRING) AS operations_total_anchorage_reward_part,
		CAST(COALESCE(SUM(operations.anchorage_usd_reward_part  ), 0) as STRING) AS operations_total_anchorage_usd_reward_part,
		CAST(COALESCE(SUM(operations.non_anchorage_reward_part ), 0) as STRING) AS operations_total_non_anchorage_reward_part,
		CAST(COALESCE(SUM(operations.non_anchorage_usd_reward_part ), 0) as STRING) AS operations_total_non_anchorage_usd_reward_part,
		CAST(COALESCE(SUM(operations.asset_usd_value ), 0) as STRING) AS operations_total_usd_value,
		CAST(COALESCE(SUM(operations.asset_quantity ), 0) as STRING) AS operations_total_asset_quantity
		FROM operations
	WHERE ((( operations.end_time  ) >= (TIMESTAMP('%s')) AND ( operations.end_time  ) < (TIMESTAMP('%s')))) AND (operations.operation_state ) = 'COMPLETE' AND (operations.type ) IN ('Delegation Reward', 'Staking Reward')
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
		13
	ORDER BY
		1,
		2,
		3,
		4 DESC,
		5,
		6,
		7,
		8,
		9,
		10,
		11,
		12,
		13
	`, day_start_str, day_end_str, day_start_str, day_end_str)

	return query
}

func StructToSlice(iter bigqueryutils.ResultIterator) ([][]string, error) {
	table := [][]string{}

	for {
		bigqueryRow := Rewards{}
		err := iter.Next(&bigqueryRow)
		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			errorMessage := "Error iterating over. " + err.Error()
			return nil, errors.New(errorMessage)
		}

		newRow := []string{}

		newRow = slices.Insert(newRow, int(rewards.ColOrganization), bigqueryRow.OPERATIONS_ORGANIZATION_NAME.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColAccount), bigqueryRow.OPERATIONS_ACCOUNT_NAME.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColOpeType), bigqueryRow.OPERATIONS_TYPE.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColAsset), bigqueryRow.OPERATIONS_ASSET_TYPE.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColAnchAssetQty), bigqueryRow.OPERATIONS_TOTAL_ASSET_QUANTITY.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColAnchValue), bigqueryRow.OPERATIONS_TOTAL_ANCHORAGE_USD_REWARD_PART.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColThirdPtQty), bigqueryRow.OPERATIONS_TOTAL_NON_ANCHORAGE_REWARD_PART.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColThirdPtValue), bigqueryRow.OPERATIONS_TOTAL_NON_ANCHORAGE_USD_REWARD_PART.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColBizDay), bigqueryRow.BUSINESS_DAY.StringVal)
		newRow = slices.Insert(newRow, int(rewards.ColAccInternalID), bigqueryRow.OPERATIONS_ACCOUNT_INTERNAL_ID.StringVal)

		table = append(table, newRow)
	}

	return table, nil
}
