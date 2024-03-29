// Copyright 2019-2023 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracing

import (
	"fmt"
	"net/http"
)

func NameHTTPRequestSpan(operation string, req *http.Request) string {
	if operation == "" {
		return fmt.Sprintf("%v %v", req.Method, req.URL.String())
	}

	return fmt.Sprintf("%v: %v %v", operation, req.Method, req.URL.String())
}
