package migration1611603404

import (
	"github.com/jinzhu/gorm"
)

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE upkeep_registrations (
			address bytea NOT NULL,
			upkeep_id bigint NOT NULL,
			check_gas_limit int NOT NULL,
			last_run_block_height bigInt DEFAULT 0 NOT NULL
		);

		CREATE UNIQUE INDEX idx_upkeep_registrations_unique_upkeep_ids_per_address ON upkeep_registrations(address, upkeep_id);
	`).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(`
		DROP TABLE IF EXISTS upkeep_registrations;
	`).Error
}