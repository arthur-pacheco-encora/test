package getcalcfromcsv

var CsvMfrString = `,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
Asset ID lookup,,,Org ID from RDB,Account ID from RDB,MSA ID & Asset ID,Anchorage Entity,Anchorage Entity ID,Org Name,Entity Legal Name,Netsuite Customer Internal ID,Billing Terms,"Agreement Term, renewal/termination conditions - RH 8/26",Status,Preparer,Prepared,Reviewer,Reviewed?,Reviewer sign off for June 2021 updates,TEMP - needs June 2021 review,Notes,Email Contact,Mailing address,Asset Type,Agreement Countersignature Date or Effective Date,Original Agreement Date,Max Delay of Fee Date,Onboarding Date,Fee Accrual Date,Minimum Charge (Customer Level),1st Tier Floor,1st Tier Rate,2nd Tier Floor,2nd Tier Rate,3rd Tier Floor,3rd Tier Rate,4th Tier Floor,4th Tier Rate,5th Tier Floor,5th Tier Rate,6th Tier Floor,6th Tier Rate,7th Tier Floor,7th Tier Rate,8th Tier Floor,8th Tier Rate,9th Tier Floor,9th Tier Rate,10th Tier Floor,10th Tier Rate,Artificial Max Floor (used for if formulas),Artificial Max Tier Rate (used for if formulas),XTZ Fee % - third party validator,Celo Fee % - Anchorage validator,Celo Fee % - third party validator,Celo credit against custody fee terms,FLOW Fee % - Anchorage validator,FLOW Fee % - third party validator,OSMO Fee % - third party validator,OSMO Staking Fee % - 100% commission validator (fee on staked balance),ROSE Fee % - Anchorage validator,ROSE Fee % - third party validator,ETH Fee % - Anchorage validator,AXL Staking Fee % - 100% commission validator (fee on staked balance),AXL Fee % - third party validator,APT Fee % - third party validator,ATOM Staking Fee % - 100% commission validator (fee on staked balance),ATOM Fee % - third party validator,HASH Staking Fee % - 100% commission validator (fee on staked balance),HASH Fee % - third party validator,EVMOS Staking Fee % - 100% commission validator (fee on staked balance),EVMOS Fee % - third party validator,APT Fee % - Anchorage validator,SOL Fee % - third party validator,SOL Staking Fee % - 100% commission validator (fee on staked balance),SUI Fee % - Anchorage validator,SUI Fee % - third party validator,SUI Staking Fee % - 100% commission validator (fee on staked balance),Staking Notes,Brokerage Client? (yes/no),Brokerage Fee,Calculation Notes,Cosmetic Notes,Notes (from prior to 1/1/2021),Reviewer Notes (from prior to 1/1/2021),Asset ID,MSA ID
BTC/ETH,0,,,,,Anchorage Digital Bank,15,Electric Sheep,ES Capital Fund,12890,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,"ALL except BTC/ETH, NFTs",,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,10%,10%,Anchorage,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,yes,Standard; varies trade by trade stated on each order,,,,,,10
BTC,1,,,,,Anchorage Digital Bank,15,Electric Sheep,ES Venture Fund,36189,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,BTC/ETH,,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,10%,10%,Anchorage,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,yes,Standard; varies trade by trade stated on each order,,,,,,0
ETH,2,,,,,Anchorage Singapore,33,Firetronics,FT Capital Fund,12567,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,"ALL except BTC/ETH, NFTs",,,,,,,,,,,,,,,,,,,,,,,,,,,,,1%,10%,1%,Anchorage,8%,1%,1%,1%,12%,1%,10%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,12%,3%,1%,6%,3%,1%,,Yes,Standard; varies trade by trade stated on each order,,,,,,10
DYDX,3,,,,,Anchorage Sinagapore,33,Firetronics,FT Venture Fund,14756,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,BTC/ETH,,,,,,,,,,,,,,,,,,,,,,,,,,,,,1%,10%,1%,Anchorage,8%,1%,1%,1%,12%,1%,10%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,12%,3%,1%,6%,3%,1%,,Yes,Standard; varies trade by trade stated on each order,,,,,,0
ROSE,4,,,,,Anchorage Digital Bank,15,Moon Enterprises,ME Capital Fund,23451,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,"ALL except BTC/ETH, NFTs",,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,6%,6%,Any,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,No,NA,,,,,,10
FIL,5,,,,,Anchorage Singapore,33,Moon Enterprises,ME Venture Fund,41962,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,BTC/ETH,,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,6%,6%,Any,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,No,NA,,,,,,0
`

var CsvInvalidMfrString = `,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
Asset ID lookup,,,Org ID from RDB,Account ID from RDB,MSA ID & Asset ID,Anchorage Entity,Anchorage Entity ID,Org Name,Entity Legal Name,Netsuite Customer Internal ID,Billing Terms,"Agreement Term, renewal/termination conditions - RH 8/26",Status,Preparer,Prepared,Reviewer,Reviewed?,Reviewer sign off for June 2021 updates,TEMP - needs June 2021 review,Notes,Email Contact,Mailing address,Asset Type,Agreement Countersignature Date or Effective Date,Original Agreement Date,Max Delay of Fee Date,Onboarding Date,Fee Accrual Date,Minimum Charge (Customer Level),1st Tier Floor,1st Tier Rate,2nd Tier Floor,2nd Tier Rate,3rd Tier Floor,3rd Tier Rate,4th Tier Floor,4th Tier Rate,5th Tier Floor,5th Tier Rate,6th Tier Floor,6th Tier Rate,7th Tier Floor,7th Tier Rate,8th Tier Floor,8th Tier Rate,9th Tier Floor,9th Tier Rate,10th Tier Floor,10th Tier Rate,Artificial Max Floor (used for if formulas),Artificial Max Tier Rate (used for if formulas),XTZ Fee % - third party validator,Celo Fee % - Anchorage validator,Celo Fee % - third party validator,Celo credit against custody fee terms,FLOW Fee % - Anchorage validator,FLOW Fee % - third party validator,OSMO Fee % - third party validator,OSMO Staking Fee % - 100% commission validator (fee on staked balance),ROSE Fee % - Anchorage validator,ROSE Fee % - third party validator,ETH Fee % - Anchorage validator,AXL Staking Fee % - 100% commission validator (fee on staked balance),AXL Fee % - third party validator,APT Fee % - third party validator,ATOM Staking Fee % - 100% commission validator (fee on staked balance),ATOM Fee % - third party validator,HASH Staking Fee % - 100% commission validator (fee on staked balance),HASH Fee % - third party validator,EVMOS Staking Fee % - 100% commission validator (fee on staked balance),EVMOS Fee % - third party validator,APT Fee % - Anchorage validator,SOL Fee % - third party validator,SOL Staking Fee % - 100% commission validator (fee on staked balance),SUI Fee % - Anchorage validator,SUI Fee % - third party validator,SUI Staking Fee % - 100% commission validator (fee on staked balance),Staking Notes,Brokerage Client? (yes/no),Brokerage Fee,Calculation Notes,Cosmetic Notes,Notes (from prior to 1/1/2021),Reviewer Notes (from prior to 1/1/2021),Asset ID,MSA ID
BTC/ETH,0,,,,,Anchorage Digital Bank,15,,ES Capital Fund,12890,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,"ALL except BTC/ETH, NFTs",,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,10%,10%,Anchorage,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,yes,Standard; varies trade by trade stated on each order,,,,,,10
BTC,1,,,,,Anchorage Digital Bank,15,,ES Venture Fund,36189,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,BTC/ETH,,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,10%,10%,Anchorage,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,yes,Standard; varies trade by trade stated on each order,,,,,,0
ETH,2,,,,,Anchorage Singapore,33,Firetronics,FT Capital Fund,12567,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,"ALL except BTC/ETH, NFTs",,,,,,,,,,,,,,,,,,,,,,,,,,,,,1%,10%,1%,Anchorage,8%,1%,1%,1%,12%,1%,10%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,12%,3%,1%,6%,3%,1%,,Yes,Standard; varies trade by trade stated on each order,,,,,,10
DYDX,3,,,,,Anchorage Sinagapore,33,Firetronics,FT Venture Fund,14756,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,BTC/ETH,,,,,,,,,,,,,,,,,,,,,,,,,,,,,1%,10%,1%,Anchorage,8%,1%,1%,1%,12%,1%,10%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,1.00%,12%,3%,1%,6%,3%,1%,,Yes,Standard; varies trade by trade stated on each order,,,,,,0
ROSE,4,,,,,Anchorage Digital Bank,15,Moon Enterprises,ME Capital Fund,23451,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,"ALL except BTC/ETH, NFTs",,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,6%,6%,Any,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,No,NA,,,,,,10
FIL,5,,,,,Anchorage Singapore,33,Moon Enterprises,ME Venture Fund,41962,Net 15,"1-year, auto-renewing, 60 days written notice required to terminate",Active,"Julie
MG (Amendment 2)","Yes
7/21/2022","Kennedy
CK (Amendment 2)","Yes
8/1/2022",,,"Fees calculated using the incremental rate at each tier, not using one rate for the whole AUC balance",,,BTC/ETH,,,,,,,,,,,,,,,,,,,,,,,,,,,,,3%,6%,6%,Any,8%,3%,3%,1%,12%,8%,10%,1.00%,3.00%,7.00%,1.00%,3.00%,1.00%,3.00%,1.00%,3.00%,12%,3%,1%,6%,3%,1%,,No,NA,,,,,,0
`

var CsvRewardsString = `,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
Operations Organization Name,Operations Account Name,Operations Vault Name,Operations Created Time,Operations End Time,Operations Type,Operations Asset Type,Operations Transaction Hash,Operations Destination Address,Operations Source Addresses,Operations Status,Operations Total Anchorage Reward Part,Operations Total Anchorage USD Reward Part,Operations Total Non Anchorage Reward Part,Operations Total Non Anchorage USD Reward Part,Operations Total USD Value,Operations Total Asset Quantity,Business day
Org Test Alpha,Test Alpha Account,Primary,2023-06-22 4:36:00,2023-06-22 4:36:00,Delegation Reward,HASH,Operations_Transaction_Hash_1,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,"1,134.45",7.94,7.94,"1,134.45",6/22/2023
Org Test Alpha,Test Alpha Account,Primary,2023-06-16 17:23:28,2023-06-16 17:23:28,Delegation Reward,HASH,Operations_Transaction_Hash_2,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,"2,356.52",16.50,16.50,"2,356.52",6/16/2023
Org Test Alpha,Test Alpha Account,Primary,2023-06-08 2:41:20,2023-06-08 2:41:20,Delegation Reward,HASH,Operations_Transaction_Hash_3,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,"1,616.13",19.39,19.39,"1,616.13",6/8/2023
Org Test Alpha,Test Alpha Account,Primary,2023-06-01 17:29:01,2023-06-01 17:29:01,Delegation Reward,HASH,Operations_Transaction_Hash_4,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,189.66,1.90,1.90,189.66,6/1/2023
Org Test Alpha,Test Alpha Account,Primary,2023-06-01 2:12:56,2023-06-01 2:12:56,Delegation Reward,HASH,Operations_Transaction_Hash_5,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,993.65,9.94,9.94,993.65,6/1/2023`

var CsvInvalidRewardsString = `,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
,,,,,,,,,,,,,,,,,
Operations Organization Name,Operations Account Name,Operations Vault Name,Operations Created Time,Operations End Time,Operations Type,Operations Asset Type,Operations Transaction Hash,Operations Destination Address,Operations Source Addresses,Operations Status,Operations Total Anchorage Reward Part,Operations Total Anchorage USD Reward Part,Operations Total Non Anchorage Reward Part,Operations Total Non Anchorage USD Reward Part,Operations Total USD Value,Operations Total Asset Quantity,Business day
,Primary,2023-06-22 4:36:00,2023-06-22 4:36:00,Delegation Reward,HASH,Operations_Transaction_Hash_1,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,"1,134.45",7.94,7.94,"1,134.45",6/22/2023
,Primary,2023-06-16 17:23:28,2023-06-16 17:23:28,Delegation Reward,HASH,Operations_Transaction_Hash_2,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,"2,356.52",16.50,16.50,"2,356.52",6/16/2023
,Primary,2023-06-08 2:41:20,2023-06-08 2:41:20,Delegation Reward,HASH,Operations_Transaction_Hash_3,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,"1,616.13",19.39,19.39,"1,616.13",6/8/2023
,Primary,2023-06-01 17:29:01,2023-06-01 17:29:01,Delegation Reward,HASH,Operations_Transaction_Hash_4,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,189.66,1.90,1.90,189.66,6/1/2023
,Primary,2023-06-01 2:12:56,2023-06-01 2:12:56,Delegation Reward,HASH,Operations_Transaction_Hash_5,Operations_Destination_Address_Alpha_1,Operations_Source_Addresses_Alpha_1,COMPLETE,0.00,0.00,993.65,9.94,9.94,993.65,6/1/2023`

var CsvUnclaimedBalancesString = `,Historic Daily Blockchain Balances Business Day Date,Historic Daily Blockchain Balances Address,Historic Daily Blockchain Balances Asset Type ID,Historic Daily Blockchain Balances Organization Key ID,Historic Daily Blockchain Balances Balance Type,Historic Daily Blockchain Balances Balance Str,Historic Daily Blockchain Balances Last Updated At Time,Account name,USD Price,USD Value
45107AXLTest Beta Account,2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_1,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585
45107AXLTest Beta Account,2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_2,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585
45107AXLTest Beta Account,2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_3,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,250274.5703,2023-06-30 23:56:39,Test Beta Account,0.3895,97481.94513
45107AXLTest Beta Account,2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_4,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,164161.7311,2023-06-30 23:56:39,Test Beta Account,0.3895,63940.99425
45107AXLTest Beta Account,2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_5,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585
45107AXLTest Beta Account,2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_6,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585`

var CsvInvalidUnclaimedBalancesString = `,Historic Daily Blockchain Balances Business Day Date,Historic Daily Blockchain Balances Address,Historic Daily Blockchain Balances Asset Type ID,Historic Daily Blockchain Balances Organization Key ID,Historic Daily Blockchain Balances Balance Type,Historic Daily Blockchain Balances Balance Str,Historic Daily Blockchain Balances Last Updated At Time,Account name,USD Price,USD Value
2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_1,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585
2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_2,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585
2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_3,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,250274.5703,2023-06-30 23:56:39,Test Beta Account,0.3895,97481.94513
2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_4,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,164161.7311,2023-06-30 23:56:39,Test Beta Account,0.3895,63940.99425
2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_5,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585
2023-06-30,Historic_Daily_Blockchain_Balances_Address_Beta_6,AXL,Historic_Daily_Blockchain_Balances_Organization_Key_ID_Beta,DELEGATION_REWARDS,0.147223,2023-06-30 23:56:39,Test Beta Account,0.3895,0.0573433585`
