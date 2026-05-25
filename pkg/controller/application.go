package controller

import (
	"context"

	"github.com/Reggles44/apperator/pkg/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Application{}).
		Complete(r)
}

func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	app := &v1alpha1.Application{}
	err := r.Get(ctx, req.NamespacedName, app)
	if err != nil {
		return ctrl.Result{}, err
	}

	subResourceLabels := map[string]string{
		"apperator.reggles44.com/managed-by": app.Name,
	}

	deployments := &appsv1.DeploymentList{}
	err = r.List(ctx, deployments, client.InNamespace(req.Namespace), )
	r.List
	if err != nil {
		return ctrl.Result{}, err
	}

	deployment := &appsv1.Deployment{}
	err = r.Create(ctx, deployment)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) createDeployment(
	app *v1alpha1.Application,
	labels map[string]string,
	ctx context.Context,
	req ctrl.Request,
) error {

	var err error
	controller := true

	err = r.Create(ctx, &appsv1.Deployment{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:        app.Spec.Name,
			Namespace:   req.Namespace,
			Labels:      labels,
			Annotations: map[string]string{},
			OwnerReferences: []v1.OwnerReference{
				{
					APIVersion:         app.APIVersion,
					Kind:               app.Kind,
					Name:               app.Name,
					UID:                app.UID,
					Controller:         &controller,
					BlockOwnerDeletion: &controller,
				},
			},
			Finalizers: []string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: new(int32),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: app.Spec.Image,
						},
					},
				},
			},
		},
	})

	return err
}

func (r *ApplicationReconciler) updateDeployment(ctx context.Context, req ctrl.Request) error {
	var err error
	return err
}

func (r *ApplicationReconciler) deleteDeployment(ctx context.Context, req ctrl.Request) error {
	var err error
	return err
}
