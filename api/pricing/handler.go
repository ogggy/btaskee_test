package pricing

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Apply(apiEngine *gin.Engine) {
	apiEngine.GET("/pricing", h.calculatePricing)
}

func (h *Handler) calculatePricing(c *gin.Context) {
	var (
		resp = CalculatePricingResp{}
		err  error
	)
	param := c.DefaultQuery("date", "")

	resp, err = h.s.CalculatePricing(c.Request.Context(), param)
	if err != nil {
		log.Println("[CalculatePricing-Service-Error]", err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{}
}
