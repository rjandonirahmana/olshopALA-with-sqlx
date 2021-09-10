package customer

// NEED RECHECK ON VALIDATION MOCK (RETURN) DATA
// Check only some key value or whole data?

import (
	"testing"

	mocks "graphql/mocks"

	"github.com/stretchr/testify/assert"
)

func InitServiceCustomer(repo Repository) CustomerInt {
	service := new(ServiceCustomer)
	service.repo = repo
	return service
}

func TestEmailAvailable(t *testing.T) {
	var testCases = []struct {
		testName    string
		email       string
		expectBool  bool
		expectError error
	}{
		{
			testName:    "success",
			email:       "jon@j.j",
			expectBool:  true,
			expectError: nil,
		},
	}

	for _, testCase := range testCases {
		mockGetCustomerByEmailrepo := new(mocks.GetCustomerByEmail)
		mockGetCustomerByEmailrepo.On("IsEmailAvailable", testCase.email).Return(testCase.expectBool, testCase.expectError)

		s := InitServiceCustomer(mockGetCustomerByEmailrepo)

		status, err := s.IsEmailAvailable(testCase.email)

		assert.Equal(t, testCase.expectBool, status)
		assert.Equal(t, testCase.expectError, err)
	}
}

func TestRegisterCustomer(t *testing.T) {
	var testCases = []struct {
		testName string
		// id int
		name            string
		email           string
		password        string
		confirmPassword string
		expectError     error
	}{
		{
			testName: "success",
			// id: 1,
			name:            "joni",
			email:           "jon@j.j",
			password:        "jojo",
			confirmPassword: "jojo",
			expectError:     nil,
		},
	}

	for _, testCase := range testCases {
		customerTest := Customer{
			Name:     testCase.name,
			Email:    testCase.email,
			Password: testCase.password,
		}
		mockRegisterUserrepo := new(mocks.RegisterUser)
		mockRegisterUserrepo.On("Register", customerTest).Return(testCase, testCase.expectError)

		s := InitServiceCustomer(mockRegisterUserrepo)

		customer, err := s.Register(customerTest)

		assert.Equal(t, testCase.name, customer.Name)
		assert.Equal(t, testCase.expectError, err)
	}

}

func TestLoginCustomer(t *testing.T) {
	var testCases = []struct {
		testName string
		// id int
		name        string
		email       string
		password    string
		expectError error
	}{
		{
			testName: "success",
			// id: 1,
			name:        "joni",
			email:       "jon@j.j",
			password:    "jojo",
			expectError: nil,
		},
	}

	for _, testCase := range testCases {
		inputTest := InputLogin{
			Email:    testCase.email,
			Password: testCase.password,
		}
		mockGetCustomerByEmailrepo := new(mocks.GetCustomerByEmail)
		mockGetCustomerByEmailrepo.On("LoginCustomer", inputTest).Return(testCase, testCase.expectError)

		s := InitServiceCustomer(mockGetCustomerByEmailrepo)

		customer, err := s.LoginCustomer(inputTest)

		assert.Equal(t, testCase.email, customer.Email)
		assert.Equal(t, testCase.expectError, err)
	}
}

func TestUpdateCustomerPhone(t *testing.T) {
	var testCases = []struct {
		testName string
		// id int
		name        string
		email       string
		phone       string
		expectError error
	}{
		{
			testName: "success",
			// id: 1,
			name:        "joni",
			email:       "jon@j.j",
			phone:       "jojo",
			expectError: nil,
		},
	}

	for _, testCase := range testCases {
		email := testCase.email
		phone := testCase.phone
		mockGetCustomerByEmailrepo := new(mocks.GetCustomerByEmail)
		mockGetCustomerByEmailrepo.On("UpdateCustomerPhone", phone, email).Return(testCase.expectError)

		s := InitServiceCustomer(mockGetCustomerByEmailrepo)

		err := s.UpdateCustomerPhone(email, phone)

		assert.Equal(t, testCase.expectError, err)
	}
}

func TestGetCustomerById(t *testing.T) {
	var testCases = []struct{
		testName string
		id int
		name string
		email string
		expectError error
	}{
		{
			testName: "success",
			id: 1,
			name: "joni",
			email: "jon@j.j",
			expectError: nil,
		},
	}

	for _, testCase := range testCases {
		id := testCase.id
		mockGetCustomerById := new(mocks.GetCustomerByID)
		mockGetCustomerById.On("GetCustomerById", id).Return(Customer, error)

		s := InitServiceCustomer(mockGetCustomerById)

		customer, err := s.GetCustomerByID(id)

		assert.Equal(t, testCase.id, customer.ID)
		assert.Equal(t, testCase.expectError, err)
	}
}

// func TestChangeProfile(t *testing.T) {
// 	var testCases = []struct{
// 		testName string
// 		id int
// 		name string
// 		profileMock string
// 	}{
// 		{
// 			testName: "success",
// 			id: 1,
// 			name: "joni",
// 			profileMock: "aaa",
// 		},
// 	}
	
// 	for _, testCase := range testCases {
// 		avatar, _ := ioutil.ReadAll(strings.NewReader(testCase.profileMock))


// 	}
// }