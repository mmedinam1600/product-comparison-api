# Prompts para la elaboracion del proyecto.

Para la realizacion de este proyecto se destacan los siguientes prompts utilizados en chatGPT y Cursor

## PROMPT #1

Necesito construir una API RESTful basica que devuelva detalles por varios elementos a comparar.
- La API debe proveer campos como:
    - Nombre del producto
    - URL de la imagen
    - Descripcion
    - Precio
    - Calificacion
    - Especificaciones (Json libre dependiendo el producto)

Puedes generarme una estructura json valida para saber que datos solicitar en el request.body y cual devolver en el response.body de la solicitud?


### Resultado

Despues de estar adecuando el objeto, se llego a la estandarizacion

POST `/api/items/compare`

Request.body (Lo que solicitamos de entrada al endpoint)
```json
{
  "ids": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3", "0ac54fa6-4521-43dc-8c0d-51a26280056b", "7b4712ce-2d55-433f-8406-19cdf71b7f3e"],
  // "fields": ["price", "rating", "specifications.buttons"] // Opcional el poder especificar que campos especificos se desean comparar, si un campo especificado aqui no es compartido, entonces no se mostrara.
}
```

Response.body (Lo que devolvera la API despues de comparar los items)
```json
{
	"data": {
    "items": [
      {
        "id": "ced9542b-8abf-40cb-bd9a-bd7483bd02f3",
        "name": "Pro Mouse X ROG",
        "image_url": "https://cdn.ecommerce-shop.com/images/pro-mouse-x-rog.png",
        "description": "Mouse with RGB and six buttons customizables!",
        "price": 39.99,
        "rating": 4.4,
        "specifications": { 
          "weight": {
          	"value": 0.1,
          	"unit": "kg"
          },
        	"sensor_dpi": 16000,
          "buttons": 6,
          "wireless": true 
        }
      },
      {
        "id": "0ac54fa6-4521-43dc-8c0d-51a26280056b",
        "name": "Pro Mouse V HP",
        "image_url": "https://cdn.ecommerce-shop.com/images/pro-mouse-v-hp.png",
        "description": "Mouse with RGB and two buttons customizables!",
        "price": 49.99,
        "rating": 4.8,
        "specifications": { 
          "weight": {
          	"value": 0.2,
          	"unit": "kg"
          },
        	"sensor_dpi": 8000,
          "buttons": 8,
          "wireless": false
        }
      },
      {
        "id": "7b4712ce-2d55-433f-8406-19cdf71b7f3e",
        "name": "UltraView 27 LG",
        "image_url": "https://cdn.ecommerce-shop.com/images/monitor-ultraview-lg.png",
        "description": "Ultra View monitor 27 inches and blue light eye care.",
        "price": 229.99,
        "rating": 4.6,
        "specifications": {
          "weight": {
          	"value": 2,
            "unit": "kg"
          },
        	"screen_size": {
          	"value": 27,
            "unit": "in"
          },
          "refresh_rate": {
          	"value": 144,
            "unit": "Hz"
          },
          "resolution": {
          	"value": "2560x1440",
            "unit": "pixels"
          }
        }
      }
    ],
    "shared_fields": ["price", "rating", "specifications.weight", "specifications.sensor_dpi", "specifications.buttons", "specifications.wireless"], // Con que el campo sea comparable con otro articulo mas, ya se toma como campo compartido. En este caso hay 2 articulos de mouse que comparten propiedades y el monitor no, pero para eso se le pone null cuando la propiedad en el item no existe
    "diff": {
      "price": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 39.99,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 49.99,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": 229.99
        },
        "metric": "lower_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
      "rating": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 4.4,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 4.8,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": 4.6
        },
        "metric": "higher_is_better",
        "best": ["0ac54fa6-4521-43dc-8c0d-51a26280056b"]
      },
      "specifications.weight": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 0.1,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 0.2,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": 2
        },
        "metric": "lower_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
      "specifications.sensor_dpi": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 16000,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 8000,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": null
        },
        "metric": "higher_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
      "specifications.buttons": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 6,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 8,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": null
        },
        "metric": "higher_is_better",
        "best": ["0ac54fa6-4521-43dc-8c0d-51a26280056b"]
      },
      "specifications.wireless": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": true,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": false,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": null
        },
        "metric": "true_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
    },
  },
  "metadata": {
  	"compare_policy": {
      // "effective_mode": "intersection",
      "effective_mode": "at_least_two",
      "comparability_score": 0.25   // 0..1: qué tanto sentido tiene esta comparación
    },
  },
  "error": null
}
```

En caso de error este sera el estandar para mostrar errores
```json
{
	"data": null,
  "metadata": null,
  "error": {
  	"error_code": "AtLeastTwoIds", // It could be AtLeastTwoIds (422), IdNotFound (404), UnknownField (422), MissingField (400)
  	"message": "At least 2 unique ids are required.", // or "Some products were not found." or "Unknown fields requested." or "Missing mandatory field 'ids'"
    //"missing_ids": ["999"], // Show when error code IdNotFound
    //"unknown_fields": ["specifications.color_depth"] // Show when error code UnknownField
  }
}
```


## Prompt # 2

Generacion del nuevo endpoint

```
Quiero que implementes un nuevo endpoint con versionado `/v1` POST `/api/v1/items/compare` siguiendo arquitectura hexagonal, principios SOLID, con resiliencia, cache por request, idempotencia opcional, y comentarios claros en la logica de casos de uso. El endpoint compara ítems (productos) leídos desde un archivo JSON que se encuentra en @items.json que actúa como “base de datos”.

## Contexto y Objetivos

Endpoint: POST /api/v1/items/compare
Request json (body)

```json
{
  "ids": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3", "0ac54fa6-4521-43dc-8c0d-51a26280056b", "7b4712ce-2d55-433f-8406-19cdf71b7f3e"],
  // "fields": ["price", "rating", "specifications.buttons"] // Opcional el poder especificar que campos especificos se desean comparar, si un campo especificado aqui no es compartido, entonces no se mostrara.
}
```

Response json esperado (body)

```json
{
	"data": {
    "items": [
      {
        "id": "ced9542b-8abf-40cb-bd9a-bd7483bd02f3",
        "name": "Pro Mouse X ROG",
        "image_url": "https://cdn.ecommerce-shop.com/images/pro-mouse-x-rog.png",
        "description": "Mouse with RGB and six buttons customizables!",
        "price": 39.99,
        "rating": 4.4,
        "specifications": { 
          "weight": {
          	"value": 0.1,
          	"unit": "kg"
          },
        	"sensor_dpi": 16000,
          "buttons": 6,
          "wireless": true 
        }
      },
      {
        "id": "0ac54fa6-4521-43dc-8c0d-51a26280056b",
        "name": "Pro Mouse V HP",
        "image_url": "https://cdn.ecommerce-shop.com/images/pro-mouse-v-hp.png",
        "description": "Mouse with RGB and two buttons customizables!",
        "price": 49.99,
        "rating": 4.8,
        "specifications": { 
          "weight": {
          	"value": 0.2,
          	"unit": "kg"
          },
        	"sensor_dpi": 8000,
          "buttons": 8,
          "wireless": false
        }
      },
      {
        "id": "7b4712ce-2d55-433f-8406-19cdf71b7f3e",
        "name": "UltraView 27 LG",
        "image_url": "https://cdn.ecommerce-shop.com/images/monitor-ultraview-lg.png",
        "description": "Ultra View monitor 27 inches and blue light eye care.",
        "price": 229.99,
        "rating": 4.6,
        "specifications": {
          "weight": {
          	"value": 2,
            "unit": "kg"
          },
        	"screen_size": {
          	"value": 27,
            "unit": "in"
          },
          "refresh_rate": {
          	"value": 144,
            "unit": "Hz"
          },
          "resolution": {
          	"value": "2560x1440",
            "unit": "pixels"
          }
        }
      }
    ],
    "shared_fields": ["price", "rating", "specifications.weight", "specifications.sensor_dpi", "specifications.buttons", "specifications.wireless"], // Con que el campo sea comparable con otro articulo mas, ya se toma como campo compartido. En este caso hay 2 articulos de mouse que comparten propiedades y el monitor no, pero para eso se le pone null cuando la propiedad en el item no existe
    "diff": {
      "price": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 39.99,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 49.99,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": 229.99
        },
        "metric": "lower_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
      "rating": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 4.4,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 4.8,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": 4.6
        },
        "metric": "higher_is_better",
        "best": ["0ac54fa6-4521-43dc-8c0d-51a26280056b"]
      },
      "specifications.weight": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 0.1,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 0.2,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": 2
        },
        "metric": "lower_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
      "specifications.sensor_dpi": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 16000,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 8000,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": null
        },
        "metric": "higher_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
      "specifications.buttons": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": 6,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": 8,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": null
        },
        "metric": "higher_is_better",
        "best": ["0ac54fa6-4521-43dc-8c0d-51a26280056b"]
      },
      "specifications.wireless": {
      	"values": {
        	"ced9542b-8abf-40cb-bd9a-bd7483bd02f3": true,
          "0ac54fa6-4521-43dc-8c0d-51a26280056b": false,
          "7b4712ce-2d55-433f-8406-19cdf71b7f3e": null
        },
        "metric": "true_is_better",
        "best": ["ced9542b-8abf-40cb-bd9a-bd7483bd02f3"]
      },
    },
  },
  "metadata": {
  	"compare_policy": {
      // "effective_mode": "intersection",
      "effective_mode": "at_least_two",
      "comparability_score": 0.25   // 0..1: qué tanto sentido tiene esta comparación
    },
  },
  "error": null
}
```

Y cuando haya algun error:

```json
{
	"data": null,
  "metadata": null,
  "error": {
  	"error_code": "AtLeastTwoIds", // It could be AtLeastTwoIds (422), IdNotFound (404), UnknownField (422), MissingField (400)
  	"message": "At least 2 unique ids are required.", // or "Some products were not found." or "Unknown fields requested." or "Missing mandatory field 'ids'"
    //"missing_ids": ["999"], // Show when error code IdNotFound
    //"unknown_fields": ["specifications.color_depth"] // Show when error code UnknownField
  }
}
```

- Modo de comparación inicial: at_least_two (comparar campo si al menos 2 ítems lo tienen; los faltantes van como null). Deja el diseño listo para soportar estrategias futuras como intersection, etc. (usa Strategy Pattern).

- Response: estructura como la que ya definí (items, shared_fields, diff, metadata, error). Mantén el orden de ids en todo.

- Cargar datos desde data/items.json. Cargar en memoria e indexar por id.

## Arquitectura (paquetes y archivos)

Te recomiendo seguir esta estructura de carpetas, es la ideal para la arquitectura hexagonal, actualmente ya existe un esqueleto, puedes ir modificandolo
```
/internal/config/config.go
/internal/http/router.go
/internal/http/middleware/idempotency.go
/internal/http/handlers/compare_handler.go
/internal/domain/item.go
/internal/domain/compare.go
/internal/service/compare_service.go
/internal/service/strategy/strategy.go
/internal/service/strategy/at_least_two.go
/internal/data/catalog_repo.go
/internal/data/catalog_loader.go
/internal/cache/request_cache.go
/internal/observability/logging.go
```

## Definicion de los dominios y contratos

- domain.Item (id, name, image_url, description, price, rating, specifications[struct con campos opcionales como en mi ejemplo]).
- domain.CompareRequest (Ids []string, Fields []string).
- domain.Metric enum: lower_is_better, higher_is_better, true_is_better.
- domain.DiffField con: Values map[string]any, Metric *Metric, Best []string.
- domain.CompareResult con: Items []Item, SharedFields []string, Diff map[string]DiffField.
- domain.Metadata con: Order []string, RequestedFields *[]string, ResolvedFields []string, ComparePolicy struct{ EffectiveMode string; ComparabilityScore float64; Warnings []string }, Currency string, Version string.
- domain.ErrorResponse con ErrorCode, Message, MissingIDs []string, UnknownFields []string.

## Strategy Pattern (comparación)

strategy.Interface
```
type Interface interface {
    Name() string // "at_least_two", "intersection", etc.
    ResolveFields(items []domain.Item, requested []string) (resolved []string)
    ComputeDiff(ctx context.Context, items []domain.Item, resolved []string) (map[string]domain.DiffField, error)
}
```

### Implementa strategy.AtLeastTwo:

- ResolveFields: conserva campos presentes con al menos 2 valores no nulos, respetando `fields` si se especifico en el request.
- ComputeDiff: para cada campo:
- Numeric/boolean: aplica métrica (price → lower; rating → higher; specifications.wireless → true_is_better; specifications.weight → lower con unit:"kg"). Para establecer las metricas me gustaria tener un mapeo de el nombre del campo y especificar cual es la metrica indicada para ese campo
- Expón ComparePolicy.EffectiveMode="at_least_two".
- Deja el arnés para futuras estrategias (intersection, etc.) registradas en un map name→strategy.

## Servicio y casos de uso

### service.CompareService con método:

```
type CompareService interface {
    Compare(ctx context.Context, req domain.CompareRequest) (domain.CompareResult, domain.Metadata, *domain.ErrorResponse)
}
```

Lógica:

1. Validar req.Ids (≥2, únicos).
2. Resolver ítems en repo; si falta alguno → 404 IdNotFound con MissingIDs.
3. Seleccionar estrategia (por ahora fija: at_least_two).
4. resolvedFields := strategy.ResolveFields(items, req.Fields)
  - Si vacío → 422 UnknownField o “No comparable fields”.
5. diff := strategy.ComputeDiff(ctx, items, resolvedFields)
6. Construir CompareResult + Metadata:
- ComparePolicy.ComparabilityScore = (len(resolved)/len(baseCandidate)) (Donde baseCandidate = el conjunto de campos “objetivo” que se pretenden comparar antes de ver los datos. Concretamente:
Si el cliente envía fields → baseCandidate = normalize(fields).
Si no envía fields → baseCandidate = es igual a los shared_fields que existan en esa comparacion)
7. Retornar 200.

## Repositorio y carga de datos

data.CatalogRepository
```
type CatalogRepository interface {
    GetByIDs(ctx context.Context, ids []string) ([]domain.Item, []string /*missing*/)
}
```

- data.FileCatalogRepo
  - Carga una vez el JSON (array) al iniciar.
  - Usa parser rápido (utiliza la libreria goccy/go-json, ya esta instalada).
  - Construye map[string]Item (índice por id).
  - Guarda en atomic.Value para lecturas lock-free.

- Config

```
type Config struct {
    DataFile string "data/products.json"
    CacheTTL time.Duration `env:"CACHE_TTL" envDefault:"60s"`
    CacheSize int `env:"CACHE_SIZE" envDefault:"1000"`
    IdempotencyTTL time.Duration `env:"IDEMPOTENCY_TTL" envDefault:"5m"`
}
```

## Cache por request (memoization)

Para el manejo de cache utilizar la libreria `github.com/dgraph-io/ristretto/v2` ya instalada
- TTL configurable. Guarda resultado completo (CompareResult+Metadata).
- Incluye headers: Cache-Status: hit|miss

## Idempotencia

- Middleware Idempotency-Key:
  - Si llega header Idempotency-Key, usa ese valor como key y compara hash del body;
    - Si ya existe con body diferente → 409 Conflict.
    - Si ya existe con mismo body → devuelve respuesta cacheada (idempotente).
  - TTL configurable (IdempotencyTTL).

## Resiliencia y buenas prácticas

- HTTP server timeouts (read/write/idle).
- Panic recovery (Gin Recovery + logs).
- Validación con go-playground/validator
- Logging estructurado con zap (ya esta instalada la libreria)

## Router y handler

Utiliza Gin, ya hay un ejemplo de una ruta para el health-check, esta en @router.go 

---

Antes de implementar el nuevo endpoint si tienes alguna duda hazmelo saber
```








