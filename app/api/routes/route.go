package route

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DmitriiKumancev/lamoda-test/api/controller"

	_ "github.com/DmitriiKumancev/lamoda-test/docs"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)


type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	r.POST("/create-warehouse", func(c *gin.Context) {
		var w controller.Warehouse
		err := c.BindJSON(&w)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse data"})
			return
		}

		err = controller.CreateWarehouse(db, &w)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": w.ID})
	})

	r.POST("/create-product", func(c *gin.Context) {
		var p controller.Product
		err := c.BindJSON(&p)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
			return
		}

		err = controller.CreateProduct(db, &p)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": p.ID})
	})

	r.DELETE("/delete-product/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		if err := controller.DeleteProduct(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	r.POST("/reserve-products", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}

		err := controller.ReserveProducts(db, productCodes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})

	r.POST("/release-products", func(c *gin.Context) {
		var productCodes []string
		if err := c.ShouldBindJSON(&productCodes); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}

		err := controller.ReleaseProducts(db, productCodes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})

	r.GET("/remaining-products/:warehouseID", func(c *gin.Context) {
		warehouseID := c.Param("warehouseID")
		var id int
		if _, err := fmt.Sscan(warehouseID, &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid warehouse ID",
			})
			return
		}

		products, err := controller.GetRemainingProducts(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, products)
	})

	return r
}
