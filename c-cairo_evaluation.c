#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <cairo/cairo.h>
#include "cairo_evaluation.h"

#define DCLAMP(comp) ((double)comp)/255.

// reference surface
cairo_surface_t *cs_ref;

/**
 * Sets the reference image to in further diffs.
 *
 * The reference image will be converted to an image of size wxh, those has to
 * be the same as the size of the image to compare with the reference.
 */
int
evaluator_init(const char* path, int w, int h)
{
  int rc = CAIRO_STATUS_SUCCESS;
  cairo_surface_t * cs;
  cs = cairo_image_surface_create_from_png(path);
  if (cs == NULL) {
    rc = cairo_surface_status(cs);
    return rc;
  }
  if (cairo_image_surface_get_format(cs) != CAIRO_FORMAT_ARGB32) {
    printf("must convert image to ARGB32 for diff, not implemented!");
  }
  cs_ref = cs;
  return rc;
}

int
evaluator_deinit() {
  if (cs_ref) {
    cairo_surface_destroy(cs_ref);
  }
  return 0;
}

int
render(const imageDNA* dna)
{
  const poly *p;
  const point *pt;
  cairo_t *cr;
  cairo_surface_t *cs;
  cairo_status_t status;
  char fn[256];

  // init
  cs = cairo_image_surface_create (CAIRO_FORMAT_ARGB32, dna->w, dna->h);
  cr = cairo_create (cs);
  for (int i=0; i < dna->npolys; ++i)
  {
    p = &dna->polys[i];
    cairo_set_source_rgba(cr, DCLAMP(p->r), DCLAMP(p->g), DCLAMP(p->b), DCLAMP(p->a));
    cairo_set_line_width(cr, 0);
    for (int j=0; j < p->npts; ++j)
    {
      pt = &p->pts[j];
      cairo_line_to(cr, pt->x, pt->y);
    }
    cairo_close_path(cr);
    cairo_stroke_preserve(cr);
    cairo_fill(cr);
  }

  // write
  /*status = cairo_surface_write_to_png(cs, "c.render.png");*/
  /*sprintf(fn, "c.render%d.png", rnd++);*/
  /*status = cairo_surface_write_to_png(cs, fn);*/

  if (status != CAIRO_STATUS_SUCCESS)
     printf("cairo error: %s\n", cairo_status_to_string(status));

  // cleanup
  cairo_surface_destroy(cs);
  cairo_destroy(cr);
}

// set `diffval` with a number representing the difference between 2 cairo
// surfaces and return 1. Both surfaces should have the same type, width, height
// and stride, or 0 will is returned and `diffval` is undefined.
int
diff_images(cairo_surface_t *cs1, cairo_surface_t *cs2, double* diffval)
{
        unsigned char *buf1, *buf2;
        int width, height, stride;
        int x, y, off;
        long diff;
        cairo_surface_type_t type;

        // flush to ensure all writing to the images are done
        cairo_surface_flush(cs1);
        cairo_surface_flush(cs2);

        // check images type
        type = cairo_surface_get_type(cs1);
        if (type != cairo_surface_get_type(cs2))
      	  return 0;

        // check images format
        if (cairo_image_surface_get_format(cs1) != CAIRO_FORMAT_ARGB32 ||
            cairo_image_surface_get_format(cs1) != CAIRO_FORMAT_ARGB32)
      	  return 0;

        // check images width
        width = cairo_image_surface_get_width(cs1);
        if (width != cairo_image_surface_get_width(cs2))
      	  return 0;

        // check images height
        height = cairo_image_surface_get_height(cs1);
        if (height != cairo_image_surface_get_height(cs2))
      	  return 0;
        //
        // check images stride
        stride = cairo_image_surface_get_stride(cs1);
        if (stride != cairo_image_surface_get_stride(cs2))
      	  return 0;

        buf1 = cairo_image_surface_get_data (cs1);
        buf2 = cairo_image_surface_get_data (cs2);

        for (y = 0; y < height; ++y)
        {
      	  for (x = 0; x < width; ++x)
      	  {
      		  off = y*stride + x*4;
      		  diff += abs((int)buf1[off+0]-(int)buf2[off+0]) +
      				  abs((int)buf1[off+1]-(int)buf2[off+1]) +
      				  abs((int)buf2[off+2]-(int)buf2[off+2]) +
      				  abs((int)buf1[off+3]-(int)buf2[off+3]);
      	  }
        }
        *diffval = (double)diff;
        return 1;
}

// int main(int argc, char *argv[])
// {
//         cairo_t *cr;
//         cairo_surface_t *cs;
//         cairo_status_t status;
//         int exit = EXIT_FAILURE;
//
//         // init
//         cs = cairo_image_surface_create (CAIRO_FORMAT_ARGB32, 512, 512);
//         cr = cairo_create (cs);
//
//         // draw
//         do_drawing(cr);
//
//         // write
//         status = cairo_surface_write_to_png(cs, "cairo.png");
//         if (status != CAIRO_STATUS_SUCCESS)
//       	  printf("cairo error: %s\n", cairo_status_to_string(status));
//         else
//       	  exit = EXIT_SUCCESS;
//
//         // test diff
//         double diff;
//         if (diff_images(cs, cs, &diff) == 0)
//         {
//       	  printf("diff error");
//         }
//         else
//         {
//       	  printf("diff: %f\n", diff);
//       	  exit = EXIT_SUCCESS;
//         }
//
//         // cleanup
//         cairo_surface_destroy(cs);
//         cairo_destroy(cr);
//         return exit;
// }

