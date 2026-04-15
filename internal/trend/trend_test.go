package trend_test

import (
	"testing"

	"portwatch/internal/trend"
)

func TestNew_ReturnsTwoMinimumWindow(t *testing.T) {
	tr := trend.New(0)
	tr.Record(5)
	tr.Record(10)
	if tr.Direction() != trend.Rising {
		t.Fatalf("expected Rising, got %s", tr.Direction())
	}
}

func TestDirection_StableWithNoSamples(t *testing.T) {
	tr := trend.New(5)
	if tr.Direction() != trend.Stable {
		t.Fatalf("expected Stable, got %s", tr.Direction())
	}
}

func TestDirection_StableWithOneSample(t *testing.T) {
	tr := trend.New(5)
	tr.Record(3)
	if tr.Direction() != trend.Stable {
		t.Fatalf("expected Stable, got %s", tr.Direction())
	}
}

func TestDirection_Rising(t *testing.T) {
	tr := trend.New(5)
	tr.Record(2)
	tr.Record(7)
	if tr.Direction() != trend.Rising {
		t.Fatalf("expected Rising, got %s", tr.Direction())
	}
}

func TestDirection_Falling(t *testing.T) {
	tr := trend.New(5)
	tr.Record(10)
	tr.Record(4)
	if tr.Direction() != trend.Falling {
		t.Fatalf("expected Falling, got %s", tr.Direction())
	}
}

func TestDirection_StableEqualValues(t *testing.T) {
	tr := trend.New(5)
	tr.Record(6)
	tr.Record(6)
	if tr.Direction() != trend.Stable {
		t.Fatalf("expected Stable, got %s", tr.Direction())
	}
}

func TestRecord_CapsAtMaxLen(t *testing.T) {
	tr := trend.New(3)
	for i := 0; i < 10; i++ {
		tr.Record(i)
	}
	samples := tr.Samples()
	if len(samples) != 3 {
		t.Fatalf("expected 3 samples, got %d", len(samples))
	}
	// Newest values: 7, 8, 9
	if samples[0].Count != 7 || samples[2].Count != 9 {
		t.Fatalf("unexpected sample values: %+v", samples)
	}
}

func TestSamples_ReturnsCopy(t *testing.T) {
	tr := trend.New(5)
	tr.Record(1)
	tr.Record(2)
	s := tr.Samples()
	s[0].Count = 999
	original := tr.Samples()
	if original[0].Count == 999 {
		t.Fatal("Samples should return an isolated copy")
	}
}

func TestReset_ClearsSamples(t *testing.T) {
	tr := trend.New(5)
	tr.Record(1)
	tr.Record(2)
	tr.Reset()
	if len(tr.Samples()) != 0 {
		t.Fatal("expected empty samples after Reset")
	}
	if tr.Direction() != trend.Stable {
		t.Fatalf("expected Stable after Reset, got %s", tr.Direction())
	}
}

func TestDirection_String(t *testing.T) {
	cases := []struct {
		d    trend.Direction
		want string
	}{
		{trend.Rising, "rising"},
		{trend.Falling, "falling"},
		{trend.Stable, "stable"},
	}
	for _, c := range cases {
		if got := c.d.String(); got != c.want {
			t.Errorf("Direction(%d).String() = %q, want %q", c.d, got, c.want)
		}
	}
}
