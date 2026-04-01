package store

import "testing"

func Test_storeFromID(t *testing.T) {
	tests := []struct {
		name          string
		in            string
		wantTenant    string
		wantPartition string
	}{
		{
			name:          "t1",
			in:            "tenant-B-89280007-8674-4aa8-9b71-8943cfaf96ff-partition-111",
			wantTenant:    "tenant-B",
			wantPartition: "89280007-8674-4aa8-9b71-8943cfaf96ff",
		},
		{
			name:          "t2",
			in:            "tenant-A-89280007-8674-4aa8-9b71-8943cfaf96ff-partition-1-whatever-test",
			wantTenant:    "tenant-A",
			wantPartition: "89280007-8674-4aa8-9b71-8943cfaf96ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant, partition := storeFromID(tt.in)
			if tenant != tt.wantTenant {
				t.Fatalf("tenant: got %q want %q", tenant, tt.wantTenant)
			}
			if partition != tt.wantPartition {
				t.Fatalf("partition: got %q want %q", partition, tt.wantPartition)
			}
		})
	}
}
