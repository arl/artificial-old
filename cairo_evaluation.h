#include <stdio.h>
#include <cairo/cairo.h>

typedef unsigned int uint32;
typedef int int32;

typedef struct
{
  int32 x, y; 
} point;

typedef struct
{
    unsigned char r, g, b, a;
    point *pts;
    uint32 npts;
} poly;

typedef struct
{
    uint32 w, h;
    uint32 npolys;
    poly *polys;
} imageDNA;

// render an image dna struct into a cairo surface
int render(const imageDNA * dna);
int render_void(const void * dna);

// set `diffval` with a number representing the difference between 2 cairo
// surfaces and return 1. Both surfaces should have the same type, width, height
// and stride, or 0 will is returned and `diffval` is undefined.
int diff_images(cairo_surface_t *cs1, cairo_surface_t *cs2, double* diffval);
