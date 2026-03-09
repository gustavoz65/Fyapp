# CSV Parser Bug Fixes - Design Document

**Date**: 2026-03-09
**Status**: Approved
**Type**: Bug Fixes + Improvements

## Problem Statement

The CSV parser has critical bugs causing data corruption in the dashboard:

1. **Encoding Issues**: BB CSV files use ISO-8859-1/Windows-1252, causing corrupted characters (`Lan�amento`, `Cart�o`)
2. **Transaction Duplication**: BB shows same transaction twice (credit entry + debit exit)
3. **Control Lines**: "Saldo Anterior", "Saldo do dia" imported as transactions
4. **Invalid Dates**: "00/00/0000" dates breaking dashboard charts
5. **Categorization Errors**: "Mercado Pago" misclassified as "Alimentação" (food)

## Solution Architecture

### 1. Encoding Detection & Conversion

**New Module**: `internal/lib/encoding/detector.go`

- Auto-detect CSV charset (UTF-8, ISO-8859-1, Windows-1252)
- Convert to UTF-8 before parsing
- Uses `golang.org/x/text/encoding`

### 2. Enhanced Filtering

**New Parser Methods**:
- `isControlLine()` - detects balance/summary lines
- `isDuplicatePayment()` - detects BB "Pagamento Pix Cartão" duplicate entries
- Execute filters BEFORE parsing attempt

**Filter Keywords**:
- "Saldo Anterior", "Saldo do dia", "S A L D O"
- "Pagamento Pix Cartão Crédito" (BB only)
- Empty "Tipo Lançamento" column (BB)

### 3. Improved Categorization

**Smart Context Matching**:
- Multi-word matches have priority over single words
- "Mercado Pago" (2 words) > "Mercado" (1 word)

**Exclusion Rules** (`rules.json`):
```json
"exclusion_patterns": [
  {
    "keywords": ["mercado pago", "pagamento mercado pago"],
    "exclude_categories": ["alimentação"],
    "reason": "payment_platform"
  }
]
```

**Payment Platform Detection**:
- Mercado Pago, PagSeguro, PicPay → "Transferências" or uncategorized
- Prevents false positives from generic words

### 4. Validation Pipeline

```
CSV Line → Detect Encoding → Convert UTF-8 →
Is Header? → Is Control Line? → Is BB Duplicate? →
Parse Date (reject 00/00/0000) → Parse Amount → Create Transaction
```

## Implementation Order

1. Create encoding detector module
2. Add control line filters
3. Fix BB duplicate detection
4. Improve categorization rules
5. Add categorization exclusions
6. Write tests for each fix
7. Add missing category rules (deferred to end)

## Testing Strategy

- Unit tests for each filter function
- Integration tests with real BB CSV samples
- Encoding detection tests (UTF-8, ISO-8859-1, Windows-1252)
- Categorization regression tests

## Success Criteria

- ✅ BB CSV imports without encoding corruption
- ✅ No duplicate transactions from BB
- ✅ Dashboard shows valid dates
- ✅ "Mercado Pago" not categorized as food
- ✅ Balance calculations correct

## Out of Scope

- Generic pattern detection engine (future work)
- Support for new banks (can be added later)
- Advanced ML categorization (future enhancement)
