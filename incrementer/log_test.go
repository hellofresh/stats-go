package incrementer

import (
	"io/ioutil"
	"testing"

	"github.com/hellofresh/stats-go/bucket"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	hook := test.NewGlobal()

	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)

	b := "foo.bar.bucket"
	n := 42

	i := Log{}
	i.Increment(b)
	i.IncrementN(b, n)

	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, "Muted stats counter increment", hook.Entries[0].Message)
	assert.Equal(t, b, hook.Entries[0].Data["metric"])
	assert.Equal(t, "Muted stats counter increment by n", hook.Entries[1].Message)
	assert.Equal(t, b, hook.Entries[1].Data["metric"])
	assert.Equal(t, n, hook.Entries[1].Data["n"])

	hook.Reset()

	bb := bucket.NewPlain("section", bucket.MetricOperation{"o1", "o2", "o3"}, true, true)
	i.IncrementAll(bb)

	assert.Equal(t, 4, len(hook.Entries))
	for j := 0; j < 4; j++ {
		assert.Equal(t, "Muted stats counter increment", hook.Entries[j].Message)
	}

	hook.Reset()

	i.IncrementAllN(bb, n)

	assert.Equal(t, 4, len(hook.Entries))
	for j := 0; j < 4; j++ {
		assert.Equal(t, "Muted stats counter increment by n", hook.Entries[j].Message)
		assert.Equal(t, n, hook.Entries[j].Data["n"])
	}
}
