/*
 * This file was last modified at 2024-06-16 14:35 by Victor N. Skurikhin.
 * reports_test.go
 * $Id$
 */

package agent

import (
	t "github.com/vskurikhin/gometrics/internal/types"
)

var enabled = []t.Name{t.TotalAlloc, t.PollCount, t.RandomValue}

//TODO
//func TestReports(t *testing.T) {
//	s := "1"
//	Storage()
//	store.Put("PollCount", &s)
//	store.Put("RandomValue", &s)
//	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
//		res.Write([]byte(""))
//	}))
//
//	a := strings.Split(testServer.URL, "://")
//	if len(a) < 2 {
//		t.Fatalf("len(%s) < 2", a)
//	}
//	t.Setenv("ADDRESS", a[1])
//	cfg := env.GetAgentConfig()
//	client := http.Client{}
//	reports(cfg, enabled, &client)
//	testServer.Close()
//	reports(cfg, enabled, &client)
//}
