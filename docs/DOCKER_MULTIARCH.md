# Docker Multi-Arch Image - Payroll v2.0

## Imagem Criada

✅ **Imagem:** `myawesomeapps/payroll:2.0`  
✅ **Tag adicional:** `myawesomeapps/payroll:latest`  
✅ **Registry:** Docker Hub  
✅ **Digest:** `sha256:a7f34f17fd40143cbd3350ba08f3bcb60abcb9ef9b372fd72e701f18a34cd52f`

## Arquiteturas Suportadas

A imagem foi compilada para as seguintes arquiteturas:

- ✅ **linux/amd64** - Processadores Intel/AMD de 64 bits
- ✅ **linux/arm64** - Processadores ARM de 64 bits (Apple Silicon M1/M2/M3, AWS Graviton, etc.)

## Comando de Build Utilizado

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t myawesomeapps/payroll:2.0 \
  -t myawesomeapps/payroll:latest \
  --push .
```

## Como Usar

### Pull da imagem

```bash
# Pull automático da arquitetura correta
docker pull myawesomeapps/payroll:2.0
```

### Run da imagem

```bash
docker run -d \
  --name payroll \
  -p 8080:8080 \
  -e INSS_RANGES='[...]' \
  -e IRRF_RANGES='[...]' \
  -e DEPENDENT_DEDUCTION_AMOUNT="189.59" \
  -e IRRF_MAX_REDUCTION_AMOUNT="312.89" \
  -e IRRF_REDUCTION_THRESHOLD="5000.00" \
  -e IRRF_REDUCTION_UPPER_LIMIT="7350.00" \
  -e IRRF_REDUCTION_CONSTANT="978.62" \
  -e IRRF_REDUCTION_MULTIPLIER="0.133145" \
  myawesomeapps/payroll:2.0
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payroll
spec:
  replicas: 3
  selector:
    matchLabels:
      app: payroll
  template:
    metadata:
      labels:
        app: payroll
    spec:
      containers:
      - name: payroll
        image: myawesomeapps/payroll:2.0
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: payroll-config
```

## Verificação da Imagem

Para verificar as arquiteturas disponíveis:

```bash
docker buildx imagetools inspect myawesomeapps/payroll:2.0
```

## Informações Técnicas

### Build Process
- **Base Builder:** golang:1.21-alpine
- **Runtime Base:** alpine:latest
- **Build Type:** Multi-stage build
- **Size:** Otimizada com Alpine Linux
- **Security:** Inclui ca-certificates para conexões HTTPS

### Funcionalidades v2.0
- ✅ Cálculo de IRRF com nova regra 2026 (Lei 15.270/2025)
- ✅ Redução automática até R$ 5.000,00
- ✅ Redução gradual entre R$ 5.000,01 e R$ 7.350,00
- ✅ Configuração completa via variáveis de ambiente
- ✅ Suporte a múltiplas arquiteturas

## Tags Disponíveis

- `myawesomeapps/payroll:2.0` - Versão específica com nova regra IRRF 2026
- `myawesomeapps/payroll:latest` - Sempre aponta para a versão mais recente

## Docker Hub

A imagem está disponível publicamente em:  
🔗 https://hub.docker.com/r/myawesomeapps/payroll

---

**Data de Build:** 01/01/2026  
**Desenvolvedor:** Sistema de Folha de Pagamento  
**Commit:** e27baca (feat: implementa nova regra de cálculo IRRF 2026)
