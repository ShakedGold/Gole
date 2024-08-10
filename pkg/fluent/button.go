package fluent

import (
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/input"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Appearance int
type Shape int

const (
	Secondary Appearance = iota
	Primary
	Outline
	Subtle
	Transparent
)

const (
	Circular Shape = iota
	Square
	Rounded
)

type ButtonStyle struct {
	// The buttons text.
	Text string
	// The buttons appearance. (Secondary, Primary, Outline, Subtle, Transparent)
	Appearance Appearance
	// The buttons shape. (Circular, Square, Rounded)
	Shape Shape
	// Icon is the buttons icon, nil if no icon is used.
	Icon *widget.Icon

	Width  unit.Dp
	Height unit.Dp

	gesture     gesture.Click
	state       state
	shaper      *text.Shaper
	font        Font
	color       color.NRGBA
	background  color.NRGBA
	inset       layout.Inset
	iconWidth   unit.Dp
	numOfClicks int
}

type state int

const (
	inactive state = iota
	disabled
	hovered
	clicked
)

// IsClicked returns true if the button is clicked.
func (b *ButtonStyle) IsClicked() bool {
	return b.state == clicked
}

// IsHovered returns true if the button is hovered.
func (b *ButtonStyle) IsHovered() bool {
	return b.state == hovered
}

// IsInactive returns true if the button is inactive.
func (b *ButtonStyle) IsInactive() bool {
	return b.state == inactive
}

func (b *ButtonStyle) IsDisabled() bool {
	return b.state == disabled
}

func (b *ButtonStyle) NumOfClicks() int {
	return b.numOfClicks
}

func (b *ButtonStyle) Update(gtx layout.Context, theme *Theme) {
	gtxDisabled := gtx.Source == (input.Source{})
	b.state = inactive
	b.numOfClicks = 0
	for {
		ev, ok := b.gesture.Update(gtx.Source)
		if !ok {
			break
		}
		switch ev.Kind {
		case gesture.KindClick:
			if !gtxDisabled {
				b.numOfClicks = ev.NumClicks
			}
		}
	}

	if !gtx.Enabled() {
		b.state = disabled
	}
	if b.gesture.Hovered() && !gtxDisabled {
		b.state = hovered
	}
	if b.gesture.Pressed() && !gtxDisabled {
		b.state = clicked
	}
}

func (b *ButtonStyle) getBackground(theme *Theme) (color.NRGBA, color.NRGBA) {
	switch b.Appearance {
	case Secondary, Outline:
		return color.NRGBA{
				R: 0xff,
				G: 0xff,
				B: 0xff,
				A: 0xff,
			},
			color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0xff,
			}
	case Primary:
		return theme.Palette.BrandBg, theme.Palette.BrandFg
	default:
		return color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0x00,
			},
			color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0xff,
			}
	}
}

func (b *ButtonStyle) getHover(theme *Theme) (color.NRGBA, color.NRGBA) {
	switch b.Appearance {
	case Secondary, Subtle:
		return color.NRGBA{
				R: 0xf5,
				G: 0xf5,
				B: 0xf5,
				A: 0xff,
			},
			color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0xff,
			}
	case Outline:
		return color.NRGBA{
				R: 0xff,
				G: 0xff,
				B: 0xff,
				A: 0xff,
			},
			color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0xff,
			}
	case Primary:
		return theme.Palette.HoverBg, theme.Palette.HoverFg
	case Transparent:
		return color.NRGBA{
			R: 0x00,
			G: 0x00,
			B: 0x00,
			A: 0x00,
		}, theme.Palette.BrandBg
	default:
		return color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0x00,
			},
			color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0xff,
			}
	}
}
func (b *ButtonStyle) Layout(gtx layout.Context, theme *Theme) layout.Dimensions {
	semantic.Button.Add(gtx.Ops)
	b.Update(gtx, theme)
	// update b.background based on the Appearance
	b.background, b.color = b.getBackground(theme)
	borderWidth := unit.Dp(0)
	if b.Appearance == Outline || b.Appearance == Secondary {
		borderWidth = unit.Dp(1)
	}
	return widget.Border{
		Color:        color.NRGBA{R: 0xd1, G: 0xd1, B: 0xd1, A: 255}, // Black color
		Width:        borderWidth,                                    // 0.5 dp wide border
		CornerRadius: unit.Dp(4),                                     // Optional: rounded corners
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				cornerRadius := unit.Dp(4)
				if b.Shape == Circular {
					cornerRadius = unit.Dp(50)
				} else if b.Shape == Square {
					cornerRadius = unit.Dp(0)
				}
				rr := gtx.Dp(cornerRadius)

				size := image.Point{
					X: gtx.Dp(b.Width),
					Y: gtx.Dp(b.Height),
				}

				defer clip.UniformRRect(image.Rectangle{Max: size}, rr).Push(gtx.Ops).Pop()

				background := b.background
				switch {
				case b.IsDisabled():
					background = theme.Palette.DisabledBg
				case b.IsHovered():
					background, b.color = b.getHover(theme)
				case b.IsClicked():
					// background = theme.Palette.PressedBg
				}
				paint.Fill(gtx.Ops, background)
				// TODO: Animation
				// for _, c := range b.click.History() {
				// 	drawInk(gtx, c)
				// }
				b.gesture.Add(gtx.Ops)
				return layout.Dimensions{Size: size}
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return b.inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Spacing:   layout.SpaceBetween,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if b.Icon == nil {
									return layout.Dimensions{}
								}
								return layout.Inset{
									Right: unit.Dp(6),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									icon := b.Icon.Layout(gtx, b.color)
									b.iconWidth = unit.Dp(icon.Size.X)
									return icon
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								colMacro := op.Record(gtx.Ops)
								paint.ColorOp{Color: b.color}.Add(gtx.Ops)
								label := widget.Label{Alignment: text.Middle, LineHeight: b.font.LineHeight}.Layout(gtx, b.shaper, b.font.Font, b.font.Size, b.Text, colMacro.Stop())
								b.Width = unit.Dp(label.Size.X+25) + b.iconWidth
								b.Height = unit.Dp(label.Size.Y + 10)
								return label
							}),
						)
					})
				})
			},
		)
	})
}

func Button(th *Theme, txt string, icon []byte) (*ButtonStyle, error) {
	var widgetIcon *widget.Icon
	var err error

	if icon != nil {
		widgetIcon, err = widget.NewIcon(icon)
	}

	return &ButtonStyle{
		Text:       txt,
		Icon:       widgetIcon,
		Appearance: Secondary,
		Width:      unit.Dp(110),
		Height:     unit.Dp(32),
		Shape:      Rounded,
		shaper:     th.Shaper,
		font:       th.Fonts.Body1,
		color:      th.Palette.BrandFg,
		background: th.Palette.BrandBg,
		inset: layout.Inset{
			Top:    unit.Dp(5),
			Right:  unit.Dp(12),
			Bottom: unit.Dp(5),
			Left:   unit.Dp(12),
		},
	}, err
}
