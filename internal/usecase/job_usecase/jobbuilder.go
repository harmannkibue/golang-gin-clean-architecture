package job_usecase

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BuildK8sJob builds a Kubernetes Job manifest with optional compute resources, image pull secrets, and ConfigMap volume
func BuildK8sJob(jobName string, image string, envMap map[string]string, backoffLimit int32, namespace string, compute *ComputeConfig, imagePullSecretName string, configMapName string) *batchv1.Job {
	envVars := []corev1.EnvVar{}
	for k, v := range envMap {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	labels := map[string]string{
		"app":        "actsml-worker",
		"job-name":   jobName,
		"managed-by": "actsml-orchestrator",
	}

	// Default resource values
	cpuRequest := "250m"
	cpuLimit := "1"
	memoryRequest := "512Mi"
	memoryLimit := "1Gi"

	// Override with compute config if provided
	if compute != nil {
		if compute.NumCPUs > 0 {
			cpuRequest = fmt.Sprintf("%dm", compute.NumCPUs*250) // 250m per CPU
			cpuLimit = fmt.Sprintf("%d", compute.NumCPUs)
		}
		if compute.Memory != "" {
			memoryRequest = compute.Memory
			// Set limit to 1.5x request if not explicitly set
			memoryLimit = compute.Memory
		}
	}

	// Configure volume mount for payload file if ConfigMap is provided
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	if configMapName != "" {
		volumeMounts = []corev1.VolumeMount{
			{
				Name:      "payload",
				MountPath: "/workspace",
				ReadOnly:  true,
			},
		}
		volumes = []corev1.Volume{
			{
				Name: "payload",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: configMapName,
						},
					},
				},
			},
		}
	}

	container := corev1.Container{
		Name:            "worker",
		Image:           image,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env:             envVars,
		VolumeMounts:    volumeMounts,
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(cpuRequest),
				corev1.ResourceMemory: resource.MustParse(memoryRequest),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(cpuLimit),
				corev1.ResourceMemory: resource.MustParse(memoryLimit),
			},
		},
	}

	// Configure image pull secrets if provided
	var imagePullSecrets []corev1.LocalObjectReference
	if imagePullSecretName != "" {
		imagePullSecrets = []corev1.LocalObjectReference{
			{Name: imagePullSecretName},
		}
	}

	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers:         []corev1.Container{container},
			ImagePullSecrets:   imagePullSecrets,
			Volumes:            volumes,
		},
	}

	ttl := int32(3600)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template:                podTemplate,
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttl,
		},
	}

	return job
}
