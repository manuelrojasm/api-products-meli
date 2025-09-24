# RUN — Cómo ejecutar este proyecto localmente

Estas instrucciones presuponen que tu máquina tiene Go instalado (versión 1.21 o superior).

---

## 1. Clona el repositorio

```bash
git clone https://github.com/manuelrojasm/api-products-meli.git
cd api-products-meli
go mod tidy
go run ./cdm

```

## 2. Prueba

```bash
curl http://localhost:8080/health
curl "http://localhost:8080/products?minPrice=50&sort=price&order=desc"
curl http://localhost:8080/products/p-001

