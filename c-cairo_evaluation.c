#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <cairo/cairo.h>
#include "cairo_evaluation.h"

#define DCLAMP(comp) ((double)comp)/255.

// reference surface
cairo_surface_t *cs_ref;
static int w_ref, h_ref;

/* forward declarations */
int diff_images_rgb24(cairo_surface_t *cs1, cairo_surface_t *cs2, double* diffval);
int draw_dna_to_surface(const imageDNA* dna, int w, int h, cairo_surface_t **cs);

/**
 * destroy dst afterwards
 */ 
int
copy_surface_to_rgb24(cairo_surface_t **dst, cairo_surface_t *org)
{
  int x, y, off;
  unsigned char *dorg, *ddst;
  int w = cairo_image_surface_get_width(org);
  int h = cairo_image_surface_get_height(org);
  int ostride = cairo_image_surface_get_stride(org);
  cairo_format_t fmt = cairo_image_surface_get_format(org);

  switch (fmt) {
    case CAIRO_FORMAT_ARGB32:
        *dst = cairo_image_surface_create(CAIRO_FORMAT_RGB24, w, h);
        int dstride = cairo_image_surface_get_stride(*dst);
        dorg = cairo_image_surface_get_data(org);
        ddst = cairo_image_surface_get_data(*dst);
        cairo_surface_flush(org);
        cairo_surface_flush(*dst);

        for (y = 0; y < h; y++) {
          uint32 *opixel = (uint32 *) (dorg + y * ostride);
          uint32 *dpixel = (uint32 *) (ddst + y * dstride);
          for (x = 0; x < w; x++, opixel++, dpixel++) {
            if ((*opixel & 0xff000000) == 0) {
              // alpha is 0
              *dpixel = 0;
            } else {
              // we can do that because in ARGB32 and RGB24, the R, G, and B
              // components occupy the same byte position
              *dpixel = (*opixel); 
            }
          }
        }
        cairo_surface_mark_dirty(*dst);
        break;

    default:
      printf("%s: unsupported image format fmt=%d\n", __func__, fmt);
  }

}

/**
 * Sets the reference image to in further diffs.
 *
 * The reference image will be converted to an image of size wxh, those has to
 * be the same as the size of the image to compare with the reference.
 */
int
evaluator_init(const char *path,
               uint32 destw, uint32 desth,
               uint32 *orgw, uint32 *orgh)
{
  int rc = CAIRO_STATUS_SUCCESS;
  cairo_surface_t * cs;

  w_ref = (int) destw;
  h_ref = (int) desth;

  cs = cairo_image_surface_create_from_png(path);
  rc = cairo_surface_status(cs);
  if (!cs) {
    return rc;
  }
  if (rc != CAIRO_STATUS_SUCCESS) {
    printf("error, surface status: %s\n", cairo_status_to_string(rc));
    return rc;
  }
  cairo_format_t fmt = cairo_image_surface_get_format(cs);

  switch (fmt) {
    case CAIRO_FORMAT_ARGB32:
      copy_surface_to_rgb24(&cs_ref, cs);
      cairo_surface_destroy(cs);
      cairo_surface_write_to_png(cs_ref, "converted.png");
      break;
    case CAIRO_FORMAT_RGB24:
      cs_ref = cs;
      break;
    case CAIRO_FORMAT_A8:
      printf("image format unsupported: A8\n");
      goto error;
    case CAIRO_FORMAT_A1:
      printf("image format unsupported: A1\n");
      goto error;
    case CAIRO_FORMAT_RGB16_565:
      printf("image format unsupported: RGB16_565\n");
      goto error;
    case CAIRO_FORMAT_RGB30:
      printf("image format unsupported: RGB30\n");
      goto error;
    case CAIRO_FORMAT_INVALID:
      printf("image format: INVALID\n");
      goto error;
  }
 
  return rc;

error:
  cairo_surface_destroy(cs);
  return 1;
}

int
evaluator_deinit() {
  if (cs_ref) {
    cairo_surface_destroy(cs_ref);
  }
  return 0;
}

int
render_and_diff(const imageDNA* dna, double *diffval, const char *dstpath)
{
  int rc = 0;
  const poly *p;
  const point *pt;
  cairo_surface_t *cs = NULL;
  cairo_status_t status;

  printf("render_and_diff dna=0x%p diffval=%f dstpath=%s\n", dna, *diffval, dstpath);

  // render the dna
  rc = draw_dna_to_surface(dna, w_ref, h_ref, &cs);
  if (rc) {
    printf("couldn't draw dna to cairo surface");
    return rc;
  }

  if (dstpath) {
    status = cairo_surface_write_to_png(cs, dstpath);
    if (status != CAIRO_STATUS_SUCCESS)
       printf("cairo error: %s\n", cairo_status_to_string(status));
  }

  rc = diff_images_rgb24(cs_ref, cs, diffval);

  // cleanup
  cairo_surface_destroy(cs);
  return rc;
}

/* Create a cairo surface and draw a DNA structure onto it. It's the caller
 * responsibility to call cairo_surface_destroy on the surface after use.
 */
int
draw_dna_to_surface(const imageDNA* dna, int w, int h, cairo_surface_t **cs)
{
  int rc = 0;
  const poly *p;
  const point *pt;
  cairo_t *cr;
  cairo_status_t status;
  char fn[256];

  if (!cs) {
    printf("draw_dna_to_surface: invalid parameter `cs`\n");
    return -1;
  }

  // init
  *cs = cairo_image_surface_create(CAIRO_FORMAT_RGB24, w, h);
  cr = cairo_create (*cs);
  for (int i=0; i < dna->npolys; ++i)
  {
    p = &dna->polys[i];
    printf("drawing poly[%d] npts=%d\n", i, p->npts);
    cairo_set_source_rgba(cr, DCLAMP(p->r), DCLAMP(p->g), DCLAMP(p->b), DCLAMP(p->a));
    cairo_set_line_width(cr, 0);
    for (int j=0; j < p->npts; ++j)
    {
      pt = &p->pts[j];
      printf("drawing pts[%d] x=%d y=%d\n", i, pt->x, pt->y);
      cairo_line_to(cr, pt->x, pt->y);
    }
    cairo_close_path(cr);
    cairo_stroke_preserve(cr);
    cairo_fill(cr);
  }

  // cleanup
  cairo_destroy(cr);
  return rc;
}


// set `diffval` with a number representing the difference between 2 cairo
// surfaces and return 1. Both surfaces should have the same type, width, height
// and stride, or 0 will is returned and `diffval` is undefined.
int
diff_images_rgb24(cairo_surface_t *cs1, cairo_surface_t *cs2, double* diffval)
{
        unsigned char *buf1, *buf2;
        int width, height, stride;
        int x, y, off;
        long diff = 0;
        cairo_surface_type_t type;

        // flush to ensure all writing to the images are done
        cairo_surface_flush(cs1);
        cairo_surface_flush(cs2);

        // check images type
        type = cairo_surface_get_type(cs1);
        if (type != cairo_surface_get_type(cs2))
      	  return 0;

        // check images format
        if (cairo_image_surface_get_format(cs1) != CAIRO_FORMAT_RGB24 ||
            cairo_image_surface_get_format(cs1) != CAIRO_FORMAT_RGB24)
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
      				  abs((int)buf1[off+2]-(int)buf2[off+2]);
      	  }
        }
        *diffval = (double)diff;
        return 1;
}

// set `diffval` with a number representing the difference between 2 cairo
// surfaces and return 1. Both surfaces should have the same type, width, height
// and stride, or 0 will is returned and `diffval` is undefined.
int
diff_images_argb32(cairo_surface_t *cs1, cairo_surface_t *cs2, double* diffval)
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
