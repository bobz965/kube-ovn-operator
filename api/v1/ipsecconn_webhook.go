/*
Copyright 2023.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var ipsecconnlog = logf.Log.WithName("ipsecconn-resource")

func (r *IpsecConn) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-vpn-gw-kube-ovn-operator-com-v1-ipsecconn,mutating=true,failurePolicy=fail,sideEffects=None,groups=vpn-gw.kube-ovn-operator.com,resources=ipsecconns,verbs=create;update,versions=v1,name=mipsecconn.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &IpsecConn{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *IpsecConn) Default() {
	ipsecconnlog.Info("default", "name", r.Name)
	if r.Spec.Auth == "" {
		r.Spec.Auth = "pubkey"
	}
	if r.Spec.IkeVersion == "" {
		r.Spec.IkeVersion = "2"
	}
	if r.Spec.Proposals == "" {
		r.Spec.Proposals = "default"
	}
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-vpn-gw-kube-ovn-operator-com-v1-ipsecconn,mutating=false,failurePolicy=fail,sideEffects=None,groups=vpn-gw.kube-ovn-operator.com,resources=ipsecconns,verbs=create;update,versions=v1,name=vipsecconn.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &IpsecConn{}

func (r *IpsecConn) validateIpsecConn() error {
	var allErrs field.ErrorList

	if r.Spec.VpnGw == "" {
		err := fmt.Errorf("ipsecConn vpn gw is required")
		e := field.Invalid(field.NewPath("spec").Child("vpnGw"), r.Spec.VpnGw, err.Error())
		allErrs = append(allErrs, e)
	}

	if r.Spec.IkeVersion != "0" && r.Spec.IkeVersion != "1" && r.Spec.IkeVersion != "2" {
		err := fmt.Errorf("ipsec connection spec ike version is invalid, ike version spec: %s", r.Spec.IkeVersion)
		e := field.Invalid(field.NewPath("spec").Child("ikeVersion"), r.Spec.IkeVersion, err.Error())
		allErrs = append(allErrs, e)
	}

	if r.Spec.Auth != "psk" && r.Spec.Auth != "pubkey" {
		err := fmt.Errorf("ipsec connection spec auth is invalid, auth spec: %s", r.Spec.Auth)
		e := field.Invalid(field.NewPath("spec").Child("auth"), r.Spec.Auth, err.Error())
		allErrs = append(allErrs, e)
	}

	if r.Spec.RemotePublicIp == "" {
		err := fmt.Errorf("ipsecConn remote public ip is required")
		e := field.Invalid(field.NewPath("spec").Child("localPublicIp"), r.Spec.RemotePublicIp, err.Error())
		allErrs = append(allErrs, e)
	}

	if r.Spec.RemotePrivateCidrs == "" {
		err := fmt.Errorf("ipsecConn remote private cidrs is required")
		e := field.Invalid(field.NewPath("spec").Child("remotePrivateCidrs"), r.Spec.RemotePrivateCidrs, err.Error())
		allErrs = append(allErrs, e)
	}

	if r.Spec.LocalPublicIp == "" {
		err := fmt.Errorf("ipsecConn localPublicIp is required")
		e := field.Invalid(field.NewPath("spec").Child("localPublicIp"), r.Spec.LocalPublicIp, err.Error())
		allErrs = append(allErrs, e)
	}

	if r.Spec.LocalPrivateCidrs == "" {
		err := fmt.Errorf("ipsecConn local private cidrs is required")
		e := field.Invalid(field.NewPath("spec").Child("localPrivateCidrs"), r.Spec.LocalPrivateCidrs, err.Error())
		allErrs = append(allErrs, e)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs.ToAggregate()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *IpsecConn) ValidateCreate() error {
	ipsecconnlog.Info("validate create", "name", r.Name)
	if err := r.validateIpsecConn(); err != nil {
		return err
	}
	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *IpsecConn) ValidateUpdate(old runtime.Object) error {
	ipsecconnlog.Info("validate update", "name", r.Name)
	if err := r.validateIpsecConn(); err != nil {
		return err
	}
	oldIpsecConn, _ := old.(*IpsecConn)
	var allErrs field.ErrorList
	if oldIpsecConn.Spec.VpnGw != "" && oldIpsecConn.Spec.VpnGw != r.Spec.VpnGw {
		err := fmt.Errorf("ipsecConn vpn gw can not be changed")
		e := field.Invalid(field.NewPath("spec").Child("vpnGw"), r.Spec.VpnGw, err.Error())
		allErrs = append(allErrs, e)
	}
	if len(allErrs) == 0 {
		return nil
	}

	// TODO(user): fill in your validation logic upon object update.
	return allErrs.ToAggregate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *IpsecConn) ValidateDelete() error {
	ipsecconnlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
