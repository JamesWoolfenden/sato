package arm

import "testing"

func Test_replaceReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "strips apiVersion arg",
			in:   "reference(resourceId('Microsoft.Network/publicIPAddresses', var.publicIpName), '2022-05-01')",
			want: "resourceId('Microsoft.Network/publicIPAddresses', var.publicIpName)",
		},
		{
			name: "strips apiVersion and Full",
			in:   "reference(var.x, '2022-05-01', 'Full').properties.y",
			want: "var.x.properties.y",
		},
		{
			name: "single arg unchanged",
			in:   "reference(azurerm_storage_account.sato0).primaryEndpoints.blob",
			want: "azurerm_storage_account.sato0.primaryEndpoints.blob",
		},
		{
			name: "nested parens in first arg",
			in:   "reference(resourceId(a, b), '2020-01-01')",
			want: "resourceId(a, b)",
		},
		{
			name: "no reference call",
			in:   "var.something",
			want: "var.something",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := replaceReference(tt.in); got != tt.want {
				t.Errorf("replaceReference(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func Test_replaceUUID_noop(t *testing.T) {
	t.Parallel()

	in := "substr(random_uuid.sato0.result, 0, 8)"
	result := map[string]interface{}{"data": map[string]interface{}{"uuid": 0}}

	got, gotResult := replaceUUID(in, result)

	if got != in {
		t.Errorf("replaceUUID() should not alter input without uuid() call: got %q", got)
	}
	if gotResult["data"].(map[string]interface{})["uuid"] != 0 {
		t.Errorf("replaceUUID() should not increment counter when no uuid() present")
	}
}
