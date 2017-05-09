/* main.go - Example for goemacs API

Copyright (C) 2016 Yann Hodique <yann.hodique@gmail.com>.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or (at
your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.  */

package main

// int plugin_is_GPL_compatible;
import "C"

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/sigma/go-emacs"
)

func init() {
	emacs.Register(initModule)
}

func initModule(env emacs.Environment) {
	env.RegisterFunction("k8s-make-client", MakeClient, 1, "k8s-make-client", nil)
	env.RegisterFunction("k8s-list-pods", ListPods, 1, "k8s-list-pods", nil)
	env.ProvideFeature("k8s")
}

func MakeClient(ctx emacs.FunctionCallContext) (emacs.Value, error) {
	env := ctx.Environment()
	kubeconfig, err := ctx.GoStringArg(0)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return env.MakeUserPointer(clientset), nil
}

func ListPods(ctx emacs.FunctionCallContext) (emacs.Value, error) {
	env := ctx.Environment()
	rawClient, ok := env.ResolveUserPointer(ctx.UserPointerArg(0))
	if !ok {
		return emacs.Error("user-ptr does not exist")
	}

	client, ok := rawClient.(*kubernetes.Clientset)
	if !ok {
		return emacs.Error("user-ptr is not a k8s client")
	}

	pods, err := client.Core().Pods("").List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podNames := make([]emacs.Value, len(pods.Items))

	for i := 0; i < len(pods.Items); i++ {
		podNames[i] = env.String(pods.Items[i].Name)
	}
	return env.StdLib().List(podNames...), nil
}

func main() {}
