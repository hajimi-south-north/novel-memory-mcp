package helpers

import (
	"errors"
	"mcpnovel/internal/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Services struct {
	DB *gorm.DB
}

func (s *Services) CreateNovel(title string, description string) (*models.Novel, error) {
	n := &models.Novel{Title: title, Description: description}
	if err := s.DB.Create(n).Error; err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Services) CreateVolume(novelID uint, title string, index int) (*models.Volume, error) {
	v := &models.Volume{NovelID: novelID, Title: title, Index: index}
	if err := s.DB.Create(v).Error; err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Services) CreateChapter(volumeID uint, title string, index int, status string) (*models.Chapter, error) {
	c := &models.Chapter{VolumeID: volumeID, Title: title, Index: index, Status: status}
	if err := s.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Services) UpsertChapterContent(chapterID uint, content string) error {
	return s.DB.Model(&models.Chapter{}).Where("id = ?", chapterID).Updates(map[string]interface{}{"content": content}).Error
}

func (s *Services) CreateWorld(name string, description string) (*models.World, error) {
	w := &models.World{Name: name, Description: description}
	if err := s.DB.Create(w).Error; err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Services) GetWorldByName(name string) (*models.World, error) {
	var w models.World
	if err := s.DB.Where("name = ?", name).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *Services) EnsureWorld(name string, description string) (*models.World, error) {
	var w models.World
	err := s.DB.Where("name = ?", name).First(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w = models.World{Name: name, Description: description}
		if e := s.DB.Create(&w).Error; e != nil {
			return nil, e
		}
		return &w, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *Services) CreatePeriod(worldID uint, name string, index int) (*models.Period, error) {
	p := &models.Period{WorldID: worldID, Name: name, Index: index}
	if err := s.DB.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Services) GetPeriodByName(worldID uint, name string) (*models.Period, error) {
	var p models.Period
	if err := s.DB.Where("world_id = ? AND name = ?", worldID, name).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Services) EnsurePeriod(worldID uint, name string, index int) (*models.Period, error) {
	var p models.Period
	err := s.DB.Where("world_id = ? AND name = ?", worldID, name).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p = models.Period{WorldID: worldID, Name: name, Index: index}
		if e := s.DB.Create(&p).Error; e != nil {
			return nil, e
		}
		return &p, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Services) CreateTimeSegment(periodID uint, name string, start time.Time, end time.Time) (*models.TimeSegment, error) {
	ts := &models.TimeSegment{PeriodID: periodID, Name: name, Start: start, End: end}
	if err := s.DB.Create(ts).Error; err != nil {
		return nil, err
	}
	return ts, nil
}

func (s *Services) GetTimeSegmentByName(periodID uint, name string) (*models.TimeSegment, error) {
	var ts models.TimeSegment
	if err := s.DB.Where("period_id = ? AND name = ?", periodID, name).First(&ts).Error; err != nil {
		return nil, err
	}
	return &ts, nil
}

func (s *Services) EnsureTimeSegment(periodID uint, name string, start time.Time, end time.Time) (*models.TimeSegment, error) {
	var ts models.TimeSegment
	err := s.DB.Where("period_id = ? AND name = ?", periodID, name).First(&ts).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ts = models.TimeSegment{PeriodID: periodID, Name: name, Start: start, End: end}
		if e := s.DB.Create(&ts).Error; e != nil {
			return nil, e
		}
		return &ts, nil
	}
	if err != nil {
		return nil, err
	}
	return &ts, nil
}

func (s *Services) CreateLocation(worldID uint, name string, description string) (*models.Location, error) {
	l := &models.Location{WorldID: worldID, Name: name, Description: description}
	if err := s.DB.Create(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func (s *Services) GetLocationByName(worldID uint, name string) (*models.Location, error) {
	var l models.Location
	if err := s.DB.Where("world_id = ? AND name = ?", worldID, name).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *Services) EnsureLocation(worldID uint, name string, description string) (*models.Location, error) {
	var l models.Location
	err := s.DB.Where("world_id = ? AND name = ?", worldID, name).First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		l = models.Location{WorldID: worldID, Name: name, Description: description}
		if e := s.DB.Create(&l).Error; e != nil {
			return nil, e
		}
		return &l, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *Services) CreateCharacter(name string, bio string) (*models.Character, error) {
	c := &models.Character{Name: name, Bio: bio}
	if err := s.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Services) GetCharacterByName(name string) (*models.Character, error) {
	var c models.Character
	if err := s.DB.Where("name = ?", name).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Services) EnsureCharacter(name string, bio string) (*models.Character, error) {
	var c models.Character
	err := s.DB.Where("name = ?", name).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c = models.Character{Name: name, Bio: bio}
		if e := s.DB.Create(&c).Error; e != nil {
			return nil, e
		}
		return &c, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Services) SetCharacterRelationship(aid uint, bid uint, rtype string, intimacy float64) (*models.CharacterRelationship, error) {
	var rel models.CharacterRelationship
	err := s.DB.Where("a_id = ? AND b_id = ?", aid, bid).First(&rel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rel = models.CharacterRelationship{AID: aid, BID: bid, Type: rtype, Intimacy: intimacy}
		if e := s.DB.Create(&rel).Error; e != nil {
			return nil, e
		}
	} else if err == nil {
		rel.Type = rtype
		rel.Intimacy = intimacy
		if e := s.DB.Save(&rel).Error; e != nil {
			return nil, e
		}
	} else {
		return nil, err
	}
	var rev models.CharacterRelationship
	e2 := s.DB.Where("a_id = ? AND b_id = ?", bid, aid).First(&rev).Error
	if errors.Is(e2, gorm.ErrRecordNotFound) {
		rev = models.CharacterRelationship{AID: bid, BID: aid, Type: rtype, Intimacy: intimacy}
		_ = s.DB.Create(&rev).Error
	} else if e2 == nil {
		rev.Type = rtype
		rev.Intimacy = intimacy
		_ = s.DB.Save(&rev).Error
	}
	return &rel, nil
}

func (s *Services) CreateEvent(chapterID uint, worldID uint, locationID uint, timeSegmentID uint, description string, characters []uint, items []uint) (*models.Event, error) {
	e := &models.Event{ChapterID: chapterID, WorldID: worldID, LocationID: locationID, TimeSegmentID: timeSegmentID, Description: description, Characters: joinIDs(characters), Items: joinIDs(items)}
	if err := s.DB.Create(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Services) CreateItem(name string, ownerID uint, locationID uint, status string) (*models.Item, error) {
	it := &models.Item{Name: name, OwnerCharacterID: ownerID, LocationID: locationID, Status: status}
	if err := s.DB.Create(it).Error; err != nil {
		return nil, err
	}
	return it, nil
}

func (s *Services) TransferItem(itemID uint, fromID uint, toID uint, eventID uint) (*models.ItemTransfer, error) {
	var it models.Item
	if err := s.DB.First(&it, itemID).Error; err != nil {
		return nil, err
	}
	t := &models.ItemTransfer{ItemID: itemID, FromCharacterID: fromID, ToCharacterID: toID, EventID: eventID}
	if err := s.DB.Create(t).Error; err != nil {
		return nil, err
	}
	it.OwnerCharacterID = toID
	if err := s.DB.Save(&it).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Services) CreateAbility(charID uint, name string, level int) (*models.Ability, error) {
	ab := &models.Ability{CharacterID: charID, Name: name, Level: level}
	if err := s.DB.Create(ab).Error; err != nil {
		return nil, err
	}
	return ab, nil
}

func (s *Services) UpgradeAbility(abilityID uint, level int) (*models.Ability, error) {
	var ab models.Ability
	if err := s.DB.First(&ab, abilityID).Error; err != nil {
		return nil, err
	}
	ab.Level = level
	if err := s.DB.Save(&ab).Error; err != nil {
		return nil, err
	}
	return &ab, nil
}

func (s *Services) UseAbility(abilityID uint, eventID uint, note string) (*models.AbilityUsage, error) {
	u := &models.AbilityUsage{AbilityID: abilityID, EventID: eventID, Note: note}
	if err := s.DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Services) CreatePlotThread(novelID uint, name string, stage string) (*models.PlotThread, error) {
	pt := &models.PlotThread{NovelID: novelID, Name: name, Stage: stage}
	if err := s.DB.Create(pt).Error; err != nil {
		return nil, err
	}
	return pt, nil
}

func (s *Services) GetNovelByTitle(title string) (*models.Novel, error) {
	var n models.Novel
	if err := s.DB.Where("title = ?", title).First(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *Services) EnsureNovel(title string, description string) (*models.Novel, error) {
	var n models.Novel
	err := s.DB.Where("title = ?", title).First(&n).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		n = models.Novel{Title: title, Description: description}
		if e := s.DB.Create(&n).Error; e != nil {
			return nil, e
		}
		return &n, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *Services) GetVolumeByTitle(novelID uint, title string) (*models.Volume, error) {
	var v models.Volume
	if err := s.DB.Where("novel_id = ? AND title = ?", novelID, title).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Services) EnsureVolume(novelID uint, title string, index int) (*models.Volume, error) {
	var v models.Volume
	err := s.DB.Where("novel_id = ? AND title = ?", novelID, title).First(&v).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		v = models.Volume{NovelID: novelID, Title: title, Index: index}
		if e := s.DB.Create(&v).Error; e != nil {
			return nil, e
		}
		return &v, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Services) GetChapterByTitle(volumeID uint, title string) (*models.Chapter, error) {
	var c models.Chapter
	if err := s.DB.Where("volume_id = ? AND title = ?", volumeID, title).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Services) EnsureChapter(volumeID uint, title string, index int, status string) (*models.Chapter, error) {
	var c models.Chapter
	err := s.DB.Where("volume_id = ? AND title = ?", volumeID, title).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c = models.Chapter{VolumeID: volumeID, Title: title, Index: index, Status: status}
		if e := s.DB.Create(&c).Error; e != nil {
			return nil, e
		}
		return &c, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type NovelContext struct {
	Novel      models.Novel
	Volumes    []VolumeContext
	Worlds     []models.World
	Locations  []models.Location
	Characters []models.Character
}

type VolumeContext struct {
	Volume   models.Volume
	Chapters []ChapterContext
}

type ChapterContext struct {
	Chapter models.Chapter
	Events  []models.Event
}

func (s *Services) GetNovelContext(novelID uint) (*NovelContext, error) {
	var n models.Novel
	if err := s.DB.First(&n, novelID).Error; err != nil {
		return nil, err
	}
	var vols []models.Volume
	if err := s.DB.Where("novel_id = ?", novelID).Order("`index` asc").Find(&vols).Error; err != nil {
		return nil, err
	}
	var vctxs []VolumeContext
	for _, v := range vols {
		var chs []models.Chapter
		_ = s.DB.Where("volume_id = ?", v.ID).Order("`index` asc").Find(&chs).Error
		var cctxs []ChapterContext
		for _, c := range chs {
			var evs []models.Event
			_ = s.DB.Where("chapter_id = ?", c.ID).Order("id asc").Find(&evs).Error
			cctxs = append(cctxs, ChapterContext{Chapter: c, Events: evs})
		}
		vctxs = append(vctxs, VolumeContext{Volume: v, Chapters: cctxs})
	}
	var worlds []models.World
	_ = s.DB.Find(&worlds).Error
	var locs []models.Location
	_ = s.DB.Find(&locs).Error
	var chars []models.Character
	_ = s.DB.Find(&chars).Error
	return &NovelContext{Novel: n, Volumes: vctxs, Worlds: worlds, Locations: locs, Characters: chars}, nil
}

func (s *Services) UpdatePlotStage(plotID uint, stage string) (*models.PlotThread, error) {
	var pt models.PlotThread
	if err := s.DB.First(&pt, plotID).Error; err != nil {
		return nil, err
	}
	pt.Stage = stage
	if err := s.DB.Save(&pt).Error; err != nil {
		return nil, err
	}
	return &pt, nil
}

func (s *Services) CreateMemory(characterID uint, eventID uint, content string, trigger string) (*models.Memory, error) {
	m := &models.Memory{CharacterID: characterID, EventID: eventID, Content: content, Trigger: trigger}
	if err := s.DB.Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Services) ExportChapter(chapterID uint) (*models.ExportResult, error) {
	var c models.Chapter
	if err := s.DB.First(&c, chapterID).Error; err != nil {
		return nil, err
	}
	return &models.ExportResult{Content: c.Content}, nil
}

func (s *Services) ExportVolume(volumeID uint) (*models.ExportResult, error) {
	var chs []models.Chapter
	if err := s.DB.Where("volume_id = ?", volumeID).Order("`index` asc").Find(&chs).Error; err != nil {
		return nil, err
	}
	var b strings.Builder
	for _, c := range chs {
		b.WriteString(c.Title)
		b.WriteString("\n")
		b.WriteString(c.Content)
		b.WriteString("\n\n")
	}
	return &models.ExportResult{Content: b.String()}, nil
}

func (s *Services) ExportNovel(novelID uint) (*models.ExportResult, error) {
	var vols []models.Volume
	if err := s.DB.Where("novel_id = ?", novelID).Order("`index` asc").Find(&vols).Error; err != nil {
		return nil, err
	}
	var b strings.Builder
	for _, v := range vols {
		res, err := s.ExportVolume(v.ID)
		if err != nil {
			return nil, err
		}
		b.WriteString(v.Title)
		b.WriteString("\n\n")
		b.WriteString(res.Content)
		b.WriteString("\n\n")
	}
	return &models.ExportResult{Content: b.String()}, nil
}

func joinIDs(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	var parts []string
	for _, id := range ids {
		parts = append(parts, strconv.FormatUint(uint64(id), 10))
	}
	return strings.Join(parts, ",")
}

func (s *Services) SetStyleRef(novelID uint, content string) (*models.StyleRef, error) {
	var sr models.StyleRef
	err := s.DB.Where("novel_id = ?", novelID).First(&sr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		sr = models.StyleRef{NovelID: novelID, Content: content}
		if e := s.DB.Create(&sr).Error; e != nil {
			return nil, e
		}
	} else if err == nil {
		sr.Content = content
		if e := s.DB.Save(&sr).Error; e != nil {
			return nil, e
		}
	} else {
		return nil, err
	}
	return &sr, nil
}

func (s *Services) GetStyleRef(novelID uint) (*models.StyleRef, error) {
	var sr models.StyleRef
	if err := s.DB.Where("novel_id = ?", novelID).First(&sr).Error; err != nil {
		return nil, err
	}
	return &sr, nil
}
