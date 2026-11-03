package svc

import (
	"context"

	pb "{{ cookiecutter.go_module_path.strip('/') }}/api/go/{{ cookiecutter.project_name.strip() }}/v1"
)

// TestAddUserRPC tests the AddUser RPC method
func (suite *SvcTestSuite) TestAddUserRPC() {
	testCases := []struct {
		name        string
		input       *pb.AddUserRequest
		wantErr     bool
		errContains string
		checkResult func(*pb.AddUserResponse)
	}{
		{
			name: "valid user",
			input: &pb.AddUserRequest{
				Name: "John Doe",
				Age:  30,
			},
			wantErr: false,
			checkResult: func(resp *pb.AddUserResponse) {
				suite.NotNil(resp.User)
				suite.Equal("John Doe", resp.User.Name)
				suite.Equal(int32(30), resp.User.Age)
				suite.NotZero(resp.User.Id)
			},
		},
		{
			name: "empty name",
			input: &pb.AddUserRequest{
				Name: "",
				Age:  30,
			},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name: "negative age",
			input: &pb.AddUserRequest{
				Name: "John Doe",
				Age:  -1,
			},
			wantErr:     true,
			errContains: "age must be positive",
		},
		{
			name: "zero age",
			input: &pb.AddUserRequest{
				Name: "John Doe",
				Age:  0,
			},
			wantErr:     true,
			errContains: "age must be positive",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Clear cache before each sub-test to ensure isolation
			if suite.cache != nil && suite.cache.GetRedisClient() != nil {
				suite.cache.GetRedisClient().FlushDB(context.Background())
			}

			resp, err := suite.svc.AddUser(context.Background(), tc.input)

			if tc.wantErr {
				suite.Error(err)
				if tc.errContains != "" && err != nil {
					suite.Contains(err.Error(), tc.errContains)
				}
				suite.Nil(resp)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
				if tc.checkResult != nil {
					tc.checkResult(resp)
				}
			}
		})
	}
}

// TestGetUserRPC tests the GetUser RPC method
func (suite *SvcTestSuite) TestGetUserRPC() {
	ctx := context.Background()

	testCases := []struct {
		name        string
		setup       func() int64 // Returns user ID to query
		input       *pb.GetUserRequest
		wantErr     bool
		errContains string
		checkResult func(*pb.GetUserResponse)
	}{
		{
			name: "existing user",
			setup: func() int64 {
				// Create a user first
				resp, err := suite.svc.AddUser(ctx, &pb.AddUserRequest{
					Name: "Alice",
					Age:  25,
				})
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				return resp.User.Id
			},
			input:   nil, // Will be set in test
			wantErr: false,
			checkResult: func(resp *pb.GetUserResponse) {
				suite.NotNil(resp.User)
				suite.Equal("Alice", resp.User.Name)
				suite.Equal(int32(25), resp.User.Age)
				suite.NotZero(resp.User.Id)
			},
		},
		{
			name: "non_existent_user",
			setup: func() int64 {
				// No setup needed for non-existent user
				return 999
			},
			input:       nil, // Will be set in test
			wantErr:     true,
			errContains: "user not found",
		},
		{
			name: "zero_id",
			setup: func() int64 {
				return 0
			},
			input:       nil, // Will be set in test
			wantErr:     true,
			errContains: "user not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Clear cache before each sub-test to ensure isolation
			if suite.cache != nil && suite.cache.GetRedisClient() != nil {
				suite.cache.GetRedisClient().FlushDB(ctx)
			}

			// Setup: create user if needed
			userId := tc.setup()

			// Set input if not provided
			if tc.input == nil {
				tc.input = &pb.GetUserRequest{
					Id: userId,
				}
			}

			resp, err := suite.svc.GetUser(ctx, tc.input)

			if tc.wantErr {
				suite.Error(err)
				if tc.errContains != "" && err != nil {
					suite.Contains(err.Error(), tc.errContains)
				}
				suite.Nil(resp)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
				if tc.checkResult != nil {
					tc.checkResult(resp)
				}
			}
		})
	}
}

// TestUpdateUserRPC tests the UpdateUser RPC method
func (suite *SvcTestSuite) TestUpdateUserRPC() {
	ctx := context.Background()

	testCases := []struct {
		name        string
		setup       func() int64                             // Returns user ID to update
		buildInput  func(userId int64) *pb.UpdateUserRequest // Builds request from user ID
		wantErr     bool
		errContains string
		checkResult func(*pb.UpdateUserResponse)
	}{
		{
			name: "valid_update",
			setup: func() int64 {
				// Create a user first
				resp, err := suite.svc.AddUser(ctx, &pb.AddUserRequest{
					Name: "Alice",
					Age:  25,
				})
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				return resp.User.Id
			},
			buildInput: func(userId int64) *pb.UpdateUserRequest {
				return &pb.UpdateUserRequest{
					Id:   userId,
					Name: "Updated Name",
					Age:  35,
				}
			},
			wantErr: false,
			checkResult: func(resp *pb.UpdateUserResponse) {
				suite.NotNil(resp.User)
				suite.Equal("Updated Name", resp.User.Name)
				suite.Equal(int32(35), resp.User.Age)
				suite.NotZero(resp.User.Id)
			},
		},
		{
			name: "non_existent_user",
			setup: func() int64 {
				// No setup needed for non-existent user
				return 999
			},
			buildInput: func(userId int64) *pb.UpdateUserRequest {
				return &pb.UpdateUserRequest{
					Id:   userId,
					Name: "Updated Name",
					Age:  35,
				}
			},
			wantErr:     true,
			errContains: "user not found",
		},
		{
			name: "update_to_empty_name",
			setup: func() int64 {
				// Create a user first
				resp, err := suite.svc.AddUser(ctx, &pb.AddUserRequest{
					Name: "Alice",
					Age:  25,
				})
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				return resp.User.Id
			},
			buildInput: func(userId int64) *pb.UpdateUserRequest {
				return &pb.UpdateUserRequest{
					Id:   userId,
					Name: "",
					Age:  35,
				}
			},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name: "negative_age",
			setup: func() int64 {
				// Create a user first
				resp, err := suite.svc.AddUser(ctx, &pb.AddUserRequest{
					Name: "Alice",
					Age:  25,
				})
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				return resp.User.Id
			},
			buildInput: func(userId int64) *pb.UpdateUserRequest {
				return &pb.UpdateUserRequest{
					Id:   userId,
					Name: "Updated Name",
					Age:  -1,
				}
			},
			wantErr:     true,
			errContains: "age must be positive",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Clear cache before each sub-test to ensure isolation
			if suite.cache != nil && suite.cache.GetRedisClient() != nil {
				suite.cache.GetRedisClient().FlushDB(ctx)
			}

			// Setup: create user if needed
			userId := tc.setup()

			// Build input from user ID
			input := tc.buildInput(userId)

			resp, err := suite.svc.UpdateUser(ctx, input)

			if tc.wantErr {
				suite.Error(err)
				if tc.errContains != "" && err != nil {
					suite.Contains(err.Error(), tc.errContains)
				}
				suite.Nil(resp)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
				if tc.checkResult != nil {
					tc.checkResult(resp)
				}
			}
		})
	}
}

// TestDeleteUserRPC tests the DeleteUser RPC method
func (suite *SvcTestSuite) TestDeleteUserRPC() {
	ctx := context.Background()

	testCases := []struct {
		name        string
		setup       func() int64 // Returns user ID to delete
		input       *pb.DeleteUserRequest
		wantErr     bool
		errContains string
		checkResult func(*pb.DeleteUserResponse)
	}{
		{
			name: "valid_delete",
			setup: func() int64 {
				// Create a user first
				resp, err := suite.svc.AddUser(ctx, &pb.AddUserRequest{
					Name: "Alice",
					Age:  25,
				})
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				return resp.User.Id
			},
			input:   nil, // Will be set in test
			wantErr: false,
			checkResult: func(resp *pb.DeleteUserResponse) {
				suite.NotNil(resp.User)
				suite.Equal("Alice", resp.User.Name)
				suite.Equal(int32(25), resp.User.Age)
				suite.NotZero(resp.User.Id)
			},
		},
		{
			name: "non_existent_user",
			setup: func() int64 {
				// No setup needed for non-existent user
				return 999
			},
			input:       nil, // Will be set in test
			wantErr:     true,
			errContains: "user not found",
		},
		{
			name: "zero_id",
			setup: func() int64 {
				return 0
			},
			input:       nil, // Will be set in test
			wantErr:     true,
			errContains: "user not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Clear cache before each sub-test to ensure isolation
			if suite.cache != nil && suite.cache.GetRedisClient() != nil {
				suite.cache.GetRedisClient().FlushDB(ctx)
			}

			// Setup: create user if needed
			userId := tc.setup()

			// Set input if not provided
			if tc.input == nil {
				tc.input = &pb.DeleteUserRequest{
					Id: userId,
				}
			}

			resp, err := suite.svc.DeleteUser(ctx, tc.input)

			if tc.wantErr {
				suite.Error(err)
				if tc.errContains != "" && err != nil {
					suite.Contains(err.Error(), tc.errContains)
				}
				suite.Nil(resp)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
				if tc.checkResult != nil {
					tc.checkResult(resp)
				}
			}
		})
	}
}
