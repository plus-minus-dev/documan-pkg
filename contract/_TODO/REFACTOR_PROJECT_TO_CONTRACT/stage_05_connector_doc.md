# Stage 5: documan-connector-ms — Document payload

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.2 (Document payload).
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Перевести сборку document payload в documan-connector-ms с `map[string]any` на `contract/document` structs.

## Сервис
- **Имя**: documan-connector-ms
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-connector-ms -name "go.mod"`

## Входы
- Найти текущий transform: `grep -rl "transform\|doc\.go" documan-connector-ms/internal/usecase/`
- Прочитать contract structs: `documan-pkg/contract/document/`
- Проверить done_log

## Что сделать

### 1. Header
**Удалить**: `source_name`, `source_date` из header, `seller_id`, `buyer_id`, `store_id`

**Добавить/перенести**:
- `source_number` в meta (MS name → source_number)
- `source_date` в meta (moment → source_date)
- `seller_name`, `seller_inn`, `seller_kpp`, `seller_address` (при возможности)
- `buyer_name`, `buyer_inn`, `buyer_kpp`, `buyer_address` (при возможности)
- `store_name`

**Новая логика**: document transform должен best-effort резолвить ERP entities по UUID из raw документа. Если lookup не удался — payload остаётся валидным с пустыми enriched fields.

### 2. Meta
Заменить legacy meta на contract meta:
- `transformed_at` → `payloaded_at`
- `updated_at` → `source_updated_at`
- `created_at` → `source_created_at`
- `deleted_at` → `source_deleted_at`
- `archived bool` → `source_status`
- Добавить: `source_account_id`, `source_number`, `source_date`, `parse_rule`, `created_at`

### 3. Summary
- `amount_with_vat` → `amount_with_tax`
- `vat` → `tax_amount`
- `amount` — только как `*int64` копейки
- Добавить: `qty`, `sku` (count unique article), `lines`
- Убрать `kopecksToRubles()` из payload

### 4. Positions
**Заменить**:
- `line_no` → `line_number` + `line_label`
- `assortment_type` → `item_type`
- `unit_price_with_vat` → `unit_price_with_tax`
- `vat_amount` → `tax_amount`
- `discount` → `discount_rate` и `discount`
- `commodity_code` → `hs_code` (если raw даёт)

**Добавить**:
- `article` из ERP article
- `item_code` из ERP code
- `variant`, `gtin`, `barcodes`, `uom_code`, `uom_name`
- `amount`, `amount_with_tax`, `tax_rate`

**Не писать больше**: `assortment_id`, `assortment_type`, `product_code`

### 5. Payload собирается через contract structs
`map[string]any` и ad-hoc JSON maps НЕ допускаются как финальное состояние.

## Критерии готовности
- [ ] Document payload собирается через `contract/document` structs
- [ ] Нет `map[string]any` для payload
- [ ] Header не содержит legacy field names
- [ ] Meta использует contract-поля
- [ ] Summary в копейках, без `kopecksToRubles()`
- [ ] Positions используют contract vocabulary
- [ ] `go build ./...` из `documan-connector-ms/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-connector-ms && go build ./...
```

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-connector-ms/`:
```
DONE: REFACTOR_CONTRACT stage 5: document payload → contract structs
```

## Stop/Continue
- Entity enrichment (best-effort lookup) — новая логика. Если lookup endpoint не готов, сделать stub и зафиксировать как `minor` (это новая функциональность, не legacy).
- Если обнаружены потребители `map[string]any` payload — обновить на текущем этапе. `map[string]any` как финальное состояние payload = `major`. Если scope за пределами connector-ms — зафиксировать как `tracked` для этапа соответствующего сервиса.
- Если contract structs не покрывают нужные поля — `major`, СТОП (может потребоваться расширение контракта).
