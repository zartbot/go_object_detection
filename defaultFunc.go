package good

import (
	"time"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/zartbot/go_utils/map2value"
	"github.com/zartbot/golap/api/datastream"
)

func DefaultModelFunc(d *datastream.DataStream, m *ModelContainer, sess *tf.Session) {

	startTime := time.Now()
	imgBytes, err := map2value.MapToBytes(d.RecordMap, "image")
	if err != nil {
		d.RecordMap["State"] = "Failed decode image"
		return
	}

	objlist, err := m.Prediction(imgBytes, sess)
	if err != nil {
		d.RecordMap["State"] = "Failed during detection"
		return
	}

	ObjectCount := make(map[string]uint32)
	for _, v := range objlist {

		item, ok := ObjectCount[v.LabelStr]
		if !ok {
			ObjectCount[v.LabelStr] = uint32(1)
		} else {
			ObjectCount[v.LabelStr] = item + 1
		}
	}
	d.RecordMap["detect_object"] = objlist
	d.RecordMap["objectCount"] = ObjectCount

	d.RecordMap["ElapsedTime_Prediction"] = time.Since(startTime)

	//if need render image

	/*
		newimg, err := RenderObject(imgBytes, objlist)
		if err != nil {
			d.RecordMap["State"] = "Failed during render"
			return
		}

			f, err := os.Create("/home/kevin/Desktop/d.jpg")
			if err != nil {
				fmt.Println(err)
				return
			}
			n2, err := f.Write(newimg)
			if err != nil {
				fmt.Println(err)
				f.Close()
				return
			}
			fmt.Println(n2, "bytes written successfully")
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
	*/

	d.RecordMap["State"] = "success"
	delete(d.RecordMap, "image")
	d.RecordMap["ElapsedTime"] = time.Since(startTime)
}
