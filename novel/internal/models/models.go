package models

import "time"

type Novel struct {
    ID uint `gorm:"primaryKey"`
    Title string
    Description string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Volume struct {
    ID uint `gorm:"primaryKey"`
    NovelID uint `gorm:"index"`
    Title string
    Index int
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Chapter struct {
    ID uint `gorm:"primaryKey"`
    VolumeID uint `gorm:"index"`
    Title string
    Index int
    Status string
    Content string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type World struct {
    ID uint `gorm:"primaryKey"`
    Name string
    Description string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Period struct {
    ID uint `gorm:"primaryKey"`
    WorldID uint `gorm:"index"`
    Name string
    Index int
    CreatedAt time.Time
    UpdatedAt time.Time
}

type TimeSegment struct {
    ID uint `gorm:"primaryKey"`
    PeriodID uint `gorm:"index"`
    Name string
    Start time.Time
    End time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Location struct {
    ID uint `gorm:"primaryKey"`
    WorldID uint `gorm:"index"`
    Name string
    Description string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Character struct {
    ID uint `gorm:"primaryKey"`
    Name string
    Bio string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type CharacterRelationship struct {
    ID uint `gorm:"primaryKey"`
    AID uint `gorm:"index"`
    BID uint `gorm:"index"`
    Type string
    Intimacy float64
    CreatedAt time.Time
    UpdatedAt time.Time
}

type LocationRelationship struct {
    ID uint `gorm:"primaryKey"`
    AID uint `gorm:"index"`
    BID uint `gorm:"index"`
    Type string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Item struct {
    ID uint `gorm:"primaryKey"`
    Name string
    OwnerCharacterID uint `gorm:"index"`
    LocationID uint `gorm:"index"`
    Status string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ItemTransfer struct {
    ID uint `gorm:"primaryKey"`
    ItemID uint `gorm:"index"`
    FromCharacterID uint
    ToCharacterID uint
    EventID uint
    CreatedAt time.Time
}

type Ability struct {
    ID uint `gorm:"primaryKey"`
    CharacterID uint `gorm:"index"`
    Name string
    Level int
    CreatedAt time.Time
    UpdatedAt time.Time
}

type AbilityUsage struct {
    ID uint `gorm:"primaryKey"`
    AbilityID uint `gorm:"index"`
    EventID uint
    Note string
    CreatedAt time.Time
}

type PlotThread struct {
    ID uint `gorm:"primaryKey"`
    NovelID uint `gorm:"index"`
    Name string
    Stage string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Event struct {
    ID uint `gorm:"primaryKey"`
    ChapterID uint `gorm:"index"`
    WorldID uint `gorm:"index"`
    LocationID uint `gorm:"index"`
    TimeSegmentID uint `gorm:"index"`
    Description string
    Characters string
    Items string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Memory struct {
    ID uint `gorm:"primaryKey"`
    CharacterID uint `gorm:"index"`
    EventID uint `gorm:"index"`
    Content string
    Trigger string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ExportResult struct {
    Content string
}

type Conflict struct {
    Type string
    Detail string
}

type StyleRef struct {
    ID uint `gorm:"primaryKey"`
    NovelID uint `gorm:"index"`
    Content string
    CreatedAt time.Time
    UpdatedAt time.Time
}