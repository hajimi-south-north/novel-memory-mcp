package outline

import (
	"mcpnovel/internal/models"
	"strings"

	"gorm.io/gorm"
)

type Generator struct {
	DB *gorm.DB
}

func (g *Generator) ChapterOutline(chapterID uint) (string, error) {
	var e []models.Event
	if err := g.DB.Where("chapter_id = ?", chapterID).Order("id asc").Find(&e).Error; err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("章节细纲\n")
	for i, ev := range e {
		b.WriteString("事件")
		b.WriteString(intToString(i + 1))
		b.WriteString(": ")
		b.WriteString(ev.Description)
		b.WriteString("\n")
	}
	return b.String(), nil
}

func (g *Generator) VolumeOutline(volumeID uint) (string, error) {
	var chs []models.Chapter
	if err := g.DB.Where("volume_id = ?", volumeID).Order("`index` asc").Find(&chs).Error; err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("分卷总纲\n")
	for _, c := range chs {
		co, err := g.ChapterOutline(c.ID)
		if err != nil {
			return "", err
		}
		b.WriteString(c.Title)
		b.WriteString("\n")
		b.WriteString(co)
		b.WriteString("\n")
	}
	return b.String(), nil
}

func (g *Generator) NovelOutline(novelID uint) (string, error) {
	var vols []models.Volume
	if err := g.DB.Where("novel_id = ?", novelID).Order("`index` asc").Find(&vols).Error; err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("小说总纲\n")
	for _, v := range vols {
		vo, err := g.VolumeOutline(v.ID)
		if err != nil {
			return "", err
		}
		b.WriteString(v.Title)
		b.WriteString("\n")
		b.WriteString(vo)
		b.WriteString("\n")
	}
	return b.String(), nil
}

func intToString(i int) string {
	s := ""
	if i == 0 {
		return "0"
	}
	for i > 0 {
		d := i % 10
		s = string('0'+d) + s
		i = i / 10
	}
	return s
}
