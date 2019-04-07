package good

import (
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/zartbot/golap/api/datastream"
)

type ModelFuncType func(*datastream.DataStream, *ModelContainer, *tf.Session)

func modelCallBackWorker(m *ModelContainer, wg *sync.WaitGroup) {
	session, err := tf.NewSession(m.graph, nil)
	if err != nil {
		logrus.Warn("Create session failed")
		wg.Done()
	}
LOOP:
	for {
		data := <-m.Input
		if m.CopyMode {
			dataOutput := data.Copy()
			if data.ID >= datastream.MaxReservedID {
				m.ModelFunc(dataOutput, m, session)
			}
			m.Output <- dataOutput
		} else {
			if data.ID >= datastream.MaxReservedID {
				m.ModelFunc(data, m, session)
			}
			m.Output <- data
		}

		if atomic.LoadInt32(&m.StopSignal) == 1 {
			logrus.Warning("Stop Signal Recieved.")
			session.Close()
			break LOOP
		}
	}
	wg.Done()
}

type GraphOperation struct {
	inputOp     *tf.Operation
	outputBox   *tf.Operation
	outputScore *tf.Operation
	outputLabel *tf.Operation
}

type ModelContainer struct {
	ID          uint32
	Name        string
	ModelPath   string
	LabelPath   string
	ModelFunc   ModelFuncType
	model       []byte
	graph       *tf.Graph
	graphOp     *GraphOperation
	Session     *tf.Session
	Label       []string
	Input       <-chan *datastream.DataStream
	Output      chan<- *datastream.DataStream
	Parallelism int
	CopyMode    bool
	StopSignal  int32
}

type ModelCfg struct {
	ID          uint32
	Name        string
	ModelPath   string
	LabelPath   string
	ModelFunc   ModelFuncType
	Input       <-chan *datastream.DataStream
	Output      chan<- *datastream.DataStream
	Parallelism int
	CopyMode    bool
}

//New is used to create and load model
func New(cfg *ModelCfg) *ModelContainer {
	m := &ModelContainer{
		ID:          cfg.ID,
		Name:        cfg.Name,
		ModelPath:   cfg.ModelPath,
		LabelPath:   cfg.LabelPath,
		ModelFunc:   cfg.ModelFunc,
		Input:       cfg.Input,
		Output:      cfg.Output,
		Parallelism: cfg.Parallelism,
		CopyMode:    cfg.CopyMode,
	}

	if cfg.ModelFunc == nil {
		m.ModelFunc = DefaultModelFunc
	}

	var err error

	//Load Label
	if m.LabelPath == "COCO" {
		m.Label, _ = LoadCOCOLabel()
	} else {
		m.Label, err = LoadLabel(m.LabelPath)
		if err != nil {
			logrus.Fatal("Load Label failed:", err)
		}
	}

	//Load Model
	m.model, err = ioutil.ReadFile(m.ModelPath)
	if err != nil {
		logrus.Fatal("Load model failed:", err)
	}

	//Load Graph
	m.graph = tf.NewGraph()
	err = m.graph.Import(m.model, "")
	if err != nil {
		logrus.Fatal("Load graph failed:", err)
	}

	m.graphOp = &GraphOperation{
		inputOp:     m.graph.Operation("image_tensor"),
		outputBox:   m.graph.Operation("detection_boxes"),
		outputScore: m.graph.Operation("detection_scores"),
		outputLabel: m.graph.Operation("detection_classes"),
	}

	return m
}

func (m *ModelContainer) Run() {
	var wg sync.WaitGroup
	wg.Add(m.Parallelism)
	logrus.WithFields(logrus.Fields{
		"Parallelism": m.Parallelism,
	}).Info(fmt.Sprintf("[model-ObjectDetection]%s Started..", m.Name))
	for id := 0; id < m.Parallelism; id++ {
		go modelCallBackWorker(m, &wg)
	}
	wg.Wait()
	logrus.WithFields(logrus.Fields{
		"Parallelism": m.Parallelism,
	}).Info(fmt.Sprintf("[model-ObjectDetection]%s Stoped..", m.Name))
}

//Stop is used stop the polling goroutines
func (m *ModelContainer) Stop() {
	atomic.StoreInt32(&m.StopSignal, 1)
}
