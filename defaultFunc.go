package good

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/zartbot/go_utils/map2value"
	"github.com/zartbot/golap/api/datastream"
)

func DefaultModelFunc(d *datastream.DataStream, m *ModelContainer, sess *tf.Session) {
	imgBytes, err := map2value.MapToBytes(d.RecordMap, "image")
	if err != nil {
		return
	}

	obj, err := m.Prediction(imgBytes, sess)
	if err != nil {
		return
	}

}
