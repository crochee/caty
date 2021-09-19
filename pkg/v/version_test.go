package v

import "testing"

func TestCompareVersion(t *testing.T) {
	list := []struct {
		name     string
		version1 string
		version2 string
		want     int
		fail     bool
	}{
		{
			"==1",
			"1.0.0",
			"1.0.0",
			0,
			false,
		},
		{
			"==2",
			"1.0.0",
			"1.0.s",
			1,
			true,
		},
		{
			"==3",
			"1.0.s",
			"1.0.1",
			-1,
			true,
		},
		{
			"<1",
			"1.0.0",
			"1.0.1",
			-1,
			false,
		},
		{
			"<2",
			"1.0.0",
			"1.0.001",
			-1,
			false,
		},
		{
			">1",
			"1.0.10",
			"1.0.001",
			1,
			false,
		},
	}
	for _, input := range list {
		t.Run(input.name, func(t *testing.T) {
			got, err := CompareVersion(input.version1, input.version2)
			if input.fail {
				return
			}
			if err != nil {
				t.Errorf("got err %v", err)
				return
			}
			if input.want != got {
				t.Errorf("want %d got %d", input.want, got)
			}
		})

	}
}
