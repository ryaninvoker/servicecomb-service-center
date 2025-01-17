/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package datasource

import (
	"fmt"

	"github.com/apache/servicecomb-service-center/pkg/log"
)

type dataSourceEngine func(opts Options) (DataSource, error)

var (
	plugins        = make(map[Kind]dataSourceEngine)
	dataSourceInst DataSource
)

// load plugins configuration into plugins
func Install(pluginImplName string, engineFunc dataSourceEngine) {
	plugins[Kind(pluginImplName)] = engineFunc
}

// Init construct storage plugin instance
// invoked by sc main process
func Init(opts Options) error {
	if opts.Kind == "" {
		return nil
	}

	dataSourceEngine, ok := plugins[opts.Kind]
	if !ok {
		return fmt.Errorf("plugin implement not supported [%s]", opts.Kind)
	}
	var err error
	dataSourceInst, err = dataSourceEngine(opts)
	if err != nil {
		return err
	}
	log.Info(fmt.Sprintf("datasource plugin [%s] enabled", opts.Kind))
	return nil
}

func GetSCManager() SCManager {
	return dataSourceInst.SCManager()
}
func GetMetadataManager() MetadataManager {
	return dataSourceInst.MetadataManager()
}
func GetSystemManager() SystemManager {
	return dataSourceInst.SystemManager()
}
func GetRoleManager() RoleManager {
	return dataSourceInst.RoleManager()
}
func GetAccountManager() AccountManager {
	return dataSourceInst.AccountManager()
}
func GetAccountLockManager() AccountLockManager {
	return dataSourceInst.AccountLockManager()
}
func GetDependencyManager() DependencyManager {
	return dataSourceInst.DependencyManager()
}
func GetMetricsManager() MetricsManager {
	return dataSourceInst.MetricsManager()
}
