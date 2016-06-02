package serialBP

import (
	"testing"
)

func TestSerializeSliceFunc(t *testing.T) {
	tp := TypeString("[]int")
	sf := Serialize(tp)

	if sf.Marshal("test", "test") != "gothicHelpers.MarshalSliceint(test)" {
		t.Error("Expected: gothicHelpers.MarshalSliceint(test)")
	}
	if sf.Unmarshal("test", "test") != "gothicHelpers.UnmarshalSliceint(test)" {
		t.Error("Expected: gothicHelpers.UnmarshalSliceint(test)")
	}

	if sf.Marshal("test", "gothicHelpers") != "MarshalSliceint(test)" {
		t.Error("Expected: MarshalSliceint(test)")
	}
	if sf.Unmarshal("test", "gothicHelpers") != "UnmarshalSliceint(test)" {
		t.Error("Expected: UnmarshalSliceint(test)")
	}
}
