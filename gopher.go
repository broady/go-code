package main

import (
  "image"
  "tour/pic"
  "math"
)

type Image struct {}

type Point struct {
  x, y int
}

func (i Image) ColorModel() image.ColorModel {
  return image.RGBAColorModel
}

func (i Image) Bounds() image.Rectangle {
    return image.Rect(0, 0, 512, 512)
}

var (
  WHITE = &image.RGBAColor{ 255, 255, 255, 255 }
  TRANSPARENT = &image.RGBAColor { 0,0,0,0}
  BLACK = &image.RGBAColor{ 0, 0, 0, 255 }
  CENTER = Point {256, 256}
  LEFTEYE = Point {160, 109}
  LEFTPUPIL = Point {188, 97}
  RIGHTEYE = Point {255, 56}
  RIGHTPUPIL = Point {285, 49}
  NOSE = Point {228, 124}
  LEFTEAR = Point {102, 142}
  RIGHTEAR = Point {329, 41}
  shapes = []Shape{
    Circle{LEFTPUPIL, 14},
    Circle{RIGHTPUPIL, 14},
    StrokedShape{Circle{LEFTEYE, 45}, 2},
    StrokedShape{Circle{RIGHTEYE, 45}, 2},
/*Nose*/
    RotatedShape{Ellipse{NOSE, 22, .3}, 0.5, NOSE},
    StrokedShape{
      CompositeShape{
        RotatedShape{Ellipse{Point{213,144}, 20, .4}, 0.5, NOSE},
        RotatedShape{Ellipse{Point{243,144}, 20, .4}, 0.5, NOSE},
      },
    2},
/*Teeth*/
    StrokedShape{
      CompositeShape{
        RotatedShape{Rectangle{Point{213,144}, Point{227, 185}}, 0.5, NOSE},
        RotatedShape{Rectangle{Point{230,144}, Point{244, 185}}, 0.5, NOSE},
      },
    2},
/*Body*/
    StrokedShape{
      RotatedShape {
        CompositeShape{
          Ellipse{Point{237, 435}, 95, .25},
          SubtractedShape{
            Rectangle{Point{75, 182}, Point{261, 310}},
            Ellipse{Point{100, 238}, 60, 2.2},
          },
          RotatedShape{
            Ellipse{Point{334, 325}, 85, 3}, -.4, Point{334,325},
          },
          Ellipse{Point{195, 379}, 90, 1},
          Ellipse{Point{294, 70}, 100, .3},
          Ellipse{Point{232, 162}, 87, 2.2},
          Ellipse{Point{349, 162}, 87, 2.2},
      }, .46, CENTER,
     }, 2},
/* tail */
     StrokedShape {
       Circle{Point{175, 438}, 7},
     2},
/* ears */
     Ellipse{LEFTEAR, 18, .7},
     StrokedShape{
       Ellipse{LEFTEAR, 40, .8},
     2},
     Ellipse{RIGHTEAR, 20, .6},
     StrokedShape{
       Ellipse{RIGHTEAR, 40, .8},
     2},
/* rear foot */
     
  }
)

type SubtractedShape struct {
  base, subtracted Shape
}

func (ss SubtractedShape) Color(x, y int) image.Color {
  c := ss.base.Color(x,y)
  if nil != c && nil == ss.subtracted.Color(x,y) {
    return c
  }
  return nil
}

type CompositeShape []Shape

func (cs CompositeShape) Color(x, y int) image.Color {
  for _, shape := range cs {
    c := shape.Color(x, y)
    if nil != c {
      return c
    }
  }
  return nil
}

type Rectangle struct {
  nw, se  Point
}

func (r Rectangle) Color(x, y int) image.Color {
  if x >= r.nw.x && x <= r.se.x && y >= r.nw.y && y <= r.se.y {
    return BLACK
  }

  return nil
}

type RotatedShape struct {
  s Shape
  r float64
  c Point
}

func (rs RotatedShape) Color(x, y int) image.Color {
  x -= rs.c.x
  y -= rs.c.y
  xr := int(float64(x) * math.Cos(rs.r) - float64(y) * math.Sin(rs.r))
  yr := int(float64(x) * math.Sin(rs.r) + float64(y) * math.Cos(rs.r))
  xr += rs.c.x
  yr += rs.c.y
  return rs.s.Color(xr, yr)
}

type StrokedShape struct {
  i       Shape
  stroke  int
}

func (s StrokedShape) Color(x, y int) image.Color {
  if nil != s.i.Color(x, y) {
    return WHITE
  }
    
  for i := x - s.stroke; i < x + 1 + s.stroke; i++ {
    for j := y - s.stroke; j < y + 1 + s.stroke; j++ {
      var c = s.i.Color(i, j)
      if nil != c {
        return c
      }
    }
  }
  return nil
}

type Circle struct {
  l Point
  r int
}

func (c Circle) Color(x, y int) image.Color {
  xa := c.l.x - x
  yb := c.l.y - y
  if (xa * xa) + (yb * yb) < c.r * c.r {
    return BLACK
  }
  return nil
}

type Ellipse struct {
  l   Point
  r   int
  sy  float64
}

func (e Ellipse) Color(x, y int) image.Color {
  xa := e.l.x - x
  yb := e.l.y - y
  if (float64(xa * xa) / 1) + (float64(yb * yb) / e.sy) < float64(e.r * e.r) {
    return BLACK
  }
  return nil
}

type Shape interface {
  Color(x, y int) image.Color
}

func (i Image) At(x, y int) image.Color {
  for _, Shape := range shapes {
    var c = Shape.Color(x, y)
    if nil != c {
      return c
    }
  }
  return TRANSPARENT
}

func myImage(dx, dy int) image.Image {
  return Image{}
}

func main() {
  pic.ServeImage(myImage)
}
