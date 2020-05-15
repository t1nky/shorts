package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"shorts/database"
	h "shorts/helper"
	"shorts/models"
	"shorts/router"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func getEmptyStringMap() map[string]string {
	return map[string]string{}
}

func encodeCredentials(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func DeleteCreatedEntities(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"

	// Setup the onCreate Hook
	db.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		fmt.Printf("Inserted entities of %s with %s=%v\n", scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
		entries = append(entries, entity{table: scope.TableName(), keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
	})
	return func() {
		// Remove the hook once we're done
		defer db.Callback().Create().Remove(hookName)
		// Find out if the current db object is already a transaction
		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}
		// Loop from the end. It is important that we delete the entries in the
		// reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
}

func performRequest(r http.Handler, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	var err error
	if len(body) > 0 {
		req, err = http.NewRequest(method, path, bytes.NewBufferString(body))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}

	if err != nil {
		fmt.Println("Cannot perform request: " + err.Error())
		return nil
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func testJSONUnMarshalling(t *testing.T, w *httptest.ResponseRecorder, record interface{}) bool {
	err := json.Unmarshal([]byte(w.Body.String()), record)
	return assert.Nil(t, err)
}

func testKeyAndValueEquality(t *testing.T, sourceMap map[string]interface{}, key string, value interface{}) bool {
	v, exists := sourceMap[key]
	if assert.True(t, exists) && assert.Equal(t, value, v) {
		return true
	}

	return false
}

func testKeyAndValueNotEmpty(t *testing.T, sourceMap map[string]interface{}, key string) bool {
	value, exists := sourceMap[key]
	if assert.True(t, exists) && assert.NotEmpty(t, value) {
		return true
	}

	return false
}

func testFailedResponse(t *testing.T, w *httptest.ResponseRecorder, expectedResponseCode int) bool {
	if !assert.Equal(t, expectedResponseCode, w.Code) {
		fmt.Println(w)
		return false
	}
	// Convert the JSON response to a map
	var response map[string]interface{}
	if !testJSONUnMarshalling(t, w, &response) {
		return false
	}

	if testKeyAndValueEquality(t, response, "result", "error") {
		return testKeyAndValueNotEmpty(t, response, "error")
	}

	return false
}

func testSuccessfulResponse(t *testing.T, w *httptest.ResponseRecorder, expectedResponseCode int) bool {
	if !assert.Equal(t, expectedResponseCode, w.Code) {
		fmt.Println(w)
		return false
	}
	// Convert the JSON response to a map
	var response map[string]interface{}
	if !testJSONUnMarshalling(t, w, &response) {
		return false
	}

	return testKeyAndValueEquality(t, response, "result", "ok")
}

func testDataResponse(t *testing.T, w *httptest.ResponseRecorder, expectedResponseCode int, response interface{}) bool {
	if !assert.Equal(t, expectedResponseCode, w.Code) {
		fmt.Println(w)
		return false
	}
	// Convert the JSON response to a map
	if !testJSONUnMarshalling(t, w, &response) {
		return false
	}

	return true
}

func testProtectedRouteResponse(t *testing.T, w *httptest.ResponseRecorder) {
	// expectedResponse := gin.H{
	// 	"result": "error",
	// 	"error":  "Authentication required",
	// }

	testFailedResponse(t, w, http.StatusUnauthorized)
}

func testFailedRegistrationResponse(t *testing.T, w *httptest.ResponseRecorder) {
	// expectedResponse := gin.H{
	// 	"result": "error",
	// 	"error":  "Error text",
	// }
	testFailedResponse(t, w, http.StatusBadRequest)
}

func testRegistrationResponse(t *testing.T, w *httptest.ResponseRecorder) (res bool) {
	// expectedResponse := gin.H{
	// 	"result": "ok",
	// }

	return testSuccessfulResponse(t, w, http.StatusCreated)
}

func testAuthenticationResponse(t *testing.T, w *httptest.ResponseRecorder, username string) (res bool) {
	var response models.UserResponse

	testDataResponse(t, w, http.StatusOK, &response)
	res = assert.Equal(t, username, response.Data.Name)

	return
}

func TestProtectedRoutesError(t *testing.T) {
	// Init local env
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	db, err := InitDatabase()
	if !assert.Nil(t, err) {
		fmt.Println("Cannot connect to the database:" + err.Error())
		return
	}
	defer db.Close()

	cleaner := DeleteCreatedEntities(database.DB)
	defer cleaner()

	r := router.SetupRouter()
	// Test that protected routes are actually protected
	testProtectedRouteResponse(t, performRequest(r, "GET", "/v1/me", "", getEmptyStringMap()))
	testProtectedRouteResponse(t, performRequest(r, "GET", "/v1/shorts", "", getEmptyStringMap()))
	testProtectedRouteResponse(t, performRequest(r, "GET", "/v1/shorts/1", "", getEmptyStringMap()))
	testProtectedRouteResponse(t, performRequest(r, "POST", "/v1/shorts", `{"full":"test.com/test"}`, getEmptyStringMap()))
	testProtectedRouteResponse(t, performRequest(r, "DELETE", "/v1/shorts/1", "", getEmptyStringMap()))
}

func TestMain(t *testing.T) {

	const USER_NAME = "Test Test"
	const USER_PASSWORD = "testPassword123"
	const FULL_LINK = "https://google.com"

	// Init local env
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	db, err := InitDatabase()
	if !assert.Nil(t, err) {
		fmt.Println("Cannot connect to the database:" + err.Error())
		return
	}
	defer db.Close()

	cleaner := DeleteCreatedEntities(database.DB)
	defer cleaner()

	// Initialize WebServer
	r := router.SetupRouter()

	// Register
	reg := performRequest(r, "POST", "/v1/users", `{"name": "`+USER_NAME+`", "password": "`+USER_PASSWORD+`"}`, getEmptyStringMap())
	if testRegistrationResponse(t, reg) {
		encodedCredentials := map[string]string{
			"Authorization": "Basic " + encodeCredentials(USER_NAME, USER_PASSWORD),
		}

		randomCredentials := map[string]string{
			"Authorization": "Basic " + encodeCredentials("123456", "123456"),
		}

		// Authenticate
		if !testFailedResponse(t, performRequest(r, "GET", "/v1/me", "", randomCredentials), http.StatusUnauthorized) {
			return
		}

		// Authenticate
		if !testAuthenticationResponse(t, performRequest(r, "GET", "/v1/me", "", encodedCredentials), USER_NAME) {
			return
		}

		// Should result in an error. As id is not provided, there is no such route
		testFailedResponse(t, performRequest(r, "DELETE", "/v1/shorts/", "", encodedCredentials), http.StatusNotFound)
		// Should result in an error, because we did not create any short link
		testFailedResponse(t, performRequest(r, "DELETE", "/v1/shorts/1", "", encodedCredentials), http.StatusNotFound)
		testFailedResponse(t, performRequest(r, "GET", "/v1/shorts/1", "", encodedCredentials), http.StatusNotFound)

		var shortsResponse models.ShortlinksResponse
		if testDataResponse(t, performRequest(r, "GET", "/v1/shorts", "", encodedCredentials), http.StatusOK, &shortsResponse) { // should be empty list
			assert.Empty(t, shortsResponse.Data)
		}

		var shortlinkResponse models.ShortlinkFullResponse
		if testDataResponse(t, performRequest(r, "POST", "/v1/shorts", `{"full":"`+FULL_LINK+`"}`, encodedCredentials), http.StatusCreated, &shortlinkResponse) { // should be empty list
			if !assert.Equal(t, shortlinkResponse.Data.Full, FULL_LINK) {
				return
			}
			shortlinkID := strconv.FormatUint(shortlinkResponse.Data.ID, 10)

			if testDataResponse(t, performRequest(r, "GET", "/v1/shorts", "", encodedCredentials), http.StatusOK, &shortsResponse) { // should return 1 record
				_ = assert.Len(t, shortsResponse.Data, 1) && assert.Equal(t, shortsResponse.Data[0].Full, FULL_LINK)
			}

			if testDataResponse(t, performRequest(r, "GET", "/v1/shorts/"+shortlinkID, "", encodedCredentials), http.StatusOK, &shortlinkResponse) { // should return information about that record with 0 uses
				assert.Equal(t, shortlinkResponse.Data.Full, FULL_LINK)
			}

			redirect := performRequest(r, "GET", "/v1/s/"+shortlinkResponse.Data.Short, "", encodedCredentials) // should redirect to a full link
			assert.Equal(t, FULL_LINK, redirect.HeaderMap.Get("Location"))

			if testDataResponse(t, performRequest(r, "GET", "/v1/shorts/"+shortlinkID, "", encodedCredentials), http.StatusOK, &shortlinkResponse) { // should return information about that record with 1 use
				_ = assert.Equal(t, shortlinkResponse.Data.Full, FULL_LINK) && assert.Len(t, shortlinkResponse.Data.Uses, 1)
			}

			testSuccessfulResponse(t, performRequest(r, "DELETE", "/v1/shorts/"+shortlinkID, "", encodedCredentials), http.StatusOK)

			testFailedResponse(t, performRequest(r, "GET", "/v1/logout", "", encodedCredentials), http.StatusUnauthorized)
		}
	}
}

func TestValidation(t *testing.T) {
	const USER_NAME = "Test Test"
	const USER_PASSWORD = "testPassword123"

	// Init local env
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	db, err := InitDatabase()
	if !assert.Nil(t, err) {
		fmt.Println("Cannot connect to the database:" + err.Error())
		return
	}
	defer db.Close()

	cleaner := DeleteCreatedEntities(database.DB)
	defer cleaner()

	// Initialize WebServer
	r := router.SetupRouter()

	// Test registration validation
	testFailedRegistrationResponse(t, performRequest(r, "POST", "/v1/users", `{"name": "`+USER_NAME+`"}`, getEmptyStringMap()))
	testFailedRegistrationResponse(t, performRequest(r, "POST", "/v1/users", `{"password": "`+USER_PASSWORD+`"}`, getEmptyStringMap()))
	testFailedRegistrationResponse(t, performRequest(r, "POST", "/v1/users", `{"name": "", "password": ""}`, getEmptyStringMap()))
	testFailedRegistrationResponse(t, performRequest(r, "POST", "/v1/users", "{}", getEmptyStringMap()))

	// Register
	reg := performRequest(r, "POST", "/v1/users", `{"name": "`+USER_NAME+`", "password": "`+USER_PASSWORD+`"}`, getEmptyStringMap())
	if testRegistrationResponse(t, reg) {
		encodedCredentials := map[string]string{
			"Authorization": "Basic " + encodeCredentials(USER_NAME, USER_PASSWORD),
		}

		// Authenticate
		if !testAuthenticationResponse(t, performRequest(r, "GET", "/v1/me", "", encodedCredentials), USER_NAME) {
			return
		}

		// Test POST /shorts validation
		testFailedResponse(t, performRequest(r, "POST", "/v1/shorts", `{}`, encodedCredentials), http.StatusBadRequest)
		testFailedResponse(t, performRequest(r, "POST", "/v1/shorts", `{""}`, encodedCredentials), http.StatusBadRequest)
		testFailedResponse(t, performRequest(r, "POST", "/v1/shorts", ``, encodedCredentials), http.StatusBadRequest)
		testFailedResponse(t, performRequest(r, "POST", "/v1/shorts", `{"full":}`, encodedCredentials), http.StatusBadRequest)
		testFailedResponse(t, performRequest(r, "POST", "/v1/shorts", `{`, encodedCredentials), http.StatusBadRequest)
		testFailedResponse(t, performRequest(r, "POST", "/v1/shorts", `}`, encodedCredentials), http.StatusBadRequest)
	}
}

func TestStats(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	const USER_NAME = "Test Test"
	const USER_PASSWORD = "testPassword123"
	var FULL_LINKS = [25]string{
		"https://google.com/?q=test",
		"https://placeholder.com/",
		"https://postman.com/",
		"https://youtube.com/",
		"https://facebook.com/",
		"https://wikipedia.org/",
		"https://reddit.com/",
		"https://netflix.com/",
		"https://www.amazon.com/This-Goodbye-IMMINENCE/dp/B01N4RYAJQ/ref=sr_1_2?dchild=1&keywords=Imminence&qid=1588013859&sr=8-2",
		"https://www.instagram.com/thisisbillgates",
		"https://napster.com/",
		"https://google.com/?q=something%20else",
		"https://www.w3.org/",
		"https://stackoverflow.com/questions/11227809/why-is-processing-a-sorted-array-faster-than-processing-an-unsorted-array/",
		"https://github.com/",
		"https://trello.com/",
		"https://notion.so/",
		"https://gitlab.com/",
		"https://procatinator.com/",
		"http://www.clicktoremove.com/",
		"http://www.skype.com/",
		"http://make-everything-ok.com/",
		"https://tour.golang.org/flowcontrol/4",
		"https://www.cam.ac.uk/",
		"https://jsonplaceholder.typicode.com/",
	}
	const ADD_LINKS_COUNT = 1000
	const ADD_USES_UP_TO = 100

	// Init local env
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	db, err := InitDatabase()
	if !assert.Nil(t, err) {
		fmt.Println("Cannot connect to the database:" + err.Error())
		return
	}
	defer db.Close()

	cleaner := DeleteCreatedEntities(database.DB)
	defer cleaner()

	// Initialize WebServer
	r := router.SetupRouter()

	// Register
	reg := performRequest(r, "POST", "/v1/users", `{"name": "`+USER_NAME+`", "password": "`+USER_PASSWORD+`"}`, getEmptyStringMap())
	if testRegistrationResponse(t, reg) {
		encodedCredentials := map[string]string{
			"Authorization": "Basic " + encodeCredentials(USER_NAME, USER_PASSWORD),
		}

		var response models.UserResponse

		if testDataResponse(t, performRequest(r, "GET", "/v1/me", "", encodedCredentials), http.StatusOK, &response) {
			if !assert.Equal(t, USER_NAME, response.Data.Name) {
				return
			}

			topDomainsExpected := make(map[string]uint64)
			domainGraphExcepted := make(models.ShortlinksGraphResponseData)
			for i, link := range FULL_LINKS {
				parsedURL, err := url.Parse(link)
				if err != nil {
					continue
				}
				websiteHost := parsedURL.Host

				var shortlinkResponse models.ShortlinkFullResponse
				if testDataResponse(t, performRequest(r, "POST", "/v1/shorts", `{"full":"`+link+`"}`, encodedCredentials), http.StatusCreated, &shortlinkResponse) { // should be empty list
					if !assert.Equal(t, shortlinkResponse.Data.Full, link) {
						return
					}
					shortlinkID := strconv.FormatUint(shortlinkResponse.Data.ID, 10)

					var shortlink models.Shortlink
					if dbc := database.DB.Find(&shortlink, shortlinkID); !assert.Nil(t, dbc.Error) {
						return
					}

					shortlinkUse := models.ShortlinkUse{LinkID: shortlinkResponse.Data.ID, UseTime: time.Date(2020, 02, i, i, i, 02, 02, time.UTC)}
					usesCount := rand.Intn(ADD_USES_UP_TO-1) + 1
					for use := 0; use <= usesCount; use++ {
						shortlink.Uses = append(shortlink.Uses, shortlinkUse)

						topDomainsExpected[websiteHost]++

						models.AddUseToGraph(&domainGraphExcepted, shortlinkUse)
					}

					if dbc := database.DB.Save(&shortlink); !assert.Nil(t, dbc.Error) {
						return
					}
				}
			}

			sortedMap := h.GetTopDomains(topDomainsExpected, 20)

			var topDomains models.TopDomainsResponse
			if testDataResponse(t, performRequest(r, "GET", "/v1/stats/top", "", getEmptyStringMap()), http.StatusOK, &topDomains) {
				// check that only top 20 returned
				if assert.Len(t, topDomains.Data, 20) {
					for _, domain := range topDomains.Data {
						if _, exists := sortedMap[domain.Website]; assert.True(t, exists) {
							if !assert.Equal(t, sortedMap[domain.Website], domain.UsesCount) {
								break
							}
						} else {
							break
						}
					}
				}
			}

			var usesGraph models.ShortlinksGraphResponse
			if testDataResponse(t, performRequest(r, "GET", "/v1/stats/graph", "", getEmptyStringMap()), http.StatusOK, &usesGraph) {
				assert.Equal(t, domainGraphExcepted, usesGraph.Data)
			}
		}
	}
}
