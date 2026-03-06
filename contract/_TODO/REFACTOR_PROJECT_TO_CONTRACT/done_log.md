# Done Log: REFACTOR_PROJECT_TO_CONTRACT

> Формат записи — строго по шаблону ниже.
> Каждый этап добавляет один блок.

---

<!-- TEMPLATE (не удалять):
## Stage N: <title>
- **DateTime**: YYYY-MM-DD HH:MM
- **Completed**: <что сделано>
- **Artifacts**: <созданные артефакты>
- **Gate**: pass | fail
- **Actual files created**: [список]
- **Actual files modified**: [список]
- **Deviations**: none | <описание>
- **Deviation level**: minor | major
- **Next stage notes**: <заметки для следующего этапа>
- **Commit**: `<service>/ <hash> DONE: REFACTOR_CONTRACT stage N: <title>`
-->

## Stage 1: documan-pkg — Proto schemas + codegen
- **DateTime**: 2026-03-06 04:30
- **Completed**: Обновлены etl.proto (убраны erp_*, file_* из transport), docmatching.proto (amount_with_vat → amount_with_tax), regenerate gen/*. queries.proto не требовал изменений — projection fields уже в целевом состоянии.
- **Artifacts**: обновлённые proto и gen файлы
- **Gate**: pass
- **Actual files created**: нет
- **Actual files modified**: [services/proto/core/v1/etl.proto, services/proto/core/v1/docmatching.proto, services/gen/core/v1/etl.pb.go, services/gen/core/v1/docmatching.pb.go]
- **Deviations**: none
- **Deviation level**: —
- **Next stage notes**: documan-core (этап 2) сломается на компиляции — ETL controllers используют удалённые proto-поля ErpUpdatedAt/ErpCreatedAt/ErpDeletedAt/ErpArchived в ETL и FileName/FileExt/FileHash в Ingest. Это ожидаемо.
- **Commit**: ожидает
