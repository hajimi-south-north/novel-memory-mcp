package conflict

import (
    "errors"
    "fmt"
    "time"
    "gorm.io/gorm"
    "mcpnovel/internal/models"
)

type Detector struct {
    DB *gorm.DB
}

func (d *Detector) DetectAll() ([]models.Conflict, error) {
    var out []models.Conflict
    c1, _ := d.TimeOrderConflicts()
    out = append(out, c1...)
    c2, _ := d.EventPresenceConflicts()
    out = append(out, c2...)
    c3, _ := d.CharacterStateConflicts()
    out = append(out, c3...)
    c4, _ := d.LocationStateConflicts()
    out = append(out, c4...)
    c5, _ := d.ReferenceIntegrityConflicts()
    out = append(out, c5...)
    c6, _ := d.RelationshipLogicConflicts()
    out = append(out, c6...)
    c7, _ := d.StatusConsistencyConflicts()
    out = append(out, c7...)
    c8, _ := d.ItemAbilityConflicts()
    out = append(out, c8...)
    c9, _ := d.PlotThreadConflicts()
    out = append(out, c9...)
    c10, _ := d.CharacterLocationConflicts()
    out = append(out, c10...)
    return out, nil
}

func (d *Detector) TimeOrderConflicts() ([]models.Conflict, error) {
    var segs []models.TimeSegment
    if err := d.DB.Order("period_id asc, start asc, end asc").Find(&segs).Error; err != nil {
        return nil, err
    }
    var out []models.Conflict
    for i := 1; i < len(segs); i++ {
        prev := segs[i-1]
        cur := segs[i]
        if prev.PeriodID == cur.PeriodID && cur.Start.Before(prev.End) {
            out = append(out, models.Conflict{Type: "时间冲突", Detail: fmt.Sprintf("时间段重叠 %d-%d", prev.ID, cur.ID)})
        }
    }
    var evs []models.Event
    if err := d.DB.Find(&evs).Error; err != nil {
        return nil, err
    }
    for _, e := range evs {
        var ts models.TimeSegment
        if err := d.DB.First(&ts, e.TimeSegmentID).Error; err == nil {
            if ts.End.Before(ts.Start) {
                out = append(out, models.Conflict{Type: "时间冲突", Detail: fmt.Sprintf("时间段无效 %d", ts.ID)})
            }
        }
    }
    return out, nil
}

func (d *Detector) EventPresenceConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var evs []models.Event
    if err := d.DB.Find(&evs).Error; err != nil {
        return nil, err
    }
    for _, e := range evs {
        if e.LocationID == 0 {
            out = append(out, models.Conflict{Type: "事件冲突", Detail: fmt.Sprintf("事件地点缺失 %d", e.ID)})
        }
        if e.WorldID == 0 {
            out = append(out, models.Conflict{Type: "事件冲突", Detail: fmt.Sprintf("事件世界缺失 %d", e.ID)})
        }
    }
    return out, nil
}

func (d *Detector) CharacterStateConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var chars []models.Character
    if err := d.DB.Find(&chars).Error; err != nil {
        return nil, err
    }
    for _, c := range chars {
        if c.Name == "" {
            out = append(out, models.Conflict{Type: "人物冲突", Detail: fmt.Sprintf("人物名称缺失 %d", c.ID)})
        }
    }
    return out, nil
}

func (d *Detector) LocationStateConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var locs []models.Location
    if err := d.DB.Find(&locs).Error; err != nil {
        return nil, err
    }
    for _, l := range locs {
        if l.WorldID == 0 {
            out = append(out, models.Conflict{Type: "地点冲突", Detail: fmt.Sprintf("地点未关联世界 %d", l.ID)})
        }
    }
    return out, nil
}

func (d *Detector) ReferenceIntegrityConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var evs []models.Event
    if err := d.DB.Find(&evs).Error; err != nil {
        return nil, err
    }
    for _, e := range evs {
        var ch models.Chapter
        if err := d.DB.First(&ch, e.ChapterID).Error; err != nil {
            out = append(out, models.Conflict{Type: "引用完整性", Detail: fmt.Sprintf("事件章节不存在 %d", e.ID)})
        }
    }
    return out, nil
}

func (d *Detector) RelationshipLogicConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var rels []models.CharacterRelationship
    if err := d.DB.Find(&rels).Error; err != nil {
        return nil, err
    }
    for _, r := range rels {
        if r.AID == r.BID {
            out = append(out, models.Conflict{Type: "关系逻辑", Detail: fmt.Sprintf("自我关系 %d", r.ID)})
        }
    }
    return out, nil
}

func (d *Detector) StatusConsistencyConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var chs []models.Chapter
    if err := d.DB.Find(&chs).Error; err != nil {
        return nil, err
    }
    for _, c := range chs {
        if c.Status != "开始" && c.Status != "进行中" && c.Status != "结束" && c.Status != "草稿" && c.Status != "完成" {
            out = append(out, models.Conflict{Type: "状态一致性", Detail: fmt.Sprintf("章节状态非法 %d", c.ID)})
        }
    }
    return out, nil
}

func (d *Detector) ItemAbilityConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var transfers []models.ItemTransfer
    if err := d.DB.Find(&transfers).Error; err != nil {
        return nil, err
    }
    for _, t := range transfers {
        var it models.Item
        if err := d.DB.First(&it, t.ItemID).Error; err != nil {
            out = append(out, models.Conflict{Type: "物品能力冲突", Detail: fmt.Sprintf("物品不存在 %d", t.ItemID)})
        }
    }
    var usages []models.AbilityUsage
    if err := d.DB.Find(&usages).Error; err != nil {
        return nil, err
    }
    for _, u := range usages {
        var ab models.Ability
        if err := d.DB.First(&ab, u.AbilityID).Error; err != nil {
            out = append(out, models.Conflict{Type: "物品能力冲突", Detail: fmt.Sprintf("能力不存在 %d", u.AbilityID)})
        }
    }
    return out, nil
}

func (d *Detector) PlotThreadConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var pts []models.PlotThread
    if err := d.DB.Find(&pts).Error; err != nil {
        return nil, err
    }
    for _, p := range pts {
        if p.Stage == "" {
            out = append(out, models.Conflict{Type: "线索冲突", Detail: fmt.Sprintf("线索阶段缺失 %d", p.ID)})
        }
    }
    return out, nil
}

func (d *Detector) CharacterLocationConflicts() ([]models.Conflict, error) {
    var out []models.Conflict
    var evs []models.Event
    if err := d.DB.Find(&evs).Error; err != nil {
        return nil, err
    }
    for _, e := range evs {
        var loc models.Location
        var w models.World
        le := d.DB.First(&loc, e.LocationID).Error
        we := d.DB.First(&w, e.WorldID).Error
        if errors.Is(le, gorm.ErrRecordNotFound) || errors.Is(we, gorm.ErrRecordNotFound) {
            out = append(out, models.Conflict{Type: "人物地点关系冲突", Detail: fmt.Sprintf("事件引用缺失 %d", e.ID)})
        }
    }
    return out, nil
}

func ValidTimeRange(start time.Time, end time.Time) bool {
    return !end.Before(start)
}