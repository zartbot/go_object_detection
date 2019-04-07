package good

import (
	"bufio"
	"os"
)

func LoadLabel(lablefilepath string) ([]string, error) {
	var labelList = make([]string, 0, 1)
	file, err := os.Open(lablefilepath)
	if err != nil {
		return labelList, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		labelList = append(labelList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return labelList, err
	}
	return labelList, nil
}

func LoadCOCOLabel() ([]string, error) {
	var labels = []string{
		"unlabeled",
		"person",
		"bicycle",
		"car",
		"motorcycle",
		"airplane",
		"bus",
		"train",
		"truck",
		"boat",
		"traffic light",
		"fire hydrant",
		"street sign",
		"stop sign",
		"parking meter",
		"bench",
		"bird",
		"cat",
		"dog",
		"horse",
		"sheep",
		"cow",
		"elephant",
		"bear",
		"zebra",
		"giraffe",
		"hat",
		"backpack",
		"umbrella",
		"shoe",
		"eye glasses",
		"handbag",
		"tie",
		"suitcase",
		"frisbee",
		"skis",
		"snowboard",
		"sports ball",
		"kite",
		"baseball bat",
		"baseball glove",
		"skateboard",
		"surfboard",
		"tennis racket",
		"bottle",
		"plate",
		"wine glass",
		"cup",
		"fork",
		"knife",
		"spoon",
		"bowl",
		"banana",
		"apple",
		"sandwich",
		"orange",
		"broccoli",
		"carrot",
		"hot dog",
		"pizza",
		"donut",
		"cake",
		"chair",
		"couch",
		"potted plant",
		"bed",
		"mirror",
		"dining table",
		"window",
		"desk",
		"toilet",
		"door",
		"tv",
		"laptop",
		"mouse",
		"remote",
		"keyboard",
		"cell phone",
		"microwave",
		"oven",
		"toaster",
		"sink",
		"refrigerator",
		"blender",
		"book",
		"clock",
		"vase",
		"scissors",
		"teddy bear",
		"hair drier",
		"toothbrush",
		"hair brush",
		"banner",
		"blanket",
		"branch",
		"bridge",
		"building-other",
		"bush",
		"cabinet",
		"cage",
		"cardboard",
		"carpet",
		"ceiling-other",
		"ceiling-tile",
		"cloth",
		"clothes",
		"clouds",
		"counter",
		"cupboard",
		"curtain",
		"desk-stuff",
		"dirt",
		"door-stuff",
		"fence",
		"floor-marble",
		"floor-other",
	}
	return labels, nil
}
