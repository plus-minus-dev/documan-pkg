# Stage 9: documan-bff — DTO + query mapping + matching

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.5.
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Обновить BFF под обновлённые `queries.proto` и `docmatching.proto`.

## Сервис
- **Имя**: documan-bff
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-bff -name "go.mod"`

## Входы
- Найти DTO: `grep -rl "dto\|types\.go\|document_details\|erp_documents" documan-bff/internal/`
- Найти matching handlers: `grep -rl "matching\|docmatching" documan-bff/internal/`
- Прочитать обновлённые gen: `documan-pkg/services/gen/`
- Проверить done_log (отклонения этапов 3, 4)

## Что сделать

### 1. Query DTO
Найти и обновить:
- `internal/dto/types.go`
- `document_details.go`
- `documents.go`
- `erp_documents.go`

### 2. Matching DTO
- `amountWithVat` → `amountWithTax`
- Обновить matching handlers и response mapping

### 3. Response mapping для details/list
- Обновить query structs, зависящие от обновлённых proto
- Projection fields `file_*`, `original_file_*`, `erp_*`, `archived` — сохранить отдельными
- Для ERP details/list полагаться на уже обогащённые ответы `core`; BFF не должен сам резолвить `seller_id`, `buyer_id`, `store_id`, `item_id`

### 4. Что вероятно НЕ меняется
- BFF может проксировать payload как `json.RawMessage` ТОЛЬКО если он НЕ конструирует, НЕ модифицирует и НЕ интерпретирует payload fields. Если BFF читает или пишет отдельные payload fields — ОБЯЗАН использовать contract structs.

## Критерии готовности
- [ ] DTO обновлены под новые proto gen
- [ ] Matching DTO использует `amountWithTax`
- [ ] Response mapping работает с новыми query proto
- [ ] Projection fields корректно маппятся
- [ ] `go build ./...` из `documan-bff/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-bff && go build ./...
```

## Grep-аудит
```bash
grep -rn "AmountWithVat\|amountWithVat\|amount_with_vat" documan-bff/internal/ --include="*.go"
```
Результат должен быть пустым.

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-bff/`:
```
DONE: REFACTOR_CONTRACT stage 9: BFF DTO + query mapping + matching
```

## Stop/Continue
- BFF проксирует payload как raw JSON без парсинга — это архитектурно допустимо, НЕ является отклонением.
- Если обнаружены зависимости от legacy payload field names — обновить на этом этапе. Оставить legacy reference = `major`.
- Если proto gen breaking change не покрывается DTO — `major`, СТОП.
