/*
Copyright 2021 Syntasso.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	platformv1alpha1 "github.com/syntasso/kratix/api/v1alpha1"
)

// WorkReconciler reconciles a Work object
type WorkReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=platform.kratix.io,resources=works,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=platform.kratix.io,resources=works/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=platform.kratix.io,resources=works/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Work object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *WorkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("work", req.NamespacedName)

	r.Log.Info("Reconciling Work " + req.Name)

	work := &platformv1alpha1.Work{}
	err := r.Client.Get(context.Background(), req.NamespacedName, work)
	if err != nil {
		r.Log.Error(err, "Error getting Work")
		return ctrl.Result{Requeue: false}, err
	}

	r.Log.Info("Decorating Work " + req.Name + " with Label: cluster=worker")
	labels := map[string]string{}
	labels["cluster"] = "worker"
	work.SetLabels(labels)

	err = r.Client.Update(context.Background(), work)
	if err != nil {
		r.Log.Error(err, "Error updating Work with new Labels")
		return ctrl.Result{Requeue: false}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.Work{}).
		Complete(r)
}
