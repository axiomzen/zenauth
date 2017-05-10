[![Build Status](https://travis-ci.com/axiomzen/compare.svg?token=TSvXZ1trTYUyHzWyBZYj&branch=master)](https://travis-ci.com/axiomzen/compare)

# compare

A more useful replacement for `reflect.DeepCompare`.  The main functionality of `compare` is to use the `DeepEquals()` method to compare two structs for equality.  `DeepCompare` will automatically traverse arbitrarily complex structs for you and generally do the sensible thing.

TODO: more unit tests

#### Example Usage ####

```
// Ping is our basic health check response
type Ping struct {
	Ping string `form:"ping"                             json:"ping" lorem:",pong"`
}

ginkgo.It("should be able to ping our service", func() {
	// get response
	var pingBack models.Ping
	statusCode, err := TestRequest().Get(routes.ResourcePing).ResponseBody(&pingBack).Do()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	// check status code
	gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
	// check result (could also use lorem.Fill(&expectedPing))
	expectedPing := models.Ping{Ping:"pong"}
	// compares all fields in pingBack to expectedPing (just the Ping field currently)
	gomega.Ω(compare.New().DeepEquals(pingBack, expectedPing, "Ping")).Should(gomega.Succeed())
})
```

#### Tolerances ####

`compare` allows you to specify tolerances for `time.Time` and floats.

```
compare.New().WithFloatEpsilon(0.00001).Float64(1.001, 1.0 , "Float Test")
```

will return:

```
"Float Test: Floats should match: 1.001, 1.0"
```

#### Ignore Fields ####

One can also pass in fields to ignore to `DeepEquals`:

```
// will ignore the ID field of the `User` struct
compare.New().Ignore(".ID").DeepEquals(&requestUser, &responseUser, "User")
```

#### Custom Comparison ####

`compare` also allows you to implement the `Valuable` interface in case you want to return a different value for comparison (useful for null* types)

```
type TreatMeAsAString struct {
	LookHere string
	NothingToSeeHere int
}

func (b *TreatMeAsAString) GetValue() reflect.Value {
  return reflect.ValueOf(b.LookHere)
}

```

#### Time Comparison ####

Times will be compared location independently so that an instant is an instant regardless of time zone. Tolerances can also be set for the nanoseconds as usually precision is typically lost through transmission.
