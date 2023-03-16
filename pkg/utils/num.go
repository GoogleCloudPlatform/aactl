// Copyright 2023 Google LLC
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

package utils

func ToFloat32(v interface{}) float32 {
	if v == nil {
		return 0
	}

	switch v := v.(type) {
	case float32:
		return v
	case float64: // TODO: handle overflow
		return float32(v)
	case int:
		return float32(v)
	case int32:
		return float32(v)
	case int64:
		return float32(v)
	case uint:
		return float32(v)
	case uint32:
		return float32(v)
	case uint64:
		return float32(v)
	}
	return 0
}
