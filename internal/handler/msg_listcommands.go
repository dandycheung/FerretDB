// Copyright 2021 FerretDB Inc.
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

package handler

import (
	"context"
	"maps"
	"slices"

	"github.com/FerretDB/wire/wirebson"

	"github.com/FerretDB/FerretDB/v2/internal/handler/middleware"
	"github.com/FerretDB/FerretDB/v2/internal/util/must"
)

// msgListCommands implements `listCommands` command.
//
// The passed context is canceled when the client connection is closed.
func (h *Handler) msgListCommands(connCtx context.Context, req *middleware.Request) (*middleware.Response, error) {
	doc := req.Document()

	if _, _, err := h.s.CreateOrUpdateByLSID(connCtx, doc); err != nil {
		return nil, err
	}

	cmdList := wirebson.MustDocument()

	for _, name := range slices.Sorted(maps.Keys(h.commands)) {
		help := h.commands[name].Help
		if help == "" {
			continue
		}

		must.NoError(cmdList.Add(name, wirebson.MustDocument(
			"help", help,
		)))
	}

	return middleware.ResponseDoc(req, wirebson.MustDocument(
		"commands", cmdList,
		"ok", float64(1),
	))
}
