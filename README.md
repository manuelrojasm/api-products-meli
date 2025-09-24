# api-products-meli

API REST para gestionar productos con arquitectura modular en Go.

---

## 📦 Estructura del proyecto

- `cmd/` — punto de entrada de la aplicación (main, wiring).  
- `internal/` — código interno (dominio, casos de uso, adaptadores).  
- `products.json` / `products.csv` — datos locales de ejemplo.  
- `go.mod`, `go.sum` — módulos y dependencias.  

---

## 🚀 Características

- Arquitectura modular / limpia (hexagonal / ports & adapters).  
- Soporta lectura de productos desde JSON y/o CSV.  
- Endpoints REST para listar productos, filtrar, obtener por ID.  
- Potencial para extender con otros recursos (como usuarios).  
- Tests y estructura desacoplada.

---

## 🛠️ Requisitos

- Go versión 1.21+  
- (Opcional) `curl` o cliente HTTP para probar la API  

---

## 🏃 Cómo ejecutar

Ver `run.md` para instrucciones detalladas.

---

## 📖 Swagger

La API expone su documentación OpenAPI/Swagger en:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Ahí 

---

## 📋 Endpoints principales

| Verbo | Ruta                  | Descripción                               |
|-------|------------------------|--------------------------------------------|
| GET   | `/health`               | Verifica que la API está viva              |
| GET   | `/products`             | Lista productos (con filtros: q, minPrice, maxPrice, sort, order) |
| GET   | `/products/:id`          | Devuelve detalle de un producto por ID     |

Ejemplos:

```bash
curl http://localhost:8080/health
curl "http://localhost:8080/products?minPrice=50&sort=price&order=asc"
curl http://localhost:8080/products/p-001
