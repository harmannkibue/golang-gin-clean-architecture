package job_route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity"
	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
	"github.com/harmannkibue/actsml-jobs-orchestrator/pkg/logger"
)

type JobRoute struct {
	u intfaces.IntJobUsecase
	l logger.Interface
}

func NewJobRoute(handler *gin.RouterGroup, t intfaces.IntJobUsecase, l logger.Interface) {
	r := &JobRoute{t, l}

	h := handler.Group("/jobs")
	{
		h.POST("/create", r.createJob)
		h.GET("/:id/status", r.getJobStatus)
	}
}

type createJobResponse struct {
	JobID       string `json:"job_id"`
	Status      string `json:"status"`
	K8sJobName  string `json:"k8s_job_name"`
	SubmittedAt string `json:"submitted_at"`
	ProjectID   string `json:"project_id,omitempty"`
	ExperimentID string `json:"experiment_id,omitempty"`
}

type getJobStatusResponse struct {
	JobID       string      `json:"job_id"`
	Status      string      `json:"status"`
	K8sConditions interface{} `json:"k8s_conditions"`
}

// @Summary     Create a Job
// @Description Create a Kubernetes training job
// @ID          Create a Job
// @Tags        Job
// @Accept      json
// @Produce     json
// @Param       request body object true "Job payload (raw JSON)"
// @Success     202 {object} createJobResponse
// @Failure     400 {object} entity.HTTPError
// @Failure     500 {object} entity.HTTPError
// @Router      /jobs/create [post]
func (route *JobRoute) createJob(ctx *gin.Context) {
	var raw json.RawMessage

	if err := ctx.ShouldBindJSON(&raw); err != nil {
		route.l.Error(err, "http - v1 - create job route")
		ctx.JSON(entity.GetStatusCode(err), entity.ErrorCodeResponse(err))
		return
	}

	res, err := route.u.CreateJob(ctx.Request.Context(), raw)
	if err != nil {
		route.l.Error(err, "http - v1 - create job usecase")
		ctx.JSON(entity.GetStatusCode(err), entity.ErrorCodeResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, createJobResponse{
		JobID:        res.JobID,
		Status:       res.Status,
		K8sJobName:   res.K8sJobName,
		SubmittedAt:  res.SubmittedAt.Format("2006-01-02T15:04:05Z"),
		ProjectID:    res.ProjectID,
		ExperimentID: res.ExperimentID,
	})
}

// @Summary     Get Job Status
// @Description Get status of a Kubernetes training job
// @ID          Get Job Status
// @Tags        Job
// @Produce     json
// @Param       id path string true "Job ID (actsml-job-{uuid})"
// @Success     200 {object} getJobStatusResponse
// @Failure     404 {object} entity.HTTPError
// @Failure     500 {object} entity.HTTPError
// @Router      /jobs/{id}/status [get]
func (route *JobRoute) getJobStatus(ctx *gin.Context) {
	jobID := ctx.Param("id")
	
	// Convert job ID to full job name (actsml-job-{uuid})
	// If the ID doesn't start with "actsml-job-", prepend it
	jobName := jobID
	if len(jobID) < 11 || jobID[:11] != "actsml-job-" {
		jobName = fmt.Sprintf("actsml-job-%s", jobID)
	}

	jobObj, err := route.u.GetJobStatus(ctx.Request.Context(), jobName)
	if err != nil {
		route.l.Error(err, "http - v1 - get job status")
		ctx.JSON(entity.GetStatusCode(err), entity.ErrorCodeResponse(err))
		return
	}

	// Determine status from job conditions
	status := "unknown"
	if len(jobObj.Status.Conditions) > 0 {
		latestCondition := jobObj.Status.Conditions[len(jobObj.Status.Conditions)-1]
		if latestCondition.Type == "Complete" && latestCondition.Status == "True" {
			status = "completed"
		} else if latestCondition.Type == "Failed" && latestCondition.Status == "True" {
			status = "failed"
		} else {
			status = "running"
		}
	} else if jobObj.Status.Active > 0 {
		status = "running"
	} else if jobObj.Status.Succeeded > 0 {
		status = "completed"
	} else if jobObj.Status.Failed > 0 {
		status = "failed"
	}

	// Extract job_id (UUID) from job name
	// Job name format: actsml-job-{uuid}
	actualJobID := jobID
	if strings.HasPrefix(jobObj.Name, "actsml-job-") {
		actualJobID = strings.TrimPrefix(jobObj.Name, "actsml-job-")
	}

	ctx.JSON(http.StatusOK, getJobStatusResponse{
		JobID:        actualJobID,
		Status:       status,
		K8sConditions: jobObj.Status.Conditions,
	})
}
