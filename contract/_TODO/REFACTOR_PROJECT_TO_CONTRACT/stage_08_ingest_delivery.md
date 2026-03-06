# Stage 8: documan-ingest — Delivery + rename

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.3 (Delivery, Repayload) + "Терминологический rename".
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Обновить delivery в documan-ingest. Выполнить терминологический rename `clear → payload`, `reclear → repayload`.

## Сервис
- **Имя**: documan-ingest
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-ingest -name "go.mod"`

## Входы
- Найти delivery: `grep -rl "delivery\|Delivery" documan-ingest/internal/`
- Найти model: `grep -rl "model\.go" documan-ingest/internal/usecase/delivery/`
- Найти все `clear/reclear`: `grep -rn "clear\|reclear" documan-ingest/internal/ --include="*.go"`
- Проверить done_log (отклонения этапа 7)

## Что сделать

### 1. Delivery
- Перестать класть `original_file` и `file_batch` внутрь `payload_meta`
- `original_file` — оставить только как transport field
- `s3_bucket/s3_key` — писать в `payload_meta`
- `file_name/file_ext/file_hash` — оставить в `payload_meta`

### 2. Терминологический rename
- `clear` → `payload` (пакеты, функции, переменные, config keys, worker names, лог-сообщения)
- `reclear` → `repayload`
- S3 path `transformed/` → `payload/`
- `documan-bff` НЕ затрагивается этим rename

### 3. Repayload workflow
- Переименовать workflow `reclear` → `repayload`
- Обновить workflow пересборки payload

## Критерии готовности
- [ ] Delivery не кладёт `original_file`/`file_batch` в `payload_meta`
- [ ] `s3_bucket/s3_key` в `payload_meta`
- [ ] Терминология `clear` заменена на `payload`
- [ ] S3 paths обновлены
- [ ] Repayload workflow переименован
- [ ] `go build ./...` из `documan-ingest/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-ingest && go build ./...
```

## Grep-аудит rename (весь сервис, все типы файлов)
```bash
grep -rn "\"clear\"\|reclear\|transformed/" documan-ingest/ \
  --include="*.go" --include="*.yaml" --include="*.yml" \
  --include="*.sh" --include="*.md" \
  --include="Makefile" --include="Dockerfile*" \
  | grep -v "_archive_del" | grep -v "vendor/"

# Отдельно .env*
find documan-ingest/ -name ".env*" -exec grep -Hn "\"clear\"\|reclear\|transformed/" {} \;
```
Результат должен быть пустым. Любое нахождение = **major**.

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-ingest/`:
```
DONE: REFACTOR_CONTRACT stage 8: delivery + rename clear→payload
```

## Stop/Continue
- Config/env ключи ДОЛЖНЫ быть переименованы на этом этапе. Дублирование старых и новых config/env keys запрещено. Docker-compose/k8s файлы обновить тут же.
- Тесты (`_test.go`) тоже должны использовать новую терминологию.
- Если rename ломает integration с core — `major`, СТОП.
