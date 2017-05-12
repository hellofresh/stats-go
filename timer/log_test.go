package timer

import (
	"io/ioutil"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	hook := test.NewGlobal()

	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)

	b := "foo.bar.bucket"

	tr := &Log{}
	tr.Start().Finish(b)

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, "Muted stats timer send", hook.Entries[0].Message)
	assert.Equal(t, b, hook.Entries[0].Data["bucket"])
}
