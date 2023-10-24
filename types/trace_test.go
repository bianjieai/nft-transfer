package types

import (
	"reflect"
	"testing"
)

func TestIsAwayFromOrigin(t *testing.T) {
	type args struct {
		sourcePort    string
		sourceChannel string
		fullClassPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"transfer forward by origin chain", args{"port-1", "channel-1", "kitty"}, true},
		{"transfer forward by relay chain", args{"port-3", "channel-3", "port-2/channel-2/kitty"}, true},
		{"transfer forward by relay chain", args{"port-5", "channel-5", "port-4/channel-4/port-2/channel-2/kitty"}, true},
		{"transfer back by relay chain", args{"port-6", "channel-6", "port-6/channel-6/port-4/channel-4/port-2/channel-2/kitty"}, false},
		{"transfer back by relay chain", args{"port-4", "channel-4", "port-4/channel-4/port-2/channel-2/kitty"}, false},
		{"transfer back by relay chain", args{"port-2", "channel-2", "port-2/channel-2/kitty"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAwayFromOrigin(tt.args.sourcePort, tt.args.sourceChannel, tt.args.fullClassPath); got != tt.want {
				t.Errorf("IsAwayFromOrigin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseClassTrace(t *testing.T) {
	type args struct {
		rawClassID string
	}
	tests := []struct {
		name string
		args args
		want ClassTrace
	}{
		{"native class", args{"kitty"}, ClassTrace{Path: "", BaseClassId: "kitty"}},
		{"transfer to (port-2,channel-2)", args{"port-2/channel-2/kitty"}, ClassTrace{Path: "port-2/channel-2", BaseClassId: "kitty"}},
		{"transfer to (port-4,channel-4)", args{"port-4/channel-4/port-2/channel-2/kitty"}, ClassTrace{Path: "port-4/channel-4/port-2/channel-2", BaseClassId: "kitty"}},
		{"transfer to (port-6,channel-6)", args{"port-6/channel-6/port-4/channel-4/port-2/channel-2/kitty"}, ClassTrace{Path: "port-6/channel-6/port-4/channel-4/port-2/channel-2", BaseClassId: "kitty"}},
		{"native class with /", args{"cat/kitty"}, ClassTrace{Path: "", BaseClassId: "cat/kitty"}},
		{"transfer to (port-2,channel-2) with /", args{"port-2/channel-2/cat/kitty"}, ClassTrace{Path: "port-2/channel-2", BaseClassId: "cat/kitty"}},
		{"transfer to (port-4,channel-4) with /", args{"port-4/channel-4/port-2/channel-2/cat/kitty"}, ClassTrace{Path: "port-4/channel-4/port-2/channel-2", BaseClassId: "cat/kitty"}},
		{"transfer to (port-6,channel-6) with /", args{"port-6/channel-6/port-4/channel-4/port-2/channel-2/cat/kitty"}, ClassTrace{Path: "port-6/channel-6/port-4/channel-4/port-2/channel-2", BaseClassId: "cat/kitty"}},
	}
	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			if got := ParseClassTrace(tests[i].args.rawClassID); !reflect.DeepEqual(got, tests[i].want) {
				t.Errorf("ParseClassTrace() = %v, want %v", got, tests[i].want)
			}
		})
	}
}

func TestClassTrace_GetFullClassPath(t *testing.T) {
	tests := []struct {
		name string
		ct   ClassTrace
		want string
	}{
		{"native class", ClassTrace{Path: "", BaseClassId: "kitty"}, "kitty"},
		{"first  tranfer", ClassTrace{Path: "port-2/channel-2", BaseClassId: "kitty"}, "port-2/channel-2/kitty"},
		{"second tranfer", ClassTrace{Path: "port-4/channel-4/port-2/channel-2", BaseClassId: "kitty"}, "port-4/channel-4/port-2/channel-2/kitty"},
		{"third  tranfer", ClassTrace{Path: "port-6/channel-6/port-4/channel-4/port-2/channel-2", BaseClassId: "kitty"}, "port-6/channel-6/port-4/channel-4/port-2/channel-2/kitty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.GetFullClassPath(); got != tt.want {
				t.Errorf("ClassTrace.GetFullClassPath() = %v, want %v", got, tt.want)
			}
		})
	}
}