package tests

import (
	"testing"

	"github.com/ocfl-archive/dlza-manager-storage-handler/service"
)

func sptr(s string) *string { return &s }

func TestTrimKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "t1", in: "89280007-8674-4aa8-9b71-8943cfaf96ff-partition-1/tenantname-89280007-8674-4aa8-9b71-8943cfaf96ff-filename-test.zip", want: "89280007-8674-4aa8-9b71-8943cfaf96ff-partition-1/filename-test.zip"},
		{name: "t2", in: "archive/tenantname-89280007-8674-4aa8-9b71-8943cfaf96ff-filename.zip", want: "archive/filename.zip"},
		{name: "t3", in: "89280007-8674-4aa8-9b71-8943cfaf96ff-partition-1/file/tenantname-89280007-8674-4aa8-9b71-8943cfaf96ff-filename.zip", want: "89280007-8674-4aa8-9b71-8943cfaf96ff-partition-1/filename.zip"},
		{name: "t4", in: "archive/file/tenantname-89280007-8674-4aa8-9b71-8943cfaf96ff-filename-test.zip", want: "archive/filename-test.zip"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := service.TrimKey(sptr(tt.in))
			if got == nil {
				t.Fatalf("got nil, want %q", tt.want)
			}
			if *got != tt.want {
				t.Fatalf("got %q, want %q", *got, tt.want)
			}
		})
	}
}
