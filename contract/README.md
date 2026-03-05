# Contract — Payload Schema

> Единый источник правды для структуры payload документов и справочников DocuMan.

## 1. Канонические типы документов

Это vocabulary для слоя классификации в `core`, а не значения для `payload_header.doc_type`.

```
type:      order | shipment | invoice | payment
direction: forward | reverse
```

| Type | Что это | Источники |
|------|---------|-----------|
| `order` | Заказ | МС: purchaseorder |
| `shipment` | Отгрузка / приёмка | МС: supply; Ingest: УПД, ТОРГ-12, накладная |
| `invoice` | Счёт | МС: invoicein; Ingest: счет, счет-оферта, СЧФ |
| `payment` | Платёж | МС: paymentout, cashout |

`direction: reverse` — возвраты, кредит-ноты (зарезервировано).

## 2. Пакеты

```
contract/
├── document/   — payload документов (Header, Meta, Summary, Position)
└── entity/     — payload справочников (Party, Item, Store)
```

### document

| Struct | Назначение |
|--------|-----------|
| `Header` | Реквизиты документа (стороны, валюта, склад) |
| `Meta` | Служебная информация (источник, файл, S3, парсинг) |
| `Summary` | Итоговые суммы и агрегаты |
| `Position` | Строка товара/услуги |

### entity

| Struct | Назначение |
|--------|-----------|
| `Party` | Контрагент / организация |
| `Item` | Товар / услуга / комплект / вариант |
| `Store` | Склад |

## 3. Conventions

- Суммы — `*int64` в копейках (nil = нет данных, 0 = валидный ноль)
- Количества и процентные значения — decimal string (`"1"`, `"0.500"`, `"10.5"`, dot separator)
- JSON — `snake_case`
- Go-поля короткие, JSON сохраняет контекст (`Type` -> `"doc_type"`, `Number` -> `"doc_number"`)
- `payload_header.doc_type` — оригинальный тип из источника (`supply`, `УПД`, `счет`...), не canonical type
- `doc_date`, `source_date` — формат `YYYY-MM-DD`
- Timestamps в `payload_meta` и `entity.Meta` — RFC3339 UTC
- `line_number` — нормализованный 0-based индекс строки в payload
- `line_label` — исходное значение номера строки из источника (`"1"`, `"1а"`, `"I"` ...)
- Налог: `tax_rate`, `tax_amount`, `amount_with_tax` (не `vat_*`)
- Единицы: `uom_code`, `uom_name` (не `unit_*`)
- Артикул в contract field names хранится как `article`
- `summary.sku` — осознанное имя агрегата, означает количество уникальных артикулов (`article`)
- Термины `sku` и `article` в проекте считаются согласованными: `article` используется как имя payload-поля для артикула, `sku` осознанно сохранён только для агрегата `summary.sku` и связанных с ним производных агрегатных имен
- Таможня: `hs_code` (не `commodity_code`)
- `item_code` — бизнесовый код товара/позиции из источника; не `article`, не `gtin`, не `hs_code`
- `source_code` — код справочника в системе-источнике; используется только там, где это отдельный source-specific атрибут entity payload
- `source_status` вычисляется по приоритету: `deleted` -> `archived` -> `active`

### Payload vs Transport

- Payload contract описывает только содержимое `payload_meta`, `payload_header`, `payload_summary`, `payload_positions`, `payload_entity`
- Transport/query-layer может содержать дополнительные поля вне payload, например `original_file`
- Хранение проекций payload-полей в отдельных колонках БД допустимо, если это нужно для запросов, сортировки или индексов

### Enforcement

- **Producer**: payload через contract struct (compile-time) -> неизвестное поле = compile error
- **Consumer**: payload читает contract field names и не должен вводить собственные альтернативные схемы

## 4. Разделение ответственности

**Producer** (ingest, connector-ms, connector-1c):
- Формирует payload строго через contract structs
- Записывает оригинальный `doc_type` источника (`supply`, `УПД`, `счет`...)
- Нормализует даты, decimal string и денежные поля до contract format

**Core**:
- Классифицирует `doc_type -> canonical type` (одно место, один switch)
- Может хранить canonical `type` и `direction` отдельно от payload
- Валидирует бизнес-инварианты payload
