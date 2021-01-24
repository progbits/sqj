package main

import (
	"bytes"
	"fmt"
	"github.com/progbits/sqjson/internal/vtable"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	query    string
	expected string
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestCmd_Numerical_Formatting(t *testing.T) {
	// Arrange.
	json := `
		{
				"α": 0.0072973525693,
				"γ": 0.5772156649015328606065120900824024310421,
				"δ": 4.669201609102990671853203820466,
				"ϵ": 8.854187812813e12,
				"ζ": 1.202056903159594285399738161511449990764986292,
				"θ": 90,
				"μ": 1.2566370614E-6,
				"ψ": 3.359885666243177553172011302918927179688905133732
		}
	`

	testCases := []TestCase{
		{"α", "0.0072973525693"},
		{"γ", "0.5772156649015329"},
		{"δ", "4.66920160910299"},
		{"ϵ", "8854187812813"},
		{"ζ", "1.2020569031595942"},
		{"θ", "90"},
		{"μ", "1.2566370614e-06"},
		{"ψ", "3.3598856662431777"},
	}

	for i := 0; i < len(testCases); i++ {
		vtable.Driver = fmt.Sprintf("TestCmd_StdIn_Select_Single_Element_%d", i)
		ioIn = bytes.NewReader([]byte(json))
		ioOut = bytes.NewBuffer(nil)
		ioErr = bytes.NewBuffer(nil)

		// Act.
		os.Args = []string{"./sqj", fmt.Sprintf("SELECT %s FROM []", testCases[i].query), "-"}
		main()

		// Assert
		result := ioOut.(*bytes.Buffer).String()
		if strings.Trim(result, "\n") != testCases[i].expected {
			t.Error("unexpected result")
		}
	}
}

func TestCmd_StdIn_NestedObject(t *testing.T) {
	// Arrange.
	vtable.Driver = "TestCmd_StdIn_NestedObject"
	json := `
		[
		  {
			"id": "5ff8d1fbdc7d0c09c7138193",
			"guid": "dc640ea9-28f8-4ff6-8609-cfc6157e52bc",
			"isActive": true,
			"about": {
			  "registered": "2014-11-04T10:37:35 -00:00",
			  "metric": -24.04675
			}
		  },
		  {
			"id": "5ff8d1fbd468a44558c9e928",
			"guid": "e35d7b31-77fe-428c-bf25-1f19a08e058b",
			"isActive": true,
			"about": {
			  "registered": "2016-03-21T09:47:51 -00:00",
			  "metric": 39.39742
			}
		  },
		  {
			"id": "5ff8d1fbe962cc214df87658",
			"guid": "0223e51e-c4a1-4bc0-8539-2472ca5c3b5d",
			"isActive": false,
			"about": {
			  "registered": "2015-03-11T08:47:56 -00:00",
			  "metric": -70.302563
			}
		  },
		  {
			"id": "5ff8d1fb95ec8b0fac4174bc",
			"guid": "6b9cae6a-4721-4aa4-987f-243127b09dbd",
			"isActive": true,
			"about": {
			  "registered": "2018-04-02T05:26:50 -01:00",
			  "metric": -42.721089
			}
		  },
		  {
			"id": "5ff8d1fbea8b81e2c427b2e9",
			"guid": "88fc8c94-f60c-46ea-a1e4-5becfb06f19d",
			"isActive": true,
			"about": {
			  "registered": "2015-12-17T09:14:19 -00:00",
			  "metric": 41.919667
			}
		  },
		  {
			"id": "5ff8d1fbd17bdee1c0768755",
			"guid": "44bc11db-e936-41af-a1a3-11ad348f4cd6",
			"isActive": false,
			"about": {
			  "registered": "2015-11-06T05:15:14 -00:00",
			  "metric": -67.967849
			}
		  }
		]
	`
	ioIn = bytes.NewReader([]byte(json))
	ioOut = bytes.NewBuffer(nil)
	ioErr = bytes.NewBuffer(nil)

	// Act.
	os.Args = []string{"./sqj", "SELECT about$metric FROM []", "-"}
	main()

	// Assert
	result := ioOut.(*bytes.Buffer).String()
	result = strings.Trim(result, "\n")
	expected := []string{
		"-24.04675",
		"39.39742",
		"-70.302563",
		"-42.721089",
		"41.919667",
		"-67.967849",
	}

	splitResult := strings.Split(result, "\n")
	if len(splitResult) != len(expected) {
		t.Error("unexpected number of values")
	}

	for i, value := range splitResult {
		if strings.Trim(value, "\n") != expected[i] {
			t.Error("unexpected values")
		}
	}
}

func TestCmd_StdIn_SelectFromSubArray(t *testing.T) {
	// Arrange.
	vtable.Driver = "TestCmd_StdIn_SelectFromSubArray"
	json := `
		[
		  {
			"guid": "9b19d50a-cec1-42e9-b7e3-8899d426a541",
			"isActive": false,
			"metric": [
			  {
				"lag": 20.67871,
				"skew": -147.55678
			  },
			  {
				"lag": -33.50249,
				"skew": 96.342544
			  },
			  {
				"lag": -78.999041,
				"skew": -73.063277
			  }
			]
		  },
		  {
			"guid": "2e7a1b41-306f-4fad-86b4-d07e49cd6e4f",
			"isActive": true,
			"metric": [
			  {
				"lag": -10.764641,
				"skew": 129.430546
			  },
			  {
				"lag": -84.682348,
				"skew": 129.620258
			  },
			  {
				"lag": -61.955773,
				"skew": -104.713877
			  }
			]
		  },
		  {
			"guid": "c73211b3-65b4-4d0c-8806-3c3b4dfecff0",
			"isActive": true,
			"metric": [
			  {
				"lag": -60.446643,
				"skew": 109.276407
			  },
			  {
				"lag": 52.830741,
				"skew": 54.130786
			  },
			  {
				"lag": 56.008626,
				"skew": -26.937118
			  }
			]
		  }
		]
	`
	ioIn = bytes.NewReader([]byte(json))
	ioOut = bytes.NewBuffer(nil)
	ioErr = bytes.NewBuffer(nil)

	// Act.
	os.Args = []string{"./sqj", "SELECT lag FROM metric WHERE skew > 0", "-"}
	main()

	// Assert
	result := ioOut.(*bytes.Buffer).String()
	result = strings.Trim(result, "\n")
	expected := []string{
		"-33.50249",
		"-10.764641",
		"-84.682348",
		"-60.446643",
		"52.830741",
	}

	splitResult := strings.Split(result, "\n")
	if len(splitResult) != len(expected) {
		t.Error("unexpected number of values")
	}

	for i, value := range splitResult {
		if strings.Trim(value, "\n") != expected[i] {
			t.Error("unexpected values")
		}
	}
}

func TestCmd_StdIn_SelectFromSubArray_DuplicateColumns(t *testing.T) {
	// Arrange.
	vtable.Driver = "TestCmd_StdIn_SelectFromSubArray_DuplicateColumns"
	json := `
		[
		  {
			"guid": "9b19d50a-cec1-42e9-b7e3-8899d426a541",
			"isActive": false,
			"metric": [
			  {
				"lag": 20.67871,
				"skew": -147.55678
			  },
			  {
				"lag": -33.50249,
				"skew": 96.342544
			  },
			  {
				"lag": -78.999041,
				"skew": -73.063277
			  }
			]
		  },
		  {
			"guid": "2e7a1b41-306f-4fad-86b4-d07e49cd6e4f",
			"isActive": true,
			"metric": [
			  {
				"lag": -10.764641,
				"skew": 129.430546
			  },
			  {
				"lag": -84.682348,
				"skew": 129.620258
			  },
			  {
				"lag": -61.955773,
				"skew": -104.713877
			  }
			]
		  },
		  {
			"guid": "c73211b3-65b4-4d0c-8806-3c3b4dfecff0",
			"isActive": true,
			"metric": [
			  {
				"lag": -60.446643,
				"skew": 109.276407
			  },
			  {
				"lag": 52.830741,
				"skew": 54.130786
			  },
			  {
				"lag": 56.008626,
				"skew": -26.937118
			  }
			]
		  }
		]
	`
	ioIn = bytes.NewReader([]byte(json))
	ioOut = bytes.NewBuffer(nil)
	ioErr = bytes.NewBuffer(nil)

	// Act.
	os.Args = []string{"./sqj", "SELECT skew, lag FROM metric WHERE skew > 0", "-"}
	main()

	// Assert
	result := ioOut.(*bytes.Buffer).String()
	result = strings.Trim(result, "\n")
	expected := []string{
		"96.342544",
		"-33.50249",
		"129.430546",
		"-10.764641",
		"129.620258",
		"-84.682348",
		"109.276407",
		"-60.446643",
		"54.130786",
		"52.830741",
	}

	splitResult := strings.Split(result, "\n")
	if len(splitResult) != len(expected) {
		t.Error("unexpected number of values")
	}

	for i, value := range splitResult {
		if strings.Trim(value, "\n") != expected[i] {
			t.Errorf("unexpected value %s", value)
		}
	}
}

func TestCmd_StdIn_JoinSubArrays(t *testing.T) {
	// Arrange.
	json := `
	  {
		"a": [
			{"id": 1},
			{"id": 2},
			{"id": 3},
			{"id": 4}
		],
		"b": [
			{"value": 3},
			{"value": 4},
			{"value": 5},
			{"value": 6}
		]
	  }
	`

	type TestCase struct {
		statement string
		expected  []string
	}
	cases := []TestCase{
		{
			"SELECT a.id, b.value FROM a, b;",
			[]string{
				"1", "3", "1", "4", "1", "5", "1", "6",
				"2", "3", "2", "4", "2", "5", "2", "6",
				"3", "3", "3", "4", "3", "5", "3", "6",
				"4", "3", "4", "4", "4", "5", "4", "6",
			},
		},
		/*		{
				"SELECT a.id, b.value FROM a JOIN b ON a.id == b.value;",
				[]string{
					"3", "4", "3", "4",
				},
			},*/
	}

	for i, test := range cases {
		vtable.Driver = fmt.Sprintf("TestCmd_StdIn_JoinSubArrays_%d", i)
		ioIn = bytes.NewReader([]byte(json))
		ioOut = bytes.NewBuffer(nil)
		ioErr = bytes.NewBuffer(nil)

		// Act.
		os.Args = []string{"./sqj", test.statement, "-"}
		main()

		// Assert.
		result := ioOut.(*bytes.Buffer).String()
		result = strings.Trim(result, "\n")

		splitResult := strings.Split(result, "\n")
		if len(splitResult) != len(test.expected) {
			t.Error("unexpected number of values")
		}

		for i, value := range splitResult {
			if strings.Trim(value, "\n") != test.expected[i] {
				t.Error("unexpected values")
			}
		}
	}
}
