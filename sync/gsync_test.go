package sync

import (
	"testing"
)

func TestPutfile(t *testing.T) {
}

func TestListfile(t *testing.T) {
}

func TestSyncDir(t *testing.T) {
	g := NewGsync("","","", []string{`/private/tmp/.*/417695_100921\.txt`, `Sourcetree\.app/.*`})
	b := g.isIgnorePath("/private/tmp/2018-04-13/417695_100921.txt")
	if !b {
		t.Error("匹配失败")
	}
	b = g.isIgnorePath("/private/tmp/Sourcetree.app/Contents/Frameworks/Sparkle.framework/Versions/A/Resources/ru.lproj/SUUpdateAlert.nib")
	if !b {
		t.Error("匹配失败")
	}


}
