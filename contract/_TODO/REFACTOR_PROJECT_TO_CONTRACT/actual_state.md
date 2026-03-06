# Актуальное состояние: REFACTOR_PROJECT_TO_CONTRACT

> Обновляется после каждого этапа.

## Текущий статус
- **Последний завершённый этап**: 3 (documan-core — payload consumers)
- **Текущий этап**: 4 (ожидает запуска)
- **Открытые отклонения**: none

## Состояние по сервисам

| Сервис | Статус | Последний этап | Заметки |
|--------|--------|---------------|---------|
| documan-pkg | proto + codegen готовы | 1 | etl.proto, docmatching.proto обновлены, queries.proto без изменений |
| documan-core | ETL + consumers обновлены | 3 | ETL transport, docmatching features/candidates, controller mapping — всё на contract vocabulary |
| documan-connector-ms | не начат | — | сломается на build — ожидаемо (этап 5) |
| documan-ingest | не начат | — | сломается на build — ожидаемо (этап 7) |
| documan-bff | не начат | — | — |
| documan-fe-my | не начат | — | — |
| documan-bruno | не начат | — | — |

## Блокеры
none
