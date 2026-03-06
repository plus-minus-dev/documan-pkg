# План перевода всего проекта на финальный `contract` относительно текущего кода

> Дата: 2026-03-06
> Основано на фактическом коде репозиториев `documan-*`, а не только на целевом контракте.
> Источник истины по целевой схеме: `documan-pkg/contract/document/*`, `documan-pkg/contract/entity/*`, `documan-pkg/contract/README.md`.

## Обязательные правила выполнения

- Без обратной совместимости: запрещены dual-read, dual-write, fallback на legacy payload fields и временное сосуществование old/new payload schema.
- Payload собирается только через contract structs. `map[string]any` и ad-hoc JSON maps не допускаются как финальное состояние producers.
- Любые legacy payload field names (`amount_with_vat`, `line_no`, `source_name`, `source_date` в header, `commodity_code`, `unit_price_with_vat`, `legalTitle` и т.п.) должны быть удалены из payload. Оставлять их «временно» нельзя.
- `seller_id`, `buyer_id`, `store_id` в `payload_header` и `item_id` в `payload_positions` являются допустимыми contract ref fields для ERP/source linkage.
- Projection fields вне payload допустимы. Поля хранения/индексации/query layer (`file_*`, `original_file_*`, `erp_*`, `archived`) не являются частью payload contract и могут сохраняться в БД/DTO отдельно.
- Если projection field дублирует payload/meta по смыслу, payload остаётся источником истины, а projection field считается производным read-model полем.
- Для ERP-документов source of truth для linkage seller/buyer/store/item — flat ref fields: `payload_header.seller_id`, `payload_header.buyer_id`, `payload_header.store_id`, `payload_positions[*].item_id`.
- Бизнес-поля `payload_header.seller_*`, `payload_header.buyer_*`, `payload_header.store_name`, `payload_positions[*].name/article/item_code` не являются reference fields и не должны содержать ids.
- В выбранной модели ERP connectors могут не писать business snapshot в `payload_header` / `payload_positions`; enrichment этих данных выполняет `core` по flat ref fields.
- `.proto` пересобираются с чистого листа под финальную схему. Обратная совместимость protobuf wire format не требуется.
- Сохранять старые номера protobuf-полей не нужно. `reserved` для удалённых legacy-полей не использовать. Нумерацию полей можно выстроить заново под финальную модель сообщений.
- Изменения `.proto` выполняются только атомарно: правка schema + regenerate + обновление producer/consumer кода в том же проходе. Состояние «.proto уже изменён, gen-код ещё старый» не допускается как этап миграции.
- Не оставлять bridge-адаптеры, compatibility wrappers, alias-поля, deprecated endpoints, дублирующие env/config keys и временные rename-слои ради “плавного перехода”. Если старое имя больше не нужно вне payload contract, оно должно быть удалено, а не алиаситься.
- Для rename `transform/clear -> payload` конечное состояние должно быть полным: без legacy route names, worker names, config/env names, S3 prefixes и лог-сообщений в затронутых сервисах, если только старое имя не сохранено осознанно вне scope этой миграции.
- Любые остаточные legacy payload references в коде, proto, DTO, UI, тестах и grep-аудите к финалу миграции считаются `major`, а не `minor`.
- Любое отклонение от этих правил требует явного согласования, а не локального компромисса внутри этапа.

## Исходные условия

- Релиза не было, обратная совместимость не требуется.
- Перевод делаем сразу по всему контуру, который реально зависит от payload / proto / query / DTO.
- Все SQL-изменения упаковываются в одну основную миграцию.
- Задача: привести текущий код проекта к финальному контракту.
- Сопутствующие задачи не делаем молча: либо согласуем, либо выносим в TODO.

## Терминологический rename внутри этой миграции

В рамках перевода `connector-ms` и `ingest` на контракт одновременно делаем механический rename legacy-терминов:

- `transform` -> `payload`
- `clear` -> `payload`
- `reclear` -> `repayload`
- S3 path `transformed/` -> `payload/`

Правила:
- rename делается только вместе с contract migration в тех же затронутых местах, не отдельной предварительной фазой
- логика не меняется, меняется только терминология
- `core.ClearGroupID` не трогаем
- `documan-bff` не затрагивается этим rename напрямую
- существующие test objects в старом S3 path не мигрируем: writers переключаются на `payload/`, старый volume удаляется и тестовые данные пересобираются заново


---

## 1. Что есть в коде сейчас

## 1.1. `documan-connector-ms`

### Document transform (`internal/usecase/transform/doc.go`)

Сейчас документ собирается в `map[string]any`, не в contract structs.

#### Header сейчас пишет legacy-поля
- `doc_number`, `doc_date`, `doc_type`
- `source_date` в header
- `source_name` в header
- `seller_id`
- `buyer_id`
- `store_id`
- `currency_code`, `currency_name`

#### Summary сейчас пишет legacy-поля и в рублях `float`
- `amount_with_vat`
- `vat`
- `amount`

Суммы строятся через `kopecksToRubles()`.

#### Positions сейчас пишут legacy-поля
- `line_no`
- `assortment_id`
- `assortment_type`
- `quantity` как `float`
- `unit_price_with_vat`
- `discount` как процент
- `vat_amount`
- `tracking_type`, `identification_codes`
- `customs_gtd`, `country_name`

Нет:
- `article`
- `item_code`
- `line_label`
- `gtin`, `barcodes`
- `uom_code`, `uom_name`
- `amount`, `amount_with_tax`, `tax_rate`, `tax_amount`
- `hs_code`

#### Meta сейчас пишет legacy-поля
- `source`
- `transformed_at`
- `updated_at`
- `created_at`
- `archived`
- `deleted_at`

Нет:
- `source_account_id`
- `source_number`
- `source_date` в meta
- `source_status`
- `parse_rule`
- `payloaded_at`

### Entity transform (`internal/usecase/transform/entity.go`)

Сейчас entity тоже собираются в `map[string]any` через field maps.

#### Product / Service
Сейчас используются поля:
- `name`
- `code`
- `article`
- `description`
- `barcodes`
- `pathName`

#### Counterparty / Organization
Сейчас используются поля:
- `name`
- `code`
- `inn`
- `kpp`
- `legalTitle`
- `companyType`

#### Store
Сейчас используются поля:
- `name`
- `code`
- `pathName`

Нет contract meta embedding, нет `address`, нет `item_code`, нет `variant`, нет `legal_title`, нет `source_status` и т.д.

### Delivery
Delivery уже возит payload как `bytes`, но transport для ERP по-прежнему несёт отдельные поля:
- `erp_updated_at`
- `erp_created_at`
- `erp_deleted_at`
- `erp_archived`

---

## 1.2. `documan-ingest`

### Parser (`internal/adapter/parser/parser.go`, `sheet.go`, extractors)

Сейчас clear payload строится не через contract structs.

#### Meta сейчас
- `doc_type`
- `format`
- `parse_elapsed_ms`

`parse_elapsed_ms` рассматривается как техническая диагностика пайплайна, а не часть payload schema. При миграции убираем его из payload/meta и оставляем только в логах/наблюдаемости.

#### Header сейчас
Берётся как `map[string]string` из `ParsedDoc.Header`.
В parse-failure `parse_error` кладётся в header.

#### Summary сейчас
- `amount`
- `vat`
- `amount_with_vat`

Денежные поля переводятся в копейки, но legacy names остаются.

#### Positions сейчас
- `line_no`
- `article`
- `name`
- `unit`
- `unit_code`
- `quantity` как `float64`
- `price`
- `price_with_vat`
- `amount`
- `excise`
- `vat_rate`
- `vat`
- `amount_with_vat`
- `commodity_code`
- `country_code`
- `country_name`
- `customs_gtd`

Нет:
- `line_number`, `line_label`
- `item_code`
- `uom_name`, `uom_code`
- `tax_amount`, `tax_rate`, `amount_with_tax`
- `hs_code`
- `discount_rate`, `discount`
- `variant`, `gtin`, `barcodes`

### Delivery (`internal/usecase/delivery/delivery.go`, `model.go`)

Сейчас delivery:
- обогащает `payload_meta` нестандартными полями `original_file` и `file_batch`
- отдельно в transport везёт `file_name`, `file_ext`, `file_hash`
- отдельно везёт `original_file_bucket`, `original_file_key`

---

## 1.3. `documan-core`

### ETL transport (`internal/controller/grpc/v1/notify_*.go`, `internal/usecase/etl/*`)

Сейчас `NotifyERPDataReady` и `NotifyIngestDataReady` принимают transport, завязанный на текущие proto:

#### ERP transport
Отдельно принимает:
- `erp_updated_at`
- `erp_created_at`
- `erp_deleted_at`
- `erp_archived`

#### Ingest transport
Отдельно принимает:
- `file_name`
- `file_ext`
- `file_hash`
- `original_file`

### Domain / storage model

#### `original_documents`
В БД и домене сейчас есть отдельные поля:
- `doc_type`
- `file_name`
- `file_ext`
- `file_hash`
- `original_file_bucket`
- `original_file_key`
- `payload_meta`, `payload_header`, `payload_summary`, `payload_positions`
- `payload_hash`, `payloaded_at`

#### `erp_documents`
В БД и домене сейчас есть отдельные поля:
- `erp_updated_at`
- `erp_created_at`
- `erp_deleted_at`
- `erp_archived`
- `archived`
- `payload_meta`, `payload_header`, `payload_summary`, `payload_positions`

#### `erp_entities`
Аналогично сейчас есть:
- `erp_updated_at`
- `erp_created_at`
- `erp_deleted_at`
- `erp_archived`
- `archived`
- `payload_entity`

### Query API (`queries.proto`, gRPC query controllers)

Сейчас query-layer отдаёт наружу отдельные legacy transport/query поля:
- для original docs: `file_name`, `file_ext`, `file_hash`, `original_file`
- для ERP docs/entities: `erp_updated_at`, `erp_created_at`, `erp_deleted_at`, `erp_archived`, `archived`

### Payload consumers

#### `dl_candidates.go`
Сейчас читает old schema:
- `source_name` как fallback для `doc_number`
- `source_date` из header как fallback для `doc_date`
- `store`
- `amount_with_vat`
- `seller_id` / `buyer_id` с enrichment через `erp_entities`
- quantity суммируется как `float64` из positions
- `SKUCount` считается как `len(positions)`

#### `resolve_field.go`
Сейчас читает поля из `payload_entity` по произвольным string-keys, ориентирован на старые имена вроде `legalTitle`, `code`, `article`.

#### `docmatching`
`docmatching.proto`, `domain/docmatching/features.go`, gRPC mapping используют legacy vocabulary:
- `amount_with_vat`
- `sku_count`
- `quantity_sum`
- `store`

### Delta
`dl_delta.go` generic, но выводит фактические field names, которые сейчас legacy.

---

## 1.4. `documan-pkg`

### `contract/*`
Контракт уже финализирован:
- `article` вместо line/item `sku`
- `item_code`
- `line_number` + `line_label`
- `amount_with_tax`, `tax_amount`
- `hs_code`
- `variant`
- `source_number`, `source_status`, `payloaded_at`
- `document.Header` с flat ref fields `seller_id`, `buyer_id`, `store_id`
- `document.Position` с flat ref field `item_id`
- `entity.Item` использует `article` + `item_code`
- `entity.Party` / `entity.Store` используют `source_code`

### `services/proto/core/v1/etl.proto`
Сейчас ещё legacy transport:
- `ERPDataItem` содержит `erp_*`
- `IngestDataItem` содержит `file_name/file_ext/file_hash`
- `original_file` уже есть

### `services/proto/core/v1/queries.proto`
Сейчас query API ещё legacy:
- original details/items содержат file fields отдельно
- ERP details/items/entities содержат `erp_*`, `archived`

### `services/proto/core/v1/docmatching.proto`
Сейчас matching proto ещё legacy:
- `amount_with_vat`
- `sku_count`

---

## 1.5. `documan-bff`

BFF payload сам не пересобирает, но жёстко зависит от текущего query/proto слоя.

### Что делает сейчас
- листинги и details проксируются из core gRPC в HTTP DTO
- payload JSON отдаётся как `json.RawMessage`
- current DTO содержат отдельные query fields:
  - original: `fileName`, `fileExt`, `fileHash`, `originalFile`
  - ERP: `erpUpdatedAt`, `erpCreatedAt`
- matching DTO содержат `AmountWithVat`

Вывод: BFF затронут не contract structs напрямую, а через изменения `queries.proto` и `docmatching.proto`.

---

## 1.6. `documan-fe-my`

Это основной frontend consumer старой payload-схемы.

### Реально зашитые legacy-поля

#### Original document pages
Файлы:
- `app/composables/useDocworker.ts`
- `app/pages/documents/[id].vue`
- `app/pages/documents/new/[id].vue`
- `app/pages/documents/new2/[id].vue`

Сейчас используются:
- `source_name`
- `line_no`
- `article`
- `commodity_code`
- `amount_with_vat`
- `vat`, `vat_rate`, `vat_amount`
- `unit_name`, `unit_price_with_vat`

#### ERP document page
Файл:
- `app/pages/documents_erp/[entityType]/[id].vue`

Сейчас дополнительно используются:
- `seller_id`, `buyer_id`
- `source_name`
- допущение, что ERP summary хранится в рублях, а не в копейках

#### Matching UI
Файл:
- `app/components/matching/DeltaViewer.vue`

Сейчас label mapping знает только `amount_with_vat`, `sku_count` и пр.

Вывод: `fe-my` нужно переводить на новый contract vocabulary и на новую semantics денег для ERP.

---

## 1.7. `documan-bruno`

По текущему дереву Bruno в основном содержит:
- smoke health checks
- login / connection flows
- local-debug orchestrator tests

Явных payload-specific API assertions по document details сейчас почти нет.

Вывод: сильного объёма работ тут, вероятно, нет, но smoke/debug коллекции надо прогнать после миграции и при необходимости добавить/обновить document-detail checks.

---

## 1.8. Репозитории без прямой зависимости от payload contract

По текущему коду прямой зависимости от payload field names не видно у:
- `documan-fe-public`
- `documan-gateway`
- `documan-fe-ms`
- `documan-auth`

Их не включаем в активную часть миграции, если не всплывёт скрытая связь через API responses.

---

## 2. Что нужно изменить, чтобы привести текущий код к контракту

## 2.1. `documan-pkg`

### Contract
- Базовый контракт уже принят. Новых правок не требуется, кроме возможных точечных clarifications по ходу миграции.

### `etl.proto`
Нужно удалить legacy transport дубли:

#### `ERPDataItem`
Убрать:
- `erp_updated_at`
- `erp_created_at`
- `erp_deleted_at`
- `erp_archived`

Оставить payload как источник ERP-meta:
- `payload_entity`
- `payload_meta`
- `payload_header`
- `payload_summary`
- `payload_positions`
- `payload_hash`

#### `IngestDataItem`
Убрать:
- `file_name`
- `file_ext`
- `file_hash`

Оставить:
- payload bytes
- `payload_hash`
- `original_file`

### `queries.proto`
Нужно привести query-layer к финальной модели.

#### Original documents
Оставить как отдельные query/storage projection fields:
- `original_file`
- `file_name`, `file_ext`, `file_hash`

#### ERP documents / entities
Оставить как отдельные query/storage projection fields:
- `erp_updated_at`, `erp_created_at`, `erp_deleted_at`, `erp_archived`, `archived`

Эти поля остаются не как часть payload contract, а как query/storage projections для list/detail, filter, sort и index-friendly access.

### `docmatching.proto`
Переименовать matching vocabulary:
- `amount_with_vat` -> `amount_with_tax`
- `sku_count` оставить как осознанное имя matching-агрегата
- `store` без изменений

Правило словаря:
- `summary.sku` — поле payload summary
- `matching.sku_count` — агрегат matching-домена

Это осознанно разные контексты, не требующие механического выравнивания по имени.

### Codegen
- regenerate `services/gen/*`
- затем пересобрать все сервисы, завязанные на `documan-pkg`

---

## 2.2. `documan-connector-ms`

### Document payload
Перевести legacy `transform/doc.go` c `map[string]any` на `contract/document` structs и одновременно переименовать terminology `transform -> payload`.

#### Header
Удалить:
- `source_name`
- `source_date` из header
- `seller_id`
- `buyer_id`
- `store_id`

Добавить / перенести:
- `source_number` в meta (`MS name -> source_number`)
- `source_date` в meta (`moment -> source_date`)
- seller/buyer/store business fields в header НЕ считать обязательным snapshot-слоем для ERP-документов

Новая логика:
- connector-ms document payload для ERP-документов пишется ref-based
- source entity ids продавца/покупателя/склада пишутся в `payload_header.seller_id/buyer_id/store_id`
- enrichment seller/buyer/store business fields на стороне connector-ms НЕ требуется; этим занимается `core`

#### Meta
Заменить legacy meta на contract meta:
- `transformed_at` -> `payloaded_at`
- `updated_at` -> `source_updated_at`
- `created_at` -> `source_created_at`
- `deleted_at` -> `source_deleted_at`
- `archived bool` -> `source_status`
- добавить `source_account_id`
- добавить `source_number`
- добавить `source_date`
- добавить `parse_rule`
- добавить `created_at`

#### Summary
Заменить:
- `amount_with_vat` -> `amount_with_tax`
- `vat` -> `tax_amount`
- `amount` оставить, но только как `*int64` копейки

Добавить:
- `qty`
- `sku` как count unique `article`
- `lines`

Убрать `kopecksToRubles()` из payload.

#### Positions
Заменить и добавить:
- `line_no` -> `line_number` + `line_label` при наличии
- `assortment_type` -> `item_type`
- `unit_price_with_vat` -> `unit_price_with_tax`
- `vat_amount` -> `tax_amount`
- `discount` -> `discount_rate` и при возможности `discount`
- `customs_gtd`, `country_name` сохранить
- `commodity_code` / equivalent -> `hs_code`, если raw даёт

Добавить новые business fields:
- `article` из ERP `article`
- `item_code` из ERP `code`
- `variant`
- `gtin`
- `barcodes`
- `uom_code`, `uom_name`
- `amount`, `amount_with_tax`, `tax_rate`

Не писать больше:
- `assortment_id`
- `assortment_type`
- `product_code`

### Entity payload
Перевести legacy `transform/entity.go` на `contract/entity` structs и одновременно переименовать terminology `transform -> payload`.

#### Party
Заменить:
- `legalTitle` -> `legal_title`
- `code` -> `source_code`
- добавить `address`, если raw даёт
- добавить embedded `entity.Meta`

Удалить:
- `companyType`

#### Store
Заменить:
- `code` -> `source_code`
- добавить `address`, если raw даёт
- добавить embedded `entity.Meta`

Удалить:
- `pathName`

#### Item
Сильнее всего меняется:
- не использовать `source_code`
- использовать `article`
- использовать `item_code` из ERP `code`
- добавить `variant`
- добавить embedded `entity.Meta`
- отделить `gtin` от `barcodes`

Удалить legacy-поля:
- `code`
- `description`
- `pathName`

### Delivery
После обновления proto:
- перестать передавать `erp_*` как отдельные transport fields
- брать ERP-meta только из payload/meta или payload_entity.Meta

### Reset / repayload
- bump payload versions (вместо legacy transform versions)
- переименовать `reclear` -> `repayload` в коде, воркерах, командах и конфиге
- прогнать repayload всех entity и documents после выкатки

---

## 2.3. `documan-ingest`

### Parser / payload build
Перевести `parser.go` с legacy clear-building на `contract/document` structs и одновременно переименовать terminology `clear -> payload`.

#### Header
- `doc_type` оставить source-specific
- `parse_error` убрать из header, перенести в meta
- `shipper_*`, `receiver_*` сохранить как максимально полный контракт, но заполнять частично, если parser не умеет надёжно разбирать

#### Meta
Удалить legacy-поля:
- `doc_type`
- `format`
- `parse_elapsed_ms`

`parse_elapsed_ms` не переносим в контракт и `payload_meta`; оставляем только в логах/метриках ingest.

Добавить:
- `source = "ingest"`
- `file_name`
- `file_ext`
- `file_hash`
- `file_modified_at`, если реально известен
- `s3_bucket`, `s3_key`
- `parse_rule`
- `parse_error`
- `created_at`
- `payloaded_at`

#### Summary
Заменить:
- `vat` -> `tax_amount`
- `amount_with_vat` -> `amount_with_tax`

Добавить:
- `qty`
- `sku` как count unique `article`
- `lines`

#### Positions
Заменить:
- `line_no` -> `line_number` + `line_label`
- `unit` -> `uom_name`
- `unit_code` -> `uom_code`
- `price` -> `unit_price`
- `price_with_vat` -> `unit_price_with_tax`
- `vat` -> `tax_amount`
- `vat_rate` -> `tax_rate`
- `amount_with_vat` -> `amount_with_tax`
- `commodity_code` -> `hs_code`
- `quantity float64` -> decimal string

Оставить / добавить:
- `article`
- `item_code`, если source document реально содержит отдельную графу `Код товара/работ, услуг`
- `variant`, если extractor может выделить
- `gtin`, `barcodes`, если source реально содержит
- `discount_rate`, `discount` — только при наличии

Правила:
- `line_label` хранит исходную графу `п/п`
- `line_number` — 0-based индекс payload
- `article` и `item_code` не дублировать автоматически

### Delivery
После обновления proto и contract meta:
- перестать класть `original_file` и `file_batch` внутрь `payload_meta`
- `original_file` оставить только transport field
- `s3_bucket/s3_key` писать в `payload_meta`
- `file_name/file_ext/file_hash` оставить в payload_meta

### Repayload
- переименовать workflow `reclear` -> `repayload`
- обновить workflow пересборки payload
- после cutover пересобрать все ingest documents

---

## 2.4. `documan-core`

### ETL transport and use case model
После обновления proto нужно синхронно обновить:
- `internal/controller/grpc/v1/notify_erp_data.go`
- `internal/controller/grpc/v1/notify_ingest_data.go`
- `internal/usecase/etl/model.go`
- `internal/usecase/etl/notify_*.go`

#### ERP
Убрать из transport/model:
- `ERPUpdatedAt`
- `ERPCreatedAt`
- `ERPDeletedAt`
- `ERPArchived`

Projection fields остаются отдельными.
При upsert core должен продолжать наполнять их отдельно от payload, используя новый payload/meta как источник данных.

#### Ingest
Убрать из transport/model:
- `FileName`
- `FileExt`
- `FileHash`

Оставить:
- `OriginalFileBucket`
- `OriginalFileKey` или новый transport `original_file`

### Domain / persistence model
Нужно привести доменные модели и SQL schema к одному решению.

#### `OriginalDocument`
Сейчас содержит:
- `FileName`, `FileExt`, `FileHash`
- `OriginalFileBucket`, `OriginalFileKey`

В final state эти поля остаются query/storage projection columns.

#### `ERPDocument` / `ERPEntity`
Сейчас содержат:
- `ERPUpdatedAt`, `ERPCreatedAt`, `ERPDeletedAt`, `ERPArchived`, `Archived`

В final state эти поля остаются query/index projection columns.

Пока по текущему коду видно, что query/list/count endpoints используют эти поля. В этой миграции их не удаляем.

### Payload consumers

#### `dl_candidates.go`
Переписать под новую схему:
- читать `payload_meta` тоже, не только header/summary/positions
- `doc_number` брать из `header.doc_number`, fallback в `meta.source_number`
- `doc_date` брать из `header.doc_date`, fallback в `meta.source_date`
- `amount_with_tax` вместо `amount_with_vat`
- quantity суммировать как decimal string -> numeric parse
- для ERP-документов читать `payload_header.seller_id`, `payload_header.buyer_id`, `payload_header.store_id`
- если seller/buyer/store business fields отсутствуют в header, делать enrichment по `header.seller_id`, `header.buyer_id`, `header.store_id` через `erp_entities`
- fallback через legacy ad-hoc ids вне contract fields не использовать

#### `resolve_field.go`
`resolve_field.go` остаётся generic resolver. Специально переписывать его только ради rename field names не требуется.

Нужно обеспечить:
- callers enrichment не передают legacy field names
- для новых consumers используются contract field names (`inn`, `name`, `legal_title`, `item_code`, `address`)
- payload_entity, публикуемый producers, уже соответствует новому vocabulary

#### `docmatching`
Переименовать match-layer vocabulary:
- `AmountWithVat` -> `AmountWithTax`
- `SKUCount` оставить без переименования
- обновить proto, domain model, gRPC mapping, BFF DTO и FE labels

### Query API
Переписать:
- query controllers
- DTO/proto mapping
- list/detail endpoints

Здесь учитываем как уже принятое решение: transport/query projection fields остаются отдельными.
Для ERP documents query-layer обязан выполнять enrichment seller/buyer/store по `payload_header.seller_id/buyer_id/store_id`, если эти business fields отсутствуют в raw `payload_header`.

### SQL migration
Одна основная миграция должна покрыть:
- `original_documents`
- `erp_documents`
- `erp_entities`
- связанные indexes / constraints / query expectations
- canonical `type` как отдельное internal/DB поле
- `direction` только как internal/DB concern при наличии реального правила, без вывода в внешний API по умолчанию

Конкретно в этой миграции:
- имена колонок `payload_*` не меняем; меняется содержимое JSON payload
- query/storage projection columns `file_*`, `original_file_*`, `erp_*`, `archived` сохраняем
- все SQL- и Go-читатели legacy JSON field names должны быть найдены и обновлены под новый contract vocabulary
- отдельно проверить запросы и выражения вида `payload_header->>'amount_with_vat'`, `payload_summary->>'vat'`, `payload_positions` по старым именам и заменить их
- canonical `type` добавляем/наполняем как internal DB поле
- `direction` не добавляем и не заполняем, если на момент миграции нет готового правила вычисления и реального потребителя

---

## 2.5. `documan-bff`

BFF нужно менять после обновления `queries.proto` и `docmatching.proto`.

### Query DTO
Проверить и обновить:
- `internal/dto/types.go`
- `document_details.go`
- `documents.go`
- `erp_documents.go`
- matching handlers / DTO

### Что именно поменяется
- matching DTO field names (`amountWithVat` -> `amountWithTax` и т.п.)
- response mapping для details/list
- query structs, которые зависят от обновлённых proto, при сохранении отдельных projection fields `file_*`, `original_file_*`, `erp_*`, `archived`

### Что, вероятно, не меняется
- payload JSON как `json.RawMessage` можно оставить
- BFF не обязан разбирать payload по contract structs, если он просто проксирует bytes/json

---

## 2.6. `documan-fe-my`

Это активный scope, потому что фронт сейчас жёстко ожидает legacy payload names.

### Original document pages
Перевести на:
- `amount_with_tax`
- `tax_amount`
- `line_number`, `line_label`
- `article`
- `item_code`
- `uom_name`, `uom_code`
- `hs_code`
- `source_number`
- при необходимости продолжать использовать отдельные query DTO fields `file_*`, так как projection fields остаются

### ERP document pages
Перевести на:
- новые payload names
- убрать ожидание raw `source_name`; `seller_id` / `buyer_id` остаются допустимыми ref fields, но UI не должен использовать raw ids как business fields и не должен сам их резолвить
- seller/buyer/store получать из уже обогащённого query/API ответа `core`, а не из raw ERP payload header
- убрать допущение “ERP summary = рубли”; в новом контракте деньги должны читаться как копейки

### Matching UI
Обновить label mappings и client types под новый `docmatching` vocabulary.

### Composables
Переписать:
- `useDocworker.ts`
- `useErpDocuments.ts`
- `useMatching.ts`, если меняется matching API

---

## 2.7. `documan-bruno`

Минимальный обязательный объём:
- прогнать smoke/debug коллекции после миграции
- обновить коллекции, если есть запросы/проверки на old query/matching response fields

По текущему дереву payload-specific coverage почти нет, значит здесь работы немного.

---

## 3. Репозитории вне активного scope

Пока не планируем правки в:
- `documan-fe-public`
- `documan-gateway`
- `documan-fe-ms`
- `documan-auth`

Основание: в текущем коде нет прямой зависимости от payload contract, legacy payload field names или query details по документам.

---

## 4. Порядок внедрения, исходя из текущего кода

### Шаг 1. Зафиксировать final contract
- `documan-pkg/contract/*` уже baseline
- новых полей не добавлять без отдельного решения

### Шаг 2. Обновить `documan-pkg` proto и regenerate
- `etl.proto`
- `queries.proto`
- `docmatching.proto`
- regenerate `services/gen/*`

### Шаг 3. Привести `documan-core` к новым proto/model/SQL
Без этого producers и BFF потом не соберутся.

Включает:
- ETL transport
- domain models
- SQL migration
- query controllers
- docmatching layer
- payload consumers

### Шаг 4. Перевести `documan-connector-ms`
- document payload
- entity payload
- delivery
- terminology rename `transform/clear -> payload`
- repayload/reset versioning

### Шаг 5. Перевести `documan-ingest`
- parser/build payload
- delivery
- terminology rename `clear -> payload`, `reclear -> repayload`
- repayload workflow

### Шаг 6. Перевести `documan-bff`
- DTO
- query mapping
- matching responses

### Шаг 7. Перевести `documan-fe-my`
- composables
- details pages
- ERP details pages
- matching UI

### Шаг 8. Обновить `documan-bruno`
- smoke/debug collections по необходимости

### Шаг 9. Coordinated cutover
- применить SQL migration
- выкатить согласованные версии сервисов
- удалить test volume со старым `transformed/`/legacy payload содержимым
- прогнать `repayload` и пересобрать тестовые данные заново без миграции старых S3 objects

### Шаг 10. Verification
- build / typecheck всех затронутых репозиториев
- выборочные API checks
- UI manual smoke
- проверка payload в S3/DB

---

## 5. Основные риски, которые видны по текущему коду

1. `core` сейчас использует отдельные DB/query поля `erp_*`, `file_*`, `original_file_*`; их нельзя “незаметно” удалить без решения, где они будут жить после миграции.
2. `fe-my` не просто показывает payload, а содержит assumptions о semantics денег ERP (`rubles`, а не `kopecks`). Это behavioural change, не только rename.
3. `docmatching` сейчас жёстко завязан на old vocabulary `amount_with_vat`; его нельзя забыть вне общего плана.
4. При выбранной flat ref-based модели `core` обязан корректно enrich-ить seller/buyer/store по `payload_header.seller_id/buyer_id/store_id`, а item-поля по `payload_positions[*].item_id` там, где это требуется; иначе ERP details/matching останутся без бизнес-полей.
5. Rename `transform/clear -> payload` затрагивает S3 keys, worker names, config/env и workflow vocabulary; его нужно делать синхронно с contract migration, а не отдельным подготовительным шагом.
6. `ingest` сейчас пишет `original_file` внутрь `payload_meta`; при переходе это поведение надо осознанно убрать и не потерять скачивание исходника.

---

## 6. Принятые решения для старта реализации

1. Projection fields остаются отдельными:
- `file_*`
- `original_file_*`
- `erp_*`
- `archived`

Это query/storage projections, а не часть payload contract.

2. Matching vocabulary:
- `amount_with_vat` -> `amount_with_tax`
- `sku_count` оставить

Правило словаря:
- `summary.sku` — поле payload summary
- `matching.sku_count` — агрегат matching-домена

3. Canonical classification:
- `type` — internal/DB сейчас
- `direction` — не тащить в внешний API без отдельной потребности

4. ERP document refs strategy:
- `seller_id`, `buyer_id`, `store_id` в `payload_header` и `item_id` в `payload_positions` являются допустимыми contract ref fields
- `connector-ms` пишет ref-based ERP document payload без обязательного business snapshot в header/positions
- `core` выполняет enrichment seller/buyer/store по `payload_header.*_id` и использует `payload_positions[*].item_id` там, где это требуется для query/matching/read models
- `ingest` обычно оставляет ref fields пустыми

Эти решения считаются зафиксированными и не открывают новый round обсуждения внутри самой миграции.
