# Extending Kubernetes - Part 1 - Custom Operator

![Cover Image](https://raw.githubusercontent.com/krvarma/pulsar-consumer-operator/master/images/extending_kubernetes_operators.png?token=AA46XG3QPTTBSKWV63LBDF27E3FK4)

Kubernetes is an open-source container orchestration project and is one of the most successful projects in the Cloud Native era. Kubernetes started as an internal project at Google. Google open-sourced Kubernetes in 2014. Since then, it has emerged as the most popular application platform in the cloud and an integral part of Cloud Native development.

Extending Kubernetes is a series of articles that explore how to extend the Kubernetes system. Starting with Operators, we will explore many other ways to extend the functionalities of the Kubernetes system.

# Extending Kubernetes
Many features make Kubernetes great. Some of them are, the ability to automate various manual processes, deploy and scale your application seamlessly, roll out & roll back application updates, manage secrets, etc. 

Another great feature of Kubernetes that makes it great is its Extensibility. Almost everything in Kubernetes is extensible. The [official documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/) describes the following seven different extension points in the Kubernetes system.

1.  [Kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/):- Kubectl plugins are nothing but an executable file with name starts with _kubectl-_
2.  [API Access Extensions](https://kubernetes.io/docs/concepts/extend-kubernetes/#api-access-extensions):- API Access Extensions are extensions that extend different stages of API Server. We can use this extension to implement custom authentication, automatic sidecar injection, etc.
3.  [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/#user-defined-types):- Kubernetes has many resource types like Pods, Services, Deployments, etc. We can implement custom resource types. Custom Resources commonly combine with custom controllers.
4.  [Scheduler Extensions](https://kubernetes.io/docs/concepts/extend-kubernetes/#scheduler-extensions):- Kubernetes scheduler decides which node to use to deploy pods. We can extend this scheduler by writing a custom scheduler extension and implement our algorithms.
5.  [Custom Controllers](https://kubernetes.io/docs/concepts/architecture/controller/):- Custom controllers used along with custom resources, which is called Operator Pattern.
6.  [Network Plugins](https://kubernetes.io/docs/concepts/overview/extending#network-plugins):- Network Plugins are plugins that extend the pod networking.
7.  [Storage Plugins](https://kubernetes.io/docs/concepts/extend-kubernetes/#storage-plugins):- Storage Plugins are plugins that extend the types of storage.

# Custom Resource Definition

A resource in Kubernetes is an endpoint in Kubernetes API. Resource stores collection of API Objects belongs to a particular kind. There are many built-in resources in Kubernetes like Pods, Deployments, Services, etc. A Custom Resource is one that is not available in vanilla Kubernetes. Cluster admins can add or delete a custom resource anytime. We can manage custom resources using kubectl. An example of CRD is defining an organization-wide SSL configuration, and another example would be an application config CRD. You can use a ConfigMap instead of CRD in many cases. The [official documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#should-i-use-a-configmap-or-a-custom-resource) describes when you have to choose ConfigMap or CRD. For more information about extending Kubernetes using Custom Resource Definition, refer to [this link](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/).

# Custom Controller

In Kubernetes, a controller watches the state of the cluster and make necessary changes to meet the desired state. It is similar to a control loop in automation. A control loop is a non-terminating loop that keeps the desired state of the system. A controller in Kubernetes keeps track of at least one resource type. There many built-in controllers in Kubernetes, replication controller, namespace controller, service account controller, etc.

Custom Resource definition along with Custom Controller makes the [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

# Operator Pattern

The Operator is a pattern introduced by CoreOS in 2016. This pattern enables you to package, deploy, and manage your application without human intervention. Without Operator Pattern, humans perform these tasks. Operator Pattern is an implementation of the concept, infrastructure as software. Using this pattern, you can automate the deployment, management, etc. of your application.

In Kubernetes, the Operator is a software extension that makes use of Custom Resource to manage an application and its components. Operators are clients of Kubernetes API that controls the custom resource. An Operator is an application-specific controller that manages the state of a custom resource.

For example, you have an application that connects to a database and store/retrieve data and performs some business logic. To deploy this application, you have the deploy the database, the app, and it's different components. Usually, engineers perform these tasks. We can automate these tasks by writing an Operator.

# Reconciliation Loop

As explained in the previous paragraph, a custom controller manages the associated custom resource and is a client of Kubernetes API Server. When a new CR object created or modified, the API Server notifies our Operator. Then the Operator starts running a loop that watches the resource for any change in the actual and desired state. This loop is called the reconciliation loop.

Kubernetes watches the current state, and if there any change in the desired state, then try to reconcile the object's state.

# Writing an Operator

There are two frameworks to write operators, [operator-sdk](https://github.com/operator-framework/operator-sdk) and [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder). Both are Golang based tools and use [controller-tools](https://github.com/kubernetes-sigs/controller-tools) and [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) libraries internally. There is also a [work in progress to combine](https://github.com/kubernetes-sigs/kubebuilder/projects/7) both SDK to create a single one. Both SDKs are almost similar with some minimal differences. This [issue discusses](https://github.com/operator-framework/operator-sdk/issues/1758) the main difference between the tools.

Both these frameworks generate lots of boilerplate code for creating a CRD. The generated code is somewhat similar. There are a couple of significant differences; the operator-sdk supports [Helm](https://docs.openshift.com/container-platform/4.2/operators/operator_sdk/osdk-helm.html) and [Ansible](https://docs.openshift.com/container-platform/4.2/operators/operator_sdk/osdk-ansible.html) operators.

These tools generate lots of code. In most of the cases, what we need to do is to change the Spec to represent the desired state, change the status to represent the observed state of the resource, and modify the reconcile loop to include the logic maintain the desired state.

The Spec structure declares the desired state of our Kind, and the State structure declares the observed state of our Kind. All CRD include Spec and State.

Inside the reconcile loop, we implement our logic to maintain the state of the resource. In a typical scenario, we create the object, check whether it is in the desired state, modify the object to match the desired state if necessary.

# Pulsar Consumer Operator

In this article, we will try to create a basic application and create a CRD to deploy it on Kubernetes. We will reuse a Pulsar Consumer from my previous article [Creating an External Scaler for KEDA](https://medium.com/@krvarma/creating-an-external-scaler-for-keda-31b314b5c4a3). The consumer is a basic one that will consume a Pulsar topic and log the message to the console. In the real world, it will be much more complicated, but for the sake of simplicity, we will use a simple application.

For this purpose, we will create a custom resource **PulsarConsumer**. This CRD can be used to deploy a Pulsar Consumer application to the Kubernetes.

To generate this CRD, we will use the kubebuilder framework. But the steps described here also apply to operator-sdk.

Let's create a basic project. The following command generates a new project.

    kubebuilder init --domain pulsarconsumer.krvarma.com

The --domain specifies the domain of our project. Every API group we define will be under this domain.

This command will generate a lot of code. The folder structure looks like this.

![Folder Structure](https://raw.githubusercontent.com/krvarma/pulsar-consumer-operator/master/images/1.png?token=AA46XGZ7D7SCOXABC43RTX27EWJK2)

The generated project has the main.go file, Dockerfile, Makefile, config folder, etc. The entry point of the project is in main.go file. Let's take a look at the main.go file.

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
    
    package main
    
    import (
    	"flag"
    	"os"
    
    	"k8s.io/apimachinery/pkg/runtime"
    	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
    	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
    	ctrl "sigs.k8s.io/controller-runtime"
    	"sigs.k8s.io/controller-runtime/pkg/log/zap"
    	// +kubebuilder:scaffold:imports
    )
    
    var (
    	scheme   = runtime.NewScheme()
    	setupLog = ctrl.Log.WithName("setup")
    )
    
    func init() {
    	_ = clientgoscheme.AddToScheme(scheme)
    
    	// +kubebuilder:scaffold:scheme
    }
    
    func main() {
    	var metricsAddr string
    	var enableLeaderElection bool
    	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
    	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
    		"Enable leader election for controller manager. "+
    			"Enabling this will ensure there is only one active controller manager.")
    	flag.Parse()
    
    	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
    
    	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
    		Scheme:             scheme,
    		MetricsBindAddress: metricsAddr,
    		Port:               9443,
    		LeaderElection:     enableLeaderElection,
    		LeaderElectionID:   "ba84c7bc.pulsarconsumer.krvarma.com",
    	})
    	if err != nil {
    		setupLog.Error(err, "unable to start manager")
    		os.Exit(1)
    	}
    
    	// +kubebuilder:scaffold:builder
    
    	setupLog.Info("starting manager")
    	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
    		setupLog.Error(err, "problem running manager")
    		os.Exit(1)
    	}
    }

As you can see, the generated code use [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) to create and start a new manager. The manager is responsible for running our controller, webhooks, etc.

The config folder contains all the configuration YAML files to launch our CRD. The framework uses [kustomize](https://github.com/kubernetes-sigs/kustomize) YAML definitions, such as files for Certificate Manager, RABC, Prometheus, Webhooks, etc.

The generated code does not contain our Kind and the controller that manages it. The following command will create the controller.

    kubebuilder create api --group consumer --version v1 --kind PulsarConsumer

As you can see from the above command, to create the API, we need to specify API Group and Kind.

![Folder Structure](https://raw.githubusercontent.com/krvarma/pulsar-consumer-operator/master/images/2.png?token=AA46XG4XX6BTTEKEVYUCMCK7EWJ2U)

An API Group is a collection of related functionalities and has a version associated with it. Each version API Group has one or more API types called Kind. In our case, the Group is Consumer and Kind PulsarConsumer.

The above command will create some additional files and modify the main.go file. The new api folder contains all the API definitions and Group definitions. As explained earlier, the API definition includes the Spec, Status, and Schema and also a structure to hold a list of the schema.

The file `<kind>_types.go` contains the API definitions. The structure `<Kind>Spec` contains the spec fields. As you can see, the framework inserts a dummy field for reference, which we should replace with actual fields. The `<Kind>Status` structure contains the status fields.

The structure `<Kind>` contains the schema of the CRD; it contains TypeMeta, ObjectMeta, Spec, and Status fields. The `<Kind>List` is a structure to hold a list of specs.

We will add the following fields to the Spec.

 - ServerAddress - Address of the Pulsar server 
 - Topic - Name of the topic 
 - Subscription - Name of the subscription 
 - Replicas - Number of replicas

Also following are the fields of the structure Status.

 - ServerAddress - Address of the Pulsar server 
 - Topic - Name of the topic 
 - Subscription - Name of the subscription 
 - Replicas - Number of replicas

The next important file is the `<Kind>_controller.go` file. The controller-runtime framework uses the Reconciler interface to implement the reconciling of a specific Kind. Here is our generated controller file.

As you can see, the file defines our PulsarConsumer reconciler definition. The most crucial piece of code is the Reconcile function. The controller-runtime calls the Reconcile method whenever there is a change in the state of a single named object of our Kind. We will implement the logic to reconcile the object's state.

Basically what reconcile logic is:

1.  *Query the named object*: The Reconcile function receives a parameter of type Request. The request parameter contains the namespaced name of the specified object. We query the system to get the PulsarConsumer with the specified name using the `Get` method.
2.  *Retrieve the object*: Once we have the PulsarConsumer with the specified name. We will check whether we already have an object in the system.
3.  *Create if it is not present*: If the client returns a not found error, it means an object with the specified name is not there in the system. So we need to create it and update the status. If the error is something else, we should gracefully return.
4.  *Check current state*: If there is no error, then it means the specified object is present in the system. We should check the current state and the desired state and see whether it is equal or not. This part is a little bit tricky. Many blogs suggest using `reflect.DeepEqual`, but this will not work since the deployment controller or other Kubernetes components will add some default fields to the Spec object, which will result in always false situations while using `reflect.DeepEqual`. After going through different blogs and kubebuilder issues, I found [this particular issue](https://github.com/kubernetes-sigs/kubebuilder/issues/592). [One of the comments](https://github.com/kubernetes-sigs/kubebuilder/issues/592#issuecomment-625738183) suggests using `equality.Semantic.DeepDerivative` since it will compare only non-zero fields on the structure. Using the proposed solution worked without any issues.
5.  Update the state if necessary: If there is a difference in the present and desired state, then we should update the state and update the status. If there is no difference, do nothing.

Here is the complete source code of the controller.

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

# Idempotency of the Operator

One of the critical requirements of an operator is that it should be idempotent. The Kubernetes system relies on controllers to reconcile the state of the resource. The Operator determines how to maintain the state.

The Operator should always expect multiple calls to reconcile unchanged resources. The Operator should make sure to produce the same result consistently.

# Marker Comments

While going through the code, you might have noticed the comments start with + symbol. These are particular comments called Marker comments. Marker comments always start with + followed by marker name and optional parameters.

The kubebuilder tool uses `controller-gen` for generating utility codes and YAML files. The kubebuilder generates a Makefile to build and run the Operator. The `controller-gen` framework sees the marker comments; the framework parses the marker comments and generates code based on the marker.

You can see the marker comments in the Spec definition. The following marker comment ensures that the particular field is mandatory.

    // +kubebuilder:validation:Required

Likewise, the marker comment above the Reconcile function ensures the proper RBAC roles.

    // +kubebuilder:rbac:groups=pulsar.pulsarconsumer.krvarma.com,resources=pulsarconsumers,verbs=get;list;watch;create;update;patch;delete
    // +kubebuilder:rbac:groups=pulsar.pulsarconsumer.krvarma.com,resources=pulsarconsumers/status,verbs=get;update;patch

For detailed documentation of marker comments, refer to [this link](https://book.kubebuilder.io/reference/markers.html).

# Building and deploying the Operator

The kubebuilder uses make utility to build and deploy the Operator. There are many make targets.

1.  `make run`:- Run the on the default Kubernetes cluster
2.  `make install`:- Install the CRD into the cluster
3.  `make uninstall`:- Uninstall the CRD
4.  `make deploy`:- Deploy the Operator into the cluster
5.  `make manifests`:- Generate the YAML files
6.  `make generate`:- Generate source codes
7.  `make docker-build`:- Build the docker image
8.  `make docker-push`:- Push the docker image to the specified registry


To install the CRD and run the Operator, run the following commands.

    make install
    make run

These commands will install the CRD to the cluster and run the Operator locally. Now we can create a resource of kind PulsarConsumer. When we add an API, the kubebuilder will generate a sample YAML file of our Kind under the folder config/samples. We can create a resource of our Kind using this file.

    kubectl apply -f config/samples/pulsar_v1_pulsarconsumer.yaml

If everything goes well, you can see the pod running by issuing the `kubectl get pods` command.

![pods](https://raw.githubusercontent.com/krvarma/pulsar-consumer-operator/master/images/3.png?token=AA46XGYJHPVR4UHX5AWUBFC7EW3E6)

# Deploying the Operator

To deploy the Operator to the cluster, we need to build and push the Operator image to the Docker registry.

    make docker-build docker-push IMG=<registry>/<user>/pulsar-operator

You should supply the values for the registry and user. This command will build and push the image to the specified image registry.

To deploy and runt the Operator into the local cluster, run the following command:

    make deploy IMG=<registry>/<user>/pulsar-operator

Again if everything goes well, you can see the controller pod running. Please note that the system will create the pod in the namespace specified.

    kubectl get pods -n pulsarconsumercrd-system

![pods](https://raw.githubusercontent.com/krvarma/pulsar-consumer-operator/master/images/4.png?token=AA46XG5QTPVVEIKUTRBZHE27EW3JM)

To test whether the consumer is running or not, you can send a message to the topic specified. I have included a sample pulsar producer written in Golang.

To send a message, run the following command.

    env PULSAR_SERVER=pulsar://localhost:6650 env PULSAR_MESSAGE="Sample Message" env PULSAR_TOPIC="my-topic" pulsar-producer

You can see the message logged to the console by running kubectl get logs command.

I hope this article is helpful to kickstart writing Kubernetes Operator. In the upcoming part of the article, we will explore other ways to extend the Kubernetes system. Till then, Happy Coding!
