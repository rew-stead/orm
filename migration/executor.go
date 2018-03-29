package migration

import (
	"bytes"
	"fmt"
	"time"
)

// Executor provides a group of operations that works with migrations.
type Executor struct {
	// Provider provides all migrations for the current project.
	Provider ItemProvider
	// Runner runs or reverts migrations for the current project.
	Runner ItemRunner
	// Generator generates a migration file.
	Generator FileGenerator
	// OnRunFn is executed when a migration is executed.
	OnRunFn ItemFn
	// OnRevertFn is executed when a migration is reverted.
	OnRevertFn ItemFn
}

// Setups setups the current project for database migrations by creating
// migration directory and related database.
func (m *Executor) Setup() error {
	migration := &Item{
		Id:          min.Format(format),
		Description: "setup",
		CreatedAt:   time.Now(),
	}

	up := &bytes.Buffer{}
	fmt.Fprintln(up, "CREATE TABLE migrations (")
	fmt.Fprintln(up, " id          TEXT      NOT NULL PRIMARY KEY,")
	fmt.Fprintln(up, " description TEXT      NOT NULL,")
	fmt.Fprintln(up, " created_at  TIMESTAMP NOT NULL")
	fmt.Fprintln(up, ");")

	down := bytes.NewBufferString("DROP TABLE IF EXISTS migrations;")

	content := &Content{
		UpCommand:   up,
		DownCommand: down,
	}

	if err := m.Generator.Write(migration, content); err != nil {
		return err
	}

	return m.Runner.Run(migration)
}

// Create creates a migration script successfully if the project has already
// been setup, otherwise returns an error.
func (m *Executor) Create(name string) (string, error) {
	timestamp := time.Now()

	migration := &Item{
		Id:          timestamp.Format(format),
		Description: name,
		CreatedAt:   timestamp,
	}

	return m.Generator.Create(migration)
}

// Run runs a pending migration for given count. If the count is negative number, it
// will execute all pending migrations.
func (m *Executor) Run(step int) error {
	migrations, err := m.Migrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if step == 0 {
			return nil
		}

		timestamp, err := time.Parse(format, migration.Id)
		if err != nil {
			return err
		}

		if !migration.CreatedAt.IsZero() || timestamp == min {
			continue
		}

		op := migration

		if m.OnRunFn != nil {
			m.OnRunFn(&op)
		}

		if err := m.Runner.Run(&op); err != nil {
			return err
		}

		step = step - 1
	}

	return nil
}

// RunAll runs all pending migrations.
func (m *Executor) RunAll() error {
	return m.Run(-1)
}

// Revert reverts an applied migration for given count. If the count is
// negative number, it will revert all applied migrations.
func (m *Executor) Revert(step int) error {
	migrations, err := m.Migrations()
	if err != nil {
		return err
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]

		if step == 0 {
			return nil
		}

		if migration.CreatedAt.IsZero() {
			continue
		}

		timestamp, err := time.Parse(format, migration.Id)
		if err != nil || timestamp == min {
			return err
		}

		op := migration

		if m.OnRevertFn != nil {
			m.OnRevertFn(&op)
		}

		if err := m.Runner.Revert(&op); err != nil {
			return err
		}

		step = step - 1
	}

	return nil
}

// RevertAll reverts all applied migrations.
func (m *Executor) RevertAll() error {
	return m.Revert(-1)
}

// Migrations returns all migrations.
func (m *Executor) Migrations() ([]Item, error) {
	return m.Provider.Migrations()
}
