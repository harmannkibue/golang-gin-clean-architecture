package mpesa_ratiba_route

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/entity/intfaces"
	"github.com/harmannkibue/golang_gin_clean_architecture/pkg/logger"
	_ "github.com/swaggo/swag/example/celler/httputil"
	"net/http"
)

type MRatibaRoute struct {
	u intfaces.IntMRatibaUsecase
	l logger.Interface
}

// NewMRatibaRoute Initialises a new http router for the blogs -.
func NewMRatibaRoute(handler *gin.RouterGroup, t intfaces.IntMRatibaUsecase, l logger.Interface) {
	r := &MRatibaRoute{t, l}

	h := handler.Group("/mratiba")
	{
		h.POST("/standing-orders", r.createMpesaRatibaStandingOrder)
	}
}

// @Summary     Create M-Pesa Ratiba standing order
// @Description Create M-Pesa Ratiba standing order
// @ID          Create M-Pesa Ratiba standing order
// @Tags  	    Mpesa Ratiba
// @Accept      json
// @Produce     json
// @Param       request body intfaces.MpesaRatibaRequestBody true "Create Mpesa Ratiba standing order for subscription management"
// @Success     201 {object} intfaces.MpesaRatibaRequestResponseBody
// @Failure     400 {object} httputil.HTTPError
// @Failure     500 {object} httputil.HTTPError
// @Router      /mratiba/standing-orders [post]
func (route *MRatibaRoute) createMpesaRatibaStandingOrder(ctx *gin.Context) {
	var body intfaces.MpesaRatibaRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		route.l.Error(err, "http - v1 - create M-Ratiba standing order")
		ctx.JSON(entity.GetStatusCode(err), entity.ErrorCodeResponse(err))
		return
	}

	standingOrder, err := route.u.CreateMpesaStandingOrder(ctx, body)

	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), entity.ErrorCodeResponse(err))
	}

	ctx.JSON(http.StatusCreated, standingOrder)
}
