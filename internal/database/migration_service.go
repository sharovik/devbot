package database

import (
	"github.com/sharovik/devbot/internal/log"
)

// BaseMigrationInterface the interface for all migration files
type BaseMigrationInterface interface {
	GetName() string
	Execute() error
}

// MigrationService the service for migrations based on GoLang run.
// This can be useful if you want to use abstraction of our SQL client in your migration
type MigrationService struct {
	Logger     log.LoggerInstance
	Dictionary BaseDatabaseInterface
	migrations []BaseMigrationInterface
}

// SetMigration method loads the migration into memory. Use this method to prepare your migration for execution
func (s *MigrationService) SetMigration(newMigration BaseMigrationInterface) {
	if s.migrations == nil {
		s.migrations = []BaseMigrationInterface{}
	}

	for _, m := range s.migrations {
		if m.GetName() == newMigration.GetName() {
			return
		}
	}

	s.migrations = append(s.migrations, newMigration)
}

func (s *MigrationService) finish() {
	s.migrations = nil
}

// RunMigrations method will run all migrations, which are set to migrations variable
func (s MigrationService) RunMigrations() error {
	for _, migration := range s.migrations {
		isMigrationAlreadyExecuted, err := s.Dictionary.IsMigrationAlreadyExecuted(migration.GetName())
		if err != nil {
			return err
		}

		if isMigrationAlreadyExecuted {
			continue
		}

		s.Logger.Info().Str("migration_name", migration.GetName()).Msg("Running migration")
		if err := migration.Execute(); err != nil {
			log.Logger().AddError(err).Str("migration_name", migration.GetName()).Msg("Failed to execute the migration")
			continue
		}

		if err := s.Dictionary.MarkMigrationExecuted(migration.GetName()); err != nil {
			return err
		}
	}

	s.finish()
	return nil
}
