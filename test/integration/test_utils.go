package integration

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/axiomzen/authentication/config"
	"github.com/axiomzen/authentication/constants"
	"github.com/axiomzen/authentication/data"
	"github.com/axiomzen/authentication/helpers"
	"github.com/axiomzen/authentication/routes"
	"github.com/axiomzen/envconfig"
	"github.com/axiomzen/yawgh"
	"github.com/joho/godotenv"
	"github.com/onsi/gomega"
	"github.com/twinj/uuid"
	"os"
	"os/exec"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"errors"
	"net/http"
	"time"
)

const (
	TesterToken = ""
)

var (
	fakeUUID = uuid.NewV4().String()
	theApp   *exec.Cmd
	theConf  *config.AUTHENTICATIONConfig
	// Do() does create a new request each time
	// but we may not want to pollute the settings across requests
	marshaler marshalerFunc = func(v interface{}, contentType string) ([]byte, error) {
		return helpers.Marshal(v, contentType)
	}

	unmarshaler unmarshalerFunc = func(data []byte, v interface{}, contentType string) error {
		return helpers.Unmarshal(data, v, contentType)
	}

	// lets you print out the request as it goes out the door
	printRequest requestIntFunc = func(r *http.Request, body []byte, err error) error {
		fmt.Printf("\nHTTP Request\n---------\nMethod: %s\nURL: %s\nBody: %s\nHeaders: %v\n", r.Method, r.URL.String(), string(body), r.Header)
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

	locationChecker responseIntFunc = func(r *http.Response, body []byte, contentType string) error {
		if r.StatusCode == http.StatusCreated && r.Header.Get("Location") == "" {
			return errors.New("Location should be present for 201 requests")
		}
		return nil
	}
)

// declare a func type that implements yawgh.Marshaler
type marshalerFunc func(v interface{}, contentType string) ([]byte, error)

// implememts yawgh.Marshaler
func (m marshalerFunc) Marshal(v interface{}, contentType string) ([]byte, error) {
	return m(v, contentType)
}

// declare a func type that implements yawgh.Unmarshaler
type unmarshalerFunc func(data []byte, v interface{}, contentType string) error

// implements yawgh.Unmarshaler
func (u unmarshalerFunc) Unmarshal(data []byte, v interface{}, contentType string) error {
	return u(data, v, contentType)
}

// interceptors
type requestIntFunc func(r *http.Request, body []byte, err error) error

func (rf requestIntFunc) InterceptRequest(r *http.Request, body []byte, err error) error {
	return rf(r, body, err)
}

type responseIntFunc func(r *http.Response, body []byte, contentType string) error

func (rf responseIntFunc) InterceptResponse(r *http.Response, body []byte, contentType string) error {
	return rf(r, body, contentType)
}

func testRequest() *yawgh.Request {
	return yawgh.New().
		Transport("http").
		DomainHost(theConf.TestDomainHost).
		Port(uint(theConf.Port)).
		Marshaler(marshaler).
		Unmarshaler(unmarshaler).
		Header(theConf.APITokenHeader, theConf.APIToken).
		ResponseInterceptor(locationChecker)
}

// TestRequestV1 gets a new configured TestRequest using API v1
func TestRequestV1() *yawgh.Request {
	return testRequest().URLComponent(routes.V1)
}

func getTempConf() *config.AUTHENTICATIONConfig {
	// create a new (temp) conf to call
	// Our DAL, then we can create and setup the DB
	// TODO: we could use a "fill defaults" call on
	// env config that would populate all the config
	// vars with their default values
	tempConf := &config.AUTHENTICATIONConfig{}
	// postgres
	tempConf.PostgreSQLHost = "localhost"
	tempConf.PostgreSQLUsername = "postgres"
	tempConf.PostgreSQLPassword = ""
	tempConf.PostgreSQLDatabase = "template1"
	f := false
	tempConf.PostgreSQLSSL = &f
	tempConf.PostgreSQLPort = 5432
	tempConf.PostgreSQLRetryNumTimes = 10
	tempConf.PostgreSQLRetrySleepTime = time.Second * 30

	return tempConf
}

func createDatabase() {
	// connect to a database (temporarily)
	testDAL, err := data.CreateProvider(getTempConf())
	defer func(dal data.AUTHENTICATIONProvider) {
		if dal != nil {
			dal.Close()
		}
	}(testDAL)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(testDAL.Create()).To(gomega.Succeed(), "it should be able to create the database")
}

func dropDatabase() {
	testDAL, err := data.CreateProvider(getTempConf())
	defer func(dal data.AUTHENTICATIONProvider) {
		if dal != nil {
			dal.Close()
		}
	}(testDAL)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(testDAL.Drop()).To(gomega.Succeed(), "it should be able to drop the database")
}

func setupDatabase() {

	// postgres
	// were going to use the migrate tool
	// we can use this for mongodb as soon as https://github.com/mattes/migrate/pull/118 is merged
	// for tests we should have the files here
	// TODO: once a number of migrations happen, we will need to merge
	// them into a single file to reduce time taken to bring a database up from scratch
	// the way to do that is a full schema dump, but we wouldn't have the schema dump nessesarily

	// perform all up migrations
	m, err := migrate.New(
		"file://../../data/migrations",
		"postgres://postgres@localhost:5432/dulpitr9o7a88d?sslmode=disable")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	defer m.Close()
	err = m.Up()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

}

func teardownDatabase() {
	// run all down migrations
	// postgres
	// this is really just a check that the migrations work in both directions
	// which might be useful for development, or a huge pain, not sure yet.
	m, err := migrate.New(
		"file://../../data/migrations",
		"postgres://postgres@localhost:5432/dulpitr9o7a88d?sslmode=disable")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	defer m.Close()
	err = m.Down()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// do nothing for mongodb
}

func fireUpApp() error {

	// (OPTIONAL) set our custom time format to something else for testing
	// TODO: investigate database precision loss with this (pg!)
	//null.SetFormat(constants.TimeFormat)

	fmt.Println("Setting up database...")

	// create (with template1)
	createDatabase()
	// setup (on real database)
	setupDatabase()

	fmt.Println("Running go install...")

	install := exec.Command("go", "install")
	install.Dir = "../.."

	if out, err := install.CombinedOutput(); err != nil {
		return err
	} else if len(out) != 0 {
		// Successful go install is silent
		fmt.Println(string(out))
		return fmt.Errorf(string(out))
	}

	fmt.Println("Firing up the app...")

	// populated from Hatch
	theApp = exec.Command("authentication")
	// so lets force the tests to be called from
	// `make test` which will boot up a db on the fly
	// so we can then hard code these values based on hatch
	// and it will be identical for CI
	// because we never want someone running the tests (locally)
	// against non-local databases or caches, etc

	// how docker plays into this I am not sure yet

	// this is the total env; it won't see anything else
	// which is what we want for local tests and CI
	//
	// so we need to somehow make sure these are all set
	// we need a lib or function that takes a config var and writes out this array
	theConf = &config.AUTHENTICATIONConfig{}
	// let the  defaults take care of most of it
	// the others are hardcoded
	theConf.APIToken = "123"
	theConf.HashSecret = "secret"
	theConf.Transport = "http"
	theConf.DomainHost = "localhost"
	theConf.Environment = constants.EnvironmentTest
	theConf.LogLevel = log.InfoLevel.String()

	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println(err.Error())
	}

	// anything with `envconfig` needs to be set as well
	//conf.AccessorServiceFQDN = "ignore" // can't be "" as that is zero value for string
	//conf.AccessorPort = "5000"
	// database stuff
	// postgres
	theConf.PostgreSQLHost = "localhost"
	theConf.PostgreSQLUsername = "postgres"
	theConf.PostgreSQLPassword = ""
	theConf.PostgreSQLDatabase = "dulpitr9o7a88d"
	f := false
	theConf.PostgreSQLSSL = &f
	//theConf.PostgreSQLPort = default is ok

	// compute dependent variables
	gomega.Expect(theConf.ComputeDependents()).To(gomega.Succeed())

	// we want export and fill in defaults
	// populated from Hatch
	envi, err := envconfig.Export("AXIOMZEN_AUTHENTICATION", theConf, true)
	if err != nil {
		return err
	}
	// debugging:
	// for _, s := range envi {
	// 	fmt.Println(s)
	// }
	theApp.Env = envi
	// If you want to see the printouts to travis logs, uncomment below
	theApp.Stdout = os.Stdout
	theApp.Stderr = os.Stderr
	theApp.Args = []string{"-race"}
	err = theApp.Start()

	if err != nil {
		return err
	}

	// for posterity: sleep command is not cross-platform
	//sleep := exec.Command("sleep", "1")

	fmt.Println("Sleeping 1 second...")
	time.Sleep(1 * time.Second)

	return nil
}

func killApp() error {

	fmt.Println("Dropping the database...")

	var err error

	fmt.Println("Killing the app...")
	if theApp.Process != nil {
		err = theApp.Process.Kill()
		//gomega.Expect(err).ToNot(gomega.HaveOccurred())

	} else {
		err = errors.New("App Process was nil!")
	}

	// teardown database
	// runs the down migrations
	teardownDatabase()

	// drop the database
	dropDatabase()

	return err
}
