# Актуальное состояние: REFACTOR_PROJECT_TO_CONTRACT

> Обновляется после каждого этапа.

## Текущий статус
- **Последний завершённый этап**: 11 (verification + cutover preparation)
- **Текущий этап**: ЗАВЕРШЕНО
- **Открытые отклонения**: none

## Состояние по сервисам

| Сервис | Статус | Последний этап | Заметки |
|--------|--------|---------------|---------|
| documan-pkg | proto + codegen + contract готовы | 1+ | etl.proto, docmatching.proto обновлены; contract/document и contract/entity обновлены (flat ref-модель) |
| documan-core | полностью обновлён | 4 | ETL + consumers + query enrichment; flat ref-модель; EnrichERPDocPayload для query layer |
| documan-connector-ms | полностью обновлён | 6 | contract structs (doc + entity), terminology rename (transform→payload), S3 path payload/, DocPayloadVersion 2.0, EntityPayloadVersion 2.0 |
| documan-ingest | полностью обновлён | 8 | contract structs + full terminology rename (clear→payload, reclear→repayload, ClearOutput→PayloadOutput, clearedAt→payloadedAt), S3 path payload/ |
| documan-bff | полностью обновлён | 9 | AmountWithVat→AmountWithTax, go.mod@acde77acd872, payload proxied as json.RawMessage |
| documan-fe-my | полностью обновлён | 10 | contract field names, ERP kopecks semantics, reclear→repayload |
| documan-bruno | проверен | 11 | нет legacy references |

## Блокеры
none
