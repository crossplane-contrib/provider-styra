package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	xpcontroller "github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
)

// Put interface here to avoid cyclic dependency in pkg/controller

// Options to initialize a new controller.
type Options struct {
	xpcontroller.Options

	ConnectionPublisher []managed.ConnectionPublisher
}

// Initializer sets up and starts a new controller.
type Initializer interface {
	Setup(ctrl.Manager, Options)
}
