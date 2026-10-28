# Changelog

## 1.1.0 (unreleased)

### Breaking changes

- `graalsystems_workspace`: `type` accepts `jupyter`, `vscode` and `rstudio`, the workspaces
  offered by the platform (`metabase`, `superset` and `zeppelin` were removed).
- `graalsystems_group`: the API cannot edit groups, so changing `name` or `description` now
  recreates the group instead of silently doing nothing.

### Fixes

- Updates now send valid JSON Patch documents (the `op` member was missing) and apply every
  changed attribute: previously only `name` was updated and other changes were ignored.
- `graalsystems_job` and `graalsystems_workflow` no longer crash when `schedule` is omitted.
- `graalsystems_user` was unusable (the schema and the code used different attribute names).
- Data sources can look resources up by `name`; without an id they used to call the API
  with an empty id.
- `max_retries = 0` is accepted.
- Authentication errors are reported instead of failing later with an empty client, and
  tokens keep refreshing after the provider is configured.
- Labels updates no longer build wrong JSON Pointer paths.

### Improvements

- Provider settings can come from `GRAAL_API_URL`, `GRAAL_KEYCLOAK_URL`, `GRAAL_TENANT`,
  `GRAAL_CLIENT_ID`, `GRAAL_CLIENT_SECRET`, `GRAAL_USERNAME`, `GRAAL_PASSWORD` (the `GS_*`
  variables still work). `password` and `application_secret` are marked sensitive.
- `auth_mode` is inferred from the credentials that are set; new `realm` argument; the realm
  falls back to the tenant when it cannot be discovered.
- `graalsystems_job`: `docker_image` is optional (the installation's runtime image is used),
  `secrets` are sent, more attributes are read back to detect drift.
- `graalsystems_workspace`: changing `type`, `infrastructure_id` or `instance_type` recreates it.
- Dependencies updated (terraform-plugin-sdk v2.40, oauth2, Go 1.26), GoReleaser v2, CI with
  tests and golangci-lint, offline unit tests.
- Removed the `main.tf` that contained real credentials, unused code and the documentation of
  the removed firewall rules.
