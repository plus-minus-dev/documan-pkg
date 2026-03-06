# Stage 4: documan-core — Query API controllers

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.4 (Query API).
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Обновить query controllers и DTO/proto mapping documan-core под обновлённые `queries.proto`, включая ERP enrichment по flat ids в `payload_header` / `payload_positions`.

## Сервис
- **Имя**: documan-core
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-core -name "go.mod"`

## Входы
- Найти query controllers: `grep -rl "queries\|Query\|List\|Detail" documan-core/internal/controller/grpc/`
- Прочитать обновлённые gen: `documan-pkg/services/gen/` (queries)
- Проверить done_log на отклонения этапов 2 и 3

## Что сделать

### 1. Query controllers
Найти все gRPC query handlers через grep и обновить:
- DTO/proto mapping для original documents
- DTO/proto mapping для ERP documents/entities
- List/detail endpoints

### 2. Учесть принятое решение
Transport/query projection fields остаются отдельными:
- `file_name`, `file_ext`, `file_hash`, `original_file` — для original docs
- `erp_updated_at`, `erp_created_at`, `erp_deleted_at`, `erp_archived`, `archived` — для ERP

Эти поля продолжают маппиться как отдельные query/storage projections, НЕ как часть payload.

Для ERP details/list query-layer обязан возвращать seller/buyer/store business fields в обогащённом виде, если raw `payload_header` их не содержит. Источник обогащения — `payload_header.seller_id/buyer_id/store_id` + `erp_entities`; line-level enrichment при необходимости использует `payload_positions[*].item_id`.

### 3. Grep-аудит
Выполнить grep по documan-core на предмет оставшихся legacy payload field names в query layer:
```bash
grep -rn "amount_with_vat\|source_name\|line_no\|commodity_code\|unit_price_with_vat\|legalTitle" documan-core/internal/controller/
```

## Критерии готовности
- [ ] Query controllers компилируются с новыми `queries.proto` gen
- [ ] DTO mapping обновлён
- [ ] Projection fields корректно маппятся
- [ ] Grep-аудит по query layer чист от legacy payload names
- [ ] ERP query responses умеют enrich-ить seller/buyer/store по `payload_header.*_id`
- [ ] `go build ./...` из `documan-core/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-core && go build ./...
```

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-core/`:
```
DONE: REFACTOR_CONTRACT stage 4: query API controllers
```

## Stop/Continue
- Если query handler зависит от payload field names (не projection) — переписать сразу.
- Если обнаружен неожиданный consumer за пределами query layer — обновить на текущем этапе. Если scope за пределами core — зафиксировать как `tracked` для этапа соответствующего сервиса. Не закрыто к stage 11 = `major`.
- Если proto gen не совместим с текущими query structures — `major`, СТОП.
