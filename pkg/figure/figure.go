package figure

import (
	"github.com/common-nighthawk/go-figure"
)

func Draw(word, font string) {
	myFigure := figure.NewFigure(word, font, true)
	myFigure.Print()
}
