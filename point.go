package game

// Point on plane
type Point struct {
	// X, Y the position
	X, Y int
}

// Add points: p and q
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Minus point
func (p Point) Minus(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul point
func (p Point) Mul(x int) Point {
	return Point{p.X * x, p.Y * x}
}

// Less whether p less than q
func (p Point) Less(q Point) bool {
	return p.X < q.X || (p.X == q.X && p.Y < q.Y)
}
