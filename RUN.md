# Guía para levantar el proyecto

Instrucciones paso a paso para levantar la aplicación en local con Docker.

---

## Inicio

### Verificar prerequisitos

- Docker y docker compose

### Clonar y levantar

```bash
# Levantar todos los servicios
docker compose up -d --build
```

### Verificar que funciona

```bash
# Health check de la API
curl http://localhost:8080/api/health-check

# Respuesta esperada: {"status":"OK"}
```

---

## Acceso a Servicios

Una vez levantado el stack, ya se puede acceder a:

| Servicio | URL | Credenciales |
|----------|-----|--------------|
| **API** | http://localhost:8080 | - |
| **Grafana** | http://localhost:3000 | user: root / pass: toor |
| **Prometheus** | http://localhost:9090 | - |
| **Loki** | http://localhost:3100 | - |
| **cAdvisor** | http://localhost:8081 | - |

---

## Probar la API

Revisa la documentacion detallada en
[ECHOAPI_TEST.md](ECHOAPI_TEST.md)

---

## Grafana Dashboard

### 1. Acceder a Grafana

Url expuesta: **http://localhost:3000**

**Login:**
- Usuario: `root`
- Contraseña: `toor`

### 2. Ver Dashboard

1. Ve a **Dashboards**
2. Busca: `Product Comparison API - Overview`
3. Verás:
   - Gráfico de CPU
   - Gráfico de Memoria
   - Logs en tiempo real

### 3. Ver Logs

1. Ve a **Explore**
2. Selecciona datasource: **Loki**
3. Query: `{container="product-comparison-api"}`
4. Click **Run query**
