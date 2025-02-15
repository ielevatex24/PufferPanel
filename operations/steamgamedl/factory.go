/*
 Copyright 2022 PufferPanel
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

package steamgamedl

import (
	"github.com/pufferpanel/pufferpanel/v2"
	"github.com/spf13/cast"
)

type OperationFactory struct {
	pufferpanel.OperationFactory
}

func (of OperationFactory) Key() string {
	return "steamgamedl"
}
func (of OperationFactory) Create(op pufferpanel.CreateOperation) (pufferpanel.Operation, error) {
	o := SteamGameDl{
		AppId:     cast.ToString(op.OperationArgs["appId"]),
		Username:  cast.ToString(op.OperationArgs["username"]),
		Password:  cast.ToString(op.OperationArgs["password"]),
		ExtraArgs: cast.ToStringSlice(op.OperationArgs["extraArgs"]),
	}
	return o, nil
}

var Factory OperationFactory
