# Stage 7: documan-ingest — Parser/payload build

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.3 (Parser/payload build).
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Перевести parser/clear-building в documan-ingest на `contract/document` structs.

## Сервис
- **Имя**: documan-ingest
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-ingest -name "go.mod"`

## Входы
- Найти parser: `grep -rl "parser\|Parser\|clear" documan-ingest/internal/adapter/`
- Найти sheet/extractors: `grep -rl "sheet\|extract" documan-ingest/internal/`
- Прочитать contract structs: `documan-pkg/contract/document/`
- Проверить done_log

## Что сделать

### 1. Header
- `doc_type` — оставить source-specific
- `parse_error` — убрать из header, перенести в meta
- `shipper_*`, `receiver_*` — сохранить, заполнять частично если parser не может надёжно

### 2. Meta
**Удалить**:
- `doc_type`, `format`, `parse_elapsed_ms`

`parse_elapsed_ms` не переносим в контракт; оставляем только в логах/метриках ingest.

**Добавить**:
- `source = "ingest"`
- `file_name`, `file_ext`, `file_hash`
- `file_modified_at` (если реально известен)
- `s3_bucket`, `s3_key`
- `parse_rule`
- `parse_error`
- `created_at`, `payloaded_at`

### 3. Summary
- `vat` → `tax_amount`
- `amount_with_vat` → `amount_with_tax`
- Добавить: `qty`, `sku` (count unique article), `lines`

### 4. Positions
**Заменить**:
- `line_no` → `line_number` + `line_label`
- `unit` → `uom_name`
- `unit_code` → `uom_code`
- `price` → `unit_price`
- `price_with_vat` → `unit_price_with_tax`
- `vat` → `tax_amount`
- `vat_rate` → `tax_rate`
- `amount_with_vat` → `amount_with_tax`
- `commodity_code` → `hs_code`
- `quantity float64` → decimal string

**Правила**:
- `line_label` хранит исходную графу `п/п`
- `line_number` — 0-based индекс payload
- `article` и `item_code` не дублировать автоматически

**Добавить** (при наличии в source):
- `article`, `item_code`
- `variant`, `gtin`, `barcodes`
- `discount_rate`, `discount`

### 5. Payload через contract structs
Перестать строить payload как ad-hoc maps. Использовать `contract/document` structs.

## Критерии готовности
- [ ] Parser строит payload через `contract/document` structs
- [ ] Header не содержит `parse_error`
- [ ] Meta содержит все contract-поля
- [ ] Summary использует `amount_with_tax`, `tax_amount`
- [ ] Positions используют contract vocabulary
- [ ] `quantity` — decimal string, не float64
- [ ] `go build ./...` из `documan-ingest/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-ingest && go build ./...
```

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-ingest/`:
```
DONE: REFACTOR_CONTRACT stage 7: parser → contract document structs
```

## Stop/Continue
- Если parser extractors не могут выделить новое поле из source — оставить nil/пустым, `minor`.
- Если `quantity` переход с float64 на string ломает downstream — `major`, СТОП.
