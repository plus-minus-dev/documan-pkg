# Stage 2: documan-core — ETL transport + domain models + SQL migration

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.4 (ETL transport, Domain/persistence, SQL migration).
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Привести ETL transport, доменные/persistence модели и SQL schema documan-core к новым proto-определениям из этапа 1 и подготовить flat ref-based ERP document model (`payload_header.seller_id/buyer_id/store_id`, `payload_positions[*].item_id`).

## Сервис
- **Имя**: documan-core
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-core -name "go.mod"`

## Входы
- Найти ETL контроллеры: `grep -rl "NotifyERP\|NotifyIngest" documan-core/internal/`
- Найти domain models: `grep -rl "OriginalDocument\|ERPDocument\|ERPEntity" documan-core/internal/`
- Найти SQL миграции: `find documan-core/ -path "*/migration*" -name "*.sql" | sort`
- Прочитать обновлённые proto gen: `documan-pkg/services/gen/`

## Что сделать

### 1. ETL transport — gRPC controllers
Файлы (искать через grep, не хардкод):
- `notify_erp_data.go` — контроллер
- `notify_ingest_data.go` — контроллер

#### ERP — убрать из transport/mapping:
- `ERPUpdatedAt`, `ERPCreatedAt`, `ERPDeletedAt`, `ERPArchived`

#### Ingest — убрать из transport/mapping:
- `FileName`, `FileExt`, `FileHash`

Оставить: `OriginalFileBucket`, `OriginalFileKey` или `original_file`.

### 2. ETL use case model
Найти: `grep -rl "etl" documan-core/internal/usecase/`

Синхронно обновить model structs под новые proto — убрать те же поля.

### 3. Domain / persistence model
- `payload_header.seller_id/buyer_id/store_id` и `payload_positions[*].item_id` должны сохраняться без потерь как часть payload JSON. Это source of truth для linkage ERP-документов и строк ERP-позиций.
#### `OriginalDocument`
- `FileName`, `FileExt`, `FileHash`, `OriginalFileBucket`, `OriginalFileKey` — ОСТАВИТЬ как query/storage projection columns.

#### `ERPDocument` / `ERPEntity`
- `ERPUpdatedAt`, `ERPCreatedAt`, `ERPDeletedAt`, `ERPArchived`, `Archived` — ОСТАВИТЬ как query/index projection columns.
- При upsert core наполняет их отдельно от payload, используя payload/meta как источник данных.
- На этом этапе не требуется business snapshot seller/buyer/store/item в raw payload: ERP documents могут полагаться на flat ref fields в `payload_header` / `payload_positions`.

### 4. SQL migration
Одна основная миграция. Покрыть:
- `original_documents`
- `erp_documents`
- `erp_entities`
- Связанные indexes / constraints
- Canonical `type` как internal/DB поле
- `direction` НЕ добавлять (нет готового правила и потребителя)

Конкретно:
- Имена колонок `payload_*` НЕ менять (меняется содержимое JSON payload)
- Query/storage projection columns `file_*`, `original_file_*`, `erp_*`, `archived` — сохранить
- Найти и обновить SQL-запросы с legacy JSON field names: `payload_header->>'amount_with_vat'`, `payload_summary->>'vat'` и т.п.

## Критерии готовности
- [ ] ETL controllers компилируются с новыми proto
- [ ] ETL model structs не содержат удалённых transport полей
- [ ] Domain models сохраняют projection fields
- [ ] SQL migration создана и покрывает все таблицы
- [ ] Все SQL/Go-readers legacy JSON field names обновлены
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
DONE: REFACTOR_CONTRACT stage 2: ETL transport + domain models + SQL migration
```

## Stop/Continue
- Query controllers и payload consumers НЕ трогать на этом этапе (этапы 3 и 4).
- Если SQL запросы с legacy JSON field names обнаружены за пределами ETL/domain — зафиксировать как `tracked`, обязательно обновить на этапах 3/4. Если не закрыто к stage 11 — автоэскалация в `major`.
- Если обнаружены неожиданные зависимости от удалённых proto полей — `major`, СТОП.
