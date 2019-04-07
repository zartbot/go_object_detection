package good

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func HLine(image *image.RGBA, y, x1, x2 int, c color.Color) {
	for i := x1; i < x2; i++ {
		image.Set(i, y, c)
	}
}

func VLine(image *image.RGBA, x, y1, y2 int, c color.Color) {
	for i := y1; i < y2; i++ {
		image.Set(x, i, c)
	}
}

func Rectangle(image *image.RGBA, x1, x2, y1, y2, width int, c color.Color) {
	for w := 0; w < width; w++ {
		HLine(image, y1+w, x1, x2, c)
		HLine(image, y2+w, x1, x2, c)
		VLine(image, x1+w, y1, y2, c)
		VLine(image, x2+w, y1, y2, c)
	}
}

func renderLabel(img *image.RGBA, x, y, label int, renderstr string) {
	col := colornames.Map[colornames.Names[label]]
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(colornames.Black),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	Rectangle(img, x, (x + len(renderstr)*7), y-13, y-6, 7, col)
	d.DrawString(renderstr)
}

func (m *ModelContainer) RenderObject(inputBytes []byte, objectList []*Object) ([]byte, error) {
	var output bytes.Buffer
	outputWriter := bufio.NewWriter(&output)

	//decode image from bytes
	img, err := DecodeImage(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("Decode image error: %v", err)
	}
	bounds := img.Bounds()
	imgRGBA := image.NewRGBA(bounds)
	draw.Draw(imgRGBA, bounds, img, bounds.Min, draw.Src)

	//render
	for _, item := range objectList {
		x1 := int(float32(bounds.Max.X) * item.Box[1])
		y1 := int(float32(bounds.Max.Y) * item.Box[0])
		x2 := int(float32(bounds.Max.X) * item.Box[3])
		y2 := int(float32(bounds.Max.Y) * item.Box[2])

		Rectangle(imgRGBA, x1, x2, y1, y2, 3, colornames.Map[colornames.Names[int(item.Label)]])
		labelinRender := fmt.Sprintf("%s (%2.0f%%)", item.LabelStr, item.Prob*100.0)
		renderLabel(imgRGBA, x1, y1, item.Label, labelinRender)

	}
	err = jpeg.Encode(outputWriter, img, &jpeg.Options{Quality: 75})
	return output.Bytes(), err
}
