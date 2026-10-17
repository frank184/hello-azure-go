# GitHub Actions deployment

This repository uses GitHub Actions to deploy container images to Azure Container Apps.

## Environments

| Environment | Azure resource group | Container App |
|---|---|---|
| `dev` | `tg-hello-azure-dev` | `hello-azure-go-dev` |

## GitHub environment configuration

Open: **Repository → Settings → Environments**

Create an environment using the same name as the deployment target, such as `dev`.

Environment protection rules determine which branches, tags, and reviewers may deploy there.

## Required environment variables

Each GitHub Environment must define these variables:

| Variable | Purpose |
|---|---|
| `AZURE_CLIENT_ID` | Entra application client ID used by GitHub OIDC |
| `AZURE_TENANT_ID` | Entra tenant containing the deployment identity |
| `AZURE_SUBSCRIPTION_ID` | Azure subscription containing the target resources |
| `AZURE_RESOURCE_GROUP` | Resource group for the target environment |
| `ACR_NAME` | Azure Container Registry name |
| `CONTAINER_APP_NAME` | Azure Container App name |

These values are non-secret configuration, so they should be added as **Environment variables**, not secrets.

## OIDC requirements

Each environment requires its own Entra federated identity credential.

The credential must match the GitHub subject for that environment:

`repo:OWNER@ORGANIZATION_ID/REPOSITORY@REPOSITORY_ID:environment:ENVIRONMENT_NAME`

Use:

- Issuer: `https://token.actions.githubusercontent.com`
- Audience: `api://AzureADTokenExchange`

The deployment identity should receive only the Azure RBAC roles required for its environment.

The current deployment roles are:

- `Reader` on the environment resource group
- `AcrPush` on the environment registry
- `Container Apps Contributor` on the environment Container App

## Deployment workflow

The CD workflow selects a GitHub Environment and reads its variables through the `vars` context.

The deployment job must bind to the selected environment:

`environment: ${{ inputs.environment }}`

Manual deployments should deploy the selected branch or commit. Automatic deployments should explicitly select the intended environment and require successful CI first.

## Image versioning

Images use the GitHub commit SHA as the deployment tag:

`<registry>.azurecr.io/<repository>:<commit-sha>`

The bootstrap `dev` tag was used only to create the initial Container App. Normal deployments should use immutable commit-based tags. Future promotion between environments should reuse the tested image digest rather than rebuild the image.

## Deployment checklist

For each environment:

1. Create the GitHub Environment.
2. Add all required environment variables.
3. Create the matching Entra federated identity credential.
4. Assign scoped Azure RBAC roles.
5. Confirm the resource group, registry, and Container App exist.
6. Confirm CI passed for the source commit.
7. Confirm the image tag or digest exists in ACR.