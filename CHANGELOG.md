# Changelog

## 1.3.0 (2026-10-29)

- `graalsystems_user` can create users. It now has `email` (required: the platform refuses a user
  without it), `first_name`, `last_name` and `password` (optional, sensitive, set at creation
  only). Before, it only sent `name` and `description`, and every creation was refused.
- `graalsystems_user.description` is deprecated and ignored: the platform does not store a
  description for users.

Requires a graal-api that accepts users without password or first name (fixed on 2026-10-29); older
versions refuse them.

## 1.2.2 (2026-10-29)

Fixes found by applying the example configuration against a real graal installation.

- Every update failed with HTTP 406: the JSON Patch request asked for `application/json`, and the
  API only answers with the resource's vendor media type
  (`application/vnd.graal.systems.v1.<type>+json`). Creating resources was not affected.
- Updates use the JSON Patch `add` operation, which sets a field whether the API returned it or
  not; `replace` was refused (`INVALID_PATCH`) on fields the API leaves out when empty, such as
  the description of a user.
- `graalsystems_job`: `timeout_seconds` and `max_retries` are computed when omitted. The API sets
  a default (3600 s), and every plan tried to set it back to null.

## 1.2.1 (2026-10-29)

- `graalsystems_workspace`: `type` accepts every workspace type the platform serves again:
  `jupyter`, `vscode`, `rstudio`, `grafana`, `superset`, `metabase` and `zeppelin`. 1.2.0 had
  wrongly limited it to `jupyter`, `vscode` and `rstudio`; that breaking change is withdrawn.

## 1.2.0 (2026-10-29)

### Breaking changes

- `graalsystems_workspace`: `type` accepts `jupyter`, `vscode` and `rstudio` (withdrawn in 1.2.1:
  the platform also serves `grafana`, `superset`, `metabase` and `zeppelin`).
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
- `graalsystems_job`: new `spark` type, `file` and `main_class_name` options; `docker_image` is optional (the installation's runtime image is used),
  `secrets` are sent, more attributes are read back to detect drift.
- `graalsystems_workspace`: changing `type`, `infrastructure_id` or `instance_type` recreates it.
- Dependencies updated (terraform-plugin-sdk v2.40, oauth2, Go 1.26), GoReleaser v2, CI with
  tests and golangci-lint, offline unit tests.
- Removed the `main.tf` that contained real credentials, unused code and the documentation of
  the removed firewall rules.
