package hst

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadFilesReturnsResponseOnHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":"FAILED","msg":"rejected"}`))
	}))
	defer server.Close()

	client, err := NewHst(&Option{BaseURL: server.URL, ChannelId: "channel-1"})
	if err != nil {
		t.Fatalf("NewHst() error = %v", err)
	}

	bizData, uploadResp, err := client.UploadFiles(context.Background(), NewUploadFilesDto("token-1"))
	if err == nil {
		t.Fatal("UploadFiles() error = nil, want HTTP error")
	}
	if bizData != nil {
		t.Fatalf("UploadFiles() bizData = %#v, want nil", bizData)
	}
	if uploadResp == nil || uploadResp.Code != "FAILED" || uploadResp.Msg != "rejected" {
		t.Fatalf("UploadFiles() response = %#v", uploadResp)
	}
}
