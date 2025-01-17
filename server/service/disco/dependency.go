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

package disco

import (
	"context"

	"github.com/apache/servicecomb-service-center/datasource"
	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/server/service/validator"
	pb "github.com/go-chassis/cari/discovery"
)

func (s *MicroServiceService) GetProviderDependencies(ctx context.Context,
	in *pb.GetDependenciesRequest) (*pb.GetProDependenciesResponse, error) {
	err := validator.Validate(in)
	if err != nil {
		log.Error("GetProviderDependencies failed for validating parameters failed", err)
		return &pb.GetProDependenciesResponse{
			Response: pb.CreateResponse(pb.ErrInvalidParams, err.Error()),
		}, nil
	}

	return datasource.GetDependencyManager().SearchProviderDependency(ctx, in)
}

func (s *MicroServiceService) GetConsumerDependencies(ctx context.Context, in *pb.GetDependenciesRequest) (*pb.GetConDependenciesResponse, error) {
	err := validator.Validate(in)
	if err != nil {
		log.Error("GetConsumerDependencies failed for validating parameters failed", err)
		return &pb.GetConDependenciesResponse{
			Response: pb.CreateResponse(pb.ErrInvalidParams, err.Error()),
		}, nil
	}

	return datasource.GetDependencyManager().SearchConsumerDependency(ctx, in)
}
