# Stage 11: Verification + Bruno + cutover preparation

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — разделы 2.7, 4 (шаги 9-10), 5 (риски).
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Cross-repo verification, обновление Bruno, grep-аудит legacy references, подготовка cutover.

## Сервисы
Все затронутые:
- `documan-pkg`
- `documan-core`
- `documan-connector-ms`
- `documan-ingest`
- `documan-bff`
- `documan-fe-my`
- `documan-bruno`

## Входы
- Полный done_log с отклонениями всех этапов
- ТЗ как reference для финального аудита

## Что сделать

### 1. Cross-repo build
```bash
cd /Users/av/Desktop/WORK/documan/documan-pkg && go build ./...
cd /Users/av/Desktop/WORK/documan/documan-core && go build ./...
cd /Users/av/Desktop/WORK/documan/documan-connector-ms && go build ./...
cd /Users/av/Desktop/WORK/documan/documan-ingest && go build ./...
cd /Users/av/Desktop/WORK/documan/documan-bff && go build ./...
cd /Users/av/Desktop/WORK/documan/documan-fe-my && npx nuxi typecheck
```

### 2. Grep-аудит legacy payload references (ВСЕ сервисы)
```bash
# Legacy payload field names — ДОЛЖНЫ быть полностью удалены
grep -rn "amount_with_vat\|source_name\|line_no\|commodity_code\|unit_price_with_vat\|legalTitle\|companyType\|pathName\|assortment_id\|assortment_type\|product_code" \
  documan-pkg/ documan-core/ documan-connector-ms/ documan-ingest/ documan-bff/ documan-fe-my/ \
  --include="*.go" --include="*.proto" --include="*.vue" --include="*.ts" \
  | grep -v "node_modules" | grep -v ".gen." | grep -v "_archive_del"
```

Любые найденные legacy payload references к финалу миграции = **major**.

### 2b. Grep-аудит refs-context для seller_id / buyer_id / store_id
```bash
# Эти идентификаторы разрешены только в contract/document/refs.go,
# document.Meta.Refs и коде core/connector, который читает или пишет payload_meta.refs.
grep -rn "seller_id\|buyer_id\|store_id" \
  documan-pkg/ documan-core/ documan-connector-ms/ documan-ingest/ documan-bff/ documan-fe-my/ \
  --include="*.go" --include="*.proto" --include="*.vue" --include="*.ts" \
  | grep -v "node_modules" | grep -v ".gen." | grep -v "_archive_del"
```

Ручная проверка результатов обязательна:
- Разрешено: `contract/document/refs.go`, `document.Meta.Refs`, код `connector-ms`/`core`, который пишет или читает `payload_meta.refs.*`
- Запрещено: `Header`, top-level payload fields, ожидания UI/BFF от raw payload `seller_id/buyer_id/store_id`
- Любое нарушение refs-context = **major**

### 3. Grep-аудит терминологического rename (все типы файлов)
```bash
# Код + конфиг + скрипты + документация
grep -rn "\"transform\"\|\"clear\"\|\"reclear\"\|transformed/" \
  documan-connector-ms/ documan-ingest/ \
  --include="*.go" --include="*.yaml" --include="*.yml" \
  --include="*.sh" --include="*.md" \
  --include="Makefile" --include="Dockerfile*" \
  | grep -v "ClearGroupID" | grep -v "_archive_del" | grep -v "vendor/"

# .env файлы
find documan-connector-ms/ documan-ingest/ -name ".env*" \
  -exec grep -Hn "transform\|clear\|reclear\|transformed/" {} \;
```
Любое нахождение legacy terminology = **major**.

### 4. Bruno
- Прогнать smoke/debug коллекции
- Обновить коллекции, если есть запросы/проверки на old query/matching response fields
- По текущему дереву payload-specific coverage почти нет — работы минимум

### 5. Аудит относительно ТЗ
Пройти по каждому разделу ТЗ (1.1–1.8 и 2.1–2.7) и подтвердить:
- Все перечисленные изменения выполнены
- Нет пропущенных файлов
- Нет оставшихся legacy references
- `seller_id` / `buyer_id` / `store_id` используются только в допустимом refs-context

### 6. Проверка proto/gen синхронизации
```bash
cd documan-pkg/services && buf generate && git diff --exit-code gen/
```
Если diff не пустой — proto и gen-код рассинхронизированы = **major**.

### 7. Аудит на запрещённые паттерны
Любое нахождение по любой из проверок ниже = **major**.

#### 7a. Deprecated endpoints и маркеры
```bash
grep -rn "deprecated\|@deprecated\|// removed\|// legacy\|// old\|// compat" \
  documan-pkg/ documan-core/ documan-connector-ms/ documan-ingest/ documan-bff/ \
  --include="*.go" --include="*.proto" \
  | grep -v "vendor/" | grep -v "_archive_del" | grep -v ".gen."
```

#### 7b. Duplicate config/env keys (старые + новые имена одновременно)
```bash
# Проверить, что старые env/config names не сосуществуют с новыми
# Искать legacy names в config-файлах всех затронутых сервисов
grep -rn "TRANSFORM\|CLEAR_\|RECLEAR\|transformed" \
  documan-connector-ms/ documan-ingest/ \
  --include="*.yaml" --include="*.yml" --include="*.env*" \
  --include="Makefile" --include="Dockerfile*" --include="*.sh" \
  | grep -v "_archive_del"

find documan-connector-ms/ documan-ingest/ -name ".env*" \
  -exec grep -Hin "TRANSFORM\|CLEAR_\|RECLEAR\|transformed" {} \;
```

#### 7c. Bridge-адаптеры, compatibility wrappers, alias-поля
```bash
grep -rn "bridge\|compat\|wrapper\|alias\|legacy\|backward\|fallback.*legacy\|dual.write\|dual.read" \
  documan-pkg/ documan-core/ documan-connector-ms/ documan-ingest/ documan-bff/ \
  --include="*.go" \
  | grep -v "vendor/" | grep -v "_archive_del" | grep -v "_test.go"
```
**Примечание**: результаты требуют ручной оценки — слова `bridge`/`wrapper` могут использоваться в не-legacy контексте. Но каждое нахождение должно быть явно проверено и подтверждено как не-нарушение.

### 8. Cutover preparation checklist
- [ ] SQL migration готова
- [ ] Все сервисы собираются
- [ ] Config/env ключи обновлены (дублирование old+new запрещено)
- [ ] S3 paths обновлены на `payload/`
- [ ] Repayload workflow готов для connector-ms и ingest
- [ ] Bruno smoke checks проходят
- [ ] Proto и gen-код синхронизированы

## Критерии готовности
- [ ] Cross-repo build — все pass
- [ ] Grep-аудит legacy payload — чисто
- [ ] Refs-context audit для `seller_id` / `buyer_id` / `store_id` пройден
- [ ] Grep-аудит rename — чисто
- [ ] Bruno обновлён
- [ ] Аудит по ТЗ — все разделы покрыты
- [ ] Cutover checklist заполнен

## Gate check
Cross-repo build (см. выше) + grep-аудит.

## После завершения
1. Записать финальный блок в `done_log.md`.
2. Обновить `actual_state.md` — статус "ЗАВЕРШЕНО".
3. Подготовить финальный отчёт соответствия ТЗ.
4. Коммиты по каждому затронутому сервису:
```
DONE: REFACTOR_CONTRACT stage 11: verification + cutover prep
```

## Stop/Continue
- Если grep-аудит находит legacy payload references — это **major**. Исправить до закрытия этапа.
- Если build fails — это **major**. Исправить до закрытия.
- Если Bruno выявляет legacy field names в запросах/ответах — обновить коллекции на этом этапе. Legacy references в тестах = **major**.

## Финал проекта
После этого этапа:
1. Провести аудит относительно ТЗ.
2. Все `tracked` deviations с промежуточных этапов проверяются здесь. Если не закрыты — автоэскалация в `major`.
3. Если отклонения minor (НЕ связанные с legacy references) — исправить корректирующими этапами (12+) с логом и коммитами. Любое отклонение требует явного согласования с пользователем.
4. Если major — СТОП и эскалация.
5. Подготовить финальный отчёт соответствия ТЗ.
