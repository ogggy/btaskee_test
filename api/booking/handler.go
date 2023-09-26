package booking

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) Apply(apiEngine *gin.Engine) {
	apiEngine.POST("/booking/new", h.createNewBooking)
	apiEngine.GET("/booking", h.getBookingInfo)
}

func (h *Handler) createNewBooking(c *gin.Context) {

	var (
		req  CreateBookingReq
		resp CreateBookingResp
	)

	if err := c.BindJSON(&req); err != nil {
		resp.ErrCode = -1
		resp.ErrMessage = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.s.CreateNewBooking(c.Request.Context(), &req)
	if err != nil {
		println("[CreateNewBooking-Service-Error] ", err)
	}
	// println("[CreateNewBooking-Service-Response] ", resp)

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getBookingInfo(c *gin.Context) {

	var (
		resp = GetBookingInfoResp{}
		err  error
	)
	param := c.DefaultQuery("id", "")
	bookingID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		resp.ErrCode = 40
		resp.ErrMessage = "invalid booking id"
	}
	resp, err = h.s.GetBookingInfo(c.Request.Context(), bookingID)
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
	return &Handler{s: s}
}
