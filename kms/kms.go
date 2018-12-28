// Copyright 2018 Neudesic LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package kms contains utility functions to invoke KMS APIs for encryption and
// decryption
package kms

import (
	"context"

	api "cloud.google.com/go/kms/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// Decrypt a byte array using the KMS key provided.
// Note that key format is of the form projects/{projectId}/locations/{location}/keyRings/{keyRingName}/cryptoKeys/{cryptoKeyName}
func Decrypt(ctx context.Context, key string, ciphertext []byte) ([]byte, error) {
	client, err := api.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &pb.DecryptRequest{
		Name:       key,
		Ciphertext: ciphertext,
	}

	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}

// Encrypt a byte array using the KMS key provided.
// Note that key format is of the form projects/{projectId}/locations/{location}/keyRings/{keyRingName}/cryptoKeys/{cryptoKeyName}
func Encrypt(ctx context.Context, key string, plaintext []byte) ([]byte, error) {
	client, err := api.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &pb.EncryptRequest{
		Name:      key,
		Plaintext: plaintext,
	}

	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}
