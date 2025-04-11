/*
Copyright 2025.

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
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	monitoringv1 "github.com/as960408/oswatcher-operator/api/v1"
)

// OSStatusReconciler reconciles a OSStatus object
type OSStatusReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// +kubebuilder:rbac:groups=monitoring.oswatcher.io,resources=osstatuses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.oswatcher.io,resources=osstatuses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.oswatcher.io,resources=osstatuses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OSStatus object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile

//webserver

func startWebServer(r client.Client, log logr.Logger) {
	http.HandleFunc("/report", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()

		var incoming monitoringv1.OSStatus
		if err := json.Unmarshal(body, &incoming); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 기본 네임스페이스 설정 (없을 경우 default)
		if incoming.Namespace == "" {
			incoming.Namespace = "default"
		}

		ctx := context.Background()
		var existing monitoringv1.OSStatus
		err = r.Get(ctx, client.ObjectKey{
			Name:      incoming.Name,
			Namespace: incoming.Namespace,
		}, &existing)

		if err != nil {
			// 없는 경우: 새로 생성
			if err := r.Create(ctx, &incoming); err != nil {
				http.Error(w, fmt.Sprintf("Failed to create object: %v", err), http.StatusInternalServerError)
				return
			}
			log.Info("OSStatus created", "name", incoming.Name)
		} else {
			// 있는 경우: Update
			existing.Spec = incoming.Spec
			if err := r.Update(ctx, &existing); err != nil {
				http.Error(w, fmt.Sprintf("Failed to update object: %v", err), http.StatusInternalServerError)
				return
			}
			log.Info("OSStatus updated", "name", incoming.Name)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Info("Starting API server at :8080")
	http.ListenAndServe(":8080", nil)
}

//webserverend

func (r *OSStatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var status monitoringv1.OSStatus
	if err := r.Get(ctx, req.NamespacedName, &status); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{RequeueAfter: time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OSStatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Log = ctrl.Log.WithName("controllers").WithName("OSStatus")
	go startWebServer(r.Client, r.Log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1.OSStatus{}).
		Named("osstatus").
		Complete(r)
}
