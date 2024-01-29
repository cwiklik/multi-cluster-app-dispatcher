/*
Copyright 2019, 2021 The Multi-Cluster App Dispatcher Authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SchedulingSpecPlural is the plural of SchedulingSpec
const SchedulingSpecPlural = "schedulingspecs"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SchedulingSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec SchedulingSpecTemplate `json:"spec,omitempty" protobuf:"bytes,1,rep,name=spec"`
}

type SchedulingSpecTemplate struct {
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,1,rep,name=nodeSelector"`
	// Expected number of pods in running and/or completed state.
	// Requeuing is triggered when the number of running/completed pods is not equal to this value.
	// When not specified, requeuing is disabled and no check is performed.
	MinAvailable int `json:"minAvailable,omitempty" protobuf:"bytes,2,rep,name=minAvailable"`
	// Specification of the requeuing strategy based on waiting time.
	// Values in this field control how often the pod check should happen,
	// and if requeuing has reached its maximum number of times.
	Requeuing RequeuingTemplate `json:"requeuing,omitempty" protobuf:"bytes,1,rep,name=requeuing"`
	ClusterScheduling ClusterSchedulingSpec `json:"clusterScheduling,omitempty"`
	// Wall clock duration time of appwrapper in seconds.
	DispatchDuration DispatchDurationSpec `json:"dispatchDuration,omitempty"`
}

type RequeuingTemplate struct {
	// Value to keep track of the initial wait time.
	// Users cannot set this as it is taken from 'timeInSeconds'.
	InitialTimeInSeconds int `json:"initialTimeInSeconds,omitempty" protobuf:"bytes,1,rep,name=initialTimeInSeconds"`
	// Initial waiting time before requeuing conditions are checked. This value is
	// specified by the user, but it may grow as requeuing events happen.
	// +kubebuilder:default=300
	TimeInSeconds int `json:"timeInSeconds,omitempty" protobuf:"bytes,2,rep,name=timeInSeconds"`
	// Maximum waiting time for requeuing checks.
	// +kubebuilder:default=0
	MaxTimeInSeconds int `json:"maxTimeInSeconds,omitempty" protobuf:"bytes,3,rep,name=maxTimeInSeconds"`
	// Growth strategy to increase the waiting time between requeuing checks.
	// The values available are 'exponential', 'linear', or 'none'.
	// For example, 'exponential' growth would double the 'timeInSeconds' value
	// every time a requeuing event is triggered.
	// If the string value is misspelled or not one of the possible options,
	// the growth behavior is defaulted to 'none'.
	// +kubebuilder:default=exponential
	GrowthType string `json:"growthType,omitempty" protobuf:"bytes,4,rep,name=growthType"`
	// Field to keep track of how many times a requeuing event has been triggered.
	// +kubebuilder:default=0
	NumRequeuings int `json:"numRequeuings,omitempty" protobuf:"bytes,5,rep,name=numRequeuings"`
	// Maximum number of requeuing events allowed. Once this value is reached (e.g.,
	// 'numRequeuings = maxNumRequeuings', no more requeuing checks are performed and the generic
	// items are stopped and removed from the cluster (AppWrapper remains deployed).
	// +kubebuilder:default=0
	MaxNumRequeuings int `json:"maxNumRequeuings,omitempty" protobuf:"bytes,6,rep,name=maxNumRequeuings"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SchedulingSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SchedulingSpec `json:"items"`
}

type ResourceName string

type ClusterReference struct {
	Name string `json:"name"`
}

type ClusterSchedulingSpec struct {
	Clusters        []ClusterReference    `json:"clusters,omitempty"`
	ClusterSelector *metav1.LabelSelector `json:"clusterSelector,omitempty"`
	PolicyResult    ClusterDecision       `json:"policyResult,omitempty"`
}

type ClusterDecision struct {
	TargetCluster ClusterReference        `json:"targetCluster,omitempty"`
	PolicySource  []PolicySourceReference `json:"policySource,omitempty"`
}

type PolicySourceReference struct {
	// ID/Name of the policy decision maker.  Most often this will be MCAD but design can support alternatives
	Name                string           `json:"name,omitempty"`
	// The latest time this condition was updated.
	LastUpdateMicroTime metav1.MicroTime `json:"lastUpdateMicroTime,omitempty"`
	// A human readable message indicating details about the cluster decision.
	Message             string           `json:"message,omitempty"`
}

type ScheduleTimeSpec struct {
	Min     metav1.Time `json:"minTimestamp,omitempty"`
	Desired metav1.Time `json:"desiredTimestamp,omitempty"`
	Max     metav1.Time `json:"maxTimestamp,omitempty"`
}

type DispatchDurationSpec struct {
	Expected int  `json:"expected,omitempty"`
	Limit    int  `json:"limit,omitempty"`
	Overrun  bool `json:"overrun,omitempty"`
}

type DispatchingWindowSpec struct {
	Start ScheduleTimeSpec `json:"start,omitempty"`
	End   ScheduleTimeSpec `json:"end,omitempty"`
}
