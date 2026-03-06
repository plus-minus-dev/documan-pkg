# INIT: Перевод проекта на финальный contract

## Режим
Ты работаешь по пайплайну `REFACTOR_PROJECT_TO_CONTRACT`.
Это мульти-сервисная миграция всего проекта documan на финальный payload contract.

## Перед любым действием — ОБЯЗАТЕЛЬНО прочитай

1. **Мастер-план**: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/master_plan.md`
2. **Лог выполнения**: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/done_log.md`
3. **Актуальное состояние**: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/actual_state.md`
4. **ТЗ (источник истины)**: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md`
5. **Архитектурные правила**: `skills/Architecture.md`, `skills/Architecture_platform.md`
6. **Обязательные правила выполнения**: раздел "Обязательные правила выполнения" в ТЗ

## После прочтения — выведи

```
Last completed stage: <N> или "none"
Open deviations: <список> или "none"
Current stage goal: <цель следующего этапа>
```

## Затем

1. Найди stage prompt для текущего этапа: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/stage_NN_*.md`
2. Прочитай его и выполняй.

## Директория артефактов
`/Users/av/Desktop/WORK/documan/documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`

## Сервисы проекта (все в `/Users/av/Desktop/WORK/documan/`)

| Сервис | Путь | Язык | Gate check |
|--------|------|------|------------|
| documan-pkg | `documan-pkg/` | Go | `go build ./...` |
| documan-core | `documan-core/` | Go | `go build ./...` |
| documan-connector-ms | `documan-connector-ms/` | Go | `go build ./...` |
| documan-ingest | `documan-ingest/` | Go | `go build ./...` |
| documan-bff | `documan-bff/` | Go | `go build ./...` |
| documan-fe-my | `documan-fe-my/` | TypeScript/Vue | `npx nuxi typecheck` |
| documan-bruno | `documan-bruno/` | Bruno | ручной прогон коллекций |

## Формат коммита
В коммите вместо обычного префикса ставить:
```
DONE: REFACTOR_CONTRACT stage N: <краткое описание>
```
Коммит делается в директории затронутого сервиса (каждый сервис — отдельный git repo).
