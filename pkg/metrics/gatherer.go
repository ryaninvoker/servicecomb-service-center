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

package metrics

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/apache/servicecomb-service-center/pkg/gopool"
	"github.com/apache/servicecomb-service-center/pkg/log"
	"github.com/apache/servicecomb-service-center/pkg/prometheus"
)

// EmptyGather just active when metrics disabled
var EmptyGather = &Gather{
	Records: NewMetrics(),
	closed:  false,
}

func NewGatherer(opts Options) *Gather {
	return &Gather{
		Interval: opts.Interval,
		Records:  NewMetrics(),
		closed:   true,
	}
}

type Gather struct {
	Records  *Metrics
	Interval time.Duration

	lock   sync.Mutex
	closed bool
}

func (mm *Gather) Start() {
	mm.lock.Lock()
	if !mm.closed {
		mm.lock.Unlock()
		return
	}
	mm.closed = false

	gopool.Go(mm.loop)

	mm.lock.Unlock()
}

func (mm *Gather) loop(ctx context.Context) {
	ticker := time.NewTicker(mm.Interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := mm.Collect(); err != nil {
				log.Error("metrics collect failed", err)
				return
			}

			Report()
		}
	}
}

func (mm *Gather) Collect() error {
	mfs, err := prometheus.Gather()
	if err != nil {
		return err
	}

	records := NewMetrics()
	for _, mf := range mfs {
		name := mf.GetName()
		if _, ok := SysMetrics.Get(name); strings.Index(name, familyNamePrefix) == 0 || ok {
			if d := Calculate(mf); d != nil {
				records.put(strings.TrimPrefix(name, familyNamePrefix), d)
			}
		}
	}
	// clean the old cache here
	mm.Records = records
	return nil
}
