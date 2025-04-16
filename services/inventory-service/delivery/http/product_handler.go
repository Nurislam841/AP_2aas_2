package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inventory-service/domain"
	"inventory-service/internal/usecase"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(router *gin.Engine, usecase *usecase.ProductUsecase) {
	handler := &ProductHandler{usecase: usecase}
	router.GET("/products/:id", handler.GetProduct)
	router.POST("/products", handler.CreateProduct)
	router.PATCH("/products/:id", handler.UpdateProduct)
	router.DELETE("/products/:id", handler.DeleteProduct)
	router.GET("/products", handler.GetAllProducts)
}
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idInt, _ := strconv.Atoi(c.Param("id"))
	product, err := h.usecase.GetProduct(int32(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.usecase.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idInt, _ := strconv.Atoi(c.Param("id"))
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = int(idInt)
	err := h.usecase.UpdateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idInt, _ := strconv.Atoi(c.Param("id"))
	err := h.usecase.DeleteProduct(int32(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	category, _ := strconv.Atoi(c.DefaultQuery("category", ""))
	fmt.Println(category)
	page := c.DefaultQuery("page", "0")
	pageSize := c.DefaultQuery("page_size", "10")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}
	products, err := h.usecase.GetAllProducts(name, category, pageSizeInt, pageInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
