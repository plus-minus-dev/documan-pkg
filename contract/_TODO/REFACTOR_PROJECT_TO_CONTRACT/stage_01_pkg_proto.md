# Stage 1: documan-pkg — Proto schemas + codegen

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.1.
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Обновить contract docs и proto-файлы под финальный контракт, затем пересобрать codegen.

## Сервис
- **Имя**: documan-pkg
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-pkg -name "go.mod"` и работать из этого корня

## Входы
- Найти proto-файлы: `grep -rl "\.proto$" documan-pkg/services/proto/`
- Прочитать текущие: `etl.proto`, `queries.proto`, `docmatching.proto`
- Прочитать финальный контракт: `documan-pkg/contract/document/`, `documan-pkg/contract/entity/`
- Прочитать `documan-pkg/contract/README.md`
- Прочитать ТЗ раздел 2.1 — что именно менять

## Что сделать

### 0. Contract docs
- Добавить и задокументировать flat ref fields в `contract/document/header.go` (`seller_id`, `buyer_id`, `store_id`)
- Добавить и задокументировать flat ref field `item_id` в `contract/document/position.go`
- Обновить `contract/README.md` под flat ref-based ERP model

### 1. `etl.proto`
#### `ERPDataItem` — убрать:
- `erp_updated_at`
- `erp_created_at`
- `erp_deleted_at`
- `erp_archived`

Оставить: `payload_entity`, `payload_meta`, `payload_header`, `payload_summary`, `payload_positions`, `payload_hash`.

#### `IngestDataItem` — убрать:
- `file_name`
- `file_ext`
- `file_hash`

Оставить: payload bytes, `payload_hash`, `original_file`.

### 2. `queries.proto`
Привести к финальной модели. Projection fields (`file_*`, `original_file_*`, `erp_*`, `archived`) оставить как отдельные query/storage поля.

### 3. `docmatching.proto`
- `amount_with_vat` → `amount_with_tax`
- `sku_count` — оставить без переименования

### 4. Codegen
- Regenerate `services/gen/*`
- Убедиться, что `.proto` пересобраны с чистого листа (новая нумерация полей, без `reserved`)

## Правила (из ТЗ)
- `.proto` пересобираются с чистого листа. Обратная совместимость protobuf wire format НЕ требуется.
- `reserved` для удалённых legacy-полей НЕ использовать.
- Нумерацию полей можно выстроить заново.
- Изменения `.proto` АТОМАРНО: правка schema + regenerate + обновление кода в том же проходе.

## Критерии готовности
- [ ] Contract docs описывают flat ref model (`seller_id`, `buyer_id`, `store_id`, `item_id`)
- [ ] `etl.proto` не содержит legacy transport полей
- [ ] `queries.proto` приведён к финальной модели с сохранением projection fields
- [ ] `docmatching.proto` использует `amount_with_tax`
- [ ] `services/gen/*` пересобраны
- [ ] `go build ./...` из `documan-pkg/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-pkg && go build ./...
```

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-pkg/`:
```
DONE: REFACTOR_CONTRACT stage 1: proto schemas + codegen
```

## Stop/Continue
- Если codegen ломает импорты в pkg — исправить В ТОМ ЖЕ проходе до gate check. Proto changes АТОМАРНЫ: schema + regenerate + обновление кода в одном проходе. Сломанные импорты на момент gate check = gate FAIL.
- Если обнаружены неожиданные зависимости от proto в pkg — `major`, СТОП.
