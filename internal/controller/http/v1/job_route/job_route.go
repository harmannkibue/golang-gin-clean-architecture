package job_route

import (
	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
	"net/http"
)

type JobRoute struct {
	u intfaces.IntJobUsecase
	l logger.Interface
}

// NewJobRoute Initialises a new http router for the blogs -.
func NewJobRoute(handler *gin.RouterGroup, t intfaces.IntJobUsecase, l logger.Interface) {
	r := &JobRoute{t, l}

	h := handler.Group("/blogs")
	{
		h.POST("/create-job/", r.createJob)
	}
}

type createJobRequestBody struct {
	Description string `json:"description"`
}

type createJobResponse struct {
	Description string `json:"message"`
	CreatedAt   string `json:"created_at"`
}

// @Summary     Create a Job
// @Description Create a Job
// @ID          Create a Job
// @Tags  	    Job
// @Accept      json
// @Produce     json
// @Param       request body createJobRequestBody true "Create blog request body"
// @Success     201 {object} createJobResponse
// @Failure     400 {object} entity.HTTPError
// @Failure     500 {object} entity.HTTPError
// @Router      /jobs/create-job/ [post]
func (route *JobRoute) createJob(ctx *gin.Context) {
	var body createJobRequestBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		route.l.Error(err, "http - v1 - create a job route")
		ctx.JSON(entity.GetStatusCode(err), entity.ErrorCodeResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}
