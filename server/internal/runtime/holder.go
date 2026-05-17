// Package runtime owns the live, hot-swappable view of master data.
//
// The Holder atomically swaps a *Catalogs aggregate every time the operator
// asks the server to re-read assets/release/20240404193219.bin.e (typically via
// the admin webhook in cmd/lunar-tear/admin.go). gRPC services hold a *Holder
// and call Get() at the start of each RPC, so they always see a consistent
// snapshot.
package runtime

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"lunar-tear/server/internal/gacha"
	"lunar-tear/server/internal/masterdata"
	"lunar-tear/server/internal/masterdata/memorydb"
	"lunar-tear/server/internal/model"
	"lunar-tear/server/internal/questflow"
	"lunar-tear/server/internal/store"
)

// Catalogs is an immutable snapshot of every catalog and catalog-derived
// handler the server needs at runtime. A new *Catalogs is built from scratch
// on every reload and atomically published via Holder.
type Catalogs struct {
	GameConfig        *masterdata.GameConfig
	Parts             *masterdata.PartsCatalog
	Quest             *masterdata.QuestCatalog
	GachaEntries      []store.GachaCatalogEntry
	GachaMedals       map[int32]masterdata.GachaMedalInfo
	GachaPool         *masterdata.GachaCatalog
	Shop              *masterdata.ShopCatalog
	DupExchange       map[int32][]model.DupExchangeEntry
	ConditionResolver *masterdata.ConditionResolver
	CageOrnament      *masterdata.CageOrnamentCatalog
	LoginBonus        *masterdata.LoginBonusCatalog
	CharacterViewer   *masterdata.CharacterViewerCatalog
	Omikuji           *masterdata.OmikujiCatalog
	Material          *masterdata.MaterialCatalog
	ConsumableItem    *masterdata.ConsumableItemCatalog
	Costume           *masterdata.CostumeCatalog
	Weapon            *masterdata.WeaponCatalog
	Explore           *masterdata.ExploreCatalog
	Gimmick           *masterdata.GimmickCatalog
	CharacterBoard    *masterdata.CharacterBoardCatalog
	CharacterRebirth  *masterdata.CharacterRebirthCatalog
	Companion         *masterdata.CompanionCatalog
	SideStory         *masterdata.SideStoryCatalog
	BigHunt           *masterdata.BigHuntCatalog
	Tower             *masterdata.TowerCatalog
	Labyrinth         *masterdata.LabyrinthCatalog

	// Catalog-derived handlers must rebuild on every reload because they
	// embed/cache pointers to specific catalog instances.
	QuestHandler *questflow.QuestHandler
	GachaHandler *gacha.GachaHandler
}

// Holder owns the current *Catalogs and the bin.e path. Concurrent readers
// call Get(); the single-writer Reload() rebuilds and atomically publishes.
type Holder struct {
	binPath string
	cur     atomic.Pointer[Catalogs]
}

// NewHolder reads the binary at binPath, builds the initial catalogs, and
// returns a ready-to-use Holder. Subsequent calls to Reload() re-read the
// same path.
func NewHolder(binPath string) (*Holder, error) {
	h := &Holder{binPath: binPath}
	if err := h.Reload(); err != nil {
		return nil, err
	}
	return h, nil
}

// Reload re-reads the bin.e from disk, rebuilds every catalog and handler,
// atomically publishes the new snapshot, and bumps the bin.e mtime so client
// caches invalidate (see service/data.go GetLatestMasterDataVersion).
func (h *Holder) Reload() error {
	if err := memorydb.Init(h.binPath); err != nil {
		return fmt.Errorf("memorydb.Init: %w", err)
	}
	c, err := buildCatalogs()
	if err != nil {
		return fmt.Errorf("buildCatalogs: %w", err)
	}
	h.cur.Store(c)
	now := time.Now()
	if err := os.Chtimes(h.binPath, now, now); err != nil {
		// Non-fatal: the catalogs swapped fine in-memory; clients may take
		// longer to invalidate their cached download but server-side state is
		// already coherent.
		log.Printf("[runtime] os.Chtimes(%s) failed (clients may not invalidate cache): %v", h.binPath, err)
	}
	return nil
}

// Get returns the current snapshot. Safe for concurrent callers; the returned
// pointer is stable for the duration of the caller's use.
func (h *Holder) Get() *Catalogs {
	return h.cur.Load()
}
