package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCmd_StdIn_Select_Single_Element(t *testing.T) {
	// Arrange.
	driver = "TestCmd_StdIn_Select_Single_Element"
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
	ioIn = bytes.NewReader([]byte(json))
	ioOut = bytes.NewBuffer(nil)
	ioErr = bytes.NewBuffer(nil)

	// Act.
	os.Args = []string{"./sqj", "SELECT ζ FROM []", "-"}
	main()

	// Assert
	result := ioOut.(*bytes.Buffer).String()
	if strings.Trim(result, "\n") != "1.2020569031595942" {
		t.Error("unexpected result")
	}
}

func TestCmd_StdIn_NestedObject(t *testing.T) {
	// Arrange.
	driver = "TestCmd_StdIn_NestedObject"
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
