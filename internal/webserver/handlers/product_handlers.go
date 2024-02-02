package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Guilherme-daCosta/goexpert-api/internal/dto"
	"github.com/Guilherme-daCosta/goexpert-api/internal/entity"
	"github.com/Guilherme-daCosta/goexpert-api/internal/infra/database"
	entityPkg "github.com/Guilherme-daCosta/goexpert-api/pkg/entity"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productDB database.ProductInteface
}

func NewProductHandler(db database.ProductInteface) *ProductHandler {
	return &ProductHandler{
		productDB: db,
	}
}

// CreateProduct godoc
// @Summary Create a product
// @Description Create a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param request body dto.CreateProductInput true "Create Product"
// @Success 201
// @Failure 500 {object} Error
// @Router /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.productDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary Get a product
// @Description Get a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID" Format(uuid)
// @Success 200 {object} entity.Product
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products/{id} [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.productDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID" Format(uuid)
// @Param request body dto.CreateProductInput true "product request"
// @Success 200
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products/{id} [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ID, err = entityPkg.ParseId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.productDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.productDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID" Format(uuid)
// @Success 200
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products/{id} [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := h.productDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.productDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetProducts godoc
// @Summary Get products
// @Description Get products
// @Tags products
// @Accept  json
// @Produce  json
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Success 200 {array} entity.Product
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /products [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}
	sort := r.URL.Query().Get("sort")
	products, err := h.productDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
