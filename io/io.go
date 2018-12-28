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

// This package exports a simple Read and Write function pair that hides the
// local path vs. GCS implementations.
package io

import (
	"context"
	"io/ioutil"
	"regexp"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

var gcsRegexp = regexp.MustCompile("gs://([^/]+)/(.*)")

// Returns true if the path looks like a GCS object specification,
// e.g. gs://bucket/object/path
func isGCS(path string) bool {
	return gcsRegexp.MatchString(path)
}

// Given a GCS path, split into a bucket and object identifier
func gcsSplit(path string) ([]string, error) {
	splits := gcsRegexp.FindStringSubmatch(path)
	if len(splits) == 3 {
		return splits[1:], nil
	}

	return nil, errors.Errorf("%s is not a valid GCS path", path)
}

// Read the content of a GCS object to byte array
func gcsRead(ctx context.Context, path string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	parts, err := gcsSplit(path)
	if err != nil {
		return nil, err
	}
	reader, err := client.Bucket(parts[0]).Object(parts[1]).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Write a byte array to the GCS object specified
func gcsWrite(ctx context.Context, path string, data []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	parts, err := gcsSplit(path)
	if err != nil {
		return err
	}
	writer := client.Bucket(parts[0]).Object(parts[1]).NewWriter(ctx)
	if _, err = writer.Write(data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

// Read the contents of the file or GCS object identified by path to a byte
// array
func Read(ctx context.Context, path string) ([]byte, error) {
	if isGCS(path) {
		return gcsRead(ctx, path)
	}
	return ioutil.ReadFile(path)
}

// Write the contents of a byte array to either a local file or GCS object
func Write(ctx context.Context, path string, data []byte) error {
	if isGCS(path) {
		return gcsWrite(ctx, path, data)
	}
	return ioutil.WriteFile(path, data, 0644)
}
