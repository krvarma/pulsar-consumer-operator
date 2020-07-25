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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PulsarConsumerSpec defines the desired state of PulsarConsumer
type PulsarConsumerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// Address of the pulsar server.
	ServerAddress string `json:"serverAddress,omitempty"`

	// +kubebuilder:validation:Required
	// Name of the topic to listen.
	Topic string `json:"topic,omitempty"`

	// +kubebuilder:validation:Required
	// Name of the subscripton.
	SubscriptionName string `json:"subscriptionName,omitempty"`

	// +kubebuilder:validation:Required
	// Number of replicas.
	Replicas *int32 `json:"replicas,omitempty"`
}

// PulsarConsumerStatus defines the observed state of PulsarConsumer
type PulsarConsumerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Server Address
	Server string `json:"server"`
	// Name of the pulsar topic
	Topic string `json:"topic"`
	// Name of the subscription
	Subscription string `json:"subscription"`
	// Number of replicas
	Replicas *int32 `json:"replicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.server",name="Server",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.topic",name="Topic",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.subscription",name="Subscription",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name="Replicas",type="integer"

// PulsarConsumer is the Schema for the pulsarconsumers API
type PulsarConsumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PulsarConsumerSpec   `json:"spec,omitempty"`
	Status PulsarConsumerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PulsarConsumerList contains a list of PulsarConsumer
type PulsarConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PulsarConsumer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PulsarConsumer{}, &PulsarConsumerList{})
}
