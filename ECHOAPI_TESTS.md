# Colección de Tests para EchoAPI

## Configuración Base

**URL Base:** `http://localhost:8080`

---

## Health Check

### GET `/api/health-check`

**Headers:** Ninguno

**Body:** Ninguno

**Respuesta Esperada:** 200 OK
```json
{"status":"OK"}
```

---

## Comparación Básica - 3 Mice Gaming

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "30106bcd-f425-4dfb-8ef6-055ab4744f6c",
    "90d8335d-4d95-475d-bc0e-1bfe63decdd0"
  ]
}
```

**Descripción:** Compara 3 mice gaming (Pro Mouse HP 2, Pro Mouse Logitech 3, Pro Mouse Razer 4)

**Respuesta Esperada:** 200 OK con campos compartidos: price, rating, specifications.buttons, specifications.sensor_dpi, specifications.weight, specifications.wireless

---

## Comparación con Fields Específicos

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "30106bcd-f425-4dfb-8ef6-055ab4744f6c"
  ],
  "fields": [
    "price",
    "rating",
    "specifications.sensor_dpi"
  ]
}
```

**Descripción:** Compara solo los campos especificados

**Respuesta Esperada:** 200 OK con solo los 3 campos solicitados en shared_fields

---

## Comparación entre Categorías Diferentes (Mouse vs Monitor)

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "d760a1da-604e-4721-a83c-28a9ba346770"
  ]
}
```

**Descripción:** Compara un mouse con un monitor - solo campos comunes (price, rating, weight)

**Respuesta Esperada:** 200 OK con shared_fields limitados

---

## Comparación de Monitores

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "d760a1da-604e-4721-a83c-28a9ba346770",
    "4ac438ff-4388-479e-b546-2bfbdf2de9a4",
    "fff26536-b70f-4598-bb41-f07b25de803a",
    "d344b957-b09a-42df-bf3e-85cf815a269d"
  ]
}
```

**Descripción:** Compara 4 monitores 34" de diferentes marcas (LG, AOC, ASUS, Dell)

**Respuesta Esperada:** 200 OK con campos de monitor: screen_size, refresh_rate, resolution

---

## Comparación de Teclados Mecánicos

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "006a918d-384a-4005-ab58-2b3c68fcbd09",
    "4cb1880b-ca72-46bf-83ef-5167e6d1ed03",
    "3fea096f-20da-483a-823a-1033885cc7cf"
  ],
  "fields": [
    "price",
    "rating",
    "specifications.switch_type",
    "specifications.backlit",
    "specifications.wireless"
  ]
}
```

**Descripción:** Compara 3 teclados mecánicos (Keychron, Varmilo, Logitech) con campos específicos

**Respuesta Esperada:** 200 OK

---

## Comparación de Audífonos

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "8288c3f7-e55d-40bc-a7f5-0330e3517127",
    "012326b6-49ff-4407-9433-e7a752096186",
    "1cb7d427-5775-4859-99f3-dee1b08003bb",
    "c60b85b7-a3bc-4507-8376-2a3e99ff05b5"
  ]
}
```

**Descripción:** Compara 4 audífonos (Anker, Sennheiser, Apple, Bose)

**Respuesta Esperada:** 200 OK con campos: battery_life, noise_cancelling, wireless

---

## Test de Cache (Segunda Request Idéntica)

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "ids": [
    "c17895dd-dffd-4ae7-89d0-89ef613bb219",
    "7504f71f-5dbf-4297-b46f-5e667ea98be1"
  ]
}
```

**Descripción:** Ejecuta esta request 2 veces seguidas

**Respuesta Esperada:** 
- Primera vez: Header `Cache-Status: miss`
- Segunda vez: Header `Cache-Status: hit` (más rápida)

---

## Test de Idempotencia - Primera Request

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
Idempotency-Key: test-key-12345
```

**Body:**
```json
{
  "ids": [
    "dd5fe08a-5ae3-468e-80aa-e3524212e966",
    "0d8310c0-b2f0-4bd2-80dc-0943ae40f3ee"
  ]
}
```

**Descripción:** Primera request con Idempotency-Key

**Respuesta Esperada:** 200 OK

---

## Test de Idempotencia - Request Duplicada

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
Idempotency-Key: test-key-12345
```

**Body:**
```json
{
  "ids": [
    "dd5fe08a-5ae3-468e-80aa-e3524212e966",
    "0d8310c0-b2f0-4bd2-80dc-0943ae40f3ee"
  ]
}
```

**Descripción:** Misma request con mismo Idempotency-Key (ejecutar después del test #9)

**Respuesta Esperada:** 200 OK (respuesta cacheada idempotente)

---

## Test de Idempotencia - Conflict

### POST `/api/v1/items/compare`

**Headers:**
```
Content-Type: application/json
Idempotency-Key: test-key-12345
```

**Body:**
```json
{
  "ids": [
    "725451db-e562-4c93-a2b6-21056cf58eb5",
    "fca65927-bf1c-43e8-834d-46a26f48a015"
  ]
}
```

**Descripción:** Mismo Idempotency-Key pero DIFERENTE body (ejecutar después del test #9)

**Respuesta Esperada:** 409 Conflict
```json
{
  "data": null,
  "metadata": null,
  "error": {
    "error_code": "Conflict",
    "message": "Request with same Idempotency-Key but different body already exists."
  }
}
```

---

## Casos de Error

### Error: Menos de 2 IDs

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d"
  ]
}
```

**Respuesta Esperada:** 422 Unprocessable Entity
```json
{
  "data": null,
  "metadata": null,
  "error": {
    "error_code": "AtLeastTwoIds",
    "message": "At least 2 unique ids are required."
  }
}
```

---

### Error: IDs Duplicados

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d"
  ]
}
```

**Respuesta Esperada:** 422 Unprocessable Entity (se filtran duplicados, queda solo 1 único)

---

### Error: ID No Encontrado

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "invalid-id-999",
    "another-invalid-id-888"
  ]
}
```

**Respuesta Esperada:** 404 Not Found
```json
{
  "data": null,
  "metadata": null,
  "error": {
    "error_code": "IdNotFound",
    "message": "Some products were not found.",
    "missing_ids": [
      "invalid-id-999",
      "another-invalid-id-888"
    ]
  }
}
```

---

### Error: Campo Faltante

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "fields": ["price", "rating"]
}
```

**Respuesta Esperada:** 400 Bad Request
```json
{
  "data": null,
  "metadata": null,
  "error": {
    "error_code": "MissingField",
    "message": "Missing mandatory field 'ids'"
  }
}
```

---

### Error: Body Vacío

**POST** `/api/v1/items/compare`

**Body:**
```json
{}
```

**Respuesta Esperada:** 400 Bad Request

---

### Error: Fields con Solo Campos Inexistentes

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "30106bcd-f425-4dfb-8ef6-055ab4744f6c"
  ],
  "fields": [
    "nonexistent_field",
    "another_bad_field"
  ]
}
```

**Respuesta Esperada:** 422 Unprocessable Entity
```json
{
  "data": null,
  "metadata": null,
  "error": {
    "error_code": "UnknownField",
    "message": "Unknown fields requested.",
    "unknown_fields": [
      "nonexistent_field",
      "another_bad_field"
    ]
  }
}
```

---

## Casos de Prueba Avanzados

### Comparación con Fields Parcialmente Válidos

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "30106bcd-f425-4dfb-8ef6-055ab4744f6c"
  ],
  "fields": [
    "price",
    "nonexistent_field",
    "rating"
  ]
}
```

**Descripción:** Mezcla de campos válidos e inválidos

**Respuesta Esperada:** 200 OK con solo los campos válidos (price, rating) y comparability_score < 1

---

###Comparación Máxima (Todos los Monitores)

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "d760a1da-604e-4721-a83c-28a9ba346770",
    "4ac438ff-4388-479e-b546-2bfbdf2de9a4",
    "fff26536-b70f-4598-bb41-f07b25de803a",
    "d344b957-b09a-42df-bf3e-85cf815a269d",
    "dd5fe08a-5ae3-468e-80aa-e3524212e966",
    "0d8310c0-b2f0-4bd2-80dc-0943ae40f3ee",
    "725451db-e562-4c93-a2b6-21056cf58eb5",
    "fca65927-bf1c-43e8-834d-46a26f48a015"
  ]
}
```

**Descripción:** Compara 8 monitores simultáneamente

**Respuesta Esperada:** 200 OK con comparación completa

---

### Comparación con Solo Campos de Especificaciones

**POST** `/api/v1/items/compare`

**Body:**
```json
{
  "ids": [
    "4897b2e4-fb8f-4aa3-b35a-a90594eb0d4d",
    "30106bcd-f425-4dfb-8ef6-055ab4744f6c",
    "90d8335d-4d95-475d-bc0e-1bfe63decdd0"
  ],
  "fields": [
    "specifications.sensor_dpi",
    "specifications.buttons",
    "specifications.wireless"
  ]
}
```

**Descripción:** Solicita solo campos anidados en specifications

**Respuesta Esperada:** 200 OK con solo esos campos

