# Мастер-план: Перевод проекта на финальный contract

> Источник истины (ТЗ): `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md`
> Дата создания плана: 2026-03-06
> Общее количество этапов: 11

## Обзор

Пайплайн покрывает все затронутые сервисы в порядке зависимостей:
`pkg → core → connector-ms → ingest → bff → fe-my → bruno → verification`

---

## Этап 1: documan-pkg — Contract docs + proto schemas + codegen
- **Файл**: `stage_01_pkg_proto.md`
- **Сервис**: `documan-pkg`
- **Цель**: Обновить contract docs (flat ref fields), `etl.proto`, `queries.proto`, `docmatching.proto`. Regenerate `services/gen/*`.
- **Gate**: `go build ./...` из `documan-pkg/`
- **Зависимости**: нет (первый этап)

## Этап 2: documan-core — ETL transport + domain models + SQL migration
- **Файл**: `stage_02_core_etl_domain.md`
- **Сервис**: `documan-core`
- **Цель**: Привести ETL transport, domain/persistence модели и SQL schema к новым proto и flat ref-based ERP document model (`payload_header.seller_id/buyer_id/store_id`, `payload_positions[*].item_id`).
- **Gate**: `go build ./...` из `documan-core/`
- **Зависимости**: этап 1

## Этап 3: documan-core — Payload consumers
- **Файл**: `stage_03_core_consumers.md`
- **Сервис**: `documan-core`
- **Цель**: Переписать `dl_candidates.go`, `resolve_field.go`, `docmatching`, `dl_delta.go` под новый contract vocabulary и конкретный enrichment по flat ids в `payload_header` / `payload_positions`.
- **Gate**: `go build ./...` из `documan-core/`
- **Зависимости**: этап 2

## Этап 4: documan-core — Query API controllers
- **Файл**: `stage_04_core_query.md`
- **Сервис**: `documan-core`
- **Цель**: Обновить query controllers и DTO/proto mapping под новые `queries.proto`, включая ERP enrichment по flat ids в `payload_header` / `payload_positions`.
- **Gate**: `go build ./...` из `documan-core/`
- **Зависимости**: этап 2

## Этап 5: documan-connector-ms — Document payload
- **Файл**: `stage_05_connector_doc.md`
- **Сервис**: `documan-connector-ms`
- **Цель**: Перевести `transform/doc.go` с `map[string]any` на `contract/document` structs и писать ERP flat ref fields в `payload_header` / `payload_positions` вместо snapshot business fields.
- **Gate**: `go build ./...` из `documan-connector-ms/`
- **Зависимости**: этап 1

## Этап 6: documan-connector-ms — Entity payload + delivery + rename
- **Файл**: `stage_06_connector_entity.md`
- **Сервис**: `documan-connector-ms`
- **Цель**: Перевести `transform/entity.go` на `contract/entity` structs. Обновить delivery. Rename `transform/clear → payload`, `reclear → repayload`.
- **Gate**: `go build ./...` из `documan-connector-ms/`
- **Зависимости**: этап 5

## Этап 7: documan-ingest — Parser/payload build
- **Файл**: `stage_07_ingest_parser.md`
- **Сервис**: `documan-ingest`
- **Цель**: Перевести `parser.go` с legacy clear-building на `contract/document` structs; `ingest` обычно оставляет flat ref fields пустыми.
- **Gate**: `go build ./...` из `documan-ingest/`
- **Зависимости**: этап 1

## Этап 8: documan-ingest — Delivery + rename
- **Файл**: `stage_08_ingest_delivery.md`
- **Сервис**: `documan-ingest`
- **Цель**: Обновить delivery. Rename `clear → payload`, `reclear → repayload`. Обновить S3 paths `transformed/ → payload/`.
- **Gate**: `go build ./...` из `documan-ingest/`
- **Зависимости**: этап 7

## Этап 9: documan-bff — DTO + query mapping + matching
- **Файл**: `stage_09_bff.md`
- **Сервис**: `documan-bff`
- **Цель**: Обновить DTO, query mapping и matching responses под новые proto; BFF полагается на уже обогащённые ответы `core` для ERP seller/buyer/store.
- **Gate**: `go build ./...` из `documan-bff/`
- **Зависимости**: этапы 4, 3

## Этап 10: documan-fe-my — UI migration
- **Файл**: `stage_10_fe.md`
- **Сервис**: `documan-fe-my`
- **Цель**: Обновить composables, document/ERP pages, matching UI под новый contract vocabulary и semantics; ERP UI не должен сам резолвить refs.
- **Gate**: `npx nuxi typecheck` из `documan-fe-my/`
- **Зависимости**: этап 9

## Этап 11: Verification + Bruno + cutover prep
- **Файл**: `stage_11_verification.md`
- **Сервисы**: все
- **Цель**: Cross-repo build/typecheck, обновить Bruno коллекции, grep-аудит legacy references, подготовка cutover.
- **Gate**: build всех сервисов + grep-аудит
- **Зависимости**: этапы 1–10

---

## Граф зависимостей

```
[1: pkg-proto]
    ├── [2: core-etl-domain]
    │       ├── [3: core-consumers]
    │       │       └── [9: bff] ──→ [10: fe-my]
    │       └── [4: core-query]
    │               └── [9: bff]
    ├── [5: connector-doc]
    │       └── [6: connector-entity]
    └── [7: ingest-parser]
            └── [8: ingest-delivery]
                                        └── [11: verification]
```

## Правила выполнения

### Абсолютные запреты (действуют на ВСЕХ этапах)
- **Без обратной совместимости**: запрещены dual-read, dual-write, fallback на legacy payload fields, временное сосуществование old/new payload schema.
- **Payload только через contract structs**: `map[string]any` и ad-hoc JSON maps запрещены как финальное состояние producers.
- **Все legacy payload field names должны быть удалены**: оставлять их «временно» нельзя.
- **Запрещены**: bridge-адаптеры, compatibility wrappers, alias-поля, deprecated endpoints, дублирующие env/config keys, временные rename-слои.
- **Proto changes — АТОМАРНО**: schema + regenerate + обновление producer/consumer кода в одном проходе. Состояние «proto обновлён, gen ещё старый» не допускается.
- **Rename `transform/clear → payload` — ПОЛНЫЙ**: без legacy route names, worker names, config/env names, S3 prefixes, лог-сообщений. Включая тесты.
- **Любое отклонение от этих правил** требует явного согласования с пользователем, а не локального компромисса.

### Deviation policy
- Обязательные правила из ТЗ (раздел "Обязательные правила выполнения") действуют на ВСЕХ этапах.
- Код — источник истины. Перед изменением файла — прочитай его.
- Каждый этап завершается gate check, done-log записью и коммитом.
- `minor` deviation — фиксируй и продолжай. `major` — СТОП и эскалация.
- **Автоэскалация**: любой `minor` deviation, связанный с legacy payload field names, `map[string]any`, незавершённым rename или оставшимися legacy references, автоматически становится `major` если не закрыт к gate check этапа 11.
- Любые остаточные legacy payload references в коде, proto, DTO, UI, **тестах** и grep-аудите к финалу миграции = `major`.
