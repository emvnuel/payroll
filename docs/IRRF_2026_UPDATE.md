# Atualização do Cálculo de IRRF - Lei nº 15.270/2025

## Resumo

Implementação das novas regras de cálculo do Imposto de Renda Retido na Fonte (IRRF) conforme a **Lei nº 15.270, de 26 de novembro de 2025**, válidas a partir de **1º de janeiro de 2026**.

## Principais Mudanças

### 1. Ampliação da Faixa de Isenção

**Até R$ 5.000,00/mês:**
- Redução de até **R$ 312,89** no imposto calculado
- Na prática, rendimentos até R$ 5.000,00 ficam isentos de IRRF

**Entre R$ 5.000,01 e R$ 7.350,00/mês:**
- Redução gradual calculada pela fórmula: **R$ 978,62 - (0,133145 × rendimento)**
- Quanto menor o rendimento, maior a redução

**Acima de R$ 7.350,00/mês:**
- Sem redução
- Aplicação normal da tabela progressiva

### 2. Tabela Progressiva Atualizada (2026)

| Faixa | Base de Cálculo (R$) | Alíquota | Dedução (R$) |
|-------|---------------------|----------|--------------|
| 1ª    | Até 2.428,80       | 0%       | 0,00         |
| 2ª    | 2.428,81 a 2.826,65| 7,5%     | 182,16       |
| 3ª    | 2.826,66 a 3.751,05| 15%      | 394,16       |
| 4ª    | 3.751,06 a 4.664,68| 22,5%    | 675,49       |
| 5ª    | Acima de 4.664,69  | 27,5%    | 908,73       |

## Configuração

### Variáveis de Ambiente

As seguintes variáveis foram adicionadas ao ConfigMap:

```yaml
# Tabela IRRF atualizada para 2026
IRRF_RANGES: '[{"init_value":0,"end_value":2428.80,"aliquot":0,"deduction":0},{"init_value":2428.81,"end_value":2826.65,"aliquot":0.075,"deduction":182.16},{"init_value":2826.66,"end_value":3751.05,"aliquot":0.15,"deduction":394.16},{"init_value":3751.06,"end_value":4664.68,"aliquot":0.225,"deduction":675.49},{"init_value":4664.69,"end_value":999999999.99,"aliquot":0.275,"deduction":908.73}]'

# Configurações da redução (Lei 15.270/2025)
IRRF_MAX_REDUCTION_AMOUNT: "312.89"           # Redução máxima
IRRF_REDUCTION_THRESHOLD: "5000.00"           # Limite para redução total
IRRF_REDUCTION_UPPER_LIMIT: "7350.00"         # Limite superior da redução
IRRF_REDUCTION_CONSTANT: "978.62"             # Constante da fórmula
IRRF_REDUCTION_MULTIPLIER: "0.133145"         # Multiplicador da fórmula
IRRF_SIMPLIFIED_DEDUCTION_PERCENTAGE: "0.25"  # 25% desconto simplificado
```

### Valores Padrão

Todas as novas variáveis possuem valores padrão embutidos no código, conforme a legislação vigente.

## Exemplos de Cálculo

### Exemplo 1: Rendimento de R$ 4.500,00

```
Rendimento Bruto: R$ 4.500,00
Desconto Simplificado: R$ 607,20 (25% de R$ 2.428,80)
Base de Cálculo: R$ 3.892,80
Imposto pela Tabela: R$ 267,32
Redução Aplicada: R$ 267,32 (limitado ao imposto)
IRRF Final: R$ 0,00 ✓
```

### Exemplo 2: Rendimento de R$ 6.000,00

```
Rendimento Bruto: R$ 6.000,00
Desconto Simplificado: R$ 607,20
Base de Cálculo: R$ 5.392,80
Imposto pela Tabela: R$ 577,21
Redução Gradual: R$ 978,62 - (0,133145 × 6.000) = R$ 179,75
IRRF Final: R$ 397,46
```

### Exemplo 3: Rendimento de R$ 8.000,00

```
Rendimento Bruto: R$ 8.000,00
Desconto Simplificado: R$ 607,20
Base de Cálculo: R$ 7.392,80
Imposto pela Tabela: R$ 1.125,20
Redução Aplicada: R$ 0,00 (acima de R$ 7.350,00)
IRRF Final: R$ 1.125,20
```

## Implementação Técnica

### Arquivos Modificados

1. **models/irrf_discount.go**
   - Adicionadas constantes configuráveis via ambiente
   - Implementada função `calculateReduction()` para aplicar a redução
   - Atualizado método `Value()` para incluir a redução no cálculo final
   - Adicionada função helper `getEnvOrDefault()`

2. **models/irrf_discount_test.go** (novo)
   - Testes unitários validando os exemplos da Receita Federal
   - Testes de casos extremos e limites
   - Cobertura de diferentes cenários (com/sem dependentes, com/sem INSS)

### Fluxo de Cálculo

1. Calcula a base tributável (rendimento - deduções)
2. Aplica a tabela progressiva para obter o imposto
3. Calcula a redução aplicável baseado no rendimento bruto
4. Subtrai a redução do imposto calculado
5. Retorna o valor final (nunca negativo)

## Testes

Todos os testes unitários passam com sucesso:

```bash
go test ./models -v -run TestIRRF
```

Cobertura de testes:
- ✓ Isenção até R$ 5.000,00
- ✓ Redução gradual entre R$ 5.000,01 e R$ 7.350,00
- ✓ Sem redução acima de R$ 7.350,00
- ✓ Cálculos com dependentes
- ✓ Cálculos com dedução de INSS
- ✓ Desconto simplificado

## Referências

- [Lei nº 15.270/2025](https://www.planalto.gov.br/ccivil_03/_ato2023-2026/2025/lei/l15270.htm)
- [Orientação da Receita Federal](https://www.gov.br/receitafederal/pt-br/assuntos/noticias/2025/dezembro/receita-federal-orienta-fontes-pagadoras-e-contribuintes-a-calcular-a-reducao-do-imposto-de-renda-a-partir-de-1o-de-janeiro-de-2026)

## Vigência

As novas regras são aplicáveis a partir de **1º de janeiro de 2026**.

---

**Data da Implementação:** 01/01/2026  
**Desenvolvedor:** Sistema de Folha de Pagamento
