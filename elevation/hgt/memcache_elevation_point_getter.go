package hgt

import (
	"fmt"

	"github.com/thanosKontos/gravelmap"
)

type memcacheElevationStorage struct {
	elevationGettersClosersCache map[string]gravelmap.ElevationPointGetterCloser
	logger                       gravelmap.Logger
	fileGetter                   hgtFileGetter
}

// NewMemcacheNasaElevationFileStorage instanciates a new in memory object holding ElevationPointGetterCloser objects
func NewMemcacheNasaElevationFileStorage(destinationDir, username, password string, logger gravelmap.Logger) *memcacheElevationStorage {
	fileGetter := &nasa30mFile{username, password, destinationDir}

	return &memcacheElevationStorage{
		elevationGettersClosersCache: make(map[string]gravelmap.ElevationPointGetterCloser),
		logger:                       logger,
		fileGetter:                   fileGetter,
	}
}

func (m *memcacheElevationStorage) Get(dms string) (gravelmap.ElevationPointGetter, error) {
	if g, ok := m.elevationGettersClosersCache[dms]; ok {
		return g, nil
	}

	m.logger.Info(fmt.Sprintf("Getting file: %s", dms))
	f, err := m.fileGetter.getFile(dms)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}
	m.logger.Info("Done")

	g := NewStrm1(f)
	m.elevationGettersClosersCache[dms] = g

	return g, nil
}

func (m *memcacheElevationStorage) Close() {
	for _, egc := range m.elevationGettersClosersCache {
		egc.Close()
	}
}
