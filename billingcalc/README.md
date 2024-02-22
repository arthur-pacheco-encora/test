# Description

This package is not part of the main Anchorage services.
It is a standalone project designed for Cloud Run and, in conjunction with Retool on the frontend, it serves as a tool to automate billing calculations for the Anchorage's Revenue/Accounting team.

# Other files related to this module

- `source/infrastructure/terraform/billingcalc/` - GCP deployment resources
- `source/go/lib/protobufs/billingcalcproto/` - Retool will pick up the service contract defined there for gRPC calls

# Utils

This section describe useful commands, intended to be run from the module root.

To build all the binaries: `go build -o bin/ ./...`

To build the docker image: `docker build -t billingcalc:latest .`

## Randomizer

The utility command cmd/randomizer can be used to randomize names in a CSV file.

Run `go run ./cmd/randomizer -h` for usage info.

# Testing locally

```
# First terminal:
go run ./cmd/server

# Second terminal:
curl -X POST http://localhost:8080/fees-csv \
    -H "Content-Type: multipart/form-data" \
    -F "mfr=@internal/services/static/gsheet/mfr_test_calc.csv" \
    -F "rewards=@internal/services/static/gsheet/rewards.csv" \
    -F "unclaimed=@internal/services/static/gsheet/unclaimed.csv" \
    -F "firstExternalId=1" \
    -F "invoiceDate=2023-06-15" \
    -F "debug=true"
```

# Deployment guide

Before getting started, it is important to understand the rationale behind the steps presented here.

The user logged in the Retool app must have a way to tell who he/she is to be able to get IAM access to the service.
This is where IAP comes in. It can tell the user identity through some OAuth checks.

In theory, as we are using terraform, the deployment shouldn't require manual intervention.
But GCP has some limitations when it comes to programmatically configuring IAP.

When an IAP resource is created via terraform, GCP does not allow specifying its OAuth redirect URIs,
even if tried manually in GCP console after the resource has been defined:

- https://registry.terraform.io/providers/hashicorp/google/5.11.0/docs/resources/iap_client
- https://stackoverflow.com/a/77614981

## One-time setup step-by-step

The initial setup requires some manual intervention, but after it's done there is no need to worry again.

1. Comment the IAP related definitions in the `sources/infrastructure/terraform/billingcalc/ilb.tf`:

```
resource "google_compute_region_backend_service" "default" {
  ...
  #iap {
  #  oauth2_client_id     = data.google_iap_client.billingcalc.client_id
  #  oauth2_client_secret = data.google_iap_client.billingcalc.secret
  #}
}

#data "google_iap_client" "billingcalc" {
#  brand     = var.iap_brand
#  client_id = var.iap_client_id
#}
```

2. Run `tf apply`. It must succeed.

3. Enable IAP by clicking the toggle button in the "billingcalc" entry in the IAP settings page `Security -> Identity-Aware Proxy`

4. Go to `APIs & Services -> Credentials` and look for `IAP-billingcalc` entry. Open it and use the "Client ID" value as the value for `iap_client_id` terraform variable.

5. In the same page as item (6), add `https://retool.anchorage-{environment}.com/oauth/user/oauthcallback` as an authorized callback URI, replacing `{environment}` with the proper environment name e.g. development

6. Start a cloud shell session and run `gcloud alpha iap oauth-brands list`. Use the `name` key value as the value for the `iap_brand` terraform variable.

7. Uncomment the IAP sections from item (1) and run `tf apply` again. It must succeed.

## Get IAP OAuth Brand info for terraform data sources
```
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/iap_client
# https://cloud.google.com/sdk/gcloud/reference/alpha/iap
# start a new cloud shell and run:
gcloud alpha iap oauth-brands list
gcloud alpha iap oauth-clients list BRAND
```

## Configuring Retool to connect to the service

HEADER: `Authorization`: `Bearer OAUTH2_ID_TOKEN`

Use self-signed certificates: Skip CA Certificate verification

Authentication: OAuth 2.0

authorization url: https://accounts.google.com/o/oauth2/v2/auth

access token url: https://www.googleapis.com/oauth2/v4/token

`client id` and `client secret`: Get from `APIs & Services -> Credentials` and look for `IAP-billingcalc` entry.

Scopes required: `https://www.googleapis.com/auth/userinfo.email openid`
Advanced > Access token lifespan:  `3600`

Reference: https://community.retool.com/t/how-do-i-give-retool-access-to-a-cloud-functions-endpoint-that-requires-authentication/2613/21

## Retool Drive integration

For resources which require drive integration, create an API Key under `APIs & Services -> Credentials`. Use the same client_id and client secret as "IAP-billingcalc".
