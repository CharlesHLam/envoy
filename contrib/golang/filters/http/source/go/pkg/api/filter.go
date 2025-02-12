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

package api

import "google.golang.org/protobuf/types/known/anypb"

// request
type StreamDecoderFilter interface {
	DecodeHeaders(RequestHeaderMap, bool) StatusType
	DecodeData(BufferInstance, bool) StatusType
	DecodeTrailers(RequestTrailerMap) StatusType
}

type StreamFilterConfigParser interface {
	Parse(any *anypb.Any) (interface{}, error)
	Merge(parentConfig interface{}, childConfig interface{}) interface{}
}

type StreamFilterConfigFactory func(config interface{}) StreamFilterFactory
type StreamFilterFactory func(callbacks FilterCallbackHandler) StreamFilter

type StreamFilter interface {
	// http request
	StreamDecoderFilter
	// response stream
	StreamEncoderFilter
	// destroy filter
	OnDestroy(DestroyReason)
	// TODO add more for stream complete and log phase
}

// response
type StreamEncoderFilter interface {
	EncodeHeaders(ResponseHeaderMap, bool) StatusType
	EncodeData(BufferInstance, bool) StatusType
	EncodeTrailers(ResponseTrailerMap) StatusType
}

// stream info
// refer https://github.com/envoyproxy/envoy/blob/main/envoy/stream_info/stream_info.h
type StreamInfo interface {
	GetRouteName() string
	FilterChainName() string
	// Protocol return the request's protocol.
	Protocol() (string, bool)
	// ResponseCode return the response code.
	ResponseCode() (uint32, bool)
	// ResponseCodeDetails return the response code details.
	ResponseCodeDetails() (string, bool)
	// AttemptCount return the number of times the request was attempted upstream.
	AttemptCount() uint32
	// Get the dynamic metadata of the request
	DynamicMetadata() DynamicMetadata
	// DownstreamLocalAddress return the downstream local address.
	DownstreamLocalAddress() string
	// DownstreamRemoteAddress return the downstream remote address.
	DownstreamRemoteAddress() string
	// UpstreamHostAddress return the upstream host address.
	UpstreamHostAddress() (string, bool)
	// UpstreamClusterName return the upstream host cluster.
	UpstreamClusterName() (string, bool)
	// FilterState return the filter state interface.
	FilterState() FilterState
}

type StreamFilterCallbacks interface {
	StreamInfo() StreamInfo
}

type FilterCallbacks interface {
	StreamFilterCallbacks
	// Continue or SendLocalReply should be last API invoked, no more code after them.
	Continue(StatusType)
	SendLocalReply(responseCode int, bodyText string, headers map[string]string, grpcStatus int64, details string)
	// RecoverPanic recover panic in defer and terminate the request by SendLocalReply with 500 status code.
	RecoverPanic()
	Log(level LogType, msg string)
	LogLevel() LogType
	// TODO add more for filter callbacks
}

type FilterCallbackHandler interface {
	FilterCallbacks
}

type DynamicMetadata interface {
	// TODO: Get(filterName string) map[string]interface{}
	Set(filterName string, key string, value interface{})
}

type StateType int

const (
	StateTypeReadOnly StateType = 0
	StateTypeMutable  StateType = 1
)

type LifeSpan int

const (
	LifeSpanFilterChain LifeSpan = 0
	LifeSpanRequest     LifeSpan = 1
	LifeSpanConnection  LifeSpan = 2
	LifeSpanTopSpan     LifeSpan = 3
)

type StreamSharing int

const (
	None                             StreamSharing = 0
	SharedWithUpstreamConnection     StreamSharing = 1
	SharedWithUpstreamConnectionOnce StreamSharing = 2
)

type FilterState interface {
	SetString(key, value string, stateType StateType, lifeSpan LifeSpan, streamSharing StreamSharing)
	GetString(key string) string
}
