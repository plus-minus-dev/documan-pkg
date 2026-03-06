# Stage 3: documan-core — Payload consumers

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.4 (Payload consumers).
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Переписать все payload consumers в documan-core под новый contract vocabulary и конкретный enrichment по flat ids в `payload_header` / `payload_positions`.

## Сервис
- **Имя**: documan-core
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-core -name "go.mod"`

## Входы
- Найти: `grep -rl "dl_candidates\|resolve_field\|docmatching\|dl_delta" documan-core/internal/`
- Прочитать contract: `documan-pkg/contract/document/`, `documan-pkg/contract/entity/`
- Проверить done_log на отклонения этапа 2

## Что сделать

### 1. `dl_candidates.go`
Переписать под новую схему:
- Читать `payload_meta` в дополнение к header/summary/positions
- `doc_number` брать из `header.doc_number`, fallback → `meta.source_number` (оба поля — новый контракт, это НЕ legacy fallback, а доменное правило выбора внутри новой схемы)
- `doc_date` брать из `header.doc_date`, fallback → `meta.source_date` (аналогично — оба contract fields)
- `amount_with_tax` вместо `amount_with_vat`
- quantity суммировать как decimal string → numeric parse
- Для ERP documents читать `header.seller_id`, `header.buyer_id`, `header.store_id`
- Если seller/buyer/store business fields отсутствуют в header, делать enrichment через `erp_entities` по `header.seller_id`, `header.buyer_id`, `header.store_id`
- На этом этапе восстанавливается конкретный enrichment в `dl_candidates.go`; generic enrich helper/service в эту миграцию не входит
- Legacy ad-hoc ids вне contract fields не использовать

### 2. `resolve_field.go`
Переписать ожидаемые поля:
- `legalTitle` → `legal_title`
- `article` для item — оставить, но в новой semantics
- `code` для item → `item_code`
- Добавить `address` для `Party` и `Store`

### 3. `docmatching`
Найти: `grep -rl "docmatching" documan-core/internal/`
Обновить:
- `AmountWithVat` → `AmountWithTax`
- `SKUCount` — оставить без переименования
- Proto mapping, domain model, gRPC mapping — всё синхронно
- Обновить `docmatching.proto` domain features `features.go`

### 4. `dl_delta.go`
Проверить выводимые field names. Если выводит legacy — обновить.

## Критерии готовности
- [ ] `dl_candidates.go` использует contract vocabulary и конкретный enrichment по `header.seller_id` / `header.buyer_id` / `header.store_id`
- [ ] `resolve_field.go` остаётся generic; callers enrichment используют только contract field names
- [ ] `docmatching` использует `AmountWithTax`
- [ ] `dl_delta.go` не выводит legacy field names
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
DONE: REFACTOR_CONTRACT stage 3: payload consumers
```

## Stop/Continue
- Если обнаружены ещё payload consumers, не перечисленные в ТЗ — обновить тут же в этом этапе. Если невозможно обновить — `major`, СТОП. Оставлять legacy consumer = `major`.
- Если consumer логика требует данных, которых ещё нет в contract — `major`, СТОП.
