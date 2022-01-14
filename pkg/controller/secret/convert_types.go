/*
Copyright 2021 The Crossplane Authors.

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

package secret

import (
	"crypto/sha1" // nolint:gosec
	"encoding/json"
	"fmt"
	"time"

	"github.com/mistermx/styra-go-client/pkg/client/secrets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane-contrib/provider-styra/apis/secret/v1alpha1"
	styraclient "github.com/crossplane-contrib/provider-styra/pkg/client"
)

const (
	checksumSecretName = "%s-styra-checksum"
)

func generateSecret(resp *secrets.GetSecretOK) *v1alpha1.Secret {
	lastModifiedAt := metav1.NewTime(time.Time(resp.Payload.Result.Metadata.LastModifiedAt))

	cr := &v1alpha1.Secret{
		Spec: v1alpha1.SecretSpec{
			ForProvider: v1alpha1.SecretParameters{
				Description: styraclient.StringValue(resp.Payload.Result.Description),
			},
		},
		Status: v1alpha1.SecretStatus{
			AtProvider: v1alpha1.SecretObservation{
				LastModifiedAt: &lastModifiedAt,
			},
		},
	}

	return cr
}

func generateChecksumSecretName(base string) string {
	return fmt.Sprintf(checksumSecretName, base)
}

func generateSecretChecksum(secretVal string, t time.Time) (string, error) {
	m := map[string]string{
		"secretVal": secretVal,
		"t":         t.UTC().String(),
	}

	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	sum := sha1.Sum(data) // nolint:gosec
	return string(sum[:]), nil
}

// isNotFound returns whether the given error is of type NotFound or not.
func isNotFound(err error) bool {
	_, ok := err.(*secrets.GetSecretNotFound)
	return ok
}
