package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type storageImpl struct {
	db           *gorm.DB
	transactions *transactionStore
	categories   *categoryStore
	rules        *ruleStore
	accounts     *accountStore
	requsitions  *requisitionStore
}

func New(filename string) (Storage, error) {
	shouldSeed := false
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		shouldSeed = true
	}

	// setup file logging so it doesn't mess with the UI
	// TODO: make this optional to better support non tui use.
	logFile, err := os.Create("gorm.log")
	if err != nil {
		return nil, err
	}

	fileLogger := log.New(logFile, "\r\n", log.LstdFlags)
	gormLogger := logger.New(fileLogger, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	db.AutoMigrate(Transaction{})
	db.AutoMigrate(Category{})
	db.AutoMigrate(Rule{})
	db.AutoMigrate(Account{})
	db.AutoMigrate(Requisition{})

	store := &storageImpl{db: db}

	// _ = shouldSeed
	if shouldSeed {
		store.populate()
	}

	return store, nil
}

func (s *storageImpl) DB() *gorm.DB {
	return s.db
}

func (s *storageImpl) Transactions() TransactionGetter {
	if s.transactions == nil {
		s.transactions = &transactionStore{
			db: s.db,
		}
	}

	return s.transactions
}

func (s *storageImpl) Categories() CategoryGetter {
	if s.categories == nil {
		s.categories = &categoryStore{
			db: s.db,
		}
	}

	return s.categories
}

func (s *storageImpl) Rules() RuleGetter {
	if s.rules == nil {
		s.rules = &ruleStore{
			db: s.db,
		}
	}

	return s.rules
}

func (s *storageImpl) Accounts() AccountGetter {
	if s.accounts == nil {
		s.accounts = &accountStore{
			db: s.db,
		}
	}

	return s.accounts
}

func (s *storageImpl) Requisitions() RequsitionGetter {
	if s.requsitions == nil {
		s.requsitions = &requisitionStore{
			db: s.db,
		}
	}

	return s.requsitions
}

func (s *storageImpl) populate() {
	categoryToId := map[string]uint{}

	categoryColors, _ := colorful.WarmPalette(len(categoryList))

	for index, category := range categoryList {
		model := &Category{Name: category, Color: categoryColors[index].Hex()}
		res := s.db.Save(model)
		if res.Error != nil {
			fmt.Println(res.Error)
		}
		categoryToId[category] = model.ID
	}

	for pattern, category := range rulesMap {
		res := s.db.Save(&Rule{
			Pattern:    pattern,
			CategoryID: categoryToId[category],
		})

		if res.Error != nil {
			fmt.Println(res.Error)
		}
	}
}
