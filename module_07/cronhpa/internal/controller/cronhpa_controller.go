/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	autoscalingv1 "github.com/KokoiRuby/module_07_cronhpa/api/v1"
)

// CronHPAReconciler reconciles a CronHPA object
type CronHPAReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autoscaling.aiops.com,resources=cronhpas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.aiops.com,resources=cronhpas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=autoscaling.aiops.com,resources=cronhpas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CronHPA object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *CronHPAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here

	logger.Info("Reconciling CronHPA")
	var cronHPA autoscalingv1.CronHPA
	if err := r.Get(ctx, req.NamespacedName, &cronHPA); err != nil {
		logger.Error(err, "unable to get cronHPA")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	now := time.Now()
	var earliestNextRuntime *time.Time

	// iter job to get target size to replicate
	for _, job := range cronHPA.Spec.Jobs {
		// get last run time by job name from status
		lastRuntime := cronHPA.Status.LastRunTime[job.Name]
		// cal next schedule time
		nextScheduledTime, err := r.getNextScheduledTime(job.Schedule, lastRuntime.Time)
		if err != nil {
			logger.Error(err, "unable to get next scheduled time")
			return ctrl.Result{}, err
		}
		logger.Info("Job info", "name", job.Name, "lastRuntime", lastRuntime, "nextScheduledTime", nextScheduledTime)

		// check if current time reaches or exceeds scheduled
		if now.Equal(nextScheduledTime) || now.After(nextScheduledTime) {
			// update replica
			logger.Info("update replicas", "jobs", job.Name, "targetSize", job.TargetSize)
			if err := r.updateDeploymentReplicas(ctx, &cronHPA, cronHPA.Spec.ScaleTargetRef, job); err != nil {
				logger.Error(err, "unable to update replicas")
				return ctrl.Result{}, err
			}

			// update status
			cronHPA.Status.CurrentReplicas = job.TargetSize
			cronHPA.Status.LastScaleTime = &metav1.Time{Time: now}

			// update job last run time
			if cronHPA.Status.LastRunTime == nil {
				cronHPA.Status.LastRunTime = make(map[string]metav1.Time)
			}
			cronHPA.Status.LastRunTime[job.Name] = metav1.Time{Time: now}

			// cal next run time from now
			nextRuntime, _ := r.getNextScheduledTime(job.Schedule, now)
			if earliestNextRuntime == nil || earliestNextRuntime.Before(nextRuntime) {
				earliestNextRuntime = &nextRuntime
			}
		} else {
			// if not reached, set it to next run time
			if earliestNextRuntime == nil || nextScheduledTime.Before(*earliestNextRuntime) {
				earliestNextRuntime = &nextScheduledTime
			}
		}

	}

	// update status
	if err := r.Status().Update(ctx, &cronHPA); err != nil {
		logger.Error(err, "unable to update cronHPA")
		return ctrl.Result{}, err
	}

	// if had next run time, requeue
	if earliestNextRuntime != nil {
		requeueAfter := earliestNextRuntime.Sub(time.Now())
		// past
		if requeueAfter < 0 {
			requeueAfter = time.Second * 1
		}
		logger.Info("requeue after time", "time", requeueAfter)
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	return ctrl.Result{}, nil
}

func (r *CronHPAReconciler) getNextScheduledTime(schedule string, after time.Time) (time.Time, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cronSchedule, err := parser.Parse(schedule)
	if err != nil {
		return time.Time{}, err
	}
	return cronSchedule.Next(after), nil
}

func (r *CronHPAReconciler) updateDeploymentReplicas(ctx context.Context, cronHPA *autoscalingv1.CronHPA, scaleTargetRef autoscalingv1.ScaleTargetReference, job autoscalingv1.JobSpec) error {
	logger := log.FromContext(ctx)

	deployment := &appsv1.Deployment{}
	deploymentKey := types.NamespacedName{
		Name:      scaleTargetRef.Name,
		Namespace: cronHPA.Namespace,
	}

	// get deployment
	if err := r.Get(ctx, deploymentKey, deployment); err != nil {
		logger.Error(err, "unable to get deployment")
		return err
	}

	// chk replica
	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == job.TargetSize {
		logger.Info("Deployment replicas is already in target size.", "targetSize", job.TargetSize)
		return nil
	}

	// update replica
	deployment.Spec.Replicas = &job.TargetSize

	// update deployment
	if err := r.Update(ctx, deployment); err != nil {
		logger.Error(err, "unable to update deployment")
		return err
	}

	logger.Info("Deployment replicas updated", "targetSize", job.TargetSize)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CronHPAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1.CronHPA{}).
		Complete(r)
}
