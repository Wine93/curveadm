/*
 *  Copyright (c) 2021 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: CurveAdm
 * Created Date: 2021-10-15
 * Author: Jingli Chen (Wine93)
 */

package common

import (
	"fmt"

	"github.com/opencurve/curveadm/cli/cli"
	"github.com/opencurve/curveadm/internal/configure"
	"github.com/opencurve/curveadm/internal/task/step"
	"github.com/opencurve/curveadm/internal/task/task"
	tui "github.com/opencurve/curveadm/internal/tui/common"
)

func NewRestartServiceTask(curveadm *cli.CurveAdm, dc *configure.DeployConfig) (*task.Task, error) {
	serviceId := configure.ServiceId(curveadm.ClusterId(), dc.GetId())
	containerId, err := curveadm.Storage().GetContainerId(serviceId)
	if err != nil {
		return nil, err
	} else if containerId == "" {
		return nil, fmt.Errorf("service(id=%s) not found", serviceId)
	}

	subname := fmt.Sprintf("host=%s role=%s containerId=%s",
		dc.GetHost(), dc.GetRole(), tui.TrimContainerId(containerId))
	t := task.NewTask("Restart Service", subname, dc.GetSshConfig())

	// add step
	t.AddStep(&step.RestartContainer{
		ContainerId:  containerId,
		ExecWithSudo: true,
		ExecInLocal:  false,
	})
	t.AddStep(&step2PostStart{
		ContainerId:  containerId,
		ExecWithSudo: true,
		ExecInLocal:  false,
	})

	return t, nil
}