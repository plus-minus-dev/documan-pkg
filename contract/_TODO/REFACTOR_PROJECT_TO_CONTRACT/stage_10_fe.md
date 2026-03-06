# Stage 10: documan-fe-my — UI migration

## Anti-context-loss пролог (ОБЯЗАТЕЛЬНО)
1. Прочитай `init.md`, `master_plan.md`, `done_log.md`, `actual_state.md` из директории `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT/`.
2. Прочитай ТЗ: `documan-pkg/contract/_TODO/REFACTOR_PROJECT_TO_CONTRACT.md` — раздел "Обязательные правила выполнения" + раздел 2.6.
3. Выведи: `Last completed stage`, `Open deviations`, `Current stage goal`.
4. Только после этого начинай работу.

---

## Цель
Обновить frontend (documan-fe-my) под новый contract vocabulary и semantics.

## Сервис
- **Имя**: documan-fe-my
- **Путь**: найти через `find /Users/av/Desktop/WORK/documan/documan-fe-my -name "package.json" -maxdepth 1`

## Входы
- Найти composables: `grep -rl "useDocworker\|useErpDocuments\|useMatching" documan-fe-my/app/`
- Найти document pages: `find documan-fe-my/app/pages/documents/ -name "*.vue"`
- Найти ERP pages: `find documan-fe-my/app/pages/documents_erp/ -name "*.vue"`
- Найти matching UI: `grep -rl "DeltaViewer\|matching" documan-fe-my/app/components/`
- Проверить done_log (отклонения этапа 9)

## Что сделать

### 1. Original document pages
Файлы (искать через grep/find):
- `useDocworker.ts`
- `documents/[id].vue`
- `documents/new/[id].vue`
- `documents/new2/[id].vue`

Перевести на:
- `amount_with_tax` вместо `amount_with_vat`
- `tax_amount` вместо `vat` / `vat_amount`
- `line_number`, `line_label` вместо `line_no`
- `article`, `item_code`
- `uom_name`, `uom_code`
- `hs_code` вместо `commodity_code`
- `source_number` вместо `source_name`
- Projection fields `file_*` — продолжать использовать из query DTO

### 2. ERP document pages
Файл: `documents_erp/[entityType]/[id].vue`
- Новые payload names
- Убрать ожидание raw `seller_id`, `buyer_id`, `source_name`
- Seller/buyer/store/item business fields для ERP UI брать из уже обогащённого query/API ответа `core`/`bff`, а не из raw `seller_id` / `buyer_id` / `store_id` / `item_id`
- **ВАЖНО**: убрать допущение "ERP summary = рубли". В новом контракте деньги = копейки. Это behavioural change, не только rename.

### 3. Matching UI
Файл: `DeltaViewer.vue`
- Обновить label mappings под новый `docmatching` vocabulary
- `amount_with_vat` → `amount_with_tax`

### 4. Composables
Переписать:
- `useDocworker.ts`
- `useErpDocuments.ts`
- `useMatching.ts` (если matching API изменился)

## Критерии готовности
- [ ] Нет legacy field names в Vue/TS файлах
- [ ] ERP summary читается как копейки
- [ ] Matching labels обновлены
- [ ] Composables обновлены
- [ ] `npx nuxi typecheck` из `documan-fe-my/` — pass

## Gate check
```bash
cd /Users/av/Desktop/WORK/documan/documan-fe-my && npx nuxi typecheck
```

## Grep-аудит
```bash
grep -rn "amount_with_vat\|source_name\|line_no\|commodity_code\|unit_price_with_vat\|legalTitle" documan-fe-my/app/ --include="*.vue" --include="*.ts"
```
Результат должен быть пустым.

## После завершения
1. Записать блок в `done_log.md` по шаблону.
2. Обновить `actual_state.md`.
3. Коммит в `documan-fe-my/`:
```
DONE: REFACTOR_CONTRACT stage 10: UI migration to contract vocabulary
```

## Stop/Continue
- Если UI компонент обнаружен за пределами известных файлов — обновить на этом этапе. Оставить legacy field names в UI = `major`.
- Если ERP money semantics change ломает отображение — тестировать визуально, зафиксировать. Это behavioural change (не legacy reference), допустимо как `minor`.
- Если typecheck показывает ошибки в не-затронутых файлах — `major`, СТОП.
