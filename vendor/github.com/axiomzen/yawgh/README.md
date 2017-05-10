# yawgh
Yet Another Wrapper for Golang Http Requests

TODO: some simple tests

#### Typical Usage ####

```
func TestRequest() *yawgh.Request {
	// define your typical setup
	return yawgh.New().
		Transport(theConf.Transport).
		DomainHost(theConf.DomainHost).
		Port(uint(theConf.Port)).
		Marshaler(marshaler).
		Unmarshaler(unmarshaler).
		URLComponent(theConf.Version).
		Header(theConf.APITokenHeader, theConf.APIToken).
		ResponseInterceptor(locationChecker)
}
```
...

```
var signup models.Signup
gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
var user models.User
// Will POST to /users/signup
statusCode, err := TestRequest().
      Post(routes.ResourceUsers).
      URLComponent(routes.ResourceSignup).
      RequestBody(&signup).
      ResponseBody(&user).
      Do()
gomega.Expect(err).ToNot(gomega.HaveOccurred())
gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
```

#### Response and Request Interceptors ####

When things are just not going your way, you can debug the requests and responses by adding in different `RequestInterceptor`s and `ResponseInterceptor`s respectivley:

```
// lets you print out the request as it goes out the door
printRequest requestIntFunc = func(r *http.Request, body []byte, err error) error {
	fmt.Printf("\nHTTP Request\n---------\nMethod: %s\nURL: %s\nBody: %s\n", r.Method, r.URL.String(), string(body))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	return nil
}

// lets you print out the response before it gets rendered into the interface
printResponse responseIntFunc = func(r *http.Response, body []byte, contentType string) error {
	fmt.Printf("\nHTTP Response\n---------\nStatus: %s\nContent Type: %s\nBody: %s\n", r.Status, contentType, string(body))
	return nil
}
```
...
```
var signup models.Signup
gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
var user models.User
// Will POST to /users/signup
statusCode, err := TestRequest().
      Post(routes.ResourceUsers + routes.ResourceSignup).
      RequestBody(&signup).
      ResponseBody(&user).
      RequestInterceptor(printRequest).
      ResponseInterceptor(printResponse).
      Do()
gomega.Expect(err).ToNot(gomega.HaveOccurred())
gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
```
#### Pluggable Marshal/UnMarhshal ####

Similiar to interceptors, marshalling is pluggable.  Hatch is configured to wrap the existing functionality in the app code.

#### Header and URL Param Support ####

Simply add a `.Header("key", "value")` or a `.URLParam("key", "value")` to set headers and url parameters (automatically encoded for you) respectivley.
