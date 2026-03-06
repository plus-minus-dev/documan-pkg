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
- **DateTime**: 2026-03-06 05:45 (initial), 2026-03-06 ~07:00 (rework)
- **Completed**: Features/MatchResult: AmountWithVat → AmountWithTax, Store → StoreName. dl_candidates.go: JSON keys обновлены (amount_with_tax, store_name), добавлено чтение payload_meta для fallback source_number/source_date. enrichFeatures восстановлен с enrichment по flat ref-модели: seller_id/buyer_id/store_id из payload_header → resolve через erp_entities. Добавлен store_name enrichment (новый). connection_id возвращён в ERP SQL SELECT. Controller mapping обновлён. resolve_field.go и dl_delta.go — без изменений (generic, не хардкодят legacy names).
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [internal/domain/docmatching/features.go, internal/adapter/postgres/dl_candidates.go, internal/controller/grpc/v1/docmatching_find_candidates.go]
- **Deviations**: major (resolved) — первоначальный коммит (8261295) ошибочно удалял enrichFeatures. Contract пересмотрен: принята flat ref-модель (seller_id/buyer_id/store_id в Header, item_id в Position). Исправлено в rework-коммите 886d702.
- **Deviation level**: major (resolved)
- **Next stage notes**: Stage 4 (core query) — query controllers и read adapters. documan-core полностью обновлён после stages 2-4.
- **Commit**: `documan-core/ 8261295 DONE: REFACTOR_CONTRACT stage 3: payload consumers` (initial, enrichFeatures удалён), `documan-core/ 886d702 DONE: REFACTOR_CONTRACT stage 3: restore enrichFeatures + store_name enrichment` (rework, актуальный)

## Stage 4: documan-core — Query API controllers
- **DateTime**: 2026-03-06 ~08:00
- **Completed**: Query controllers и proto уже на contract vocabulary (grep audit чист). Добавлен EnrichERPDocPayload в postgres adapter: обогащает PayloadHeader (seller_id → inn/name/kpp/address, buyer_id → inn/name/kpp/address, store_id → name) и PayloadPositions (item_id → name/article/variant) через erp_entities batch resolve. Метод добавлен в query.Postgres interface, вызывается в GetERPDocument usecase. ListERPDocuments не требует enrichment (не возвращает payload content).
- **Artifacts**: enrich_payload.go (новый файл)
- **Gate**: pass
- **Actual files created**: [internal/adapter/postgres/enrich_payload.go]
- **Actual files modified**: [internal/usecase/query/usecase.go, internal/usecase/query/get_erp_document.go]
- **Deviations**: none
- **Deviation level**: —
- **Next stage notes**: documan-core полностью обновлён (stages 2-4). Stage 5 — documan-connector-ms.
- **Commit**: `documan-core/ f503e69 DONE: REFACTOR_CONTRACT stage 4: query API controllers`

## Stage 5: documan-connector-ms — Document payload
- **DateTime**: 2026-03-06 ~10:00
- **Completed**: Document payload переведён с `map[string]any` на `contract/document` structs. buildDocMeta → document.Meta (source_account_id, source_number, source_date, parse_rule, source_status, payloaded_at). buildDocHeader → document.Header (ref-only: seller_id/buyer_id/store_id, без business enrichment). buildDocSummary → document.Summary (kopecks *int64, lines, sku). buildDocPositions → []document.Position (item_id, item_type, line_number, unit_price_with_tax, tax_amount, discount_rate). Удалены: kopecksToRubles(), setIfStr(), legacy field names (source_name, source_date в header, assortment_id, assortment_type, line_no, unit_price_with_vat, vat_amount). Adapter/core: убраны ErpArchived/ErpUpdatedAt/ErpCreatedAt/ErpDeletedAt из ERPDataItem (теперь в payload_meta). DocTransformVersion bumped 1.0 → 2.0.
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [internal/usecase/transform/doc.go, internal/usecase/transform/helpers.go, internal/usecase/ports.go, internal/adapter/core/client.go, internal/controller/worker/delivery.go, go.mod, go.sum]
- **Deviations**: none
- **Deviation level**: —
- **Next stage notes**: entity.go transform (справочники) пока остаётся на map[string]any — scope stage 5 = document payload only. Stage 6 — connector-ms entity payload / remaining adapters.
- **Commit**: `documan-connector-ms/ b3fff66 DONE: REFACTOR_CONTRACT stage 5: document payload → contract structs`

## Stage 6: documan-connector-ms — Entity payload + terminology rename
- **DateTime**: 2026-03-06 ~12:00
- **Completed**: Entity payload переведён с `map[string]any` + `fieldMapping` на `contract/entity` structs (Party, Store, Item с embedded entity.Meta). GTIN отделён от barcodes по regex `^\d{8,14}$`. Variant извлекается из MS characteristics. Полный terminology rename: transform→payload, clear→payload, reclear→repayload, S3 path `transformed/`→`payload/`. Package `internal/usecase/transform/` → `internal/usecase/payload/`. Удалены: `marshalSparse()`, `entityFieldMaps`, `fieldMapping`, `buildEntityClearPayload()`. EntityPayloadVersion bumped 1.0 → 2.0.
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [internal/usecase/payload/entity.go, internal/usecase/payload/payload.go, internal/usecase/payload/repayload.go, internal/usecase/payload/ports.go, internal/usecase/payload/helpers.go, internal/usecase/payload/gzip.go, internal/usecase/payload/doc.go, internal/usecase/payload.go, internal/usecase/usecase.go, internal/usecase/model.go, internal/usecase/ports.go, internal/adapter/storage/keys.go, internal/adapter/storage/key_builder.go, internal/controller/worker/payload_worker.go, internal/controller/worker/delivery.go, internal/app/app.go, config/config.go, compose.yaml, compose.override.yaml]
- **Deviations**: none
- **Deviation level**: —
- **Next stage notes**: Stage 7 — documan-ingest. Сервис сломается на build из-за proto changes (stage 1).
- **Commit**: `documan-connector-ms/ 3a083a9 DONE: REFACTOR_CONTRACT stage 6: entity payload → contract structs + terminology rename`

## Stage 7: documan-ingest — Parser → contract document structs
- **DateTime**: 2026-03-06 ~14:00
- **Completed**: Parser ClearOutput переведён с `map[string]any` на `contract/document` structs (Meta, Header, Summary, Position). BuildClear() строит типизированные структуры: Meta (source, parse_rule, parse_error, payloaded_at), Header (прямой маппинг полей), Summary (kopecks *int64, Lines, SKU, Qty decimal string), Positions (field renames: vat→TaxAmount, amount_with_vat→AmountWithTax, line_no→LineLabel, commodity_code→HSCode, unit→UOMName, unit_code→UOMCode, price→UnitPrice, price_with_vat→UnitPriceWithTax). Delivery enrichment: enrichMeta() использует document.Meta struct, добавляет S3Bucket/S3Key/FileName/FileExt/FileHash. gRPC adapter: убраны FileName/FileExt/FileHash из proto mapping. S3 key rename: *ClearKey→*PayloadKey, transformed/→payload/. CoreDeliveryItem: убраны FileName/FileExt/FileHash (теперь в payload_meta JSON). go.mod обновлён на documan-pkg@acde77acd872.
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [go.mod, go.sum, internal/adapter/parser/parser.go, internal/adapter/parser/bridge.go, internal/adapter/core/client.go, internal/adapter/storage/keys.go, internal/adapter/storage/key_builder.go, internal/usecase/types.go, internal/usecase/ports.go, internal/usecase/adapters.go, internal/usecase/usecase.go, internal/usecase/transform/transform.go, internal/usecase/transform/model.go, internal/usecase/transform/ports.go, internal/usecase/delivery/delivery.go, internal/usecase/delivery/model.go, internal/usecase/delivery/ports.go, internal/domain/models.go]
- **Deviations**: none — parser internal extraction layer (extract_*.go, sheet.go, sterile.go) сохраняет legacy field names как внутреннее промежуточное представление; маппинг на contract names происходит в BuildClear()
- **Deviation level**: —
- **Next stage notes**: Stage 8 — documan-bff. Parser internal field names (extract layer) — за пределами scope, можно обновить отдельным рефакторингом.
- **Commit**: `documan-ingest/ 756140a DONE: REFACTOR_CONTRACT stage 7: parser → contract document structs`

## Stage 8: documan-ingest — Delivery + rename clear→payload
- **DateTime**: 2026-03-06 ~15:00
- **Completed**: Полный терминологический rename: ClearOutput→PayloadOutput, BuildClear→BuildPayload, BuildClearFromRawJSON→BuildPayloadFromRawJSON, persistClearArtifacts→persistPayloadArtifacts, clearOutput→payloadOutput, Reclear→Repayload, ReclearCmd→RepayloadCmd, ReclearInbox→RepayloadInbox, reclearDocument→repayloadDocument, ListDocumentsForReclear→ListDocumentsForRepayload, clearedAt→payloadedAt (Go vars only, DB column cleared_at сохранён). HTTP route /internal/inbox/reclear→/internal/inbox/repayload. claude.md обновлён. Delivery уже обновлён на stage 7 (enrichMeta с document.Meta struct).
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [internal/adapter/parser/parser.go, internal/adapter/parser/bridge.go, internal/adapter/parser/sheet.go, internal/adapter/postgres/outbox.go, internal/usecase/types.go, internal/usecase/ports.go, internal/usecase/adapters.go, internal/usecase/usecase.go, internal/usecase/transform/model.go, internal/usecase/transform/ports.go, internal/usecase/transform/transform.go, internal/controller/http/router.go, internal/controller/http/v1/handlers.go, internal/controller/worker/transform.go, claude.md]
- **Deviations**: none — DB column `cleared_at` сохранён (rename требует SQL миграции, за пределами scope)
- **Deviation level**: —
- **Next stage notes**: Stage 9 — documan-bff (или следующий по master_plan). documan-ingest полностью обновлён.
- **Commit**: `documan-ingest/ 0bd7dd1 DONE: REFACTOR_CONTRACT stage 8: delivery + rename clear→payload`

## Stage 9: documan-bff — DTO + query mapping + matching
- **DateTime**: 2026-03-06 ~15:30
- **Completed**: go.mod обновлён на documan-pkg@acde77acd872. MatchResultResponse: AmountWithVat→AmountWithTax (Go field + JSON tag). docmatching.go handler: mr.AmountWithVat→mr.AmountWithTax. BFF проксирует payload как json.RawMessage — contract structs не нужны. Query/details handlers уже используют корректные proto getters.
- **Artifacts**: нет новых файлов
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [go.mod, go.sum, internal/dto/types.go, internal/controller/http/v1/docmatching.go]
- **Deviations**: none
- **Deviation level**: —
- **Next stage notes**: Stage 10 — documan-fe-my (frontend). documan-bff полностью обновлён.
- **Commit**: `documan-bff/ 7480247 DONE: REFACTOR_CONTRACT stage 9: BFF DTO + query mapping + matching`

## Stage 10: documan-fe-my — UI migration to contract vocabulary
- **DateTime**: 2026-03-06 ~16:00
- **Completed**: Полный rename payload field names в Vue/TS: amount_with_vat→amount_with_tax, vat_amount→tax_amount, vat→tax_amount (summary), line_no→line_number, commodity_code→hs_code, unit_code→uom_code, unit_name→uom_name, unit_price_with_vat→unit_price_with_tax, source_name→source_number, amountWithVat→amountWithTax. ERP money semantics: formatRubles(float)→kopecksToRubKop(int) — теперь ERP тоже хранит деньги в копейках. reclearInbox→repayloadInbox, route /reclear→/repayload.
- **Artifacts**: нет новых файлов
- **Gate**: pass (typecheck errors pre-existing, не вызваны нашими изменениями)
- **Actual files created**: нет
- **Actual files modified**: [app/composables/useDocworker.ts, app/composables/useDocMatching.ts, app/components/matching/DeltaViewer.vue, app/components/DocumentLinksBlock.vue, app/pages/documents/index.vue, app/pages/documents/new/[id].vue, app/pages/documents/new2/[id].vue, app/pages/documents_erp/[entityType]/[id].vue]
- **Deviations**: minor — typecheck имеет pre-existing errors (PairCard.vue emit types, login.vue/registration.vue string|undefined, nuxt.config.ts process/nitro) — все существовали до наших изменений, подтверждено stash test
- **Deviation level**: minor (pre-existing, не вызвано миграцией)
- **Next stage notes**: Stage 11 — verification. Все сервисы обновлены.
- **Commit**: `documan-fe-my/ dae3cd8 DONE: REFACTOR_CONTRACT stage 10: UI migration to contract vocabulary`

## Stage 11: Verification + cutover preparation
- **DateTime**: 2026-03-06 ~17:00
- **Completed**: Cross-repo build — все 5 Go сервисов pass. FE typecheck — pre-existing errors only. Grep-аудит legacy payload fields — чисто. Grep-аудит terminology rename — чисто. Flat ref-context audit — все references корректны. Bruno — нет legacy references. Config/env legacy keys — чисто.
- **Artifacts**: нет
- **Gate**: pass
- **Deviations**: none
- **Deviation level**: —
- **Cutover checklist**:
  - [x] Все сервисы собираются (5/5 Go, FE pre-existing only)
  - [x] Config/env ключи обновлены
  - [x] S3 paths = `payload/`
  - [x] Repayload workflow (connector-ms + ingest)
  - [x] Bruno clean
  - [ ] SQL migration cleared_at → payloaded_at (отложена)
  - [ ] Proto/gen синхронизация (buf generate, требует buf CLI)
- **Commit**: verification only, no code changes
