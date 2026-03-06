# Stage 6: documan-connector-ms — Entity payload + delivery + rename

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + разделы 2.2 (Entity payload, Delivery, Reset/repayload) + "Терминологический rename".
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Перевести entity payload на contract structs. Обновить delivery. Выполнить терминологический rename `transform/clear → payload`, `reclear → repayload`.

## Сервис
- **Имя**: documan-connector-ms
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-connector-ms -name "go.mod"`

## Входы
- Найти entity transform: `grep -rl "entity\|transform" documan-connector-ms/internal/usecase/`
- Найти delivery: `grep -rl "delivery\|Deliver" documan-connector-ms/internal/`
- Прочитать contract entity structs: `documan-pkg/contract/entity/`
- Проверить done_log (отклонения этапа 5)

## Что сделать

### 1. Entity payload — Party (Counterparty / Organization)
- `legalTitle` → `legal_title`
- `code` → `source_code`
- Добавить `address` (если raw даёт)
- Добавить embedded `entity.Meta`
- Удалить `companyType`

### 2. Entity payload — Store
- `code` → `source_code`
- Добавить `address` (если raw даёт)
- Добавить embedded `entity.Meta`
- Удалить `pathName`

### 3. Entity payload — Item
- Не использовать `source_code`
- Использовать `article`
- `item_code` из ERP `code`
- Добавить `variant`
- Добавить embedded `entity.Meta`
- Отделить `gtin` от `barcodes`
- Удалить: `code`, `description`, `pathName`

### 4. Delivery
После обновления proto:
- Перестать передавать `erp_*` как отдельные transport fields
- Брать ERP-meta только из payload/meta или payload_entity.Meta
- `payload_header.seller_id/buyer_id/store_id` и `payload_positions[*].item_id` в document payload должны ссылаться на те же source ids, по которым публикуются `entity.Party`, `entity.Store`, `entity.Item`

### 5. Терминологический rename (СИНХРОННО с contract migration)
Правила:
- Rename делается только вместе с contract migration в затронутых местах
- Логика НЕ меняется, меняется только терминология
- `core.ClearGroupID` НЕ трогаем

Переименовать:
- `transform` → `payload` (пакеты, функции, переменные, config keys, worker names)
- `clear` → `payload`
- `reclear` → `repayload`
- S3 path `transformed/` → `payload/`
- Существующие test objects в старом S3 path НЕ мигрируем

### 6. Reset / repayload
- Bump payload versions (вместо legacy transform versions)
- Переименовать `reclear` → `repayload` в коде, воркерах, командах и конфиге

## Критерии готовности
- [ ] Entity payload собирается через `contract/entity` structs
- [ ] Delivery не передаёт `erp_*` отдельно
- [ ] Терминология `transform/clear` заменена на `payload`
- [ ] S3 paths обновлены на `payload/`
- [ ] `core.ClearGroupID` НЕ изменён
- [ ] `go build ./...` из `documan-connector-ms/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-connector-ms && go build ./...
```

## Grep-аудит rename (весь сервис, все типы файлов)
```bash
grep -rn "transform\|reclear\|transformed/" documan-connector-ms/ \
  --include="*.go" --include="*.yaml" --include="*.yml" \
  --include="*.sh" --include="*.md" \
  --include="Makefile" --include="Dockerfile*" \
  | grep -v "ClearGroupID" | grep -v "_archive_del" | grep -v "vendor/"

# Отдельно .env* (не поддерживает --include glob)
find documan-connector-ms/ -name ".env*" -exec grep -Hn "transform\|reclear\|transformed/" {} \;
```
Результат должен быть пустым. Любое нахождение = **major**.

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-connector-ms/`:
```
DONE: REFACTOR_CONTRACT stage 6: entity payload + delivery + rename transform→payload
```

## Stop/Continue
- Config/env keys ДОЛЖНЫ быть переименованы на этом этапе. Если docker-compose/k8s файлы в отдельном месте — обновить их тут же. Дублирование старых и новых config/env keys запрещено (правило из ТЗ).
- Тесты (`_test.go`) тоже должны использовать новую терминологию. Legacy terminology в тестах = `major` к stage 11.
- Если rename ломает существующие integration points — `major`, СТОП.
