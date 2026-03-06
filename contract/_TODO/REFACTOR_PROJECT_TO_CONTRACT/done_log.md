# Done Log: REFACTOR_PROJECT_TO_CONTRACT

> Формат записи — строго по шаблону ниже.
> Каждый этап добавляет один блок.

---

<!-- TEMPLATE (не удалять):
## Stage N: <title>
- **DateTime**: YYYY-MM-DD HH:MM
- **Completed**: <что сделано>
- **Artifacts**: <созданные артефакты>
- **Gate**: pass | fail
- **Actual files created**: [список]
- **Actual files modified**: [список]
- **Deviations**: none | <описание>
- **Deviation level**: minor | major
- **Next stage notes**: <заметки для следующего этапа>
- **Commit**: `<service>/ <hash> DONE: REFACTOR_CONTRACT stage N: <title>`
-->

## Stage 1: documan-pkg — Proto schemas + codegen
- **DateTime**: 2026-03-06 04:30
- **Completed**: Обновлены etl.proto (убраны erp_*, file_* из transport), docmatching.proto (amount_with_vat → amount_with_tax), regenerate gen/*. queries.proto не требовал изменений — projection fields уже в целевом состоянии.
- **Artifacts**: обновлённые proto и gen файлы
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [services/proto/core/v1/etl.proto, services/proto/core/v1/docmatching.proto, services/gen/core/v1/etl.pb.go, services/gen/core/v1/docmatching.pb.go]
- **Deviations**: none
- **Deviation level**: —
- **Next stage notes**: documan-core (этап 2) сломается на компиляции — ETL controllers используют удалённые proto-поля ErpUpdatedAt/ErpCreatedAt/ErpDeletedAt/ErpArchived в ETL и FileName/FileExt/FileHash в Ingest. Это ожидаемо.
- **Commit**: `documan-pkg/ 25ed2c0 DONE: REFACTOR_CONTRACT stage 1: proto schemas + codegen`

## Stage 2: documan-core — ETL transport + domain models
- **DateTime**: 2026-03-06 05:15
- **Completed**: Убраны ERPUpdatedAt/CreatedAt/DeletedAt/ERPArchived из ETL ERP transport (controller + model + usecase). Убраны FileName/FileExt/FileHash из ETL Ingest transport. Domain models без изменений (projection fields сохранены). Docmatching controller: AmountWithVat → AmountWithTax (proto field rename). go.mod обновлён на documan-pkg@25ed2c0.
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [go.mod, go.sum, internal/controller/grpc/v1/notify_erp_data.go, internal/controller/grpc/v1/notify_ingest_data.go, internal/controller/grpc/v1/docmatching_find_candidates.go, internal/usecase/etl/model.go, internal/usecase/etl/notify_erp_data.go, internal/usecase/etl/notify_ingest_data.go]
- **Deviations**: tracked — `dl_candidates.go:205` JSON key `"amount_with_vat"` и `features.go` Go field `AmountWithVat` остаются legacy (за пределами ETL/domain scope). Обновить на этапах 3/4.
- **Deviation level**: tracked (auto-escalation to major at stage 11)
- **Next stage notes**: Этап 3 (core consumers) должен переименовать AmountWithVat → AmountWithTax в features.go, dl_candidates.go и всех потребителях. SQL миграция не создавалась — schema columns не менялись (payload JSON content changes, не column names).
- **Commit**: `documan-core/ 6e2f7ec DONE: REFACTOR_CONTRACT stage 2: ETL transport + domain models`

## Stage 3: documan-core — Payload consumers
- **DateTime**: 2026-03-06 05:45
- **Completed**: Features/MatchResult: AmountWithVat → AmountWithTax, Store → StoreName. dl_candidates.go: JSON keys обновлены (amount_with_tax, store_name), добавлено чтение payload_meta для fallback source_number/source_date, удалён enrichFeatures (seller_id/buyer_id fallback через erp_entities). ERP SQL: убран connection_id из SELECT. Controller mapping обновлён. resolve_field.go и dl_delta.go — без изменений (generic, не хардкодят legacy names).
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [internal/domain/docmatching/features.go, internal/adapter/postgres/dl_candidates.go, internal/controller/grpc/v1/docmatching_find_candidates.go]
- **Deviations**: none (tracked deviation из stage 2 закрыта)
- **Deviation level**: —
- **Next stage notes**: Stage 4 (core query) — query controllers и read adapters. documan-core полностью обновлён после stages 2-4.
- **Commit**: `documan-core/ 8261295 DONE: REFACTOR_CONTRACT stage 3: payload consumers`
