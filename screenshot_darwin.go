// +build darwin

package gutils

import (
	/*
		#cgo CFLAGS: -x objective-c
		#cgo LDFLAGS: -framework Cocoa -framework Foundation
		#import <Cocoa/Cocoa.h>
		#import <Foundation/Foundation.h>

		NSData *data = 0;
		int width, height, stride; // stride = bytes per row

		void
		Raw() {
			CGImageRef image = CGDisplayCreateImage(kCGDirectMainDisplay);
			width = CGImageGetWidth(image);
			height = CGImageGetHeight(image);
			stride = CGImageGetBytesPerRow(image);

			CFDataRef dataref = CGDataProviderCopyData(
				CGImageGetDataProvider(image)
			);
			data = [NSData dataWithData:(NSData *)dataref];
			CFRelease(dataref);
			CGImageRelease(image);
		}

		int
		Length() {
			return [data length];
		}

		const void*
		Data() {
			return [data bytes];
		}

		void
		Clean () {
			[data release];
		}
	*/
	"C"
	"image"
	"unsafe"
)

// Screenshot returns an image.Image object containing the current screenshot.
func Screenshot() image.Image {
	C.Raw()
	defer C.Clean()

	data := (*[1 << 30]byte)(unsafe.Pointer(C.Data()))[0:C.Length()]
	width := int(C.width)
	height := int(C.height)
	stride := int(C.stride)
	return ConvertMacBGRAToRGBA(width, height, stride, data)
}
