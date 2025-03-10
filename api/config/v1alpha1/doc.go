// Package types holds the definition of the gikopsctl configuration struct and
// supporting structs.  It's the k8s API conformant object that describes
// a set of generation and transformation operations to create and/or
// modify k8s resources.
// A gikopsctl configuration file is a serialization of this struct.
// +k8s:deepcopy-gen=package
package v1alpha1
