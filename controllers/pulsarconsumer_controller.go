/*


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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	pulsarv1 "github.com/krvarma/pulsarconsumercrd/api/v1"
)

// PulsarConsumerReconciler reconciles a PulsarConsumer object
type PulsarConsumerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *PulsarConsumerReconciler) getDeployment(pulsarcrd pulsarv1.PulsarConsumer) (appsv1.Deployment, error) {
	labels := pulsarcrd.Labels

	labels["pulsarcrd"] = pulsarcrd.Name + "-pulsarcrd"

	depl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pulsarcrd.Name + "-deployment",
			Namespace: pulsarcrd.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pulsarcrd.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "pulsarconsumer",
							Image: "krvarma/pulsarconsumer:latest",
							Env: []corev1.EnvVar{
								{Name: "PULSAR_SERVER", Value: pulsarcrd.Spec.ServerAddress},
								{Name: "PULSAR_TOPIC", Value: pulsarcrd.Spec.Topic},
								{Name: "PULSAR_SUBSCRIPTION_NAME", Value: pulsarcrd.Spec.SubscriptionName},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&pulsarcrd, &depl, r.Scheme); err != nil {
		return depl, err
	}

	return depl, nil
}

// +kubebuilder:rbac:groups=pulsar.pulsarconsumer.krvarma.com,resources=pulsarconsumers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pulsar.pulsarconsumer.krvarma.com,resources=pulsarconsumers/status,verbs=get;update;patch

func (r *PulsarConsumerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.TODO()
	log := r.Log.WithValues("pulsarconsumer", req.NamespacedName)

	var pulsarcrd pulsarv1.PulsarConsumer

	if err := r.Get(ctx, req.NamespacedName, &pulsarcrd); err != nil {
		log.Info("NO CRD Found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployment, err := r.getDeployment(pulsarcrd)
	if err != nil {
		return ctrl.Result{}, err
	}

	found := &appsv1.Deployment{}
	objKey := client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}
	err = r.Get(ctx, objKey, found)

	if err != nil {
		if !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		log.Info("Creating deployment")

		err = r.Create(ctx, &deployment)

		if err == nil {
			pulsarcrd.Status.Server = pulsarcrd.Spec.ServerAddress
			pulsarcrd.Status.Topic = pulsarcrd.Spec.Topic
			pulsarcrd.Status.Subscription = pulsarcrd.Spec.SubscriptionName
			pulsarcrd.Status.Replicas = pulsarcrd.Spec.Replicas

			err = r.Status().Update(ctx, &pulsarcrd)

			if err != nil {
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	} else {
		deplspec := deployment.Spec
		foundspec := found.Spec

		if !equality.Semantic.DeepDerivative(deplspec, foundspec) {
			found.Spec = deployment.Spec
			log.Info("Updating Deployment")

			err = r.Update(ctx, found)
			if err == nil {
				pulsarcrd.Status.Server = pulsarcrd.Spec.ServerAddress
				pulsarcrd.Status.Topic = pulsarcrd.Spec.Topic
				pulsarcrd.Status.Subscription = pulsarcrd.Spec.SubscriptionName
				pulsarcrd.Status.Replicas = pulsarcrd.Spec.Replicas

				err = r.Status().Update(ctx, &pulsarcrd)

				if err != nil {
					return reconcile.Result{}, err
				}
			}

			return reconcile.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func (r *PulsarConsumerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pulsarv1.PulsarConsumer{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
