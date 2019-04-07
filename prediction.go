package good

import (
	"fmt"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type Object struct {
	Prob     float32
	Label    int
	LabelStr string
	Box      []float32
}

func (m *ModelContainer) Prediction(inputBytes []byte, session *tf.Session) ([]*Object, error) {
	//make tensor
	tensor, err := makeTensorFromImage(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("Convert to tensor error: %v", err)
	}

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			m.graphOp.inputOp.Output(0): tensor,
		},
		[]tf.Output{
			m.graphOp.outputBox.Output(0),
			m.graphOp.outputScore.Output(0),
			m.graphOp.outputLabel.Output(0),
		}, nil)

	if err != nil {
		return nil, fmt.Errorf("Session execution error: %v", err)
	}

	//summary output
	box := output[0].Value().([][][]float32)[0]
	score := output[1].Value().([][]float32)[0]
	lables := output[2].Value().([][]float32)[0]

	objList := make([]*Object, 0, 1)

	for k, v := range score {
		//filter threshold
		if v < 0.35 {
			break
		}

		obj := &Object{
			Prob:  v,
			Label: int(lables[k]),
			Box:   box[k],
		}

		obj.LabelStr = m.Label[obj.Label]
		objList = append(objList, obj)
	}
	return objList, nil
}
