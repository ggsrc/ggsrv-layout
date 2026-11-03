package users

import (
	"context"
	"time"
)

// UsersTableSerde implements the Loader and Dumper interfaces for the users table.
type UsersTableSerde struct {
	Users *Queries
}

func (u UsersTableSerde) Load(data []byte) error {
	// Use the auto-generated Load function.
	err := u.Users.Load(context.Background(), data)
	if err != nil {
		return err
	}

	// After loading data with explicit IDs, we must reset the identity sequence
	// to prevent conflicts with future auto-generated IDs.
	// Note: users table doesn't have a RefreshID method, so we skip this step
	return nil
}

func (u UsersTableSerde) Dump() ([]byte, error) {
	// Use the auto-generated Dump function, but provide a callback to
	// normalize any "flaky" data, like timestamps, before serialization.
	// This ensures the golden files are stable.
	return u.Users.Dump(context.Background(), func(m *User) {
		m.CreatedAt = time.Unix(0, 0).UTC()
		m.UpdatedAt = time.Unix(0, 0).UTC()
	})
}
