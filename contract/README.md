# Contract — Payload Schema

> Единый источник правды для структуры payload документов и справочников DocuMan.

## 1. Канонические типы документов

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
| `Item` | Товар / услуга / комплект (предварительно) |
| `Store` | Склад |

## 3. Conventions

- Суммы — `*int64` в копейках (nil = нет данных, 0 = валидный ноль)
- Количества — decimal string (`"1"`, `"0.500"`, dot separator)
- JSON — `snake_case`
- Go-поля короткие, JSON сохраняет контекст (`Type` → `"doc_type"`, `Number` → `"doc_number"`)
- Налог: `tax_rate`, `tax_amount`, `amount_with_tax` (не `vat_*`)
- Единицы: `uom_code`, `uom_name` (не `unit_*`)
- Артикул: `sku` (не `article`)
- Таможня: `hs_code` (не `commodity_code`)

### Enforcement

- **Producer**: payload через contract struct (compile-time) → неизвестное поле = compile error
- **Consumer**: неизвестные поля — warning в лог (strict с error — для CI)

## 4. Разделение ответственности

**Producer** (ingest, connector-ms, connector-1c):
- Формирует payload строго через contract structs
- Записывает оригинальный `doc_type` источника (`supply`, `УПД`, `счет`...)

**Core**:
- Классифицирует `doc_type → canonical type` (одно место, один switch)
- Валидирует бизнес-инварианты payload
