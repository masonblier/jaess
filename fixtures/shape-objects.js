// @src https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/create

// Shape - superclass
function Shape() {
  this.x = -1;
  this.y = 5;
}

// superclass method
Shape.prototype.move = function(x, y) {
    this.x += x;
    this.y += y;
    console.info("Shape moved.");
};

// Rectangle - subclass
function Rectangle() {
  Shape.call(this); //call super constructor.
}

// subclass extends superclass
Rectangle.prototype = Object.create(Shape.prototype);
Rectangle.prototype.constructor = Rectangle;

var rect = new Rectangle();

assert(rect instanceof Rectangle);
assert(rect instanceof Shape);

rect.move(3, -4);

assert(rect.x === 2);
assert(rect.y === 1);
