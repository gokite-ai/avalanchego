// Copyright (C) 2026, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/snowtest"
	"github.com/ava-labs/avalanchego/utils/logging"
)

// The root RPC and websocket routes are a production Kite interface. The alias
// must retain upstream bootstrap rejection and apply only to the configured L1.
func TestKiteAliases(t *testing.T) {
	previous := KiteChainID
	t.Cleanup(func() { KiteChainID = previous })
	for _, tc := range []struct {
		name        string
		configured  string
		state       snow.State
		aliasStatus int
	}{
		{"configured", snowtest.CChainID.String(), snow.NormalOp, http.StatusNoContent},
		{"bootstrapping", snowtest.CChainID.String(), snow.Bootstrapping, http.StatusServiceUnavailable},
		{"disabled", "", snow.NormalOp, http.StatusNotFound},
		{"other_chain", "different-chain", snow.NormalOp, http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			KiteChainID = tc.configured
			m, err := newMetrics(prometheus.NewRegistry())
			require.NoError(t, err)
			s := &server{log: logging.NoLog{}, router: newRouter(), metrics: m}
			ctx := snowtest.ConsensusContext(snowtest.Context(t, snowtest.CChainID))
			ctx.State.Set(snow.EngineState{State: tc.state})
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
			base := "bc/" + ctx.ChainID.String()
			require.NoError(t, s.addChainRoute("kite", handler, ctx, base, "/rpc"))
			require.NoError(t, s.addChainRoute("kite", handler, ctx, base, "/ws"))
			for _, route := range []string{"/", "/ws", "/ext/bc/kite/rpc", "/ext/bc/kite/ws"} {
				w := httptest.NewRecorder()
				s.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, route, nil))
				require.Equal(t, tc.aliasStatus, w.Code, route)
			}
			canonicalStatus := http.StatusNoContent
			if tc.state != snow.NormalOp {
				canonicalStatus = http.StatusServiceUnavailable
			}
			for _, endpoint := range []string{"/rpc", "/ws"} {
				w := httptest.NewRecorder()
				s.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ext/"+base+endpoint, nil))
				require.Equal(t, canonicalStatus, w.Code, endpoint)
			}
		})
	}
}