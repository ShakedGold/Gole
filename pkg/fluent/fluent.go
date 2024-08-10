package fluent

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	// TODO: Replace with FluentUI icons.
)

// Palette contains the minimal set of colors that a widget may need to
// draw itself.
type Palette struct {
	// Bg is the background color atop which content is currently being
	// drawn.
	Bg color.NRGBA

	// Fg is a color suitable for drawing on top of Bg.
	Fg color.NRGBA

	// DisabledBg is a background color suitable for drawing disabled widgets.
	DisabledBg color.NRGBA

	// DisabledFg is a foreground color suitable for drawing disabled widgets.
	DisabledFg color.NRGBA

	// HoverBg is a color used to draw attention to widgets that are
	// currently being hovered over.
	HoverBg color.NRGBA

	// HoverFg is a color suitable for content drawn on top of HoverBg.
	HoverFg color.NRGBA

	// BrandBg is a color used to draw attention to active,
	// important, interactive widgets such as buttons.
	BrandBg color.NRGBA

	// ContrastFg is a color suitable for content drawn on top of
	// ContrastBg.
	BrandFg color.NRGBA
}

type Font struct {
	font.Font
	Size       unit.Sp
	LineHeight unit.Sp
}

type Fonts struct {
	Caption1 Font
	Caption2 Font

	Body1 Font
	Body2 Font

	Subtitle1 Font
	Subtitle2 Font

	Title1 Font
	Title2 Font
	Title3 Font
}

func NewFonts() Fonts {
	typeface := font.Typeface("Segoe UI")
	return Fonts{
		Caption1: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   100,
			},
			Size:       unit.Sp(12),
			LineHeight: unit.Sp(16),
		},
		Caption2: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   100,
			},
			Size:       unit.Sp(10),
			LineHeight: unit.Sp(14),
		},

		Body1: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   100,
			},
			Size:       unit.Sp(14),
			LineHeight: unit.Sp(20),
		},
		Body2: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   100,
			},
			Size:       unit.Sp(16),
			LineHeight: unit.Sp(22),
		},

		Subtitle1: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   600,
			},
			Size:       unit.Sp(20),
			LineHeight: unit.Sp(28),
		},
		Subtitle2: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   600,
			},
			Size:       unit.Sp(16),
			LineHeight: unit.Sp(22),
		},

		Title1: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   600,
			},
			Size:       unit.Sp(24),
			LineHeight: unit.Sp(32),
		},
		Title2: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   600,
			},
			Size:       unit.Sp(28),
			LineHeight: unit.Sp(36),
		},
		Title3: Font{
			Font: font.Font{
				Typeface: typeface,
				Weight:   600,
			},
			Size:       unit.Sp(32),
			LineHeight: unit.Sp(40),
		},
	}
}

type Theme struct {
	Shaper *text.Shaper
	Palette

	// Fonts contains the fonts for the theme.
	Fonts Fonts

	// FingerSize is the minimum touch target size.
	FingerSize unit.Dp
}

// NewTheme constructs a theme (and underlying text shaper).
func NewTheme() *Theme {
	t := &Theme{Shaper: &text.Shaper{}}
	t.Palette = Palette{
		Fg:         rgb(0x115ea3),
		Bg:         rgb(0xfafafa),
		HoverBg:    rgb(0x115ea3),
		HoverFg:    rgb(0xffffff),
		DisabledBg: rgb(0xf0f0f0),
		DisabledFg: rgb(0xbdbdbd),
		BrandBg:    rgb(0x0f6cbd),
		BrandFg:    rgb(0xffffff),
	}
	t.Fonts = NewFonts()

	// 38dp is on the lower end of possible finger size.
	t.FingerSize = 38

	return t
}

func (t Theme) WithPalette(p Palette) Theme {
	t.Palette = p
	return t
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
