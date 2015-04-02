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

		void FreeIfRequired();

		void
		Raw() {
			FreeIfRequired();

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

		long
		DataPtr() {
			return (long)data;
		}

		void
		Free (long ptr) {
			NSData *dat = (NSData*)ptr;
			[dat release];
		}

		void
		FreeIfRequired() {
			if (data) {
				[data release];
				data = 0;
			}
		}
	*/
	"C"
	"image"
	// "runtime"
	"unsafe"
)

// CleaningUpMacBGRA keeps track of underlying C data pointer, and cleans up
// using runtime.SetFinalizer
type CleaningUpMacBGRA struct {
	MacBGRA
	dataPtr uint64
}

// CleanUpCleaningUpMacBGRA cleans up a CleaningUpMacBGRA. Can also be safely
// called when you are done with CleaningUpMacBGRA, to reclaim memory sooner
// than waiting for GC to clean up.
func CleanUpCleaningUpMacBGRA(m *CleaningUpMacBGRA) {
	if m.dataPtr != 0 {
		C.Free(C.long(m.dataPtr))
		m.dataPtr = 0
	}
}

// Screenshot returns an image.Image object containing the current screenshot.
// Currently we must.
func Screenshot() image.Image {
	C.Raw()
	data := (*[1 << 30]byte)(unsafe.Pointer(C.Data()))[0:C.Length()]
	m := &CleaningUpMacBGRA{
		MacBGRA: *NewMacBGRA(
			image.Rect(0, 0, int(C.width), int(C.height)), int(C.stride), data,
		),
		dataPtr: uint64(C.DataPtr()),
	}
	// runtime.SetFinalizer(m, CleanUpCleaningUpMacBGRA)
	return m
}
